package gui

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// cmdKey is the platform-appropriate name for the primary modifier key
// (Cmd on macOS, Ctrl elsewhere), matching fyne.KeyModifierSuper.
var cmdKey = func() string {
	if runtime.GOOS == "darwin" {
		return "Cmd"
	}
	return "Ctrl"
}()

// tooltipButton is a widget.Button that displays a text tooltip on hover.
// Fyne 2.7 has no built-in tooltip support, so we provide a lightweight one
// by extending the button's hover handling to show/hide a small popup.
type tooltipButton struct {
	widget.Button
	tooltip string
	popup   *widget.PopUp
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

// MouseIn shows the tooltip in addition to the button's normal hover behavior.
func (b *tooltipButton) MouseIn(e *desktop.MouseEvent) {
	b.Button.MouseIn(e)
	b.showTooltip()
}

// MouseOut hides the tooltip and restores the button's normal state.
func (b *tooltipButton) MouseOut() {
	b.Button.MouseOut()
	b.hideTooltip()
}

// Tapped hides the tooltip before invoking the button action so it does not
// linger over a dialog opened by the action.
func (b *tooltipButton) Tapped(e *fyne.PointEvent) {
	b.hideTooltip()
	b.Button.Tapped(e)
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
