package gui

import (
	"testing"

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
