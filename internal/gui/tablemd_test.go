package gui

import (
	"reflect"
	"strings"
	"testing"
)

const sampleDoc = `intro line

| A | B |
| :--- | ---: |
| 1 | 2 |
| 3 | 4 |

outro`

func TestDetectTableAt(t *testing.T) {
	tests := []struct {
		name      string
		row       int
		wantOK    bool
		wantStart int
		wantEnd   int
	}{
		{"on header", 2, true, 2, 5},
		{"on separator", 3, true, 2, 5},
		{"on data row", 4, true, 2, 5},
		{"on last data row", 5, true, 2, 5},
		{"on intro text", 0, false, 0, 0},
		{"on blank line", 1, false, 0, 0},
		{"on outro text", 7, false, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, ok := detectTableAt(sampleDoc, tt.row)
			if ok != tt.wantOK || (ok && (start != tt.wantStart || end != tt.wantEnd)) {
				t.Errorf("detectTableAt(row=%d) = (%d,%d,%v), want (%d,%d,%v)",
					tt.row, start, end, ok, tt.wantStart, tt.wantEnd, tt.wantOK)
			}
		})
	}
}

func TestDetectTableAtRejectsPipeParagraph(t *testing.T) {
	// Pipes but no separator row -> not a table.
	text := "a | b\nc | d\n"
	if _, _, ok := detectTableAt(text, 0); ok {
		t.Error("pipe paragraph without a separator should not be detected as a table")
	}
}

func TestParseTableLines(t *testing.T) {
	lines := strings.Split(sampleDoc, "\n")[2:6]
	headers, aligns, rows := parseTableLines(lines)

	if !reflect.DeepEqual(headers, []string{"A", "B"}) {
		t.Errorf("headers = %v", headers)
	}
	if !reflect.DeepEqual(aligns, []string{"Left", "Right"}) {
		t.Errorf("aligns = %v", aligns)
	}
	if !reflect.DeepEqual(rows, [][]string{{"1", "2"}, {"3", "4"}}) {
		t.Errorf("rows = %v", rows)
	}
}

func TestParseTableLinesUnescapesAndNormalizes(t *testing.T) {
	lines := []string{
		`| a\|b | c |`,
		`| --- | :---: |`,
		`| x | y | extra |`, // ragged: truncated to header count
		`| z |`,             // ragged: padded
	}
	headers, aligns, rows := parseTableLines(lines)
	if !reflect.DeepEqual(headers, []string{"a|b", "c"}) {
		t.Errorf("headers = %v (pipe should be unescaped)", headers)
	}
	if !reflect.DeepEqual(aligns, []string{"Left", "Center"}) {
		t.Errorf("aligns = %v", aligns)
	}
	if !reflect.DeepEqual(rows, [][]string{{"x", "y"}, {"z", ""}}) {
		t.Errorf("rows = %v (should be normalized to 2 columns)", rows)
	}
}

func TestEditExistingTableRoundTrip(t *testing.T) {
	// Parse a table, load it into the editor, regenerate — alignment and
	// content (with re-escaped pipes) should round-trip.
	lines := []string{
		`| Name | Qty |`,
		`| :--- | ---: |`,
		`| a\|b | 5 |`,
	}
	headers, aligns, rows := parseTableLines(lines)
	te := NewTableEditorWithData(headers, aligns, rows, nil)

	got := strings.Split(strings.TrimRight(te.GenerateMarkdown(), "\n"), "\n")
	want := []string{
		`| Name | Qty |`,
		`| :--- | ---: |`,
		`| a\|b | 5 |`,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round-trip mismatch:\n got:  %v\n want: %v", got, want)
	}
}

func TestReplaceLines(t *testing.T) {
	e := NewMarkdownEditor(nil)
	e.SetText("line0\nOLD1\nOLD2\nline3")
	e.ReplaceLines(1, 2, "NEW")
	if got := e.GetText(); got != "line0\nNEW\nline3" {
		t.Errorf("ReplaceLines = %q", got)
	}

	// Out-of-range is a no-op.
	e.SetText("a\nb")
	e.ReplaceLines(0, 5, "X")
	if got := e.GetText(); got != "a\nb" {
		t.Errorf("out-of-range ReplaceLines changed text: %q", got)
	}
}
