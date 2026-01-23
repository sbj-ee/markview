package markdown

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestNewRenderer(t *testing.T) {
	source := []byte("test")
	renderer := NewRenderer(source)

	if renderer == nil {
		t.Fatal("NewRenderer() returned nil")
	}

	if renderer.source == nil {
		t.Error("Renderer.source is nil")
	}

	if renderer.highlighter == nil {
		t.Error("Renderer.highlighter is nil")
	}
}

func TestRender_Heading(t *testing.T) {
	md := goldmark.New()
	source := []byte("# Test Heading")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for heading")
	}
}

func TestRender_Paragraph(t *testing.T) {
	md := goldmark.New()
	source := []byte("This is a paragraph.")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for paragraph")
	}
}

func TestRender_CodeBlock(t *testing.T) {
	md := goldmark.New()
	source := []byte("```go\nfunc main() {}\n```")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for code block")
	}
}

func TestRender_List(t *testing.T) {
	md := goldmark.New()
	source := []byte("- Item 1\n- Item 2\n- Item 3")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for list")
	}
}

func TestRender_EmptyDocument(t *testing.T) {
	md := goldmark.New()
	source := []byte("")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	// Empty document should return empty Container
	if len(cont.Objects) != 0 {
		t.Errorf("Render() returned %d widgets for empty document, want 0", len(cont.Objects))
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple text",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "text with newlines",
			input: "hello\nworld",
			want:  "hello world",
		},
		{
			name:  "text with carriage returns",
			input: "hello\r\nworld",
			want:  "hello world",
		},
		{
			name:  "text with tabs",
			input: "hello\tworld",
			want:  "hello world",
		},
		{
			name:  "text with multiple spaces",
			input: "hello    world",
			want:  "hello world",
		},
		{
			name:  "text with leading/trailing whitespace",
			input: "   hello world   ",
			want:  "hello world",
		},
		{
			name:  "HTML entity - smart quote",
			input: "don&rsquo;t",
			want:  "don\u2019t",
		},
		{
			name:  "HTML entity - left double quote",
			input: "&ldquo;quoted&rdquo;",
			want:  "\u201cquoted\u201d",
		},
		{
			name:  "HTML entity - em dash",
			input: "hello&mdash;world",
			want:  "hello\u2014world",
		},
		{
			name:  "HTML entity - ampersand",
			input: "rock &amp; roll",
			want:  "rock & roll",
		},
		{
			name:  "control characters",
			input: "hello\x00world",
			want:  "hello world",
		},
		{
			name:  "mixed issues",
			input: "  hello\n\n&rsquo;world&rsquo;  \n",
			want:  "hello \u2019world\u2019",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only whitespace",
			input: "   \n\t  ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeText(tt.input)
			if got != tt.want {
				t.Errorf("normalizeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractInlineText(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantText string
	}{
		{
			name:     "simple text",
			markdown: "# Hello",
			wantText: "Hello",
		},
		{
			name:     "bold text",
			markdown: "# **Bold**",
			wantText: "Bold",
		},
		{
			name:     "italic text",
			markdown: "# *Italic*",
			wantText: "Italic",
		},
		{
			name:     "mixed formatting",
			markdown: "# **Bold** and *italic*",
			wantText: "Bold and italic",
		},
	}

	md := goldmark.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.markdown)
			reader := text.NewReader(source)
			doc := md.Parser().Parse(reader)

			renderer := NewRenderer(source)

			// Find the heading node
			heading := doc.FirstChild()
			if heading == nil {
				t.Fatal("No heading node found")
			}

			text := renderer.extractInlineText(heading)
			if !strings.Contains(text, tt.wantText) {
				t.Errorf("extractInlineText() = %q, want text containing %q", text, tt.wantText)
			}
		})
	}
}

func TestRender_Blockquote(t *testing.T) {
	md := goldmark.New()
	source := []byte("> This is a quote")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for blockquote")
	}
}

func TestRender_ThematicBreak(t *testing.T) {
	md := goldmark.New()
	source := []byte("---")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for thematic break")
	}
}

func TestRender_MultipleHeadings_WithSpacing(t *testing.T) {
	md := goldmark.New()
	source := []byte("# Heading 1\n\n## Heading 2\n\n### Heading 3")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	// Should have headings plus spacers
	// First heading has no spacer before, but has spacer after
	// Second and third headings have spacers before and after
	if len(cont.Objects) < 5 {
		t.Errorf("Render() returned %d widgets, expected at least 5 (3 headings + spacers)", len(cont.Objects))
	}
}

func TestRender_NestedList(t *testing.T) {
	md := goldmark.New()
	source := []byte("- Item 1\n  - Nested 1\n  - Nested 2\n- Item 2")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for nested list")
	}
}

func TestRender_OrderedList(t *testing.T) {
	md := goldmark.New()
	source := []byte("1. First\n2. Second\n3. Third")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)
	content := renderer.Render(doc)

	cont, ok := content.(*fyne.Container)
	if !ok {
		t.Fatal("Render() did not return Container")
	}

	if len(cont.Objects) < 1 {
		t.Error("Render() returned no widgets for ordered list")
	}
}

func TestToSubscript(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2", "₂"},
		{"0123456789", "₀₁₂₃₄₅₆₇₈₉"},
		{"H2O", "H₂O"}, // H has no subscript, stays as H
		{"abc", "ₐbc"}, // a has subscript, b and c don't
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSubscript(tt.input)
			if got != tt.want {
				t.Errorf("toSubscript(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToSuperscript(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2", "²"},
		{"0123456789", "⁰¹²³⁴⁵⁶⁷⁸⁹"},
		{"abc", "ᵃᵇᶜ"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSuperscript(tt.input)
			if got != tt.want {
				t.Errorf("toSuperscript(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestContainsInlineHTML(t *testing.T) {
	md := goldmark.New()

	tests := []struct {
		name     string
		markdown string
		want     bool
	}{
		{"plain text", "Hello world", false},
		{"subscript", "H<sub>2</sub>O", true},
		{"superscript", "x<sup>2</sup>", true},
		{"bold", "**bold**", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.markdown)
			reader := text.NewReader(source)
			doc := md.Parser().Parse(reader)

			renderer := NewRenderer(source)

			// Find the paragraph
			para := doc.FirstChild()
			if para == nil {
				t.Fatal("No paragraph node found")
			}

			got := renderer.containsInlineHTML(para)
			if got != tt.want {
				t.Errorf("containsInlineHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderInlineNode_Subscript(t *testing.T) {
	md := goldmark.New()
	source := []byte("H<sub>2</sub>O")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)

	// Get the paragraph node
	para := doc.FirstChild()
	if para == nil {
		t.Fatal("No paragraph node found")
	}

	// Render all inline nodes and collect text
	var resultText string
	for child := para.FirstChild(); child != nil; child = child.NextSibling() {
		segments := renderer.renderInlineNode(child)
		for _, seg := range segments {
			if textSeg, ok := seg.(*widget.TextSegment); ok {
				resultText += textSeg.Text
			}
		}
	}

	want := "H₂O"
	if resultText != want {
		t.Errorf("renderInlineNode subscript = %q, want %q", resultText, want)
	}
}

func TestRenderInlineNode_Superscript(t *testing.T) {
	md := goldmark.New()
	source := []byte("x<sup>2</sup>")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	renderer := NewRenderer(source)

	// Get the paragraph node
	para := doc.FirstChild()
	if para == nil {
		t.Fatal("No paragraph node found")
	}

	// Render all inline nodes and collect text
	var resultText string
	for child := para.FirstChild(); child != nil; child = child.NextSibling() {
		segments := renderer.renderInlineNode(child)
		for _, seg := range segments {
			if textSeg, ok := seg.(*widget.TextSegment); ok {
				resultText += textSeg.Text
			}
		}
	}

	want := "x²"
	if resultText != want {
		t.Errorf("renderInlineNode superscript = %q, want %q", resultText, want)
	}
}

func TestExtractInlineTextWithHTML(t *testing.T) {
	md := goldmark.New()

	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{"subscript", "H<sub>2</sub>O", "H₂O"},
		{"superscript", "x<sup>2</sup>", "x²"},
		{"mixed", "H<sub>2</sub>O and E=mc<sup>2</sup>", "H₂O and E=mc²"},
		{"no html", "plain text", "plain text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte(tt.markdown)
			reader := text.NewReader(source)
			doc := md.Parser().Parse(reader)

			renderer := NewRenderer(source)

			// Get the paragraph node
			para := doc.FirstChild()
			if para == nil {
				t.Fatal("No paragraph node found")
			}

			got := renderer.extractInlineTextWithHTML(para)
			if got != tt.want {
				t.Errorf("extractInlineTextWithHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}
