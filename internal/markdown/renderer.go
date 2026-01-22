package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/yuin/goldmark/ast"
)

// Renderer converts Goldmark AST to Fyne RichText segments
type Renderer struct {
	source      []byte
	segments    []widget.RichTextSegment
	highlighter *SyntaxHighlighter
}

// NewRenderer creates a new AST to Fyne renderer
func NewRenderer(source []byte) *Renderer {
	return &Renderer{
		source:      source,
		segments:    make([]widget.RichTextSegment, 0),
		highlighter: NewSyntaxHighlighter(),
	}
}

// Render converts the AST node to Fyne segments
func (r *Renderer) Render(node ast.Node) []widget.RichTextSegment {
	r.renderNode(node)
	return r.segments
}

// renderNode recursively renders AST nodes
func (r *Renderer) renderNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Document:
		r.renderChildren(n)

	case *ast.Heading:
		r.renderHeading(n)

	case *ast.Paragraph:
		r.renderParagraph(n)

	case *ast.FencedCodeBlock:
		r.renderFencedCodeBlock(n)

	case *ast.CodeBlock:
		r.renderCodeBlock(n)

	case *ast.Blockquote:
		r.renderBlockquote(n)

	case *ast.List:
		r.renderList(n)

	case *ast.ListItem:
		r.renderListItem(n)

	case *ast.ThematicBreak:
		r.addSegment("---", widget.RichTextStyle{})
		r.addSegment("\n\n", widget.RichTextStyle{})

	case *ast.HTMLBlock:
		// Skip HTML blocks for now
		r.renderChildren(n)

	default:
		r.renderChildren(node)
	}
}

// renderHeading renders heading nodes with visual hierarchy
func (r *Renderer) renderHeading(node *ast.Heading) {
	text := r.extractInlineText(node)

	// Add visual markers for different heading levels
	var prefix string
	switch node.Level {
	case 1:
		prefix = "# "
	case 2:
		prefix = "## "
	case 3:
		prefix = "### "
	case 4:
		prefix = "#### "
	case 5:
		prefix = "##### "
	case 6:
		prefix = "###### "
	}

	if prefix != "" {
		r.addSegment(prefix, widget.RichTextStyle{
			TextStyle: fyne.TextStyle{Bold: true},
			ColorName: widget.RichTextStyleInline.ColorName,
		})
	}

	r.addSegment(text, widget.RichTextStyle{
		TextStyle: fyne.TextStyle{Bold: true},
	})
	r.addSegment("\n\n", widget.RichTextStyle{})
}

// renderParagraph renders paragraph nodes
func (r *Renderer) renderParagraph(node *ast.Paragraph) {
	r.renderInline(node)
	r.addSegment("\n\n", widget.RichTextStyle{})
}

// renderFencedCodeBlock renders fenced code blocks with syntax highlighting
func (r *Renderer) renderFencedCodeBlock(node *ast.FencedCodeBlock) {
	var buf bytes.Buffer

	// Extract code content
	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		buf.Write(line.Value(r.source))
	}

	code := buf.String()
	// Remove trailing newline if present
	code = strings.TrimSuffix(code, "\n")

	// Get language from info string
	language := string(node.Language(r.source))
	if language == "" {
		language = "text"
	}

	// Add language label if specified
	if language != "text" && language != "" {
		r.addSegment("```"+language, widget.RichTextStyle{
			TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
			ColorName: widget.RichTextStyleInline.ColorName,
		})
		r.addSegment("\n", widget.RichTextStyle{})
	}

	// Highlight code
	highlightedSegments := r.highlighter.Highlight(code, language)
	r.segments = append(r.segments, highlightedSegments...)

	// Add closing fence marker
	if language != "text" && language != "" {
		r.addSegment("\n```", widget.RichTextStyle{
			TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
			ColorName: widget.RichTextStyleInline.ColorName,
		})
	}
	r.addSegment("\n\n", widget.RichTextStyle{})
}

// renderCodeBlock renders indented code blocks
func (r *Renderer) renderCodeBlock(node *ast.CodeBlock) {
	var buf bytes.Buffer

	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		buf.Write(line.Value(r.source))
	}

	code := buf.String()
	code = strings.TrimSuffix(code, "\n")

	r.addSegment(code, widget.RichTextStyle{
		TextStyle: fyne.TextStyle{Monospace: true},
	})
	r.addSegment("\n\n", widget.RichTextStyle{})
}

// renderBlockquote renders blockquote nodes
func (r *Renderer) renderBlockquote(node *ast.Blockquote) {
	r.addSegment("│ ", widget.RichTextStyle{
		ColorName: widget.RichTextStyleInline.ColorName,
	})

	// Render content with italic style
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if para, ok := child.(*ast.Paragraph); ok {
			text := r.extractInlineText(para)
			r.addSegment(text, widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Italic: true},
			})
		} else {
			r.renderNode(child)
		}
	}
	r.addSegment("\n\n", widget.RichTextStyle{})
}

// renderList renders list nodes
func (r *Renderer) renderList(node *ast.List) {
	r.renderChildren(node)
}

// renderListItem renders list item nodes
func (r *Renderer) renderListItem(node *ast.ListItem) {
	r.addSegment("• ", widget.RichTextStyle{})
	r.renderChildren(node)
}

// renderInline renders inline content (text, emphasis, links, etc.)
func (r *Renderer) renderInline(node ast.Node) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Text:
			r.renderText(n)

		case *ast.String:
			r.addSegment(string(n.Value), widget.RichTextStyle{})

		case *ast.Emphasis:
			r.renderEmphasis(n)

		case *ast.CodeSpan:
			r.renderCodeSpan(n)

		case *ast.Link:
			r.renderLink(n)

		case *ast.Image:
			r.renderImage(n)

		case *ast.AutoLink:
			r.renderAutoLink(n)

		default:
			// Recursively render unknown inline elements
			r.renderInline(child)
		}
	}
}

// renderText renders text nodes
func (r *Renderer) renderText(node *ast.Text) {
	text := string(node.Segment.Value(r.source))
	if node.SoftLineBreak() {
		text += " "
	} else if node.HardLineBreak() {
		text += "\n"
	}
	r.addSegment(text, widget.RichTextStyle{})
}

// renderEmphasis renders emphasis (italic/bold) nodes
func (r *Renderer) renderEmphasis(node *ast.Emphasis) {
	style := widget.RichTextStyle{}

	if node.Level == 1 {
		style.TextStyle = fyne.TextStyle{Italic: true}
	} else if node.Level == 2 {
		style.TextStyle = fyne.TextStyle{Bold: true}
	}

	text := r.extractInlineText(node)
	r.addSegment(text, style)
}

// renderCodeSpan renders inline code
func (r *Renderer) renderCodeSpan(node *ast.CodeSpan) {
	text := string(node.Text(r.source))
	r.addSegment(text, widget.RichTextStyle{
		TextStyle: fyne.TextStyle{Monospace: true},
	})
}

// renderLink renders link nodes
func (r *Renderer) renderLink(node *ast.Link) {
	text := r.extractInlineText(node)
	url := string(node.Destination)

	// For now, just render as text with the URL
	linkText := fmt.Sprintf("%s (%s)", text, url)
	r.addSegment(linkText, widget.RichTextStyle{
		ColorName: widget.RichTextStyleInline.ColorName,
	})
}

// renderImage renders image nodes
func (r *Renderer) renderImage(node *ast.Image) {
	// For now, just render alt text
	altText := r.extractInlineText(node)
	url := string(node.Destination)

	imageText := fmt.Sprintf("[Image: %s] (%s)", altText, url)
	r.addSegment(imageText, widget.RichTextStyle{
		TextStyle: fyne.TextStyle{Italic: true},
	})
	r.addSegment("\n\n", widget.RichTextStyle{})
}

// renderAutoLink renders auto-link nodes
func (r *Renderer) renderAutoLink(node *ast.AutoLink) {
	url := string(node.URL(r.source))
	r.addSegment(url, widget.RichTextStyle{
		ColorName: widget.RichTextStyleInline.ColorName,
	})
}

// renderChildren renders all child nodes
func (r *Renderer) renderChildren(node ast.Node) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		r.renderNode(child)
	}
}

// extractInlineText extracts text from inline elements
func (r *Renderer) extractInlineText(node ast.Node) string {
	var buf bytes.Buffer

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Text:
			buf.Write(n.Segment.Value(r.source))
			if n.SoftLineBreak() {
				buf.WriteString(" ")
			}
		case *ast.String:
			buf.Write(n.Value)
		default:
			// Recursively extract text from nested elements
			buf.WriteString(r.extractInlineText(child))
		}
	}

	return buf.String()
}

// addSegment adds a segment to the segments list
func (r *Renderer) addSegment(text string, style widget.RichTextStyle) {
	if text == "" {
		return
	}

	r.segments = append(r.segments, &widget.TextSegment{
		Text:  text,
		Style: style,
	})
}
