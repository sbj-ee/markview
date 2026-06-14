package gui

import (
	"regexp"
	"strings"
)

// tableSeparatorCellRe matches a single GFM table separator cell, e.g. "---",
// ":---", "---:" or ":---:".
var tableSeparatorCellRe = regexp.MustCompile(`^:?-+:?$`)

// isTableRow reports whether a line could be part of a pipe table: non-blank
// and containing at least one pipe.
func isTableRow(line string) bool {
	return strings.TrimSpace(line) != "" && strings.Contains(line, "|")
}

// isTableSeparator reports whether a line is a GFM table separator row.
func isTableSeparator(line string) bool {
	if !strings.Contains(line, "-") {
		return false
	}
	cells := splitTableCells(line)
	if len(cells) == 0 {
		return false
	}
	for _, c := range cells {
		if !tableSeparatorCellRe.MatchString(c) {
			return false
		}
	}
	return true
}

// splitTableCells splits a table row into trimmed, pipe-unescaped cell strings,
// dropping the empty fields created by the leading and trailing pipes.
func splitTableCells(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := splitUnescapedPipes(line)
	cells := make([]string, len(parts))
	for i, p := range parts {
		cells[i] = strings.TrimSpace(p)
	}
	return cells
}

// splitUnescapedPipes splits on pipes that are not backslash-escaped, and
// unescapes "\|" to "|" in the resulting fields.
func splitUnescapedPipes(s string) []string {
	var cells []string
	var cur strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == '|' {
			cur.WriteByte('|')
			i++
			continue
		}
		if s[i] == '|' {
			cells = append(cells, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteByte(s[i])
	}
	cells = append(cells, cur.String())
	return cells
}

// detectTableAt finds the line range of the GFM table containing cursorRow, if
// any. It returns the inclusive [start, end] line indices and true when the
// cursor sits on a contiguous block of pipe rows whose second line is a valid
// separator.
func detectTableAt(text string, cursorRow int) (start, end int, ok bool) {
	lines := strings.Split(text, "\n")
	if cursorRow < 0 || cursorRow >= len(lines) || !isTableRow(lines[cursorRow]) {
		return 0, 0, false
	}

	start = cursorRow
	for start > 0 && isTableRow(lines[start-1]) {
		start--
	}
	end = cursorRow
	for end < len(lines)-1 && isTableRow(lines[end+1]) {
		end++
	}

	// A table needs a header followed by a separator row; the header itself
	// must not look like a separator.
	if end-start < 1 || isTableSeparator(lines[start]) || !isTableSeparator(lines[start+1]) {
		return 0, 0, false
	}
	return start, end, true
}

// parseTableLines parses the lines of a GFM table into headers, per-column
// alignments (values from columnAlignments), and data rows. Rows are normalized
// to the header column count.
func parseTableLines(lines []string) (headers, aligns []string, rows [][]string) {
	if len(lines) < 2 {
		return nil, nil, nil
	}
	headers = splitTableCells(lines[0])
	aligns = parseAlignments(splitTableCells(lines[1]), len(headers))
	for _, l := range lines[2:] {
		rows = append(rows, normalizeRow(splitTableCells(l), len(headers)))
	}
	return headers, aligns, rows
}

// parseAlignments maps separator cells to alignment labels, padded to n columns.
func parseAlignments(sepCells []string, n int) []string {
	aligns := make([]string, n)
	for i := 0; i < n; i++ {
		align := "Left"
		if i < len(sepCells) {
			c := sepCells[i]
			left := strings.HasPrefix(c, ":")
			right := strings.HasSuffix(c, ":")
			switch {
			case left && right:
				align = "Center"
			case right:
				align = "Right"
			}
		}
		aligns[i] = align
	}
	return aligns
}

// normalizeRow pads or truncates a row to exactly n cells.
func normalizeRow(row []string, n int) []string {
	out := make([]string, n)
	copy(out, row)
	return out
}
