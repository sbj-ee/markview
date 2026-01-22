package markdown

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func init() {
	// Initialize a test app for Fyne widget tests
	test.NewApp()
}

func TestNewBlockquote(t *testing.T) {
	bq := NewBlockquote("This is a quote")

	if bq == nil {
		t.Fatal("NewBlockquote() returned nil")
	}

	if bq.label == nil {
		t.Error("Blockquote.label is nil")
	}

	if bq.padding != 12 {
		t.Errorf("Blockquote.padding = %f, want 12", bq.padding)
	}

	if bq.borderWidth != 4 {
		t.Errorf("Blockquote.borderWidth = %f, want 4", bq.borderWidth)
	}
}

func TestBlockquote_MinSize(t *testing.T) {
	bq := NewBlockquote("Test quote text")
	size := bq.MinSize()

	// MinSize should be positive and include padding + border
	if size.Width <= 0 {
		t.Errorf("MinSize().Width = %f, want > 0", size.Width)
	}
	if size.Height <= 0 {
		t.Errorf("MinSize().Height = %f, want > 0", size.Height)
	}

	// MinSize should account for padding and border
	minWidth := bq.padding*2 + bq.borderWidth
	if size.Width < minWidth {
		t.Errorf("MinSize().Width = %f, should be >= %f (padding*2 + border)", size.Width, minWidth)
	}
}

func TestBlockquote_CreateRenderer(t *testing.T) {
	bq := NewBlockquote("Test")
	renderer := bq.CreateRenderer()

	if renderer == nil {
		t.Fatal("CreateRenderer() returned nil")
	}
}

func TestBlockquoteRenderer_Objects(t *testing.T) {
	bq := NewBlockquote("Test")
	renderer := bq.CreateRenderer()

	objects := renderer.Objects()
	if len(objects) != 3 {
		t.Errorf("Objects() returned %d objects, want 3 (background + border + label)", len(objects))
	}
}

func TestBlockquoteRenderer_Layout(t *testing.T) {
	bq := NewBlockquote("Test quote")
	renderer := bq.CreateRenderer()

	// Layout should not panic
	renderer.Layout(fyne.NewSize(200, 100))
}

func TestBlockquoteRenderer_Refresh(t *testing.T) {
	bq := NewBlockquote("Test")
	renderer := bq.CreateRenderer()

	// Refresh should not panic
	renderer.Refresh()
}

func TestBlockquoteRenderer_Destroy(t *testing.T) {
	bq := NewBlockquote("Test")
	renderer := bq.CreateRenderer()

	// Destroy should not panic
	renderer.Destroy()
}

func TestBlockquote_ItalicStyle(t *testing.T) {
	bq := NewBlockquote("Test")

	if !bq.label.TextStyle.Italic {
		t.Error("Blockquote label should have italic text style")
	}
}

func TestBlockquote_WordWrap(t *testing.T) {
	bq := NewBlockquote("Test")

	if bq.label.Wrapping != fyne.TextWrapWord {
		t.Errorf("Blockquote label Wrapping = %v, want TextWrapWord", bq.label.Wrapping)
	}
}
