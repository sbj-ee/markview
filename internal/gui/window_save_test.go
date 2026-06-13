package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"go.uber.org/zap"
)

// TestFlashSavedStatus verifies the transient save confirmation appears in the
// status bar.
func TestFlashSavedStatus(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")

	w.flashSavedStatus()
	// Stop the revert timer so it doesn't fire against a torn-down app.
	if w.saveStatusTimer != nil {
		w.saveStatusTimer.Stop()
	}

	if got := w.statusBar.Text; got != "✓ Saved" {
		t.Errorf("status bar = %q, want %q", got, "✓ Saved")
	}
}
