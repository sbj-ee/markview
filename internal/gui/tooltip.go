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
// by extending the button's hover handling to show/hide a small bubble.
type tooltipButton struct {
	widget.Button
	tooltip   string
	popup     *tooltipOverlay
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

// MouseOut cancels a pending tooltip. If the bubble is already showing, leave
// it up: showing a canvas overlay steals hover from the button, so a MouseOut
// here does not mean the pointer left the icon. The overlay hides the tooltip
// once the pointer actually moves off the button (or the bubble).
func (b *tooltipButton) MouseOut() {
	b.Button.MouseOut()
	if b.popup != nil {
		return
	}
	b.hovered = false
	b.cancelTooltip()
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

	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	pos = fyne.NewPos(pos.X, pos.Y+b.Size().Height+2)
	pos = clampTipPos(c, pos, content.MinSize())

	b.popup = newTooltipOverlay(content, c, pos, b)
	b.popup.Show()
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

func clampTipPos(c fyne.Canvas, pos fyne.Position, size fyne.Size) fyne.Position {
	cs := c.Size()
	if pos.X+size.Width > cs.Width {
		pos.X = cs.Width - size.Width
		if pos.X < 0 {
			pos.X = 0
		}
	}
	if pos.Y+size.Height > cs.Height {
		pos.Y = cs.Height - size.Height
		if pos.Y < 0 {
			pos.Y = 0
		}
	}
	return pos
}

func objectContains(obj fyne.CanvasObject, abs fyne.Position) bool {
	if obj == nil || fyne.CurrentApp() == nil {
		return false
	}
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
	sz := obj.Size()
	return abs.X >= pos.X && abs.Y >= pos.Y && abs.X < pos.X+sz.Width && abs.Y < pos.Y+sz.Height
}

// tooltipOverlay draws a tooltip bubble on the canvas overlay stack.
//
// widget.PopUp cannot be used for this: it is always sized to the whole canvas
// and implements Tappable, so a click on a toolbar button is consumed to
// dismiss the popup and never reaches the button. Users have to click several
// times (each click only hides or retriggers the tooltip). Fyne also only
// hit-tests window content when no overlay is present, so a non-tappable
// overlay would swallow clicks too. This overlay covers the canvas, then on
// tap hides itself and forwards the event to the button underneath.
type tooltipOverlay struct {
	widget.BaseWidget
	canvas  fyne.Canvas
	content fyne.CanvasObject
	pos     fyne.Position
	owner   *tooltipButton
	shown   bool
}

func newTooltipOverlay(content fyne.CanvasObject, c fyne.Canvas, pos fyne.Position, owner *tooltipButton) *tooltipOverlay {
	o := &tooltipOverlay{canvas: c, content: content, pos: pos, owner: owner}
	o.ExtendBaseWidget(o)
	return o
}

func (o *tooltipOverlay) Show() {
	if o.canvas == nil {
		return
	}
	if !o.shown {
		o.canvas.Overlays().Add(o)
		o.shown = true
	}
	o.Resize(o.canvas.Size())
	o.Refresh()
	o.BaseWidget.Show()
}

func (o *tooltipOverlay) Hide() {
	if o.shown && o.canvas != nil {
		o.canvas.Overlays().Remove(o)
		o.shown = false
	}
	o.BaseWidget.Hide()
}

// Tapped dismisses the tooltip and delivers the click to the control under
// the pointer, so the first press on a toolbar icon still runs its action.
func (o *tooltipOverlay) Tapped(e *fyne.PointEvent) {
	owner := o.owner
	abs := e.AbsolutePosition
	if (abs.X == 0 && abs.Y == 0) && e.Position != (fyne.Position{}) {
		abs = e.Position
	}
	o.dismiss()
	if owner == nil || !objectContains(owner, abs) {
		return
	}
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(owner)
	owner.Tapped(&fyne.PointEvent{
		Position:         fyne.NewPos(abs.X-pos.X, abs.Y-pos.Y),
		AbsolutePosition: abs,
	})
}

// MouseIn and MouseMoved keep the tooltip visible while the pointer stays on
// the button or the bubble, and hide it once it leaves both.
func (o *tooltipOverlay) MouseIn(e *desktop.MouseEvent) {
	o.MouseMoved(e)
}

func (o *tooltipOverlay) MouseMoved(e *desktop.MouseEvent) {
	if o.pointerOverTipOrButton(e.AbsolutePosition) {
		return
	}
	o.dismiss()
}

func (o *tooltipOverlay) MouseOut() {
	o.dismiss()
}

func (o *tooltipOverlay) pointerOverTipOrButton(abs fyne.Position) bool {
	if objectContains(o.owner, abs) {
		return true
	}
	sz := o.content.Size()
	return abs.X >= o.pos.X && abs.Y >= o.pos.Y &&
		abs.X < o.pos.X+sz.Width && abs.Y < o.pos.Y+sz.Height
}

func (o *tooltipOverlay) dismiss() {
	if o.owner != nil {
		o.owner.hovered = false
		o.owner.cancelTooltip()
		o.owner.hideTooltip()
		return
	}
	o.Hide()
}

func (o *tooltipOverlay) CreateRenderer() fyne.WidgetRenderer {
	return &tooltipOverlayRenderer{overlay: o, objects: []fyne.CanvasObject{o.content}}
}

type tooltipOverlayRenderer struct {
	overlay *tooltipOverlay
	objects []fyne.CanvasObject
}

func (r *tooltipOverlayRenderer) Layout(fyne.Size) {
	sz := r.overlay.content.MinSize()
	r.overlay.content.Resize(sz)
	r.overlay.content.Move(r.overlay.pos)
}

func (r *tooltipOverlayRenderer) MinSize() fyne.Size {
	return r.overlay.content.MinSize()
}

func (r *tooltipOverlayRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *tooltipOverlayRenderer) Destroy() {}

func (r *tooltipOverlayRenderer) Refresh() {
	if r.overlay.canvas != nil && r.overlay.Size() != r.overlay.canvas.Size() {
		r.overlay.Resize(r.overlay.canvas.Size())
	}
	r.Layout(r.overlay.Size())
	r.overlay.content.Refresh()
}
