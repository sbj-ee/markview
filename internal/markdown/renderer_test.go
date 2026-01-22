package markdown

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
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
