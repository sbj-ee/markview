package markdown

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func init() {
	// Initialize a test app for Fyne widget tests
	test.NewApp()
}

func TestNewCodeBlock(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text: "func main() {}",
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		},
	}

	cb := NewCodeBlock(segments)

	if cb == nil {
		t.Fatal("NewCodeBlock() returned nil")
	}

	if cb.richText == nil {
		t.Error("CodeBlock.richText is nil")
	}

	if cb.padding != 8 {
		t.Errorf("CodeBlock.padding = %f, want 8", cb.padding)
	}
}

func TestCodeBlock_MinSize(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "test code",
			Style: widget.RichTextStyle{},
		},
	}

	cb := NewCodeBlock(segments)
	size := cb.MinSize()

	// MinSize should be positive and include padding
	if size.Width <= 0 {
		t.Errorf("MinSize().Width = %f, want > 0", size.Width)
	}
	if size.Height <= 0 {
		t.Errorf("MinSize().Height = %f, want > 0", size.Height)
	}

	// MinSize should be larger than just the text due to padding
	if size.Width < cb.padding*2 {
		t.Errorf("MinSize().Width = %f, should include padding of %f*2", size.Width, cb.padding)
	}
}

func TestCodeBlock_CreateRenderer(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "test",
			Style: widget.RichTextStyle{},
		},
	}

	cb := NewCodeBlock(segments)
	renderer := cb.CreateRenderer()

	if renderer == nil {
		t.Fatal("CreateRenderer() returned nil")
	}
}

func TestCodeBlockRenderer_Objects(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "test",
			Style: widget.RichTextStyle{},
		},
	}

	cb := NewCodeBlock(segments)
	renderer := cb.CreateRenderer()

	objects := renderer.Objects()
	if len(objects) != 2 {
		t.Errorf("Objects() returned %d objects, want 2 (background + richtext)", len(objects))
	}
}

func TestCodeBlockRenderer_Layout(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "test",
			Style: widget.RichTextStyle{},
		},
	}

	cb := NewCodeBlock(segments)
	renderer := cb.CreateRenderer()

	// Layout should not panic
	renderer.Layout(fyne.NewSize(200, 100))
}

func TestCodeBlockRenderer_Refresh(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "test",
			Style: widget.RichTextStyle{},
		},
	}

	cb := NewCodeBlock(segments)
	renderer := cb.CreateRenderer()

	// Refresh should not panic
	renderer.Refresh()
}

func TestCodeBlockRenderer_Destroy(t *testing.T) {
	segments := []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "test",
			Style: widget.RichTextStyle{},
		},
	}

	cb := NewCodeBlock(segments)
	renderer := cb.CreateRenderer()

	// Destroy should not panic
	renderer.Destroy()
}
