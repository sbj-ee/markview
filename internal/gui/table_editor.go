package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// TableEditor provides a visual interface for creating markdown tables
type TableEditor struct {
	rows      int
	cols      int
	cells     [][]*widget.Entry
	headers   []*widget.Entry
	container *fyne.Container
	onInsert  func(markdown string)
}

// NewTableEditor creates a new table editor
func NewTableEditor(onInsert func(markdown string)) *TableEditor {
	te := &TableEditor{
		rows:     2,
		cols:     3,
		onInsert: onInsert,
	}
	te.buildUI()
	return te
}

// buildUI builds the table editor UI
func (te *TableEditor) buildUI() {
	te.createCells()
	te.updateContainer()
}

// createCells creates the entry widgets for the table
func (te *TableEditor) createCells() {
	// Create header cells
	te.headers = make([]*widget.Entry, te.cols)
	for j := 0; j < te.cols; j++ {
		te.headers[j] = widget.NewEntry()
		te.headers[j].SetPlaceHolder(fmt.Sprintf("Header %d", j+1))
	}

	// Create data cells
	te.cells = make([][]*widget.Entry, te.rows)
	for i := 0; i < te.rows; i++ {
		te.cells[i] = make([]*widget.Entry, te.cols)
		for j := 0; j < te.cols; j++ {
			te.cells[i][j] = widget.NewEntry()
			te.cells[i][j].SetPlaceHolder(fmt.Sprintf("R%dC%d", i+1, j+1))
		}
	}
}

// updateContainer updates the container with current cells
func (te *TableEditor) updateContainer() {
	// Build header row
	headerRow := container.NewGridWithColumns(te.cols)
	for j := 0; j < te.cols; j++ {
		headerRow.Add(te.headers[j])
	}

	// Build data rows
	dataRows := container.NewVBox()
	for i := 0; i < te.rows; i++ {
		row := container.NewGridWithColumns(te.cols)
		for j := 0; j < te.cols; j++ {
			row.Add(te.cells[i][j])
		}
		dataRows.Add(row)
	}

	// Build size controls
	rowsLabel := widget.NewLabel(fmt.Sprintf("Rows: %d", te.rows))
	colsLabel := widget.NewLabel(fmt.Sprintf("Columns: %d", te.cols))

	addRowBtn := widget.NewButton("+Row", func() {
		te.addRow()
	})
	removeRowBtn := widget.NewButton("-Row", func() {
		te.removeRow()
	})
	addColBtn := widget.NewButton("+Col", func() {
		te.addColumn()
	})
	removeColBtn := widget.NewButton("-Col", func() {
		te.removeColumn()
	})

	controls := container.NewHBox(
		rowsLabel,
		addRowBtn,
		removeRowBtn,
		widget.NewSeparator(),
		colsLabel,
		addColBtn,
		removeColBtn,
	)

	// Combine into main container
	te.container = container.NewVBox(
		widget.NewLabel("Headers:"),
		headerRow,
		widget.NewSeparator(),
		widget.NewLabel("Data:"),
		dataRows,
		widget.NewSeparator(),
		controls,
	)
}

// addRow adds a new row to the table
func (te *TableEditor) addRow() {
	if te.rows >= 20 {
		return // Limit to 20 rows
	}
	te.rows++

	// Add new row of cells
	newRow := make([]*widget.Entry, te.cols)
	for j := 0; j < te.cols; j++ {
		newRow[j] = widget.NewEntry()
		newRow[j].SetPlaceHolder(fmt.Sprintf("R%dC%d", te.rows, j+1))
	}
	te.cells = append(te.cells, newRow)

	te.updateContainer()
}

// removeRow removes the last row from the table
func (te *TableEditor) removeRow() {
	if te.rows <= 1 {
		return // Keep at least 1 row
	}
	te.rows--
	te.cells = te.cells[:te.rows]
	te.updateContainer()
}

// addColumn adds a new column to the table
func (te *TableEditor) addColumn() {
	if te.cols >= 10 {
		return // Limit to 10 columns
	}
	te.cols++

	// Add new header cell
	newHeader := widget.NewEntry()
	newHeader.SetPlaceHolder(fmt.Sprintf("Header %d", te.cols))
	te.headers = append(te.headers, newHeader)

	// Add new cell to each row
	for i := 0; i < te.rows; i++ {
		newCell := widget.NewEntry()
		newCell.SetPlaceHolder(fmt.Sprintf("R%dC%d", i+1, te.cols))
		te.cells[i] = append(te.cells[i], newCell)
	}

	te.updateContainer()
}

// removeColumn removes the last column from the table
func (te *TableEditor) removeColumn() {
	if te.cols <= 1 {
		return // Keep at least 1 column
	}
	te.cols--
	te.headers = te.headers[:te.cols]
	for i := 0; i < te.rows; i++ {
		te.cells[i] = te.cells[i][:te.cols]
	}
	te.updateContainer()
}

// GenerateMarkdown generates markdown table from the current data
func (te *TableEditor) GenerateMarkdown() string {
	var sb strings.Builder

	// Header row
	sb.WriteString("|")
	for j := 0; j < te.cols; j++ {
		text := te.headers[j].Text
		if text == "" {
			text = fmt.Sprintf("Column %d", j+1)
		}
		sb.WriteString(" " + text + " |")
	}
	sb.WriteString("\n")

	// Separator row
	sb.WriteString("|")
	for j := 0; j < te.cols; j++ {
		sb.WriteString(" --- |")
	}
	sb.WriteString("\n")

	// Data rows
	for i := 0; i < te.rows; i++ {
		sb.WriteString("|")
		for j := 0; j < te.cols; j++ {
			text := te.cells[i][j].Text
			if text == "" {
				text = " "
			}
			sb.WriteString(" " + text + " |")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// GetContainer returns the table editor container
func (te *TableEditor) GetContainer() *fyne.Container {
	return te.container
}

// ShowTableEditorDialog shows a dialog for creating/editing tables
func ShowTableEditorDialog(parent fyne.Window, onInsert func(markdown string)) {
	te := NewTableEditor(onInsert)

	var d dialog.Dialog

	insertBtn := widget.NewButton("Insert Table", func() {
		markdown := te.GenerateMarkdown()
		if onInsert != nil {
			onInsert(markdown)
		}
		d.Hide()
	})
	insertBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("Cancel", func() {
		d.Hide()
	})

	buttons := container.NewHBox(
		cancelBtn,
		insertBtn,
	)

	content := container.NewVBox(
		te.GetContainer(),
		widget.NewSeparator(),
		buttons,
	)

	d = dialog.NewCustomWithoutButtons("Insert Table", content, parent)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
