package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

// collectText walks a canvas object tree and returns the text of every Label
// and Button found.
func collectText(obj fyne.CanvasObject) []string {
	var out []string
	switch o := obj.(type) {
	case *widget.Label:
		out = append(out, o.Text)
	case *widget.Button:
		out = append(out, o.Text)
	case *fyne.Container:
		for _, child := range o.Objects {
			out = append(out, collectText(child)...)
		}
	}
	return out
}

func TestWelcomeContentShownOnStartup(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")

	if w.scrollContent.Content == nil {
		t.Fatal("content area is empty on startup; expected welcome panel")
	}
	texts := collectText(w.scrollContent.Content)
	for _, want := range []string{"Welcome to MarkView", "Open File…", "Open Folder…"} {
		if !contains(texts, want) {
			t.Errorf("welcome panel missing %q (got %v)", want, texts)
		}
	}
}

func TestWelcomeContentListsRecentFiles(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	w := NewWindow(app, zap.NewNop(), "test")
	w.recentFiles = []string{"/tmp/alpha.md", "/tmp/beta.md"}

	texts := collectText(w.buildWelcomeContent())
	if !contains(texts, "Recent") {
		t.Errorf("expected a Recent heading, got %v", texts)
	}
	for _, want := range []string{"alpha.md", "beta.md"} {
		if !contains(texts, want) {
			t.Errorf("welcome panel missing recent file %q (got %v)", want, texts)
		}
	}
}
