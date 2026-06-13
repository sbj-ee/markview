package gui

import (
	"strings"
	"testing"
)

// TestExportConvertsMathToMathJaxDelimiters verifies LaTeX is rendered into the
// \(…\) and \[…\] delimiters the injected MathJax script understands, rather
// than passing through raw (where Markdown could mangle it).
func TestExportConvertsMathToMathJaxDelimiters(t *testing.T) {
	out := renderMarkdown("Inline $a^2+b^2$ here.\n\n$$\nE=mc^2\n$$\n")

	if !strings.Contains(out, `\(a^2+b^2\)`) {
		t.Errorf("inline math not converted to \\(…\\): %s", out)
	}
	if !strings.Contains(out, `\[`) || !strings.Contains(out, "E=mc^2") {
		t.Errorf("block math not converted to \\[…\\]: %s", out)
	}
	// The literal dollar-delimited form should no longer be present.
	if strings.Contains(out, "$a^2+b^2$") {
		t.Errorf("raw inline LaTeX left unconverted: %s", out)
	}
}
