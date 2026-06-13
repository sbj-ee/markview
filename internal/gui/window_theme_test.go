package gui

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/sbj-ee/markview/internal/themes"
	"go.uber.org/zap"
)

// TestThemeChangePreservesUnsavedEdits ensures changing the theme while editing
// re-renders without re-reading the file from disk, so unsaved edits and the
// dirty state are kept.
func TestThemeChangePreservesUnsavedEdits(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")

	tmp := filepath.Join(t.TempDir(), "doc.md")
	if err := os.WriteFile(tmp, []byte("# on disk\n"), 0644); err != nil {
		t.Fatal(err)
	}
	w.currentFile = tmp
	w.contentBuffer = "# on disk\n"

	// Simulate active editing with unsaved changes.
	w.editMode = true
	w.editor.SetText("# unsaved edits in progress")
	w.isDirty = true

	// Default theme is Dark; switch to a different one.
	w.setTheme(themes.ThemeLight)

	if got := w.editor.GetText(); got != "# unsaved edits in progress" {
		t.Errorf("theme change clobbered editor text: got %q", got)
	}
	if !w.isDirty {
		t.Error("theme change reset the unsaved-changes flag")
	}
	if w.currentTheme != themes.ThemeLight {
		t.Errorf("theme not applied: got %v", w.currentTheme)
	}
}
