package markdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Block widgets in the rendered document live in a VBox inside a scroll
// container. If any single block reports a MinSize wider than the viewport,
// the whole VBox is laid out at that width and every word-wrapped paragraph
// wraps at the wider width too, pushing text off the right edge of the
// window. The widgets in this file keep their minimum width small and adapt
// to whatever width they are given instead.

const (
	// minTableColWidth is the narrowest a table column will shrink to.
	minTableColWidth float32 = 60
	// minImageWidth is the narrowest an image will shrink to.
	minImageWidth float32 = 100
)

// fitColumnWidths scales natural column widths down so that they fit in
// avail. Columns never shrink below minTableColWidth (or their natural width,
// if that is smaller). If the columns already fit, the natural widths are
// returned unchanged.
func fitColumnWidths(natural []float32, avail float32) []float32 {
	out := make([]float32, len(natural))
	var total float32
	for _, w := range natural {
		total += w
	}
	if avail <= 0 || total <= avail {
		copy(out, natural)
		return out
	}
	scale := avail / total
	for i, w := range natural {
		out[i] = fyne.Max(w*scale, fyne.Min(w, minTableColWidth))
	}
	return out
}

// ResponsiveTable holds a grid of fixed-width cells and recalculates the
// cell widths whenever its own width changes, so the table shrinks (and its
// cells wrap) instead of overflowing the viewport.
type ResponsiveTable struct {
	widget.BaseWidget
	content *fyne.Container
	cells   [][]*FixedWidthContainer
	natural []float32
}

// NewResponsiveTable creates a table wrapper. content is the VBox of rows,
// cells are the per-row cell widgets, and natural holds each column's
// preferred width.
func NewResponsiveTable(content *fyne.Container, cells [][]*FixedWidthContainer, natural []float32) *ResponsiveTable {
	t := &ResponsiveTable{content: content, cells: cells, natural: natural}
	t.ExtendBaseWidget(t)
	return t
}

// applyWidth sets every cell's width for the given available table width.
func (t *ResponsiveTable) applyWidth(avail float32) {
	if avail > 0 && len(t.natural) > 1 {
		// Rows are HBoxes, which put theme padding between cells.
		avail -= theme.Padding() * float32(len(t.natural)-1)
	}
	widths := fitColumnWidths(t.natural, avail)
	for _, row := range t.cells {
		for i, cell := range row {
			if i < len(widths) {
				cell.width = widths[i]
			}
		}
	}
}

// MinSize returns a small minimum width and the height the table needs at
// its current width.
func (t *ResponsiveTable) MinSize() fyne.Size {
	t.applyWidth(t.Size().Width)
	var minW float32
	for _, w := range t.natural {
		minW += fyne.Min(w, minTableColWidth)
	}
	return fyne.NewSize(minW, t.content.MinSize().Height)
}

// CreateRenderer implements fyne.Widget.
func (t *ResponsiveTable) CreateRenderer() fyne.WidgetRenderer {
	return &responsiveTableRenderer{table: t}
}

type responsiveTableRenderer struct {
	table *ResponsiveTable
}

func (r *responsiveTableRenderer) Layout(size fyne.Size) {
	r.table.applyWidth(size.Width)
	r.table.content.Move(fyne.NewPos(0, 0))
	r.table.content.Resize(size)
}

func (r *responsiveTableRenderer) MinSize() fyne.Size {
	return r.table.MinSize()
}

func (r *responsiveTableRenderer) Refresh() {
	r.table.content.Refresh()
}

func (r *responsiveTableRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.table.content}
}

func (r *responsiveTableRenderer) Destroy() {}

// FitWidth shows an object at its natural size when there is room, and
// scales it down (keeping its aspect ratio) when the available width is
// narrower. Used for images.
type FitWidth struct {
	widget.BaseWidget
	obj     fyne.CanvasObject
	natural fyne.Size
}

// NewFitWidth wraps obj, whose preferred size is natural.
func NewFitWidth(obj fyne.CanvasObject, natural fyne.Size) *FitWidth {
	f := &FitWidth{obj: obj, natural: natural}
	f.ExtendBaseWidget(f)
	return f
}

// scaledSize returns the size to display obj at for the given width.
func (f *FitWidth) scaledSize(avail float32) fyne.Size {
	if avail <= 0 || f.natural.Width <= 0 || avail >= f.natural.Width {
		return f.natural
	}
	return fyne.NewSize(avail, f.natural.Height*avail/f.natural.Width)
}

// MinSize returns a small minimum width and the height needed at the
// current width.
func (f *FitWidth) MinSize() fyne.Size {
	s := f.scaledSize(f.Size().Width)
	return fyne.NewSize(fyne.Min(f.natural.Width, minImageWidth), s.Height)
}

// CreateRenderer implements fyne.Widget.
func (f *FitWidth) CreateRenderer() fyne.WidgetRenderer {
	return &fitWidthRenderer{fit: f}
}

type fitWidthRenderer struct {
	fit *FitWidth
}

func (r *fitWidthRenderer) Layout(size fyne.Size) {
	r.fit.obj.Move(fyne.NewPos(0, 0))
	r.fit.obj.Resize(r.fit.scaledSize(size.Width))
}

func (r *fitWidthRenderer) MinSize() fyne.Size {
	return r.fit.MinSize()
}

func (r *fitWidthRenderer) Refresh() {
	r.fit.obj.Refresh()
}

func (r *fitWidthRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.fit.obj}
}

func (r *fitWidthRenderer) Destroy() {}
