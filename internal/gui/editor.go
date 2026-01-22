package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// MarkdownEditor wraps a multi-line entry with markdown editing support
type MarkdownEditor struct {
	widget.BaseWidget
	entry     *widget.Entry
	OnChanged func(content string)
}

// NewMarkdownEditor creates a new markdown editor
func NewMarkdownEditor(onChanged func(content string)) *MarkdownEditor {
	editor := &MarkdownEditor{
		OnChanged: onChanged,
	}
	editor.entry = widget.NewMultiLineEntry()
	editor.entry.Wrapping = fyne.TextWrapWord
	editor.entry.OnChanged = func(text string) {
		if editor.OnChanged != nil {
			editor.OnChanged(text)
		}
	}
	editor.ExtendBaseWidget(editor)
	return editor
}

// CreateRenderer implements fyne.Widget
func (e *MarkdownEditor) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(e.entry)
}

// SetText sets the editor content
func (e *MarkdownEditor) SetText(text string) {
	e.entry.SetText(text)
}

// GetText returns the current editor content
func (e *MarkdownEditor) GetText() string {
	return e.entry.Text
}

// Focus focuses the editor for input
func (e *MarkdownEditor) Focus(canvas fyne.Canvas) {
	canvas.Focus(e.entry)
}

// TypedKey handles key events - delegates to entry
func (e *MarkdownEditor) TypedKey(key *fyne.KeyEvent) {
	e.entry.TypedKey(key)
}

// TypedRune handles rune events - delegates to entry
func (e *MarkdownEditor) TypedRune(r rune) {
	e.entry.TypedRune(r)
}

// MinSize returns the minimum size of the editor
func (e *MarkdownEditor) MinSize() fyne.Size {
	return e.entry.MinSize()
}

// WrapSelection wraps the selected text with prefix and suffix, or inserts at cursor
func (e *MarkdownEditor) WrapSelection(prefix, suffix string) {
	text := e.entry.Text
	// Get cursor position
	col := e.entry.CursorColumn
	row := e.entry.CursorRow

	// Find position in text
	pos := 0
	lines := splitLines(text)
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	pos += col

	// Check if there's a selection by looking at selected text
	selected := e.entry.SelectedText()
	if selected != "" {
		// Find selection in text and wrap it
		selStart := findSelectionStart(text, selected, pos)
		if selStart >= 0 {
			newText := text[:selStart] + prefix + selected + suffix + text[selStart+len(selected):]
			e.entry.SetText(newText)
			return
		}
	}

	// No selection, insert at cursor
	if pos > len(text) {
		pos = len(text)
	}
	newText := text[:pos] + prefix + suffix + text[pos:]
	e.entry.SetText(newText)
}

// InsertAtCursor inserts text at the current cursor position
func (e *MarkdownEditor) InsertAtCursor(insert string) {
	text := e.entry.Text
	col := e.entry.CursorColumn
	row := e.entry.CursorRow

	// Find position in text
	pos := 0
	lines := splitLines(text)
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1
	}
	pos += col

	if pos > len(text) {
		pos = len(text)
	}
	newText := text[:pos] + insert + text[pos:]
	e.entry.SetText(newText)
}

// InsertAtLineStart inserts text at the beginning of the current line
func (e *MarkdownEditor) InsertAtLineStart(prefix string) {
	text := e.entry.Text
	row := e.entry.CursorRow

	lines := splitLines(text)
	if row < len(lines) {
		lines[row] = prefix + lines[row]
	}
	e.entry.SetText(joinLines(lines))
}

// splitLines splits text into lines
func splitLines(text string) []string {
	if text == "" {
		return []string{""}
	}
	var lines []string
	start := 0
	for i, c := range text {
		if c == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}
	lines = append(lines, text[start:])
	return lines
}

// joinLines joins lines with newlines
func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

// findSelectionStart finds the start position of selected text near cursor
func findSelectionStart(text, selected string, cursorPos int) int {
	// Search backwards from cursor first
	searchStart := cursorPos - len(selected)
	if searchStart < 0 {
		searchStart = 0
	}
	for i := searchStart; i <= cursorPos && i+len(selected) <= len(text); i++ {
		if text[i:i+len(selected)] == selected {
			return i
		}
	}
	// Search forwards from cursor
	for i := cursorPos; i+len(selected) <= len(text); i++ {
		if text[i:i+len(selected)] == selected {
			return i
		}
	}
	return -1
}
