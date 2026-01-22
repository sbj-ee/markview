package markdown

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Blockquote is a custom widget that displays quoted text with a left border
type Blockquote struct {
	widget.BaseWidget
	label       *widget.Label
	padding     float32
	borderWidth float32
}

// NewBlockquote creates a new blockquote widget with the given text
func NewBlockquote(text string) *Blockquote {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle = fyne.TextStyle{Italic: true}

	bq := &Blockquote{
		label:       label,
		padding:     12,
		borderWidth: 4,
	}
	bq.ExtendBaseWidget(bq)
	return bq
}

// MinSize returns the minimum size of the blockquote
func (b *Blockquote) MinSize() fyne.Size {
	labelSize := b.label.MinSize()
	return fyne.NewSize(
		labelSize.Width+b.padding*2+b.borderWidth,
		labelSize.Height+b.padding*2,
	)
}

// CreateRenderer creates the renderer for the blockquote widget
func (b *Blockquote) CreateRenderer() fyne.WidgetRenderer {
	return &blockquoteRenderer{
		blockquote: b,
	}
}

// blockquoteRenderer is the renderer for the Blockquote widget
type blockquoteRenderer struct {
	blockquote *Blockquote
	background *canvas.Rectangle
	border     *canvas.Rectangle
}

func (r *blockquoteRenderer) Destroy() {}

func (r *blockquoteRenderer) Layout(size fyne.Size) {
	if r.background == nil {
		r.background = canvas.NewRectangle(r.getBackgroundColor())
		r.background.CornerRadius = 2
	}
	r.background.Resize(size)

	if r.border == nil {
		r.border = canvas.NewRectangle(r.getBorderColor())
		r.border.CornerRadius = 2
	}
	r.border.Move(fyne.NewPos(0, 0))
	r.border.Resize(fyne.NewSize(r.blockquote.borderWidth, size.Height))

	// Position the label with padding, accounting for border
	padding := r.blockquote.padding
	borderWidth := r.blockquote.borderWidth
	r.blockquote.label.Move(fyne.NewPos(borderWidth+padding, padding))
	r.blockquote.label.Resize(fyne.NewSize(
		size.Width-padding*2-borderWidth,
		size.Height-padding*2,
	))
}

func (r *blockquoteRenderer) MinSize() fyne.Size {
	return r.blockquote.MinSize()
}

func (r *blockquoteRenderer) Objects() []fyne.CanvasObject {
	if r.background == nil {
		r.background = canvas.NewRectangle(r.getBackgroundColor())
		r.background.CornerRadius = 2
	}
	if r.border == nil {
		r.border = canvas.NewRectangle(r.getBorderColor())
		r.border.CornerRadius = 2
	}
	return []fyne.CanvasObject{r.background, r.border, r.blockquote.label}
}

func (r *blockquoteRenderer) Refresh() {
	if r.background != nil {
		r.background.FillColor = r.getBackgroundColor()
		r.background.Refresh()
	}
	if r.border != nil {
		r.border.FillColor = r.getBorderColor()
		r.border.Refresh()
	}
	r.blockquote.label.Refresh()
}

// getBackgroundColor returns the appropriate background color based on theme
func (r *blockquoteRenderer) getBackgroundColor() color.Color {
	// Subtle background for blockquotes
	return color.RGBA{R: 50, G: 52, B: 64, A: 128}
}

// getBorderColor returns the border color for the left edge
func (r *blockquoteRenderer) getBorderColor() color.Color {
	// Use primary color for the border to make it stand out
	c := theme.Current().Color(theme.ColorNamePrimary, theme.VariantDark)
	if c == nil {
		return color.RGBA{R: 139, G: 233, B: 253, A: 255}
	}
	return c
}
