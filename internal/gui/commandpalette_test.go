package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func TestShortcutLabel(t *testing.T) {
	tests := []struct {
		name     string
		shortcut fyne.Shortcut
		want     string
	}{
		{"nil", nil, ""},
		{"non-custom", &fyne.ShortcutCopy{}, ""},
		{
			"super",
			&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierSuper},
			cmdKey + "+S",
		},
		{
			"super shift",
			&desktop.CustomShortcut{KeyName: fyne.KeyP, Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift},
			cmdKey + "+Shift+P",
		},
		{
			"control only",
			&desktop.CustomShortcut{KeyName: fyne.KeyP, Modifier: fyne.KeyModifierControl},
			"Ctrl+P",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortcutLabel(tt.shortcut); got != tt.want {
				t.Errorf("shortcutLabel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPaletteEntryNavigationKeys(t *testing.T) {
	e := newPaletteEntry()
	var up, down, accept, cancel int
	e.onUp = func() { up++ }
	e.onDown = func() { down++ }
	e.onAccept = func() { accept++ }
	e.onCancel = func() { cancel++ }

	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})

	if up != 1 || down != 1 || cancel != 1 {
		t.Errorf("nav callbacks: up=%d down=%d cancel=%d, want 1 each", up, down, cancel)
	}
	if accept != 2 {
		t.Errorf("accept callback fired %d times, want 2 (Return and Enter)", accept)
	}
}

func TestNewCommandPaletteStoresCommands(t *testing.T) {
	called := false
	cmds := []command{
		{label: "File: Save", shortcut: "Cmd+S", action: func() { called = true }},
		{label: "View: Split View"},
	}
	cp := NewCommandPalette(cmds)
	if len(cp.commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(cp.commands))
	}
	cp.commands[0].action()
	if !called {
		t.Error("stored command action was not invoked")
	}
}
