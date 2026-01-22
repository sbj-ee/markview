package gui

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ImagePaster handles pasting images from clipboard
type ImagePaster struct {
	basePath string
	assetDir string
}

// NewImagePaster creates a new image paster
func NewImagePaster(basePath string) *ImagePaster {
	return &ImagePaster{
		basePath: basePath,
		assetDir: "assets",
	}
}

// SetAssetDir sets the directory name for saving images
func (ip *ImagePaster) SetAssetDir(dir string) {
	ip.assetDir = dir
}

// SaveImage saves an image and returns the markdown reference
func (ip *ImagePaster) SaveImage(img image.Image, name string) (string, error) {
	if ip.basePath == "" {
		return "", fmt.Errorf("no base path set - save the document first")
	}

	// Create assets directory if it doesn't exist
	assetsPath := filepath.Join(ip.basePath, ip.assetDir)
	if err := os.MkdirAll(assetsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create assets directory: %w", err)
	}

	// Generate filename if not provided
	if name == "" {
		name = fmt.Sprintf("image-%s.png", time.Now().Format("20060102-150405"))
	}

	// Ensure .png extension
	if filepath.Ext(name) == "" {
		name += ".png"
	}

	// Save the image
	imagePath := filepath.Join(assetsPath, name)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Return relative markdown reference
	relativePath := filepath.Join(ip.assetDir, name)
	return fmt.Sprintf("![%s](%s)", name, relativePath), nil
}

// GetClipboardImage attempts to get an image from the clipboard
// Note: Fyne's clipboard currently only supports text, so this is a placeholder
// for future enhancement when Fyne adds image clipboard support
func (ip *ImagePaster) GetClipboardImage(clipboard fyne.Clipboard) (image.Image, error) {
	// Fyne currently doesn't support image clipboard
	// This is a placeholder for future implementation
	return nil, fmt.Errorf("clipboard image paste not yet supported by Fyne")
}

// ShowImagePasteDialog shows a dialog for pasting/uploading images
func ShowImagePasteDialog(window fyne.Window, basePath string, onInsert func(markdown string)) {
	if basePath == "" {
		dialog.ShowInformation("Paste Image", "Please save the document first before adding images.", window)
		return
	}

	paster := NewImagePaster(basePath)

	// Create name entry
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("image-name (optional)")

	// File picker for manual image selection
	selectBtn := widget.NewButton("Select Image File...", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			// Decode image
			img, _, err := image.Decode(reader)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to decode image: %w", err), window)
				return
			}

			// Get name from entry or use original filename
			name := nameEntry.Text
			if name == "" {
				name = filepath.Base(reader.URI().Path())
			}

			// Save and insert
			markdown, err := paster.SaveImage(img, name)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if onInsert != nil {
				onInsert(markdown)
			}
		}, window)
		fd.SetFilter(&imageFilter{})
		fd.Resize(fyne.NewSize(800, 600))
		fd.Show()
	})

	infoLabel := widget.NewLabel("Images will be saved to the 'assets' folder\nrelative to your document.")

	content := container.NewVBox(
		widget.NewLabel("Image Name:"),
		nameEntry,
		widget.NewSeparator(),
		selectBtn,
		widget.NewSeparator(),
		infoLabel,
	)

	d := dialog.NewCustom("Add Image", "Cancel", content, window)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

// imageFilter is a file filter for image files
type imageFilter struct{}

func (f *imageFilter) Matches(uri fyne.URI) bool {
	ext := filepath.Ext(uri.Path())
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return true
	}
	return false
}
