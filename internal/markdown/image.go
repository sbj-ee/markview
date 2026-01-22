package markdown

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ImageLoader handles loading images from local paths and URLs
type ImageLoader struct {
	basePath string
	cache    map[string]*canvas.Image
	mu       sync.RWMutex
}

// NewImageLoader creates a new image loader with the given base path
func NewImageLoader(basePath string) *ImageLoader {
	return &ImageLoader{
		basePath: basePath,
		cache:    make(map[string]*canvas.Image),
	}
}

// SetBasePath updates the base path for resolving relative image paths
func (l *ImageLoader) SetBasePath(basePath string) {
	l.basePath = basePath
}

// LoadImage loads an image from a path or URL and returns a Fyne canvas object
func (l *ImageLoader) LoadImage(src string, altText string) fyne.CanvasObject {
	// Check cache first
	l.mu.RLock()
	if cached, ok := l.cache[src]; ok {
		l.mu.RUnlock()
		return l.wrapImage(cached, altText)
	}
	l.mu.RUnlock()

	var img *canvas.Image

	if isURL(src) {
		img = l.loadFromURL(src)
	} else {
		img = l.loadFromFile(src)
	}

	if img == nil {
		// Return placeholder for failed loads
		return l.createPlaceholder(src, altText)
	}

	// Cache the loaded image
	l.mu.Lock()
	l.cache[src] = img
	l.mu.Unlock()

	return l.wrapImage(img, altText)
}

// loadFromFile loads an image from a local file path
func (l *ImageLoader) loadFromFile(src string) *canvas.Image {
	// Resolve relative paths
	path := src
	if !filepath.IsAbs(src) && l.basePath != "" {
		path = filepath.Join(l.basePath, src)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// Open and decode the image to get dimensions
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil
	}

	// Create canvas image from file
	img := canvas.NewImageFromFile(path)
	img.FillMode = canvas.ImageFillContain

	// Set reasonable size based on image dimensions
	width, height := calculateDisplaySize(imgConfig.Width, imgConfig.Height)
	img.SetMinSize(fyne.NewSize(width, height))

	return img
}

// loadFromURL loads an image from a URL
func (l *ImageLoader) loadFromURL(src string) *canvas.Image {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(src)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// Read the image data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// Create a resource from the data
	resource := fyne.NewStaticResource(filepath.Base(src), data)

	// Decode to get dimensions
	imgConfig, _, err := image.DecodeConfig(strings.NewReader(string(data)))
	if err != nil {
		// Still try to create the image even if we can't get dimensions
		img := canvas.NewImageFromResource(resource)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(400, 300))
		return img
	}

	img := canvas.NewImageFromResource(resource)
	img.FillMode = canvas.ImageFillContain

	// Set reasonable size based on image dimensions
	width, height := calculateDisplaySize(imgConfig.Width, imgConfig.Height)
	img.SetMinSize(fyne.NewSize(width, height))

	return img
}

// wrapImage wraps an image with optional alt text caption
func (l *ImageLoader) wrapImage(img *canvas.Image, altText string) fyne.CanvasObject {
	// Clone the image to avoid sharing state
	newImg := canvas.NewImageFromResource(img.Resource)
	if img.File != "" {
		newImg = canvas.NewImageFromFile(img.File)
	}
	newImg.FillMode = canvas.ImageFillContain
	newImg.SetMinSize(img.MinSize())

	if altText == "" {
		return newImg
	}

	// Add caption below image
	caption := widget.NewLabel(altText)
	caption.Alignment = fyne.TextAlignCenter
	caption.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(newImg, caption)
}

// createPlaceholder creates a placeholder widget for failed image loads
func (l *ImageLoader) createPlaceholder(src string, altText string) fyne.CanvasObject {
	text := altText
	if text == "" {
		text = "Image: " + filepath.Base(src)
	}

	label := widget.NewLabel("[" + text + "]")
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Italic: true}

	return label
}

// isURL checks if a string is a URL
func isURL(src string) bool {
	return strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")
}

// calculateDisplaySize calculates appropriate display dimensions for an image
func calculateDisplaySize(origWidth, origHeight int) (float32, float32) {
	const maxWidth float32 = 800
	const maxHeight float32 = 600

	width := float32(origWidth)
	height := float32(origHeight)

	// Scale down if too large
	if width > maxWidth {
		ratio := maxWidth / width
		width = maxWidth
		height = height * ratio
	}

	if height > maxHeight {
		ratio := maxHeight / height
		height = maxHeight
		width = width * ratio
	}

	// Ensure minimum size
	if width < 100 {
		width = 100
	}
	if height < 75 {
		height = 75
	}

	return width, height
}
