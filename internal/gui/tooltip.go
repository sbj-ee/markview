package gui

import (
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// tooltipDelay is how long the pointer must rest on a button before its tooltip
// appears, so sweeping across the toolbar doesn't flash every tooltip.
const tooltipDelay = 500 * time.Millisecond

// cmdKey is the platform-appropriate name for the primary modifier key
// (Cmd on macOS, Ctrl elsewhere), matching fyne.KeyModifierSuper.
var cmdKey = func() string {
	if runtime.GOOS == "darwin" {
		return "Cmd"
	}
	return "Ctrl"
}()

// menuTitle widens a top-level menu title with surrounding spaces so the menu
// bar entries sit further apart. Only applied where Fyne draws the menu bar
// in-window (Linux, etc.); macOS uses the native system menu bar, which must
// keep plain titles.
func menuTitle(name string) string {
	if runtime.GOOS == "darwin" {
		return name
	}
	return "  " + name + "  "
}

// tooltipButton is a widget.Button that displays a text tooltip on hover.
// Fyne 2.7 has no built-in tooltip support, so we provide a lightweight one
// by extending the button's hover handling to show/hide a small popup.
type tooltipButton struct {
	widget.Button
	tooltip   string
	popup     *widget.PopUp
	showTimer *time.Timer
	hovered   bool
}

// newTooltipButton creates an icon button that shows the given tooltip on hover.
func newTooltipButton(icon fyne.Resource, tooltip string, tapped func()) *tooltipButton {
	b := &tooltipButton{tooltip: tooltip}
	b.Icon = icon
	b.OnTapped = tapped
	b.Importance = widget.LowImportance
	b.ExtendBaseWidget(b)
	return b
}

// MouseIn schedules the tooltip after a short delay, in addition to the
// button's normal hover behavior.
func (b *tooltipButton) MouseIn(e *desktop.MouseEvent) {
	b.Button.MouseIn(e)
	b.hovered = true
	b.scheduleTooltip()
}

// MouseOut cancels any pending tooltip and hides a shown one.
func (b *tooltipButton) MouseOut() {
	b.Button.MouseOut()
	b.hovered = false
	b.cancelTooltip()
	b.hideTooltip()
}

// Tapped dismisses the tooltip before invoking the button action so it does not
// linger over a dialog opened by the action.
func (b *tooltipButton) Tapped(e *fyne.PointEvent) {
	b.hovered = false
	b.cancelTooltip()
	b.hideTooltip()
	b.Button.Tapped(e)
}

// scheduleTooltip arms a timer to show the tooltip once the pointer has rested
// on the button for tooltipDelay.
func (b *tooltipButton) scheduleTooltip() {
	if b.tooltip == "" {
		return
	}
	b.cancelTooltip()
	b.showTimer = time.AfterFunc(tooltipDelay, func() {
		// Runs on a timer goroutine; show on the UI thread, and only if the
		// pointer is still over the button.
		fyne.Do(func() {
			if b.hovered {
				b.showTooltip()
			}
		})
	})
}

// cancelTooltip stops a pending show timer.
func (b *tooltipButton) cancelTooltip() {
	if b.showTimer != nil {
		b.showTimer.Stop()
		b.showTimer = nil
	}
}

func (b *tooltipButton) showTooltip() {
	if b.tooltip == "" || b.popup != nil {
		return
	}
	c := fyne.CurrentApp().Driver().CanvasForObject(b)
	if c == nil {
		return
	}

	label := widget.NewLabel(b.tooltip)
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameOverlayBackground))
	bg.StrokeColor = theme.Color(theme.ColorNameSeparator)
	bg.StrokeWidth = 1
	bg.CornerRadius = theme.Padding()
	content := container.NewStack(bg, label)

	b.popup = widget.NewPopUp(content, c)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	b.popup.ShowAtPosition(fyne.NewPos(pos.X, pos.Y+b.Size().Height+2))
}

func (b *tooltipButton) hideTooltip() {
	if b.popup != nil {
		b.popup.Hide()
		b.popup = nil
	}
}

// setActive toggles the button's highlighted state, used to indicate that the
// panel it controls is currently visible.
func (b *tooltipButton) setActive(active bool) {
	imp := widget.LowImportance
	if active {
		imp = widget.HighImportance
	}
	if b.Importance != imp {
		b.Importance = imp
		b.Refresh()
	}
}
