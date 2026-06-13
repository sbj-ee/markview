package gui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// renderMarkdown mirrors the Goldmark configuration used by the HTML export.
func renderMarkdown(src string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Typographer),
		goldmark.WithRendererOptions(goldmarkhtml.WithHardWraps()),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return ""
	}
	return buf.String()
}

func TestConvertMermaidBlocksRendersDiagram(t *testing.T) {
	html := renderMarkdown("```mermaid\ngraph TD;\nA-->B;\n```\n")

	// Sanity-check the assumption about Goldmark's output.
	if !strings.Contains(html, `class="language-mermaid"`) {
		t.Fatalf("unexpected Goldmark output for mermaid block: %s", html)
	}

	out := convertMermaidBlocks(html)
	if !strings.Contains(out, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid container, got: %s", out)
	}
	if strings.Contains(out, `language-mermaid`) {
		t.Errorf("mermaid code wrapper not removed: %s", out)
	}
	// The diagram source survives (HTML-escaped is fine; the browser decodes it).
	if !strings.Contains(out, "graph TD;") {
		t.Errorf("diagram source lost: %s", out)
	}
}

func TestConvertMermaidBlocksLeavesOtherCodeAlone(t *testing.T) {
	html := renderMarkdown("```go\nfmt.Println(\"hi\")\n```\n")
	out := convertMermaidBlocks(html)
	if out != html {
		t.Errorf("non-mermaid code block was modified:\n got:  %s\n want: %s", out, html)
	}
}
