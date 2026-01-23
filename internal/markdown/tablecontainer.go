package markdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TableContainer is a container that wraps a table with a fixed size
type TableContainer struct {
	widget.BaseWidget
	table  *widget.Table
	width  float32
	height float32
}

// NewTableContainer creates a new table container with fixed dimensions
func NewTableContainer(table *widget.Table, width, height float32) *TableContainer {
	tc := &TableContainer{
		table:  table,
		width:  width,
		height: height,
	}
	tc.ExtendBaseWidget(tc)
	return tc
}

// CreateRenderer implements fyne.Widget
func (tc *TableContainer) CreateRenderer() fyne.WidgetRenderer {
	return &tableContainerRenderer{
		container: tc,
	}
}

// MinSize returns the minimum size of the table container
func (tc *TableContainer) MinSize() fyne.Size {
	return fyne.NewSize(tc.width, tc.height)
}

type tableContainerRenderer struct {
	container *TableContainer
}

func (r *tableContainerRenderer) Layout(size fyne.Size) {
	r.container.table.Resize(fyne.NewSize(r.container.width, r.container.height))
	r.container.table.Move(fyne.NewPos(0, 0))
}

func (r *tableContainerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.container.width, r.container.height)
}

func (r *tableContainerRenderer) Refresh() {
	r.container.table.Refresh()
}

func (r *tableContainerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container.table}
}

func (r *tableContainerRenderer) Destroy() {}
