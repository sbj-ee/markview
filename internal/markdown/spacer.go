package markdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Spacer is a custom widget that creates vertical space
type Spacer struct {
	widget.BaseWidget
	height float32
}

// NewSpacer creates a new spacer with the specified height
func NewSpacer(height float32) *Spacer {
	s := &Spacer{height: height}
	s.ExtendBaseWidget(s)
	return s
}

// MinSize returns the minimum size of the spacer
func (s *Spacer) MinSize() fyne.Size {
	return fyne.NewSize(1, s.height)
}

// CreateRenderer creates the renderer for the spacer widget
func (s *Spacer) CreateRenderer() fyne.WidgetRenderer {
	return &spacerRenderer{spacer: s}
}

// spacerRenderer is the renderer for the Spacer widget
type spacerRenderer struct {
	spacer *Spacer
}

func (r *spacerRenderer) Destroy() {}

func (r *spacerRenderer) Layout(size fyne.Size) {}

func (r *spacerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, r.spacer.height)
}

func (r *spacerRenderer) Objects() []fyne.CanvasObject {
	// Return an invisible rectangle to take up space
	rect := canvas.NewRectangle(nil)
	rect.SetMinSize(fyne.NewSize(1, r.spacer.height))
	return []fyne.CanvasObject{rect}
}

func (r *spacerRenderer) Refresh() {}
