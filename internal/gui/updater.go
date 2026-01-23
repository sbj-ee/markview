package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	githubRepoOwner = "sbj-ee"
	githubRepoName  = "markview"
	githubAPIURL    = "https://api.github.com/repos/%s/%s/releases/latest"
	updateCheckKey  = "lastUpdateCheck"
	skippedVersion  = "skippedVersion"
	checkInterval   = 24 * time.Hour // Check once per day
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
}

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
				dialog.ShowError(fmt.Errorf("Failed to check for updates: %v", err), uc.window)
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
			dialog.ShowInformation("Up to Date",
				fmt.Sprintf("You're running the latest version (%s).", uc.currentVersion),
				uc.window)
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

// showUpdateDialog displays the update available dialog
func (uc *UpdateChecker) showUpdateDialog(release *GitHubRelease) {
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	// Create release notes (truncate if too long)
	releaseNotes := release.Body
	if len(releaseNotes) > 500 {
		releaseNotes = releaseNotes[:500] + "..."
	}

	// Build content
	titleLabel := widget.NewLabel(fmt.Sprintf("A new version of MarkView is available!"))
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	versionLabel := widget.NewLabel(fmt.Sprintf("Current: %s  →  Latest: %s",
		uc.currentVersion, latestVersion))

	notesLabel := widget.NewLabel("Release Notes:")
	notesLabel.TextStyle = fyne.TextStyle{Bold: true}

	notesText := widget.NewLabel(releaseNotes)
	notesText.Wrapping = fyne.TextWrapWord

	notesScroll := container.NewVScroll(notesText)
	notesScroll.SetMinSize(fyne.NewSize(400, 150))

	content := container.NewVBox(
		titleLabel,
		versionLabel,
		widget.NewSeparator(),
		notesLabel,
		notesScroll,
	)

	// Find download URL for current platform and architecture
	downloadURL := release.HTMLURL
	var universalURL string

	// Determine which architecture suffix to look for
	archSuffix := ""
	switch runtime.GOARCH {
	case "amd64":
		archSuffix = "x86_64"
	case "arm64":
		archSuffix = "arm64"
	}

	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if !strings.HasSuffix(name, ".dmg") {
			continue
		}

		// Check for architecture-specific DMG
		if archSuffix != "" && strings.Contains(name, archSuffix) {
			downloadURL = asset.BrowserDownloadURL
			break
		}

		// Track universal DMG (no arch suffix) as fallback
		if !strings.Contains(name, "arm64") && !strings.Contains(name, "x86_64") && !strings.Contains(name, "amd64") {
			universalURL = asset.BrowserDownloadURL
		}
	}

	// Use universal DMG if no arch-specific one was found
	if downloadURL == release.HTMLURL && universalURL != "" {
		downloadURL = universalURL
	}

	// Create custom dialog with buttons
	d := dialog.NewCustomConfirm(
		"Update Available",
		"Download",
		"Later",
		content,
		func(download bool) {
			if download {
				// Open download URL in browser
				if err := openURL(downloadURL); err != nil {
					dialog.ShowError(err, uc.window)
				}
			}
		},
		uc.window,
	)

	// Add "Skip This Version" button by wrapping the dialog
	skipBtn := widget.NewButton("Skip This Version", func() {
		uc.app.Preferences().SetString(skippedVersion, latestVersion)
		d.Hide()
	})

	d.Show()

	// Show skip button in a separate small dialog after main dialog
	// This is a workaround since Fyne dialogs don't support 3 buttons easily
	go func() {
		time.Sleep(100 * time.Millisecond)
		// Add skip functionality via preferences menu instead
	}()
	_ = skipBtn // Reserved for future use
}

// openURL opens a URL in the default browser
func openURL(urlStr string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	return fyne.CurrentApp().OpenURL(parsed)
}
