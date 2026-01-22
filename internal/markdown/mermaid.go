package markdown

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MermaidBlock is a custom widget for displaying Mermaid diagram code
type MermaidBlock struct {
	widget.BaseWidget
	content string
}

// NewMermaidBlock creates a new mermaid block widget
func NewMermaidBlock(content string) *MermaidBlock {
	m := &MermaidBlock{content: content}
	m.ExtendBaseWidget(m)
	return m
}

// CreateRenderer implements fyne.Widget
func (m *MermaidBlock) CreateRenderer() fyne.WidgetRenderer {
	// Create background with slightly different color for mermaid
	bg := canvas.NewRectangle(color.NRGBA{R: 35, G: 50, B: 55, A: 255})
	bg.CornerRadius = 6

	// Create left border indicator (green for mermaid)
	border := canvas.NewRectangle(color.NRGBA{R: 80, G: 200, B: 120, A: 255})

	// Create header label
	header := widget.NewRichText(&widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: "Mermaid Diagram",
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true, Italic: true},
					ColorName: theme.ColorNameForeground,
				},
			},
		},
	})

	// Create the mermaid content label
	label := widget.NewRichText(&widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: m.content,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Monospace: true},
					ColorName: theme.ColorNameForeground,
				},
			},
		},
	})

	return &mermaidBlockRenderer{
		mermaidBlock: m,
		bg:           bg,
		border:       border,
		header:       header,
		label:        label,
	}
}

type mermaidBlockRenderer struct {
	mermaidBlock *MermaidBlock
	bg           *canvas.Rectangle
	border       *canvas.Rectangle
	header       *widget.RichText
	label        *widget.RichText
}

func (r *mermaidBlockRenderer) Layout(size fyne.Size) {
	// Background fills the entire area
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	// Left border is 4px wide
	r.border.Resize(fyne.NewSize(4, size.Height))
	r.border.Move(fyne.NewPos(0, 0))

	// Header at top
	padding := float32(16)
	headerSize := r.header.MinSize()
	r.header.Resize(fyne.NewSize(size.Width-padding-8, headerSize.Height))
	r.header.Move(fyne.NewPos(padding, 8))

	// Content below header
	r.label.Resize(fyne.NewSize(size.Width-padding-8, size.Height-headerSize.Height-24))
	r.label.Move(fyne.NewPos(padding, headerSize.Height+16))
}

func (r *mermaidBlockRenderer) MinSize() fyne.Size {
	headerMin := r.header.MinSize()
	labelMin := r.label.MinSize()
	return fyne.NewSize(
		fyne.Max(headerMin.Width, labelMin.Width)+24,
		headerMin.Height+labelMin.Height+32,
	)
}

func (r *mermaidBlockRenderer) Refresh() {
	r.bg.Refresh()
	r.border.Refresh()
	r.header.Refresh()
	r.label.Refresh()
}

func (r *mermaidBlockRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.border, r.header, r.label}
}

func (r *mermaidBlockRenderer) Destroy() {}

// IsMermaidCodeBlock checks if a code block language indicates mermaid
func IsMermaidCodeBlock(language string) bool {
	return language == "mermaid" || language == "mmd"
}

// MermaidPreviewContainer creates a container with the diagram and preview hint
func MermaidPreviewContainer(content string) *fyne.Container {
	block := NewMermaidBlock(content)
	hint := widget.NewLabel("Mermaid diagrams render in exported HTML/PDF")
	hint.TextStyle = fyne.TextStyle{Italic: true}
	return container.NewVBox(block, hint)
}
