package gui

import (
	"testing"

	"fyne.io/fyne/v2/theme"
)

// TestFileTreeNodeUsesCompactFont verifies file names render at the smaller
// caption text size.
func TestFileTreeNodeUsesCompactFont(t *testing.T) {
	n := newFileTreeNode()
	if n.label.SizeName != theme.SizeNameCaptionText {
		t.Errorf("file tree label SizeName = %q, want %q", n.label.SizeName, theme.SizeNameCaptionText)
	}
}
