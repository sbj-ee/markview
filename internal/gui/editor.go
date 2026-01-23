package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// MarkdownEditor wraps a multi-line entry with markdown editing support
type MarkdownEditor struct {
	widget.BaseWidget
	entry   *widget.Entry
	canvas  fyne.Canvas
	focused bool

	// lastInsertPos tracks the position after the last programmatic insert
	// This is used when the entry doesn't have focus and cursor position may be stale
	lastInsertPos int

	// lastSelectedText stores the selected text when focus was lost
	// This allows toolbar actions to wrap selected text even after focus shifts
	lastSelectedText string

	OnChanged func(content string)
	// OnKeyEvent is called before key events are processed.
	// Return true to consume the event (prevent default handling).
	OnKeyEvent func(key *fyne.KeyEvent) bool
}

// Ensure MarkdownEditor implements required interfaces
var _ fyne.Focusable = (*MarkdownEditor)(nil)
var _ fyne.Tappable = (*MarkdownEditor)(nil)
var _ desktop.Keyable = (*MarkdownEditor)(nil)

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

// GetCursorPosition returns the current cursor row and column (0-indexed)
func (e *MarkdownEditor) GetCursorPosition() (row, col int) {
	return e.entry.CursorRow, e.entry.CursorColumn
}

// Focus focuses the editor for input
func (e *MarkdownEditor) Focus(canvas fyne.Canvas) {
	e.canvas = canvas
	canvas.Focus(e.entry)
}

// FocusGained is called when the editor receives focus
func (e *MarkdownEditor) FocusGained() {
	e.focused = true
	e.entry.FocusGained()
	// Sync lastInsertPos with cursor position when focus is gained
	e.SyncLastInsertPos()
}

// FocusLost is called when the editor loses focus
func (e *MarkdownEditor) FocusLost() {
	// Save cursor position and selection before losing focus
	// This allows toolbar actions to work correctly even after focus shifts
	e.SyncLastInsertPos()
	e.lastSelectedText = e.entry.SelectedText()

	e.focused = false
	e.entry.FocusLost()
}

// Focused returns whether the editor has focus
func (e *MarkdownEditor) Focused() bool {
	return e.focused
}

// TypedKey handles key events
func (e *MarkdownEditor) TypedKey(key *fyne.KeyEvent) {
	// Check if the key event should be intercepted
	if e.OnKeyEvent != nil && e.OnKeyEvent(key) {
		return // Event was consumed
	}
	e.entry.TypedKey(key)
}

// TypedRune handles rune events - delegates to entry
func (e *MarkdownEditor) TypedRune(r rune) {
	e.entry.TypedRune(r)
}

// Tapped handles tap events to focus the entry
func (e *MarkdownEditor) Tapped(event *fyne.PointEvent) {
	if e.canvas != nil {
		e.canvas.Focus(e.entry)
	}
}

// KeyDown is called when a key is pressed (desktop.Keyable)
func (e *MarkdownEditor) KeyDown(key *fyne.KeyEvent) {
	// Check if the key event should be intercepted
	if e.OnKeyEvent != nil && e.OnKeyEvent(key) {
		return // Event was consumed
	}
	// widget.Entry handles key events through TypedKey, not KeyDown
	// This method is required by desktop.Keyable but Entry doesn't use it
}

// KeyUp is called when a key is released (desktop.Keyable)
func (e *MarkdownEditor) KeyUp(key *fyne.KeyEvent) {
	// widget.Entry doesn't use KeyUp
}

// MinSize returns the minimum size of the editor
func (e *MarkdownEditor) MinSize() fyne.Size {
	// Return a reasonable minimum size for the editor
	min := e.entry.MinSize()
	if min.Width < 200 {
		min.Width = 200
	}
	if min.Height < 100 {
		min.Height = 100
	}
	return min
}

// WrapSelection wraps the selected text with prefix and suffix, or inserts at cursor
func (e *MarkdownEditor) WrapSelection(prefix, suffix string) {
	text := e.entry.Text

	// Always try to get selection from entry first (it may persist after focus lost)
	// Fall back to saved selection only if entry has none
	selected := e.entry.SelectedText()
	if selected == "" {
		selected = e.lastSelectedText
	}

	// Get cursor position from entry (even if not focused, entry may still have valid position)
	col := e.entry.CursorColumn
	row := e.entry.CursorRow
	pos := 0
	lines := splitLines(text)
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	pos += col

	// If entry position seems invalid, use lastInsertPos
	if pos > len(text) || (pos == 0 && e.lastInsertPos > 0) {
		pos = e.lastInsertPos
	}

	if selected != "" {
		// Find selection in text and wrap it
		selStart := findSelectionStart(text, selected, pos)
		if selStart >= 0 {
			newText := text[:selStart] + prefix + selected + suffix + text[selStart+len(selected):]
			e.entry.SetText(newText)
			// Position cursor after the wrapped text
			newPos := selStart + len(prefix) + len(selected) + len(suffix)
			e.lastInsertPos = newPos
			e.lastSelectedText = "" // Clear saved selection after use
			e.setCursorPosition(newPos)
			return
		}
	}

	// No selection, insert at cursor and position between prefix and suffix
	if pos > len(text) {
		pos = len(text)
	}
	if pos < 0 {
		pos = 0
	}
	newText := text[:pos] + prefix + suffix + text[pos:]
	e.entry.SetText(newText)
	// Position cursor between prefix and suffix
	newPos := pos + len(prefix)
	e.lastInsertPos = newPos
	e.setCursorPosition(newPos)
}

// InsertAtCursor inserts text at the current cursor position
func (e *MarkdownEditor) InsertAtCursor(insert string) {
	text := e.entry.Text

	var pos int
	if e.focused {
		// Entry has focus, use its cursor position
		col := e.entry.CursorColumn
		row := e.entry.CursorRow

		// Find position in text
		pos = 0
		lines := splitLines(text)
		for i := 0; i < row && i < len(lines); i++ {
			pos += len(lines[i]) + 1
		}
		pos += col
	} else {
		// Entry doesn't have focus, use last known insert position
		pos = e.lastInsertPos
	}

	if pos > len(text) {
		pos = len(text)
	}
	if pos < 0 {
		pos = 0
	}

	newText := text[:pos] + insert + text[pos:]
	e.entry.SetText(newText)

	// Update last insert position and cursor
	newPos := pos + len(insert)
	e.lastInsertPos = newPos
	e.setCursorPosition(newPos)
}

// InsertAtLineStart inserts text at the beginning of the current line
func (e *MarkdownEditor) InsertAtLineStart(prefix string) {
	text := e.entry.Text
	row := e.entry.CursorRow
	col := e.entry.CursorColumn

	lines := splitLines(text)
	if row < len(lines) {
		lines[row] = prefix + lines[row]
	}
	e.entry.SetText(joinLines(lines))

	// Calculate new cursor position (same row, column shifted by prefix length)
	newPos := 0
	for i := 0; i < row && i < len(lines); i++ {
		newPos += len(lines[i]) + 1
	}
	newPos += col + len(prefix)
	e.setCursorPosition(newPos)
}

// setCursorPosition sets the cursor to a specific position in the text
func (e *MarkdownEditor) setCursorPosition(pos int) {
	text := e.entry.Text
	if pos > len(text) {
		pos = len(text)
	}
	if pos < 0 {
		pos = 0
	}

	// Calculate row and column from position
	row := 0
	col := 0
	for i := 0; i < pos; i++ {
		if i < len(text) && text[i] == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}

	e.entry.CursorRow = row
	e.entry.CursorColumn = col
	e.entry.Refresh()
}

// SyncLastInsertPos updates lastInsertPos from the current cursor position
func (e *MarkdownEditor) SyncLastInsertPos() {
	text := e.entry.Text
	col := e.entry.CursorColumn
	row := e.entry.CursorRow

	pos := 0
	lines := splitLines(text)
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1
	}
	pos += col

	if pos > len(text) {
		pos = len(text)
	}
	e.lastInsertPos = pos
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
	selLen := len(selected)
	textLen := len(text)

	// Check the two most likely positions first:
	// 1. Cursor is at END of selection (most common) - selection starts at cursorPos - selLen
	if cursorPos >= selLen && cursorPos <= textLen {
		start := cursorPos - selLen
		if start+selLen <= textLen && text[start:start+selLen] == selected {
			return start
		}
	}

	// 2. Cursor is at START of selection - selection starts at cursorPos
	if cursorPos >= 0 && cursorPos+selLen <= textLen {
		if text[cursorPos:cursorPos+selLen] == selected {
			return cursorPos
		}
	}

	// Expand search outward from cursor position
	// Search all positions where selection could contain or be near cursor
	for offset := 1; offset <= textLen; offset++ {
		// Search backwards - selection might start before cursor
		start := cursorPos - offset
		if start >= 0 && start+selLen <= textLen {
			if text[start:start+selLen] == selected {
				return start
			}
		}
		// Search forwards - selection might start after cursor
		start = cursorPos + offset
		if start >= 0 && start+selLen <= textLen {
			if text[start:start+selLen] == selected {
				return start
			}
		}
	}

	return -1
}
