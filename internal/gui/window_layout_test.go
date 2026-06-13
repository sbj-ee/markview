package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"go.uber.org/zap"
)

// TestTogglePanesReclaimSpace verifies that hiding the file tree and the
// TOC/outline pane causes the content area to expand into the freed space.
func TestTogglePanesReclaimSpace(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp() // reset global app afterwards

	w := NewWindow(app, zap.NewNop(), "test")
	w.fyneWindow.Resize(fyne.NewSize(1200, 800))

	contentWidth := func() float32 { return w.scrollContent.Size().Width }

	both := contentWidth()
	if both <= 0 {
		t.Fatalf("content width not laid out (got %v); cannot verify", both)
	}

	// Hide the file tree: the content pane should widen.
	w.toggleFileTree()
	noTree := contentWidth()
	if noTree <= both {
		t.Errorf("hiding file tree did not widen content: both=%v noTree=%v", both, noTree)
	}

	// Also hide the TOC/outline pane: content should widen further.
	w.toggleTOC()
	noToc := contentWidth()
	if noToc <= noTree {
		t.Errorf("hiding TOC did not widen content: noTree=%v noToc=%v", noTree, noToc)
	}

	// With both panes hidden, content should fill almost the whole width.
	if noToc < w.fyneWindow.Canvas().Size().Width*0.9 {
		t.Errorf("content did not fill width with panes hidden: content=%v window=%v",
			noToc, w.fyneWindow.Canvas().Size().Width)
	}

	// Restoring both panes should shrink the content back toward the start.
	w.toggleFileTree()
	w.toggleTOC()
	restored := contentWidth()
	if restored >= noToc {
		t.Errorf("restoring panes did not shrink content: noToc=%v restored=%v", noToc, restored)
	}

	t.Logf("content width — canvas=%.0f both=%.0f noTree=%.0f noToc=%.0f restored=%.0f",
		w.fyneWindow.Canvas().Size().Width, both, noTree, noToc, restored)
}
