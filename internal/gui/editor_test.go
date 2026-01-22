package gui

import (
	"testing"
)

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "single line",
			input:    "hello world",
			expected: []string{"hello world"},
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "trailing newline",
			input:    "line1\nline2\n",
			expected: []string{"line1", "line2", ""},
		},
		{
			name:     "empty lines",
			input:    "line1\n\nline3",
			expected: []string{"line1", "", "line3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLines(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitLines(%q) returned %d lines, expected %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("splitLines(%q)[%d] = %q, expected %q", tt.input, i, line, tt.expected[i])
				}
			}
		})
	}
}

func TestJoinLines(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: "",
		},
		{
			name:     "single line",
			input:    []string{"hello world"},
			expected: "hello world",
		},
		{
			name:     "multiple lines",
			input:    []string{"line1", "line2", "line3"},
			expected: "line1\nline2\nline3",
		},
		{
			name:     "with empty string",
			input:    []string{"line1", "", "line3"},
			expected: "line1\n\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinLines(tt.input)
			if result != tt.expected {
				t.Errorf("joinLines(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindSelectionStart(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		selected  string
		cursorPos int
		expected  int
	}{
		{
			name:      "selection at start",
			text:      "hello world",
			selected:  "hello",
			cursorPos: 5,
			expected:  0,
		},
		{
			name:      "selection at end",
			text:      "hello world",
			selected:  "world",
			cursorPos: 6,
			expected:  6,
		},
		{
			name:      "selection in middle",
			text:      "hello beautiful world",
			selected:  "beautiful",
			cursorPos: 10,
			expected:  6,
		},
		{
			name:      "not found",
			text:      "hello world",
			selected:  "xyz",
			cursorPos: 5,
			expected:  -1,
		},
		{
			name:      "cursor before selection",
			text:      "hello world",
			selected:  "world",
			cursorPos: 2,
			expected:  6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findSelectionStart(tt.text, tt.selected, tt.cursorPos)
			if result != tt.expected {
				t.Errorf("findSelectionStart(%q, %q, %d) = %d, expected %d",
					tt.text, tt.selected, tt.cursorPos, result, tt.expected)
			}
		})
	}
}

func TestSplitLinesRoundTrip(t *testing.T) {
	// Test that split followed by join returns original
	tests := []string{
		"",
		"single line",
		"line1\nline2",
		"line1\nline2\nline3",
		"has\n\nempty\n\nlines",
	}

	for _, original := range tests {
		lines := splitLines(original)
		result := joinLines(lines)
		if result != original {
			t.Errorf("Round trip failed for %q: got %q", original, result)
		}
	}
}
