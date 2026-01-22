package markdown

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func TestNewSyntaxHighlighter(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	if highlighter == nil {
		t.Fatal("NewSyntaxHighlighter returned nil")
	}

	if highlighter.style == nil {
		t.Error("Highlighter style is nil")
	}
}

func TestSyntaxHighlighter_SetStyle(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	// Test setting a valid style
	highlighter.SetStyle("github")

	// Test setting an invalid style (should not panic)
	highlighter.SetStyle("invalid-style-name")
}

func TestSyntaxHighlighter_Highlight_Go(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	code := `func main() {
	fmt.Println("Hello, World!")
}`

	segments := highlighter.Highlight(code, "go")

	if len(segments) == 0 {
		t.Fatal("No segments returned for Go code")
	}

	// Should have multiple segments for different tokens
	if len(segments) < 5 {
		t.Errorf("Expected at least 5 segments for Go code, got %d", len(segments))
	}

	// Check that at least one segment has monospace style
	hasMonospace := false
	for _, seg := range segments {
		if textSeg, ok := seg.(*widget.TextSegment); ok {
			if textSeg.Style.TextStyle.Monospace {
				hasMonospace = true
				break
			}
		}
	}

	if !hasMonospace {
		t.Error("Expected at least one segment to have monospace style")
	}
}

func TestSyntaxHighlighter_Highlight_Python(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	code := `def hello():
    print("Hello, World!")`

	segments := highlighter.Highlight(code, "python")

	if len(segments) == 0 {
		t.Fatal("No segments returned for Python code")
	}

	// Should have multiple segments
	if len(segments) < 3 {
		t.Errorf("Expected at least 3 segments for Python code, got %d", len(segments))
	}
}

func TestSyntaxHighlighter_Highlight_JavaScript(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	code := `function hello() {
	console.log("Hello, World!");
}`

	segments := highlighter.Highlight(code, "javascript")

	if len(segments) == 0 {
		t.Fatal("No segments returned for JavaScript code")
	}

	// Should have multiple segments
	if len(segments) < 5 {
		t.Errorf("Expected at least 5 segments for JavaScript code, got %d", len(segments))
	}
}

func TestSyntaxHighlighter_Highlight_UnknownLanguage(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	code := "some code here"

	segments := highlighter.Highlight(code, "unknown-language")

	if len(segments) == 0 {
		t.Fatal("No segments returned for unknown language")
	}

	// Should fall back to plain text rendering
	if len(segments) != 1 {
		t.Errorf("Expected 1 segment for unknown language (fallback), got %d", len(segments))
	}
}

func TestSyntaxHighlighter_Highlight_EmptyCode(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	code := ""

	segments := highlighter.Highlight(code, "go")

	if len(segments) != 0 {
		t.Errorf("Expected 0 segments for empty code, got %d", len(segments))
	}
}

func TestSyntaxHighlighter_HighlightAsPlainText(t *testing.T) {
	highlighter := NewSyntaxHighlighter()

	code := `func main() {
	fmt.Println("Hello")
}`

	result := highlighter.HighlightAsPlainText(code, "go")

	if result == "" {
		t.Error("HighlightAsPlainText returned empty string")
	}

	// Should return some formatted text
	if len(result) < len(code) {
		t.Errorf("Expected result to be at least as long as input, got %d chars for %d input chars", len(result), len(code))
	}
}

func TestGetAvailableStyles(t *testing.T) {
	styles := GetAvailableStyles()

	if len(styles) == 0 {
		t.Fatal("No styles available")
	}

	// Should include common styles
	hasMonokai := false
	for _, style := range styles {
		if style == "monokai" {
			hasMonokai = true
			break
		}
	}

	if !hasMonokai {
		t.Error("Expected 'monokai' to be in available styles")
	}
}

func TestCreateHighlightedLabel(t *testing.T) {
	// Skip this test as it requires Fyne app context for theme colors
	t.Skip("Skipping CreateHighlightedLabel test - requires Fyne app context")
}

func TestCodeBlockLayout_MinSize(t *testing.T) {
	layout := &codeBlockLayout{
		padding: fyne.NewSize(10, 10),
	}

	// Test with empty objects
	size := layout.MinSize([]fyne.CanvasObject{})
	if size.Width != 0 || size.Height != 0 {
		t.Errorf("Expected zero size for empty objects, got %v", size)
	}
}

func TestCodeBlockLayout_Layout(t *testing.T) {
	layout := &codeBlockLayout{
		padding: fyne.NewSize(10, 10),
	}

	// Test with empty objects (should not panic)
	layout.Layout([]fyne.CanvasObject{}, fyne.NewSize(100, 100))
}
