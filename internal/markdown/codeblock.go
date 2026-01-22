package markdown

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CodeBlock is a custom widget that displays code with a background
type CodeBlock struct {
	widget.BaseWidget
	richText *widget.RichText
	padding  float32
}

// NewCodeBlock creates a new code block widget with the given segments
func NewCodeBlock(segments []widget.RichTextSegment) *CodeBlock {
	rt := widget.NewRichText(&widget.ParagraphSegment{
		Texts: segments,
	})
	rt.Wrapping = fyne.TextWrapOff

	cb := &CodeBlock{
		richText: rt,
		padding:  8,
	}
	cb.ExtendBaseWidget(cb)
	return cb
}

// MinSize returns the minimum size of the code block
func (c *CodeBlock) MinSize() fyne.Size {
	rtSize := c.richText.MinSize()
	return fyne.NewSize(
		rtSize.Width+c.padding*2,
		rtSize.Height+c.padding*2,
	)
}

// CreateRenderer creates the renderer for the code block widget
func (c *CodeBlock) CreateRenderer() fyne.WidgetRenderer {
	return &codeBlockRenderer{
		codeBlock: c,
	}
}

// codeBlockRenderer is the renderer for the CodeBlock widget
type codeBlockRenderer struct {
	codeBlock  *CodeBlock
	background *canvas.Rectangle
}

func (r *codeBlockRenderer) Destroy() {}

func (r *codeBlockRenderer) Layout(size fyne.Size) {
	if r.background == nil {
		r.background = canvas.NewRectangle(r.getBackgroundColor())
		r.background.CornerRadius = 4
	}
	r.background.Resize(size)

	// Position the rich text with padding
	padding := r.codeBlock.padding
	r.codeBlock.richText.Move(fyne.NewPos(padding, padding))
	r.codeBlock.richText.Resize(fyne.NewSize(
		size.Width-padding*2,
		size.Height-padding*2,
	))
}

func (r *codeBlockRenderer) MinSize() fyne.Size {
	return r.codeBlock.MinSize()
}

func (r *codeBlockRenderer) Objects() []fyne.CanvasObject {
	if r.background == nil {
		r.background = canvas.NewRectangle(r.getBackgroundColor())
		r.background.CornerRadius = 4
	}
	return []fyne.CanvasObject{r.background, r.codeBlock.richText}
}

func (r *codeBlockRenderer) Refresh() {
	if r.background != nil {
		r.background.FillColor = r.getBackgroundColor()
		r.background.Refresh()
	}
	r.codeBlock.richText.Refresh()
}

// getBackgroundColor returns the appropriate background color based on theme
func (r *codeBlockRenderer) getBackgroundColor() color.Color {
	// Use InputBackground which is slightly darker than the main background
	bg := theme.Current().Color(theme.ColorNameInputBackground, theme.VariantDark)
	if bg == nil {
		// Fallback to a dark gray
		return color.RGBA{R: 40, G: 42, B: 48, A: 255}
	}
	return bg
}
