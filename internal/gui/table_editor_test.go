package gui

import (
	"strings"
	"testing"
)

// TestTableEditorAlignmentAndEscaping verifies the generated markdown uses the
// per-column alignment markers and escapes literal pipes in cell content.
func TestTableEditorAlignmentAndEscaping(t *testing.T) {
	te := NewTableEditor(nil) // 2 rows x 3 cols, all left-aligned
	te.headers[0].SetText("Name")
	te.headers[1].SetText("a|b") // pipe must be escaped
	te.headers[2].SetText("Qty")
	te.aligns = []string{"Left", "Center", "Right"}
	te.cells[0][0].SetText("Widget")
	te.cells[0][1].SetText("x|y")
	te.cells[0][2].SetText("5")

	lines := strings.Split(strings.TrimRight(te.GenerateMarkdown(), "\n"), "\n")

	if lines[0] != `| Name | a\|b | Qty |` {
		t.Errorf("header row = %q", lines[0])
	}
	if lines[1] != "| :--- | :---: | ---: |" {
		t.Errorf("separator row = %q, want left/center/right markers", lines[1])
	}
	if lines[2] != `| Widget | x\|y | 5 |` {
		t.Errorf("data row = %q", lines[2])
	}
}

// TestAlignmentSeparator covers the marker mapping directly.
func TestAlignmentSeparator(t *testing.T) {
	cases := map[string]string{
		"Left":    ":---",
		"Center":  ":---:",
		"Right":   "---:",
		"unknown": ":---", // defaults to left
	}
	for align, want := range cases {
		if got := alignmentSeparator(align); got != want {
			t.Errorf("alignmentSeparator(%q) = %q, want %q", align, got, want)
		}
	}
}
