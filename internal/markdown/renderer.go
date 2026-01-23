package markdown

import (
	"bytes"
	"fmt"
	"html"
	"image/color"
	"net/url"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

var (
	whitespaceRegex = regexp.MustCompile(`[\s\n\r\t]+`)
	// Remove ALL control characters and line breaks
	controlCharsRegex = regexp.MustCompile(`[\x00-\x1F\x7F]+`)
)

// Renderer converts Goldmark AST to Fyne widgets
type Renderer struct {
	source      []byte
	segments    []widget.RichTextSegment
	widgets     []fyne.CanvasObject
	highlighter *SyntaxHighlighter
	imageLoader *ImageLoader
}

// normalizeText aggressively removes all newlines and normalizes whitespace
func normalizeText(text string) string {
	// First remove ALL control characters (includes \n, \r, \t, and other control codes)
	text = controlCharsRegex.ReplaceAllString(text, " ")

	// Decode HTML entities (before final normalization)
	text = html.UnescapeString(text)

	// Use regex to replace ALL remaining whitespace with single space
	text = whitespaceRegex.ReplaceAllString(text, " ")

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

// NewRenderer creates a new AST to Fyne renderer
func NewRenderer(source []byte) *Renderer {
	return &Renderer{
		source:      source,
		segments:    make([]widget.RichTextSegment, 0),
		widgets:     make([]fyne.CanvasObject, 0),
		highlighter: NewSyntaxHighlighter(),
		imageLoader: NewImageLoader(""),
	}
}

// SetBasePath sets the base path for resolving relative image paths
func (r *Renderer) SetBasePath(basePath string) {
	r.imageLoader.SetBasePath(basePath)
}

// Render converts the AST node to a Fyne container with properly sized widgets
func (r *Renderer) Render(node ast.Node) fyne.CanvasObject {
	r.renderNodeAsWidget(node)
	return container.NewVBox(r.widgets...)
}

// RenderSegments converts the AST node to Fyne segments (legacy)
func (r *Renderer) RenderSegments(node ast.Node) []widget.RichTextSegment {
	r.renderNode(node)
	return r.segments
}

// renderNodeAsWidget recursively renders AST nodes as widgets
func (r *Renderer) renderNodeAsWidget(node ast.Node) {
	switch n := node.(type) {
	case *ast.Document:
		r.renderChildrenAsWidgets(n)

	case *ast.Heading:
		r.renderHeadingAsWidget(n)

	case *ast.Paragraph:
		r.renderParagraphAsWidget(n)

	case *ast.FencedCodeBlock:
		r.renderCodeBlockAsWidget(n, true)

	case *ast.CodeBlock:
		r.renderCodeBlockAsWidget(n, false)

	case *ast.Blockquote:
		r.renderBlockquoteAsWidget(n)

	case *ast.List:
		r.renderListAsWidget(n)

	case *ast.ThematicBreak:
		r.widgets = append(r.widgets, widget.NewSeparator())

	case *ast.HTMLBlock:
		// Skip HTML blocks
		return

	case *east.Table:
		r.renderTableAsWidget(n)

	default:
		// Check for math blocks
		if node.Kind() == mathjax.KindMathBlock {
			r.renderMathBlockAsWidget(node)
			return
		}
		r.renderChildrenAsWidgets(node)
	}
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
		sep := &widget.SeparatorSegment{}
		r.segments = append(r.segments, sep)

	case *ast.HTMLBlock:
		// Skip HTML blocks for now
		return

	default:
		r.renderChildren(node)
	}
}

// renderHeading renders heading nodes
func (r *Renderer) renderHeading(node *ast.Heading) {
	text := r.extractInlineText(node)

	// Add extra spacing before headings for visual separation
	r.segments = append(r.segments, &widget.TextSegment{
		Text:  "\n",
		Style: widget.RichTextStyle{},
	})

	// Create paragraph segment with bold heading
	// Use different visual markers for different levels
	var displayText string
	switch node.Level {
	case 1:
		// H1: All caps and extra spacing
		displayText = strings.ToUpper(text)
	case 2:
		// H2: Title case with visual marker
		displayText = text
	default:
		// H3-H6: Regular with marker
		displayText = text
	}

	para := &widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: displayText,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
				},
			},
		},
	}

	r.segments = append(r.segments, para)

	// Add separator line for H1 and H2
	if node.Level <= 2 {
		r.segments = append(r.segments, &widget.SeparatorSegment{})
	}
}

// renderParagraph renders paragraph nodes
func (r *Renderer) renderParagraph(node *ast.Paragraph) {
	// Don't render paragraph wrapper if inside list item
	if _, ok := node.Parent().(*ast.ListItem); ok {
		return
	}

	// Collect inline segments for this paragraph
	var texts []widget.RichTextSegment

	// Render inline content and collect segments
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		texts = append(texts, r.renderInlineNode(child)...)
	}

	if len(texts) > 0 {
		para := &widget.ParagraphSegment{
			Texts: texts,
		}
		r.segments = append(r.segments, para)
	}
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
	code = strings.TrimSuffix(code, "\n")

	// Get language from info string
	language := string(node.Language(r.source))
	if language == "" {
		language = "text"
	}

	// Build code block content
	var texts []widget.RichTextSegment

	// Add language label if specified
	if language != "text" && language != "" {
		texts = append(texts, &widget.TextSegment{
			Text: "```" + language + "\n",
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
				ColorName: widget.RichTextStyleInline.ColorName,
			},
		})
	}

	// Highlight code
	highlightedSegments := r.highlighter.Highlight(code, language)
	texts = append(texts, highlightedSegments...)

	// Add closing fence marker
	if language != "text" && language != "" {
		texts = append(texts, &widget.TextSegment{
			Text: "\n```",
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
				ColorName: widget.RichTextStyleInline.ColorName,
			},
		})
	}

	// Wrap in paragraph
	para := &widget.ParagraphSegment{
		Texts: texts,
	}
	r.segments = append(r.segments, para)
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

	para := &widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: code,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Monospace: true},
				},
			},
		},
	}
	r.segments = append(r.segments, para)
}

// renderBlockquote renders blockquote nodes
func (r *Renderer) renderBlockquote(node *ast.Blockquote) {
	var texts []widget.RichTextSegment

	// Add quote marker
	texts = append(texts, &widget.TextSegment{
		Text: "│ ",
		Style: widget.RichTextStyle{
			ColorName: widget.RichTextStyleInline.ColorName,
		},
	})

	// Render content with italic style
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if para, ok := child.(*ast.Paragraph); ok {
			text := r.extractInlineText(para)
			texts = append(texts, &widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Italic: true},
				},
			})
		}
	}

	para := &widget.ParagraphSegment{
		Texts: texts,
	}
	r.segments = append(r.segments, para)
}

// renderList renders list nodes
func (r *Renderer) renderList(node *ast.List) {
	var items []widget.RichTextSegment

	// Collect all list items
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			// Get text from list item
			var text strings.Builder
			for itemChild := listItem.FirstChild(); itemChild != nil; itemChild = itemChild.NextSibling() {
				if para, ok := itemChild.(*ast.Paragraph); ok {
					text.WriteString(r.extractInlineText(para))
				}
			}

			items = append(items, &widget.TextSegment{
				Text:  text.String(),
				Style: widget.RichTextStyle{},
			})
		}
	}

	if len(items) > 0 {
		listSeg := &widget.ListSegment{
			Items:   items,
			Ordered: node.IsOrdered(),
		}
		r.segments = append(r.segments, listSeg)
	}
}

// renderListItem renders list item nodes (called individually)
func (r *Renderer) renderListItem(node *ast.ListItem) {
	// ListItems are now handled by renderList, so this is a no-op
}

// renderInlineNode renders a single inline node and returns segments
func (r *Renderer) renderInlineNode(node ast.Node) []widget.RichTextSegment {
	var segments []widget.RichTextSegment

	switch n := node.(type) {
	case *ast.Text:
		rawText := string(n.Segment.Value(r.source))

		// Normalize the text (removes all control characters and whitespace)
		text := normalizeText(rawText)

		// For hard line break, use actual newline (rare case - two spaces at end of line)
		if n.HardLineBreak() {
			text = text + "\n"
		} else if n.SoftLineBreak() && text == "" {
			// If soft line break and text is now empty after normalization, it was just a newline
			text = " "
		}

		// Only add non-empty segments
		if text != "" {
			segments = append(segments, &widget.TextSegment{
				Text:  text,
				Style: widget.RichTextStyle{},
			})
		}

	case *ast.String:
		text := normalizeText(string(n.Value))

		if text != "" {
			segments = append(segments, &widget.TextSegment{
				Text:  text,
				Style: widget.RichTextStyle{},
			})
		}

	case *ast.Emphasis:
		style := widget.RichTextStyle{}
		if n.Level == 1 {
			style.TextStyle = fyne.TextStyle{Italic: true}
		} else if n.Level == 2 {
			style.TextStyle = fyne.TextStyle{Bold: true}
		}
		text := r.extractInlineText(n)
		segments = append(segments, &widget.TextSegment{
			Text:  text,
			Style: style,
		})

	case *ast.CodeSpan:
		// For code spans, don't normalize whitespace as much (preserve internal spacing)
		text := string(n.Text(r.source))
		// But do decode HTML entities
		text = html.UnescapeString(text)
		segments = append(segments, &widget.TextSegment{
			Text: text,
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		})

	case *ast.Link:
		text := r.extractInlineText(n)
		urlStr := string(n.Destination)
		parsedURL, err := url.Parse(urlStr)
		if err == nil {
			segments = append(segments, &widget.HyperlinkSegment{
				Text: text,
				URL:  parsedURL,
			})
		} else {
			// Fallback to plain text if URL parsing fails
			segments = append(segments, &widget.TextSegment{
				Text: text + " (" + urlStr + ")",
				Style: widget.RichTextStyle{
					ColorName: widget.RichTextStyleInline.ColorName,
				},
			})
		}

	case *ast.Image:
		altText := r.extractInlineText(n)
		url := string(n.Destination)
		segments = append(segments, &widget.TextSegment{
			Text: fmt.Sprintf("[Image: %s] (%s)", altText, url),
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Italic: true},
			},
		})

	case *ast.AutoLink:
		urlStr := string(n.URL(r.source))
		parsedURL, err := url.Parse(urlStr)
		if err == nil {
			segments = append(segments, &widget.HyperlinkSegment{
				Text: urlStr,
				URL:  parsedURL,
			})
		} else {
			segments = append(segments, &widget.TextSegment{
				Text: urlStr,
				Style: widget.RichTextStyle{
					ColorName: widget.RichTextStyleInline.ColorName,
				},
			})
		}

	default:
		// Check for inline math
		if node.Kind() == mathjax.KindInlineMath {
			mathText := r.extractInlineText(node)
			segments = append(segments, &widget.TextSegment{
				Text: "⟨" + mathText + "⟩",
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Monospace: true, Italic: true},
					ColorName: "math",
				},
			})
		} else {
			// Recursively handle children
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				segments = append(segments, r.renderInlineNode(child)...)
			}
		}
	}

	return segments
}

// renderInline renders inline content (text, emphasis, links, etc.) - deprecated
func (r *Renderer) renderInline(node ast.Node) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		segs := r.renderInlineNode(child)
		for _, seg := range segs {
			r.segments = append(r.segments, seg)
		}
	}
}

// renderChildren renders all child nodes
func (r *Renderer) renderChildren(node ast.Node) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		r.renderNode(child)
	}
}

// extractInlineText extracts text from inline elements and decodes HTML entities
func (r *Renderer) extractInlineText(node ast.Node) string {
	var buf bytes.Buffer

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Text:
			buf.WriteString(string(n.Segment.Value(r.source)))
			if n.SoftLineBreak() {
				buf.WriteString(" ")
			}
		case *ast.String:
			buf.WriteString(string(n.Value))
		default:
			// Recursively extract text from nested elements
			buf.WriteString(r.extractInlineText(child))
		}
	}

	// Normalize all text at the end
	return normalizeText(buf.String())
}

// Widget rendering methods

// getHeadingColorName returns the color name for a heading level
func getHeadingColorName(level int) fyne.ThemeColorName {
	switch level {
	case 1, 2:
		return "heading1" // Cyan/teal
	case 3, 4:
		return "heading3" // Orange/gold
	default:
		return "heading5" // Light gray
	}
}

// getHeadingSizeName returns the size name for a heading level
func getHeadingSizeName(level int) fyne.ThemeSizeName {
	switch level {
	case 1:
		return theme.SizeNameHeadingText
	case 2:
		return theme.SizeNameSubHeadingText
	case 3:
		return "heading3"
	case 4:
		return "heading4"
	default:
		return theme.SizeNameText
	}
}

// renderHeadingAsWidget renders heading with level-appropriate color and size
func (r *Renderer) renderHeadingAsWidget(node *ast.Heading) {
	text := r.extractInlineText(node)

	// Add spacing before heading (not for first element)
	if len(r.widgets) > 0 {
		r.widgets = append(r.widgets, NewSpacer(16))
	}

	// Get color and size based on heading level
	colorName := getHeadingColorName(node.Level)
	sizeName := getHeadingSizeName(node.Level)

	rt := widget.NewRichText(&widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
					ColorName: colorName,
					SizeName:  sizeName,
				},
			},
		},
	})
	rt.Wrapping = fyne.TextWrapWord

	r.widgets = append(r.widgets, rt)

	// Add horizontal rule after H1 headings for visual separation
	if node.Level == 1 {
		r.widgets = append(r.widgets, NewSpacer(4))
		r.widgets = append(r.widgets, widget.NewSeparator())
	}

	// Add spacing after heading
	r.widgets = append(r.widgets, NewSpacer(8))
}

// renderParagraphAsWidget renders paragraph as RichText
func (r *Renderer) renderParagraphAsWidget(node *ast.Paragraph) {
	// Check if paragraph contains only an image (common markdown pattern)
	if r.isImageOnlyParagraph(node) {
		r.renderParagraphImages(node)
		return
	}

	// Extract text as plain text (to avoid line break issues)
	text := r.extractInlineText(node)

	if text != "" {
		label := widget.NewLabel(text)
		label.Wrapping = fyne.TextWrapWord
		r.widgets = append(r.widgets, label)

		// Add spacing after paragraph
		r.widgets = append(r.widgets, NewSpacer(8))
	}
}

// isImageOnlyParagraph checks if a paragraph contains only image(s)
func (r *Renderer) isImageOnlyParagraph(node *ast.Paragraph) bool {
	hasImage := false
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.(type) {
		case *ast.Image:
			hasImage = true
		case *ast.Text:
			// Allow whitespace-only text nodes
			text := string(child.(*ast.Text).Segment.Value(r.source))
			if strings.TrimSpace(text) != "" {
				return false
			}
		default:
			return false
		}
	}
	return hasImage
}

// renderParagraphImages renders all images in a paragraph as block elements
func (r *Renderer) renderParagraphImages(node *ast.Paragraph) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if img, ok := child.(*ast.Image); ok {
			r.renderImageAsWidget(img)
		}
	}
}

// renderImageAsWidget renders an image as a widget
func (r *Renderer) renderImageAsWidget(node *ast.Image) {
	src := string(node.Destination)
	altText := r.extractInlineText(node)

	// Add spacing before image
	r.widgets = append(r.widgets, NewSpacer(8))

	// Load and add the image
	imgWidget := r.imageLoader.LoadImage(src, altText)
	r.widgets = append(r.widgets, imgWidget)

	// Add spacing after image
	r.widgets = append(r.widgets, NewSpacer(8))
}

// renderCodeBlockAsWidget renders code block with syntax highlighting and background
func (r *Renderer) renderCodeBlockAsWidget(node ast.Node, fenced bool) {
	var buf bytes.Buffer
	var language string

	if fenced {
		fcb := node.(*ast.FencedCodeBlock)
		for i := 0; i < fcb.Lines().Len(); i++ {
			line := fcb.Lines().At(i)
			buf.Write(line.Value(r.source))
		}
		language = string(fcb.Language(r.source))
	} else {
		cb := node.(*ast.CodeBlock)
		for i := 0; i < cb.Lines().Len(); i++ {
			line := cb.Lines().At(i)
			buf.Write(line.Value(r.source))
		}
	}

	code := strings.TrimSuffix(buf.String(), "\n")
	if language == "" {
		language = "text"
	}

	// Check for special code block types
	if IsMermaidCodeBlock(language) {
		// Render as mermaid diagram block
		mermaidBlock := NewMermaidBlock(code)
		r.widgets = append(r.widgets, NewSpacer(8))
		r.widgets = append(r.widgets, mermaidBlock)
		r.widgets = append(r.widgets, NewSpacer(8))
		return
	}

	// Get syntax highlighted segments (no fence markers)
	highlightedSegments := r.highlighter.Highlight(code, language)

	// Create code block widget with background styling
	codeBlock := NewCodeBlock(highlightedSegments)

	// Add spacing before code block
	r.widgets = append(r.widgets, NewSpacer(8))
	r.widgets = append(r.widgets, codeBlock)
	// Add spacing after code block
	r.widgets = append(r.widgets, NewSpacer(8))
}

// renderBlockquoteAsWidget renders blockquote with styled left border
func (r *Renderer) renderBlockquoteAsWidget(node *ast.Blockquote) {
	// Extract all text from blockquote
	var quoteText string
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if para, ok := child.(*ast.Paragraph); ok {
			quoteText += r.extractInlineText(para)
		}
	}

	// Create blockquote widget with styled border and background
	blockquote := NewBlockquote(quoteText)

	// Add spacing before blockquote
	r.widgets = append(r.widgets, NewSpacer(8))
	r.widgets = append(r.widgets, blockquote)
	// Add spacing after blockquote
	r.widgets = append(r.widgets, NewSpacer(8))
}

// renderListAsWidget renders list items as individual paragraphs with bullets
func (r *Renderer) renderListAsWidget(node *ast.List) {
	isOrdered := node.IsOrdered()
	itemNum := 1

	// Add spacing before list
	r.widgets = append(r.widgets, NewSpacer(4))

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			r.renderListItemAsWidget(listItem, isOrdered, itemNum, 0)
			if isOrdered {
				itemNum++
			}
		}
	}

	// Add spacing after list
	r.widgets = append(r.widgets, NewSpacer(8))
}

// renderListItemAsWidget renders a single list item with proper indentation
func (r *Renderer) renderListItemAsWidget(node *ast.ListItem, isOrdered bool, number int, depth int) {
	indent := strings.Repeat("    ", depth)

	var bullet string
	if isOrdered {
		bullet = fmt.Sprintf("%d. ", number)
	} else {
		bullet = "• "
	}

	// Build text for the list item - try direct text extraction first
	itemText := r.extractInlineText(node)

	// Render the list item
	label := widget.NewLabel(indent + bullet + itemText)
	label.Wrapping = fyne.TextWrapWord
	r.widgets = append(r.widgets, label)

	// Handle nested lists separately
	for itemChild := node.FirstChild(); itemChild != nil; itemChild = itemChild.NextSibling() {
		if nestedList, ok := itemChild.(*ast.List); ok {
			r.renderNestedList(nestedList, depth+1)
		}
	}
}

// renderNestedList renders nested lists with proper indentation
func (r *Renderer) renderNestedList(node *ast.List, depth int) {
	isOrdered := node.IsOrdered()
	itemNum := 1

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			r.renderListItemAsWidget(listItem, isOrdered, itemNum, depth)
			if isOrdered {
				itemNum++
			}
		}
	}
}

// renderChildrenAsWidgets renders all child nodes as widgets
func (r *Renderer) renderChildrenAsWidgets(node ast.Node) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		r.renderNodeAsWidget(child)
	}
}

// renderMathBlockAsWidget renders a math block with special styling
func (r *Renderer) renderMathBlockAsWidget(node ast.Node) {
	// Extract math content from the node
	var buf bytes.Buffer
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		if text, ok := c.(*ast.Text); ok {
			buf.Write(text.Segment.Value(r.source))
		}
	}

	// If no children, try to get text directly from Lines
	if buf.Len() == 0 {
		lines := node.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			buf.Write(line.Value(r.source))
		}
	}

	mathContent := strings.TrimSpace(buf.String())

	// Create math display widget
	mathWidget := NewMathBlock(mathContent)

	// Add spacing before math block
	r.widgets = append(r.widgets, NewSpacer(8))
	r.widgets = append(r.widgets, mathWidget)
	// Add spacing after math block
	r.widgets = append(r.widgets, NewSpacer(8))
}

// renderTableAsWidget renders a table from GFM extension using a grid layout
func (r *Renderer) renderTableAsWidget(node *east.Table) {
	// Collect table data
	var headers []string
	var rows [][]string

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch row := child.(type) {
		case *east.TableHeader:
			// Extract header cells
			for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
				if tableCell, ok := cell.(*east.TableCell); ok {
					headers = append(headers, r.extractInlineText(tableCell))
				}
			}
		case *east.TableRow:
			// Extract row cells
			var rowData []string
			for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
				if tableCell, ok := cell.(*east.TableCell); ok {
					rowData = append(rowData, r.extractInlineText(tableCell))
				}
			}
			rows = append(rows, rowData)
		}
	}

	numCols := len(headers)
	if numCols == 0 && len(rows) > 0 {
		numCols = len(rows[0])
	}

	if numCols == 0 {
		return // Empty table
	}

	// Add spacing before table
	r.widgets = append(r.widgets, NewSpacer(8))

	// Calculate column widths based on content
	colWidths := make([]float32, numCols)
	for col := 0; col < numCols; col++ {
		maxWidth := float32(80) // Minimum width

		// Check header width
		if col < len(headers) {
			label := widget.NewLabel(headers[col])
			label.TextStyle = fyne.TextStyle{Bold: true}
			size := label.MinSize()
			width := size.Width + 32 // Add padding
			if width > maxWidth {
				maxWidth = width
			}
		}

		// Check data widths
		for _, row := range rows {
			if col < len(row) {
				label := widget.NewLabel(row[col])
				size := label.MinSize()
				width := size.Width + 32 // Add padding
				if width > maxWidth {
					maxWidth = width
				}
			}
		}

		// Cap maximum width to allow wrapping for very long content
		if maxWidth > 600 {
			maxWidth = 600
		}

		colWidths[col] = maxWidth
	}

	// Define colors - use primary color for header to make it distinct
	headerBgColor := theme.Color(theme.ColorNamePrimary)
	borderColor := theme.Color(theme.ColorNameSeparator)

	// Build table as a VBox of rows
	tableRows := []fyne.CanvasObject{}

	// Create header row
	headerCells := []fyne.CanvasObject{}
	for col := 0; col < numCols; col++ {
		var text string
		if col < len(headers) {
			text = headers[col]
		}
		cell := r.createTableCell(text, true, headerBgColor, borderColor, colWidths[col])
		headerCells = append(headerCells, cell)
	}
	headerRow := container.NewHBox(headerCells...)
	tableRows = append(tableRows, headerRow)

	// Create data rows
	for rowIdx, row := range rows {
		rowCells := []fyne.CanvasObject{}
		var bgColor fyne.ThemeColorName
		if rowIdx%2 == 0 {
			bgColor = theme.ColorNameInputBackground
		} else {
			bgColor = theme.ColorNameBackground
		}
		rowBgColor := theme.Color(bgColor)

		for col := 0; col < numCols; col++ {
			var text string
			if col < len(row) {
				text = row[col]
			}
			cell := r.createTableCell(text, false, rowBgColor, borderColor, colWidths[col])
			rowCells = append(rowCells, cell)
		}
		dataRow := container.NewHBox(rowCells...)
		tableRows = append(tableRows, dataRow)
	}

	// Create the table container
	tableContainer := container.NewVBox(tableRows...)

	r.widgets = append(r.widgets, tableContainer)

	// Add spacing after table
	r.widgets = append(r.widgets, NewSpacer(8))
}

// createTableCell creates a single table cell with background, border, and fixed width
func (r *Renderer) createTableCell(text string, isHeader bool, bgColor, borderColor color.Color, width float32) fyne.CanvasObject {
	// Create background
	bg := canvas.NewRectangle(bgColor)

	var cellContent fyne.CanvasObject
	if isHeader {
		// Use RichText for header to allow custom text color
		rt := widget.NewRichText(&widget.TextSegment{
			Text: text,
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Bold: true},
				ColorName: "tableHeader",
			},
		})
		rt.Wrapping = fyne.TextWrapWord
		cellContent = container.NewPadded(rt)
	} else {
		label := widget.NewLabel(text)
		label.Wrapping = fyne.TextWrapWord
		cellContent = container.NewPadded(label)
	}

	// Stack background and content
	cell := container.NewStack(bg, cellContent)

	// Create a fixed-width cell using a custom container
	fixedWidthCell := NewFixedWidthContainer(cell, width)

	return fixedWidthCell
}

// FixedWidthContainer is a container that has a fixed width but flexible height
type FixedWidthContainer struct {
	widget.BaseWidget
	content fyne.CanvasObject
	width   float32
}

// NewFixedWidthContainer creates a container with fixed width
func NewFixedWidthContainer(content fyne.CanvasObject, width float32) *FixedWidthContainer {
	c := &FixedWidthContainer{
		content: content,
		width:   width,
	}
	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer implements fyne.Widget
func (c *FixedWidthContainer) CreateRenderer() fyne.WidgetRenderer {
	return &fixedWidthRenderer{container: c}
}

// MinSize returns the minimum size with fixed width
func (c *FixedWidthContainer) MinSize() fyne.Size {
	contentMin := c.content.MinSize()
	// Calculate height based on content that would wrap at our fixed width
	return fyne.NewSize(c.width, contentMin.Height)
}

type fixedWidthRenderer struct {
	container *FixedWidthContainer
}

func (r *fixedWidthRenderer) Layout(size fyne.Size) {
	r.container.content.Resize(fyne.NewSize(r.container.width, size.Height))
	r.container.content.Move(fyne.NewPos(0, 0))
}

func (r *fixedWidthRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *fixedWidthRenderer) Refresh() {
	r.container.content.Refresh()
}

func (r *fixedWidthRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container.content}
}

func (r *fixedWidthRenderer) Destroy() {}
