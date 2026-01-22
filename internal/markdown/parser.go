package markdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.uber.org/zap"
)

// Parser wraps Goldmark markdown parser and converts to Fyne widgets
type Parser struct {
	md     goldmark.Markdown
	logger *zap.Logger
}

// NewParser creates a new markdown parser with GFM extensions
func NewParser(logger *zap.Logger) *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,          // GitHub Flavored Markdown
			extension.Typographer,  // Smart quotes, dashes
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Auto-generate heading IDs
		),
	)

	return &Parser{
		md:     md,
		logger: logger,
	}
}

// Parse parses markdown content and returns a Fyne CanvasObject
func (p *Parser) Parse(content []byte) (fyne.CanvasObject, error) {
	return p.ParseWithBasePath(content, "")
}

// ParseWithBasePath parses markdown content with a base path for resolving relative images
func (p *Parser) ParseWithBasePath(content []byte, basePath string) (fyne.CanvasObject, error) {
	// Parse markdown to AST
	reader := text.NewReader(content)
	doc := p.md.Parser().Parse(reader)

	// Convert AST to Fyne widgets using renderer
	renderer := NewRenderer(content)
	if basePath != "" {
		renderer.SetBasePath(basePath)
	}
	container := renderer.Render(doc)

	return container, nil
}

// ParseLegacy parses markdown content and returns Fyne RichText segments (for compatibility)
func (p *Parser) ParseLegacy(content []byte) ([]widget.RichTextSegment, error) {
	// Parse markdown to AST
	reader := text.NewReader(content)
	doc := p.md.Parser().Parse(reader)

	// Convert AST to Fyne segments using renderer
	renderer := NewRenderer(content)
	segments := renderer.RenderSegments(doc)

	return segments, nil
}

// GetMarkdown returns the underlying Goldmark instance
func (p *Parser) GetMarkdown() goldmark.Markdown {
	return p.md
}
