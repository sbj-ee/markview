package markdown

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MathBlock is a custom widget for displaying math formulas
type MathBlock struct {
	widget.BaseWidget
	content string
}

// NewMathBlock creates a new math block widget
func NewMathBlock(content string) *MathBlock {
	m := &MathBlock{content: content}
	m.ExtendBaseWidget(m)
	return m
}

// CreateRenderer implements fyne.Widget
func (m *MathBlock) CreateRenderer() fyne.WidgetRenderer {
	// Create background with slightly different color for math
	bg := canvas.NewRectangle(color.NRGBA{R: 40, G: 45, B: 60, A: 255})
	bg.CornerRadius = 6

	// Create left border indicator (purple for math)
	border := canvas.NewRectangle(color.NRGBA{R: 147, G: 112, B: 219, A: 255})

	// Create the math content label
	label := widget.NewRichText(&widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: "∫ " + m.content,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Monospace: true},
					ColorName: theme.ColorNameForeground,
				},
			},
		},
	})

	return &mathBlockRenderer{
		mathBlock: m,
		bg:        bg,
		border:    border,
		label:     label,
	}
}

type mathBlockRenderer struct {
	mathBlock *MathBlock
	bg        *canvas.Rectangle
	border    *canvas.Rectangle
	label     *widget.RichText
}

func (r *mathBlockRenderer) Layout(size fyne.Size) {
	// Background fills the entire area
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	// Left border is 4px wide
	r.border.Resize(fyne.NewSize(4, size.Height))
	r.border.Move(fyne.NewPos(0, 0))

	// Content is padded from the left border
	padding := float32(16)
	r.label.Resize(fyne.NewSize(size.Width-padding-8, size.Height))
	r.label.Move(fyne.NewPos(padding, 8))
}

func (r *mathBlockRenderer) MinSize() fyne.Size {
	labelMin := r.label.MinSize()
	return fyne.NewSize(labelMin.Width+24, labelMin.Height+16)
}

func (r *mathBlockRenderer) Refresh() {
	r.bg.Refresh()
	r.border.Refresh()
	r.label.Refresh()
}

func (r *mathBlockRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.border, r.label}
}

func (r *mathBlockRenderer) Destroy() {}

// InlineMath creates a container for inline math display
func InlineMath(content string) *fyne.Container {
	label := widget.NewRichText(&widget.TextSegment{
		Text: "⟨" + content + "⟩",
		Style: widget.RichTextStyle{
			TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
		},
	})
	return container.NewHBox(label)
}
