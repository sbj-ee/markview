package gui

import (
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// OutlineEntry represents a heading in the document outline
type OutlineEntry struct {
	Level int
	Text  string
	Line  int
}

// Outline is a widget that shows document structure based on headings
type Outline struct {
	widget.BaseWidget
	entries  []OutlineEntry
	list     *widget.List
	onSelect func(line int)
}

// NewOutline creates a new outline widget
func NewOutline(onSelect func(line int)) *Outline {
	o := &Outline{
		entries:  make([]OutlineEntry, 0),
		onSelect: onSelect,
	}
	o.list = widget.NewList(
		func() int { return len(o.entries) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(o.entries) {
				entry := o.entries[id]
				// Indent based on heading level
				indent := strings.Repeat("  ", entry.Level-1)
				label.SetText(indent + entry.Text)
			}
		},
	)
	o.list.OnSelected = func(id widget.ListItemID) {
		if id < len(o.entries) && o.onSelect != nil {
			o.onSelect(o.entries[id].Line)
		}
	}
	o.ExtendBaseWidget(o)
	return o
}

// CreateRenderer implements fyne.Widget
func (o *Outline) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(o.list)
}

// UpdateFromText parses markdown text and updates the outline
func (o *Outline) UpdateFromText(text string) {
	o.entries = parseHeadings(text)
	o.list.Refresh()
}

// GetContainer returns the outline in a scrollable container with a header
func (o *Outline) GetContainer() *fyne.Container {
	header := widget.NewLabel("Outline")
	header.TextStyle = fyne.TextStyle{Bold: true}
	return container.NewBorder(header, nil, nil, nil, container.NewScroll(o.list))
}

// parseHeadings extracts headings from markdown text
func parseHeadings(text string) []OutlineEntry {
	var entries []OutlineEntry
	lines := strings.Split(text, "\n")

	// Match markdown headings (# to ######)
	headingRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	for i, line := range lines {
		matches := headingRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			level := len(matches[1])
			text := strings.TrimSpace(matches[2])
			entries = append(entries, OutlineEntry{
				Level: level,
				Text:  text,
				Line:  i,
			})
		}
	}

	return entries
}
