package gui

import (
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// command is a single named action that can be run from the command palette.
type command struct {
	label    string
	shortcut string // optional display hint, e.g. "Cmd+Shift+P"
	action   func()
}

// paletteEntry is a single-line search entry that forwards navigation keys
// (Up/Down/Enter/Escape) to the command palette instead of editing text.
type paletteEntry struct {
	widget.Entry
	onUp     func()
	onDown   func()
	onAccept func()
	onCancel func()
}

func newPaletteEntry() *paletteEntry {
	e := &paletteEntry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *paletteEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyUp:
		if e.onUp != nil {
			e.onUp()
		}
	case fyne.KeyDown:
		if e.onDown != nil {
			e.onDown()
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if e.onAccept != nil {
			e.onAccept()
		}
	case fyne.KeyEscape:
		if e.onCancel != nil {
			e.onCancel()
		}
	default:
		e.Entry.TypedKey(key)
	}
}

// CommandPalette presents a fuzzy-filterable list of all app commands and
// runs the chosen one. It is a keyboard-first counterpart to the menu bar.
type CommandPalette struct {
	commands []command
}

// NewCommandPalette creates a palette over the given commands.
func NewCommandPalette(commands []command) *CommandPalette {
	return &CommandPalette{commands: commands}
}

// Show displays the command palette dialog.
func (cp *CommandPalette) Show(window fyne.Window) {
	// Sorted master list (stable so equal labels keep input order).
	all := make([]command, len(cp.commands))
	copy(all, cp.commands)
	sort.SliceStable(all, func(i, j int) bool {
		return strings.ToLower(all[i].label) < strings.ToLower(all[j].label)
	})

	filtered := make([]command, len(all))
	copy(filtered, all)

	selected := 0
	// suppressRun distinguishes keyboard navigation (highlight only) from a
	// mouse click (which should run the command) inside list.OnSelected.
	suppressRun := false

	var d dialog.Dialog
	var list *widget.List

	list = widget.NewList(
		func() int { return len(filtered) },
		func() fyne.CanvasObject {
			name := widget.NewLabel("command")
			shortcut := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
			return container.NewHBox(name, layout.NewSpacer(), shortcut)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(filtered) {
				return
			}
			box := obj.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(filtered[id].label)
			box.Objects[2].(*widget.Label).SetText(filtered[id].shortcut)
		},
	)

	run := func(id int) {
		if id < 0 || id >= len(filtered) {
			return
		}
		action := filtered[id].action
		d.Hide()
		if action != nil {
			action()
		}
	}

	list.OnSelected = func(id widget.ListItemID) {
		selected = id
		if suppressRun {
			suppressRun = false
			return
		}
		run(id)
	}

	move := func(delta int) {
		if len(filtered) == 0 {
			return
		}
		next := selected + delta
		if next < 0 {
			next = 0
		}
		if next >= len(filtered) {
			next = len(filtered) - 1
		}
		selected = next
		suppressRun = true
		list.Select(next)
	}

	entry := newPaletteEntry()
	entry.SetPlaceHolder("Type a command…")
	entry.onDown = func() { move(1) }
	entry.onUp = func() { move(-1) }
	entry.onAccept = func() { run(selected) }
	entry.onCancel = func() { d.Hide() }

	entry.OnChanged = func(query string) {
		query = strings.ToLower(strings.TrimSpace(query))
		filtered = filtered[:0]
		for _, c := range all {
			if query == "" || fuzzyMatch(strings.ToLower(c.label), query) {
				filtered = append(filtered, c)
			}
		}
		selected = 0
		list.UnselectAll()
		list.Refresh()
		if len(filtered) > 0 {
			suppressRun = true
			list.Select(0)
		}
	}

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(520, 400))
	content := container.NewBorder(entry, nil, nil, nil, scroll)

	d = dialog.NewCustom("Command Palette", "Cancel", content, window)
	d.Resize(fyne.NewSize(560, 500))
	d.Show()

	// Pre-select the first command and focus the search field.
	if len(filtered) > 0 {
		suppressRun = true
		list.Select(0)
	}
	window.Canvas().Focus(entry)
}
