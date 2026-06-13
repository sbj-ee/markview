package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

// TestToggleButtonActiveState verifies the file-tree toolbar button reflects
// the panel's visibility via its importance (highlighted when visible).
func TestToggleButtonActiveState(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")

	// File tree is visible by default, so its button should start highlighted.
	if got := w.fileTreeAction.button.Importance; got != widget.HighImportance {
		t.Errorf("file tree button importance = %v, want High (active) at startup", got)
	}

	w.toggleFileTree() // hide
	if got := w.fileTreeAction.button.Importance; got != widget.LowImportance {
		t.Errorf("file tree button importance = %v, want Low (inactive) when hidden", got)
	}

	w.toggleFileTree() // show
	if got := w.fileTreeAction.button.Importance; got != widget.HighImportance {
		t.Errorf("file tree button importance = %v, want High (active) when shown", got)
	}

	// Library is off by default, so its button should not be highlighted.
	if got := w.libraryAction.button.Importance; got != widget.LowImportance {
		t.Errorf("library button importance = %v, want Low (inactive) at startup", got)
	}
}
