package gui

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"go.uber.org/zap"
)

// TestScrollBehaviorOnReload verifies the scroll position is preserved when the
// same document is reloaded (live reload / refresh) but reset when a different
// document is opened.
func TestScrollBehaviorOnReload(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	long, _ := filepath.Abs("../../testdata/theme-showcase.md")
	other, _ := filepath.Abs("../../testdata/sample.md")

	w := NewWindow(app, zap.NewNop(), "test")
	w.fyneWindow.Resize(fyne.NewSize(700, 400))

	w.loadFile(long)
	w.fyneWindow.Resize(fyne.NewSize(700, 400)) // ensure layout after content set

	// Scroll down.
	w.scrollContent.ScrollToOffset(fyne.NewPos(0, 250))
	afterScroll := w.scrollContent.Offset.Y
	t.Logf("after manual scroll: Y=%.0f", afterScroll)

	if afterScroll == 0 {
		t.Fatal("could not scroll the test content; cannot verify reload behavior")
	}

	// Reload the SAME document (simulates live reload / refresh): position kept.
	w.loadFile(long)
	afterReload := w.scrollContent.Offset.Y
	if afterReload != afterScroll {
		t.Errorf("reloading same file changed scroll: was %.0f, now %.0f (want preserved)",
			afterScroll, afterReload)
	}

	// Open a DIFFERENT document: scroll resets to the top.
	w.loadFile(other)
	afterOther := w.scrollContent.Offset.Y
	if afterOther != 0 {
		t.Errorf("opening different file did not reset scroll: Y=%.0f (want 0)", afterOther)
	}
}
