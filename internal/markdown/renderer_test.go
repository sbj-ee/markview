package markdown

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/widget"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

func TestNewRenderer(t *testing.T) {
	source := []byte("# Test")
	renderer := NewRenderer(source)

	if renderer == nil {
		t.Fatal("NewRenderer returned nil")
	}

	if renderer.source == nil {
		t.Error("Renderer source is nil")
	}

	if renderer.segments == nil {
		t.Error("Renderer segments is nil")
	}

	if renderer.highlighter == nil {
		t.Error("Renderer highlighter is nil")
	}
}

func TestRenderer_Headings(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "H1",
			input: "# Heading 1",
			want:  "Heading 1",
		},
		{
			name:  "H2",
			input: "## Heading 2",
			want:  "Heading 2",
		},
		{
			name:  "H3",
			input: "### Heading 3",
			want:  "Heading 3",
		},
		{
			name:  "H6",
			input: "###### Heading 6",
			want:  "Heading 6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.input)
			reader := text.NewReader(source)
			doc := md.Parser().Parse(reader)

			renderer := NewRenderer(source)
			segments := renderer.Render(doc)

			if len(segments) == 0 {
				t.Fatal("No segments rendered")
			}

			// Find the heading segment
			found := false
			for _, seg := range segments {
				if textSeg, ok := seg.(*widget.TextSegment); ok {
					if textSeg.Text == tt.want {
						found = true
						// Check that it's bold
						if !textSeg.Style.TextStyle.Bold {
							t.Error("Heading should be bold")
						}
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected heading text '%s' not found in segments", tt.want)
			}
		})
	}
}

func TestRenderer_Paragraphs(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	input := "This is a simple paragraph."
	source := []byte(input)
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	segments := renderer.Render(doc)

	if len(segments) == 0 {
		t.Fatal("No segments rendered")
	}

	// Should contain the paragraph text (may be across multiple segments)
	var combinedText string
	for _, seg := range segments {
		if textSeg, ok := seg.(*widget.TextSegment); ok {
			combinedText += textSeg.Text
		}
	}

	if !strings.Contains(combinedText, input) {
		t.Errorf("Paragraph text '%s' not found in combined segments: '%s'", input, combinedText)
	}
}

func TestRenderer_CodeBlocks(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	tests := []struct {
		name     string
		input    string
		wantCode string
	}{
		{
			name:     "fenced code block",
			input:    "```go\nfunc main() {}\n```",
			wantCode: "func main()",
		},
		{
			name:     "code block with language",
			input:    "```python\nprint('hello')\n```",
			wantCode: "print('hello')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.input)
			reader := text.NewReader(source)
			doc := md.Parser().Parse(reader)

			renderer := NewRenderer(source)
			segments := renderer.Render(doc)

			if len(segments) == 0 {
				t.Fatal("No segments rendered")
			}

			// Code should be in the segments (may be split across multiple segments due to highlighting)
			var combinedText string
			for _, seg := range segments {
				if textSeg, ok := seg.(*widget.TextSegment); ok {
					combinedText += textSeg.Text
				}
			}

			if !strings.Contains(combinedText, tt.wantCode) {
				t.Errorf("Code '%s' not found in combined segments: '%s'", tt.wantCode, combinedText)
			}
		})
	}
}

func TestRenderer_Lists(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	input := "- Item 1\n- Item 2\n- Item 3"
	source := []byte(input)
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	segments := renderer.Render(doc)

	if len(segments) == 0 {
		t.Fatal("No segments rendered")
	}

	// Should have list item markers
	bulletCount := 0
	for _, seg := range segments {
		if textSeg, ok := seg.(*widget.TextSegment); ok {
			if textSeg.Text == "• " {
				bulletCount++
			}
		}
	}

	if bulletCount < 3 {
		t.Errorf("Expected at least 3 bullet points, got %d", bulletCount)
	}
}

func TestRenderer_Emphasis(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	tests := []struct {
		name       string
		input      string
		wantText   string
		wantItalic bool
		wantBold   bool
	}{
		{
			name:       "italic",
			input:      "*italic*",
			wantText:   "italic",
			wantItalic: true,
			wantBold:   false,
		},
		{
			name:       "bold",
			input:      "**bold**",
			wantText:   "bold",
			wantItalic: false,
			wantBold:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.input)
			reader := text.NewReader(source)
			doc := md.Parser().Parse(reader)

			renderer := NewRenderer(source)
			segments := renderer.Render(doc)

			if len(segments) == 0 {
				t.Fatal("No segments rendered")
			}

			// Find the emphasized text
			found := false
			for _, seg := range segments {
				if textSeg, ok := seg.(*widget.TextSegment); ok {
					if textSeg.Text == tt.wantText {
						found = true
						if textSeg.Style.TextStyle.Italic != tt.wantItalic {
							t.Errorf("Expected italic=%v, got italic=%v", tt.wantItalic, textSeg.Style.TextStyle.Italic)
						}
						if textSeg.Style.TextStyle.Bold != tt.wantBold {
							t.Errorf("Expected bold=%v, got bold=%v", tt.wantBold, textSeg.Style.TextStyle.Bold)
						}
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected text '%s' not found in segments", tt.wantText)
			}
		})
	}
}

func TestRenderer_EmptyContent(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))

	source := []byte("")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	segments := renderer.Render(doc)

	// Empty content should return empty segments
	if len(segments) != 0 {
		t.Errorf("Expected 0 segments for empty content, got %d", len(segments))
	}
}
