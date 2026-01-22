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
