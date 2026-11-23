package gui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// TestTooltipScheduling verifies the hover-delay timer is armed and cancelled
// correctly, and that buttons without a tooltip never schedule one.
func TestTooltipScheduling(t *testing.T) {
	b := newTooltipButton(theme.InfoIcon(), "Bold", nil)
	b.hovered = true

	b.scheduleTooltip()
	if b.showTimer == nil {
		t.Fatal("scheduleTooltip should arm a timer")
	}
	if b.popup != nil {
		t.Error("tooltip should not be shown before the delay elapses")
	}

	b.cancelTooltip()
	if b.showTimer != nil {
		t.Error("cancelTooltip should clear the timer")
	}

	empty := newTooltipButton(theme.InfoIcon(), "", nil)
	empty.hovered = true
	empty.scheduleTooltip()
	if empty.showTimer != nil {
		t.Error("a button with no tooltip should not schedule a timer")
	}
}

func TestToolbarClickFiresWhileTooltipShown(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()

	clicked := 0
	b := newTooltipButton(theme.InfoIcon(), "Open File", func() { clicked++ })
	b.Resize(fyne.NewSize(40, 40))

	w := app.NewWindow("tooltip")
	w.SetContent(container.NewWithoutLayout(b))
	w.Resize(fyne.NewSize(240, 120))
	w.Show()

	b.hovered = true
	b.showTooltip()
	if b.popup == nil {
		t.Fatal("showTooltip should display a tooltip")
	}

	// A canvas tap at the button, with the tooltip overlay up, must fire the
	// action once. widget.PopUp is a full-canvas Tappable, so this used to
	// dismiss the tooltip and leave clicked == 0.
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	test.TapCanvas(w.Canvas(), pos.Add(fyne.NewPos(b.Size().Width/2, b.Size().Height/2)))

	if clicked != 1 {
		t.Fatalf("button action fired %d times, want 1", clicked)
	}
	if b.popup != nil {
		t.Error("tooltip should be dismissed after the click")
	}
}
