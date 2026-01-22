package gui

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// ImageInsertDialog provides a dialog for inserting images with preview
type ImageInsertDialog struct {
	window     fyne.Window
	baseDir    string
	onInsert   func(markdown string)
	pathEntry  *widget.Entry
	altEntry   *widget.Entry
	preview    *canvas.Image
	previewBox *fyne.Container
}

// NewImageInsertDialog creates a new image insert dialog
func NewImageInsertDialog(window fyne.Window, baseDir string, onInsert func(markdown string)) *ImageInsertDialog {
	d := &ImageInsertDialog{
		window:   window,
		baseDir:  baseDir,
		onInsert: onInsert,
	}
	return d
}

// Show displays the image insert dialog
func (d *ImageInsertDialog) Show() {
	// Create form fields
	d.pathEntry = widget.NewEntry()
	d.pathEntry.SetPlaceHolder("Image path or URL")
	d.pathEntry.OnChanged = func(path string) {
		d.updatePreview(path)
	}

	d.altEntry = widget.NewEntry()
	d.altEntry.SetPlaceHolder("Alt text (description)")

	// Browse button
	browseBtn := widget.NewButton("Browse...", func() {
		d.showFilePicker()
	})

	// Create preview area
	d.preview = canvas.NewImageFromFile("")
	d.preview.FillMode = canvas.ImageFillContain
	d.preview.SetMinSize(fyne.NewSize(300, 200))

	placeholder := widget.NewLabel("No image selected")
	placeholder.Alignment = fyne.TextAlignCenter
	d.previewBox = container.NewStack(placeholder)

	// Create layout
	pathRow := container.NewBorder(nil, nil, nil, browseBtn, d.pathEntry)

	form := container.NewVBox(
		widget.NewLabel("Image Path or URL:"),
		pathRow,
		widget.NewLabel("Alt Text:"),
		d.altEntry,
		widget.NewSeparator(),
		widget.NewLabel("Preview:"),
		d.previewBox,
	)

	var dlg dialog.Dialog

	insertBtn := widget.NewButton("Insert", func() {
		markdown := d.generateMarkdown()
		if d.onInsert != nil && markdown != "" {
			d.onInsert(markdown)
		}
		dlg.Hide()
	})
	insertBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("Cancel", func() {
		dlg.Hide()
	})

	buttons := container.NewHBox(
		cancelBtn,
		insertBtn,
	)

	content := container.NewVBox(
		form,
		widget.NewSeparator(),
		buttons,
	)

	dlg = dialog.NewCustomWithoutButtons("Insert Image", content, d.window)
	dlg.Resize(fyne.NewSize(500, 450))
	dlg.Show()
}

// showFilePicker shows the file picker dialog
func (d *ImageInsertDialog) showFilePicker() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		path := reader.URI().Path()

		// Convert to relative path if possible
		if d.baseDir != "" {
			relPath, err := filepath.Rel(d.baseDir, path)
			if err == nil && !strings.HasPrefix(relPath, "..") {
				path = relPath
			}
		}

		d.pathEntry.SetText(path)
		d.updatePreview(path)
	}, d.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".bmp"}))
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
}

// updatePreview updates the image preview
func (d *ImageInsertDialog) updatePreview(path string) {
	if path == "" {
		placeholder := widget.NewLabel("No image selected")
		placeholder.Alignment = fyne.TextAlignCenter
		d.previewBox.Objects = []fyne.CanvasObject{placeholder}
		d.previewBox.Refresh()
		return
	}

	// Resolve path
	fullPath := path
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "http") && d.baseDir != "" {
		fullPath = filepath.Join(d.baseDir, path)
	}

	// Check if it's a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// For URLs, show a placeholder
		label := widget.NewLabel("URL preview not available\n" + path)
		label.Alignment = fyne.TextAlignCenter
		d.previewBox.Objects = []fyne.CanvasObject{label}
		d.previewBox.Refresh()
		return
	}

	// Load local image
	img := canvas.NewImageFromFile(fullPath)
	if img.Resource != nil || img.File != "" {
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(300, 200))
		d.previewBox.Objects = []fyne.CanvasObject{img}
	} else {
		label := widget.NewLabel("Unable to load image")
		label.Alignment = fyne.TextAlignCenter
		d.previewBox.Objects = []fyne.CanvasObject{label}
	}
	d.previewBox.Refresh()
}

// generateMarkdown generates the markdown image syntax
func (d *ImageInsertDialog) generateMarkdown() string {
	path := d.pathEntry.Text
	if path == "" {
		return ""
	}

	alt := d.altEntry.Text
	if alt == "" {
		alt = "image"
	}

	return "![" + alt + "](" + path + ")"
}

// ShowImageInsertDialog is a convenience function to show the dialog
func ShowImageInsertDialog(window fyne.Window, baseDir string, onInsert func(markdown string)) {
	d := NewImageInsertDialog(window, baseDir, onInsert)
	d.Show()
}
