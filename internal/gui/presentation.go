package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// PresentationView displays markdown as slides
type PresentationView struct {
	window       fyne.Window
	slides       []string
	currentSlide int
	content      *widget.RichText
	counter      *widget.Label
	container    *fyne.Container
}

// NewPresentationView creates a new presentation view
func NewPresentationView(window fyne.Window, markdown string) *PresentationView {
	p := &PresentationView{
		window:       window,
		currentSlide: 0,
	}

	// Split markdown by horizontal rules (---)
	p.slides = splitIntoSlides(markdown)
	if len(p.slides) == 0 {
		p.slides = []string{"No content"}
	}

	p.content = widget.NewRichTextFromMarkdown("")
	p.content.Wrapping = fyne.TextWrapWord

	p.counter = widget.NewLabel("")
	p.updateSlide()

	return p
}

// splitIntoSlides splits markdown content into slides by --- separators
func splitIntoSlides(markdown string) []string {
	// Split by horizontal rules
	parts := strings.Split(markdown, "\n---\n")

	var slides []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			slides = append(slides, part)
		}
	}

	return slides
}

// updateSlide updates the current slide display
func (p *PresentationView) updateSlide() {
	if p.currentSlide >= 0 && p.currentSlide < len(p.slides) {
		p.content.ParseMarkdown(p.slides[p.currentSlide])
		p.counter.SetText(intToString(p.currentSlide+1) + " / " + intToString(len(p.slides)))
	}
}

// nextSlide moves to the next slide
func (p *PresentationView) nextSlide() {
	if p.currentSlide < len(p.slides)-1 {
		p.currentSlide++
		p.updateSlide()
	}
}

// prevSlide moves to the previous slide
func (p *PresentationView) prevSlide() {
	if p.currentSlide > 0 {
		p.currentSlide--
		p.updateSlide()
	}
}

// Show displays the presentation in a dialog
func (p *PresentationView) Show() {
	scroll := container.NewScroll(p.content)
	scroll.SetMinSize(fyne.NewSize(700, 450))

	prevBtn := widget.NewButton("Previous", func() {
		p.prevSlide()
	})
	nextBtn := widget.NewButton("Next", func() {
		p.nextSlide()
	})

	controls := container.NewHBox(
		prevBtn,
		p.counter,
		nextBtn,
	)

	content := container.NewBorder(
		nil,
		container.NewCenter(controls),
		nil,
		nil,
		scroll,
	)

	d := dialog.NewCustom("Presentation Mode", "Exit", content, p.window)
	d.Resize(fyne.NewSize(800, 600))
	d.Show()
}

// ShowPresentationMode shows the markdown as a presentation
func ShowPresentationMode(window fyne.Window, markdown string) {
	p := NewPresentationView(window, markdown)
	p.Show()
}
