package gui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/sbj-ee/markview/internal/update"
)

const (
	githubRepoOwner = "sbj-ee"
	githubRepoName  = "markview"
	githubAPIURL    = "https://api.github.com/repos/%s/%s/releases/latest"
	updateCheckKey  = "lastUpdateCheck"
	skippedVersion  = "skippedVersion"
	checkInterval   = 24 * time.Hour // Check once per day
	// checkAtStartupKey holds the Help > Check for Updates at Startup toggle.
	checkAtStartupKey = "checkUpdatesAtStartup"
)

// GitHubRelease represents the GitHub API response for a release
type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpdateChecker handles checking for application updates
type UpdateChecker struct {
	currentVersion string
	window         fyne.Window
	app            fyne.App

	// RestartState, when set, reports the file to reopen after a restart
	// and whether there are unsaved changes that a restart would lose.
	RestartState func() (file string, dirty bool)
}

// installedBinary is where the .deb installs MarkView. Linux updates only
// install over this copy, and a restart starts it (the new version).
const installedBinary = "/usr/bin/markview"

// NewUpdateChecker creates a new update checker
func NewUpdateChecker(currentVersion string, window fyne.Window, app fyne.App) *UpdateChecker {
	return &UpdateChecker{
		currentVersion: currentVersion,
		window:         window,
		app:            app,
	}
}

// CheckForUpdates checks GitHub for a newer version
func (uc *UpdateChecker) CheckForUpdates(silent bool) {
	// In silent mode, check if we should skip this check
	if silent {
		lastCheck := uc.app.Preferences().String(updateCheckKey)
		if lastCheck != "" {
			lastTime, err := time.Parse(time.RFC3339, lastCheck)
			if err == nil && time.Since(lastTime) < checkInterval {
				return // Skip check, too recent
			}
		}
	}

	go func() {
		release, err := uc.fetchLatestRelease()
		if err != nil {
			if !silent {
				fyne.Do(func() { dialog.ShowError(fmt.Errorf("failed to check for updates: %v", err), uc.window) })
			}
			return
		}

		// Save check time
		uc.app.Preferences().SetString(updateCheckKey, time.Now().Format(time.RFC3339))

		latestVersion := strings.TrimPrefix(release.TagName, "v")
		currentVersion := strings.TrimPrefix(uc.currentVersion, "v")

		if uc.isNewerVersion(latestVersion, currentVersion) {
			// Check if user skipped this version
			if silent {
				skipped := uc.app.Preferences().String(skippedVersion)
				if skipped == latestVersion {
					return // User skipped this version
				}
			}
			uc.showUpdateDialog(release)
		} else if !silent {
			fyne.Do(func() {
				dialog.ShowInformation("Up to Date",
					fmt.Sprintf("You're running the latest version (%s).", uc.currentVersion),
					uc.window)
			})
		}
	}()
}

// fetchLatestRelease fetches the latest release from GitHub API
func (uc *UpdateChecker) fetchLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf(githubAPIURL, githubRepoOwner, githubRepoName)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "MarkView-UpdateChecker")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found - you may be running a development version")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// isNewerVersion compares two semantic versions
// Returns true if latest > current
func (uc *UpdateChecker) isNewerVersion(latest, current string) bool {
	latestParts := parseVersion(latest)
	currentParts := parseVersion(current)

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	// If all compared parts are equal, longer version is newer
	return len(latestParts) > len(currentParts)
}

// parseVersion parses a version string into integer parts
func parseVersion(v string) []int {
	// Remove any prefix like "v"
	v = strings.TrimPrefix(v, "v")

	// Split by dots and dashes
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == '.' || r == '-'
	})

	result := make([]int, 0, len(parts))
	for _, p := range parts {
		// Try to parse as integer, skip non-numeric parts
		if n, err := strconv.Atoi(p); err == nil {
			result = append(result, n)
		}
	}

	return result
}

// showUpdateDialog displays the update available dialog: release notes and
// Install / Skip This Version / Later. Install follows zwriter's flow:
// download SHA256SUMS + SHA256SUMS.minisig, verify the minisign signature
// with the embedded public key, download the package, check its SHA-256,
// then install it (Linux: pkexec apt) or open it (macOS: the .dmg).
func (uc *UpdateChecker) showUpdateDialog(release *GitHubRelease) {
	fyne.Do(func() { uc.buildUpdateDialog(release).Show() })
}

func (uc *UpdateChecker) buildUpdateDialog(release *GitHubRelease) dialog.Dialog {
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	releaseNotes := release.Body
	if len(releaseNotes) > 1500 {
		releaseNotes = releaseNotes[:1500] + "..."
	}

	titleLabel := widget.NewLabel(fmt.Sprintf("MarkView %s is available", latestVersion))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	versionLabel := widget.NewLabel(fmt.Sprintf("You have %s.", uc.currentVersion))
	notes := widget.NewRichTextFromMarkdown(releaseNotes)
	notes.Wrapping = fyne.TextWrapWord
	notesScroll := container.NewVScroll(notes)
	notesScroll.SetMinSize(fyne.NewSize(460, 180))

	pkgName := update.PackageName(latestVersion, runtime.GOOS, runtime.GOARCH)
	rel := toUpdateRelease(release)
	_, hasPkg := rel.Find(pkgName)
	_, hasSums := rel.Find("SHA256SUMS")
	_, hasSig := rel.Find("SHA256SUMS.minisig")
	how := fmt.Sprintf("Install downloads %s and the signed SHA256SUMS, checks the signature and the SHA-256, then ", pkgName)
	if runtime.GOOS == "darwin" {
		how += "replaces MarkView.app in place."
	} else {
		how += "runs pkexec apt install (you'll be asked for your password)."
	}
	canInstall := pkgName != "" && hasPkg && hasSums && hasSig
	if !canInstall {
		how = "This release has no signed package for this system, so it can't be installed from here. Use the release page."
	}
	howLabel := widget.NewLabel(how)
	howLabel.Wrapping = fyne.TextWrapWord

	var d dialog.Dialog
	install := widget.NewButton("Install", func() {
		d.Hide()
		uc.installRelease(rel, pkgName)
	})
	install.Importance = widget.HighImportance
	if !canInstall {
		install.Disable()
	}
	skip := widget.NewButton("Skip This Version", func() {
		uc.app.Preferences().SetString(skippedVersion, latestVersion)
		d.Hide()
	})
	later := widget.NewButton("Later", func() { d.Hide() })
	page := widget.NewHyperlink("Release page", nil)
	if u, err := url.Parse(release.HTMLURL); err == nil {
		page.SetURL(u)
	}

	content := container.NewBorder(
		container.NewVBox(titleLabel, container.NewHBox(versionLabel, page), widget.NewSeparator()),
		container.NewVBox(widget.NewSeparator(), howLabel, container.NewHBox(layout.NewSpacer(), install, skip, later)),
		nil, nil, notesScroll)
	d = dialog.NewCustomWithoutButtons("Update Available", content, uc.window)
	d.Resize(fyne.NewSize(560, 440))
	return d
}

func toUpdateRelease(r *GitHubRelease) *update.Release {
	out := &update.Release{TagName: r.TagName, Name: r.Name, Body: r.Body, HTMLURL: r.HTMLURL, PublishedAt: r.PublishedAt}
	for _, a := range r.Assets {
		out.Assets = append(out.Assets, update.Asset{Name: a.Name, URL: a.BrowserDownloadURL})
	}
	return out
}

// installRelease downloads and verifies the package, then installs or opens it.
func (uc *UpdateChecker) installRelease(rel *update.Release, pkgName string) {
	bar := widget.NewProgressBar()
	status := widget.NewLabel("Downloading and verifying " + pkgName + "…")
	prog := dialog.NewCustomWithoutButtons("Updating MarkView", container.NewVBox(status, bar), uc.window)
	prog.Show()
	go func() {
		dir, err := os.MkdirTemp("", "markview-update-")
		if err != nil {
			fyne.Do(func() { prog.Hide(); dialog.ShowError(err, uc.window) })
			return
		}
		c := &update.Client{
			HTTP:      &http.Client{Timeout: 10 * time.Minute},
			UserAgent: "MarkView/" + uc.currentVersion,
			OnProgress: func(_ string, got, total int64) {
				if total > 0 {
					fyne.Do(func() { bar.SetValue(float64(got) / float64(total)) })
				}
			},
		}
		path, sum, err := c.DownloadVerifiedSum(context.Background(), rel, pkgName, dir)
		if err != nil {
			os.RemoveAll(dir)
			fyne.Do(func() {
				prog.Hide()
				dialog.ShowError(fmt.Errorf("the update was not installed: %v", err), uc.window)
			})
			return
		}
		// Hand the installer a private copy, re-verified against the signed
		// hash at the last moment.
		stageDir, err := os.MkdirTemp("", "markview-install-")
		if err == nil {
			var staged string
			if staged, err = update.StageVerified(path, sum, stageDir); err == nil {
				path = staged
			} else {
				os.RemoveAll(stageDir)
			}
		}
		if err != nil {
			os.RemoveAll(dir)
			fyne.Do(func() {
				prog.Hide()
				dialog.ShowError(fmt.Errorf("the update was not installed: %v", err), uc.window)
			})
			return
		}
		fyne.Do(func() { status.SetText("Verified. Installing…") })
		err = uc.install(path, stageDir)
		fyne.Do(func() {
			prog.Hide()
			if err != nil {
				uc.fallback(err, path)
				return
			}
			os.RemoveAll(dir)
			os.RemoveAll(stageDir)
			uc.offerRestart(rel.Version())
		})
	}()
}

// install puts the verified package in place: on macOS by replacing the
// running MarkView.app from the disk image, on Linux with pkexec apt.
func (uc *UpdateChecker) install(pkg, workDir string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if resolved, rerr := filepath.EvalSymlinks(exe); rerr == nil {
		exe = resolved
	}
	if runtime.GOOS == "darwin" {
		bundle := update.RunningBundle(exe)
		if bundle == "" {
			return errors.New("this MarkView isn't running from an app bundle, so it can't replace itself")
		}
		return update.ReplaceBundleFromDMG(pkg, bundle, workDir, nil)
	}
	if exe != installedBinary {
		return fmt.Errorf("this MarkView (%s) wasn't installed from the .deb, so the package can't be installed in its place automatically", exe)
	}
	cmdline := update.InstallCommand(runtime.GOOS, pkg)
	if cmdline == nil {
		return fmt.Errorf("installing updates isn't supported on %s", runtime.GOOS)
	}
	if out, err := exec.Command(cmdline[0], cmdline[1:]...).CombinedOutput(); err != nil {
		return fmt.Errorf("the installer failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// fallback explains why the update wasn't installed automatically and keeps
// the verified package so it can be installed by hand.
func (uc *UpdateChecker) fallback(why error, pkg string) {
	how := "Install it with: sudo apt install " + pkg
	if runtime.GOOS == "darwin" {
		how = "Open it and drag MarkView.app to Applications."
		_ = exec.Command("/usr/bin/open", pkg).Start()
	}
	dialog.ShowError(fmt.Errorf("%v\n\nThe verified package is at:\n%s\n\n%s", why, pkg, how), uc.window)
}

// offerRestart asks whether to restart into the newly installed version.
func (uc *UpdateChecker) offerRestart(version string) {
	file, dirty := "", false
	if uc.RestartState != nil {
		file, dirty = uc.RestartState()
	}
	text := "MarkView " + version + " was installed. Restart MarkView now to use it?"
	if file != "" {
		text += "\n\n" + filepath.Base(file) + " will be reopened."
	}
	if dirty {
		text += "\n\nWarning: you have unsaved changes, and restarting now will lose them. Choose Later and save first to keep them."
	}
	msg := widget.NewLabel(text)
	msg.Wrapping = fyne.TextWrapWord
	var d dialog.Dialog
	restart := widget.NewButton("Restart Now", func() {
		d.Hide()
		uc.restart(file)
	})
	later := widget.NewButton("Later", func() { d.Hide() })
	if dirty {
		later.Importance = widget.HighImportance
	} else {
		restart.Importance = widget.HighImportance
	}
	d = dialog.NewCustomWithoutButtons("Restart MarkView",
		container.NewVBox(msg, container.NewHBox(layout.NewSpacer(), restart, later)), uc.window)
	d.Resize(fyne.NewSize(460, 200))
	d.Show()
}

// restart starts the installed MarkView (reopening file, if any) and quits.
func (uc *UpdateChecker) restart(file string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		exe, _ := os.Executable()
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}
		args := []string{"-n", update.RunningBundle(exe)}
		if file != "" {
			args = append(args, "--args", file)
		}
		cmd = exec.Command("/usr/bin/open", args...)
	} else {
		var args []string
		if file != "" {
			args = append(args, file)
		}
		cmd = exec.Command(installedBinary, args...)
	}
	if err := cmd.Start(); err != nil {
		dialog.ShowError(fmt.Errorf("could not restart MarkView: %v\nQuit and reopen it to use the new version", err), uc.window)
		return
	}
	_ = cmd.Process.Release()
	uc.app.Quit()
}

// CheckAtStartupEnabled reports the Help > Check for Updates at Startup choice (default on).
func CheckAtStartupEnabled(app fyne.App) bool {
	return app.Preferences().BoolWithFallback(checkAtStartupKey, true)
}

// SetCheckAtStartupEnabled stores the Help > Check for Updates at Startup choice.
func SetCheckAtStartupEnabled(app fyne.App, on bool) {
	app.Preferences().SetBool(checkAtStartupKey, on)
}

// openURL opens a URL in the default browser
func openURL(urlStr string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	return fyne.CurrentApp().OpenURL(parsed)
}
