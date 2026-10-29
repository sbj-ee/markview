package gui

import (
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/sbj-ee/markview/internal/update"
)

func sampleRelease(withSigned bool) *GitHubRelease {
	r := &GitHubRelease{TagName: "v1.1.0", HTMLURL: "https://github.com/sbj-ee/markview/releases/tag/v1.1.0",
		Body: "### Added\n\n- Help > Check for Updates installs minisign-verified packages\n- Check for Updates at Startup toggle"}
	add := func(name string) {
		r.Assets = append(r.Assets, struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{name, "https://github.com/sbj-ee/markview/releases/download/v1.1.0/" + name})
	}
	add(update.PackageName("1.1.0", runtime.GOOS, runtime.GOARCH))
	if withSigned {
		add("SHA256SUMS")
		add("SHA256SUMS.minisig")
	}
	return r
}

// The dialog offers Install only for a signed release; optionally saves a
// screenshot (MARKVIEW_SCREENSHOT_DIR) rendered with Fyne's test driver.
func TestUpdateDialog(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	w := a.NewWindow("MarkView")
	w.Resize(fyne.NewSize(640, 520))
	uc := NewUpdateChecker("1.0.22", w, a)

	d := uc.buildUpdateDialog(sampleRelease(true))
	d.Show()
	if dir := os.Getenv("MARKVIEW_SCREENSHOT_DIR"); dir != "" {
		f, err := os.Create(filepath.Join(dir, "update-dialog.png"))
		if err != nil {
			t.Fatal(err)
		}
		png.Encode(f, w.Canvas().Capture())
		f.Close()
	}
	d.Hide()

	unsigned := uc.buildUpdateDialog(sampleRelease(false))
	unsigned.Show()
	if dir := os.Getenv("MARKVIEW_SCREENSHOT_DIR"); dir != "" {
		f, _ := os.Create(filepath.Join(dir, "update-dialog-unsigned.png"))
		png.Encode(f, w.Canvas().Capture())
		f.Close()
	}
	unsigned.Hide()
}

func TestCheckAtStartupPreference(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()
	if !CheckAtStartupEnabled(a) {
		t.Fatal("default should be on")
	}
	SetCheckAtStartupEnabled(a, false)
	if CheckAtStartupEnabled(a) {
		t.Fatal("toggle off not stored")
	}
}

// Renders the Help menu's update items (Fyne's own menu widget) for the PR.
func TestHelpMenuUpdateItemsScreenshot(t *testing.T) {
	dir := os.Getenv("MARKVIEW_SCREENSHOT_DIR")
	if dir == "" {
		t.Skip("set MARKVIEW_SCREENSHOT_DIR to render")
	}
	a := test.NewApp()
	defer test.NewApp()
	w := a.NewWindow("MarkView")
	mw := &Window{app: a, fyneWindow: w}
	m := fyne.NewMenu("Help",
		fyne.NewMenuItem("Keyboard Shortcuts", func() {}),
		fyne.NewMenuItem("Help", func() {}),
		fyne.NewMenuItem("Check for Updates…", func() {}),
		mw.checkAtStartupMenuItem(),
		fyne.NewMenuItem("About MarkView", func() {}))
	menu := widget.NewMenu(m)
	w.SetContent(container.NewWithoutLayout(menu))
	menu.Resize(menu.MinSize())
	w.Resize(menu.MinSize())
	f, _ := os.Create(filepath.Join(dir, "help-menu.png"))
	png.Encode(f, w.Canvas().Capture())
	f.Close()
}
