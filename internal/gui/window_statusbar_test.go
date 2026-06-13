package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"go.uber.org/zap"
)

// TestStatusBarDirtyIndicator verifies the status bar reflects unsaved changes
// while in edit mode.
func TestStatusBarDirtyIndicator(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")
	w.editMode = true

	w.isDirty = true
	w.updateStatusBar()
	if got := w.statusBar.Text; got != "● Unsaved" {
		t.Errorf("status bar = %q, want %q when dirty", got, "● Unsaved")
	}

	w.isDirty = false
	w.updateStatusBar()
	if got := w.statusBar.Text; got != "Edit Mode" {
		t.Errorf("status bar = %q, want %q when clean", got, "Edit Mode")
	}
}
