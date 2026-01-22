package markdown

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yuin/goldmark/ast"
)

// Renderer converts Goldmark AST to Fyne widgets
type Renderer struct {
	source      []byte
	segments    []widget.RichTextSegment
	widgets     []fyne.CanvasObject
	highlighter *SyntaxHighlighter
}

// NewRenderer creates a new AST to Fyne renderer
func NewRenderer(source []byte) *Renderer {
	return &Renderer{
		source:      source,
		segments:    make([]widget.RichTextSegment, 0),
		widgets:     make([]fyne.CanvasObject, 0),
		highlighter: NewSyntaxHighlighter(),
	}
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

	default:
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
		text := string(n.Segment.Value(r.source))
		if n.SoftLineBreak() {
			text += " "
		} else if n.HardLineBreak() {
			text += "\n"
		}
		segments = append(segments, &widget.TextSegment{
			Text:  text,
			Style: widget.RichTextStyle{},
		})

	case *ast.String:
		segments = append(segments, &widget.TextSegment{
			Text:  string(n.Value),
			Style: widget.RichTextStyle{},
		})

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
		text := string(n.Text(r.source))
		segments = append(segments, &widget.TextSegment{
			Text: text,
			Style: widget.RichTextStyle{
				TextStyle: fyne.TextStyle{Monospace: true},
			},
		})

	case *ast.Link:
		text := r.extractInlineText(n)
		url := string(n.Destination)
		segments = append(segments, &widget.TextSegment{
			Text: text + " (" + url + ")",
			Style: widget.RichTextStyle{
				ColorName: widget.RichTextStyleInline.ColorName,
			},
		})

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
		url := string(n.URL(r.source))
		segments = append(segments, &widget.TextSegment{
			Text: url,
			Style: widget.RichTextStyle{
				ColorName: widget.RichTextStyleInline.ColorName,
			},
		})

	default:
		// Recursively handle children
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			segments = append(segments, r.renderInlineNode(child)...)
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

	// Decode HTML entities (e.g., &ldquo; → ", &rdquo; → ", &mdash; → —)
	return html.UnescapeString(buf.String())
}

// Widget rendering methods

// renderHeadingAsWidget renders heading with blue color and larger size
func (r *Renderer) renderHeadingAsWidget(node *ast.Heading) {
	text := r.extractInlineText(node)

	// Use hyperlink color which is a nice shade of blue
	rt := widget.NewRichText(&widget.ParagraphSegment{
		Texts: []widget.RichTextSegment{
			&widget.TextSegment{
				Text: text,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
					ColorName: theme.ColorNameHyperlink,
					SizeName:  theme.SizeNameHeadingText,
				},
			},
		},
	})
	rt.Wrapping = fyne.TextWrapWord

	r.widgets = append(r.widgets, rt)
}

// renderParagraphAsWidget renders paragraph as RichText
func (r *Renderer) renderParagraphAsWidget(node *ast.Paragraph) {
	// Collect inline segments for this paragraph
	var texts []widget.RichTextSegment

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		texts = append(texts, r.renderInlineNode(child)...)
	}

	if len(texts) > 0 {
		rt := widget.NewRichText(&widget.ParagraphSegment{
			Texts: texts,
		})
		rt.Wrapping = fyne.TextWrapWord
		r.widgets = append(r.widgets, rt)
	}
}

// renderCodeBlockAsWidget renders code block with syntax highlighting
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

	// Get syntax highlighted segments (no fence markers)
	highlightedSegments := r.highlighter.Highlight(code, language)

	rt := widget.NewRichText(&widget.ParagraphSegment{
		Texts: highlightedSegments,
	})

	// Don't wrap code - let it scroll horizontally
	rt.Wrapping = fyne.TextWrapOff

	r.widgets = append(r.widgets, rt)
}

// renderBlockquoteAsWidget renders blockquote
func (r *Renderer) renderBlockquoteAsWidget(node *ast.Blockquote) {
	var texts []widget.RichTextSegment

	texts = append(texts, &widget.TextSegment{
		Text: "│ ",
		Style: widget.RichTextStyle{
			ColorName: widget.RichTextStyleInline.ColorName,
		},
	})

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

	rt := widget.NewRichText(&widget.ParagraphSegment{
		Texts: texts,
	})
	rt.Wrapping = fyne.TextWrapWord
	r.widgets = append(r.widgets, rt)
}

// renderListAsWidget renders list items as individual paragraphs with bullets
func (r *Renderer) renderListAsWidget(node *ast.List) {
	isOrdered := node.IsOrdered()
	itemNum := 1

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			r.renderListItemAsWidget(listItem, isOrdered, itemNum, 0)
			if isOrdered {
				itemNum++
			}
		}
	}
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

	// Build the complete text for the list item
	var itemText strings.Builder
	itemText.WriteString(indent)
	itemText.WriteString(bullet)

	// Extract content from the list item
	for itemChild := node.FirstChild(); itemChild != nil; itemChild = itemChild.NextSibling() {
		switch child := itemChild.(type) {
		case *ast.Paragraph:
			// Extract text from paragraph
			text := r.extractInlineText(child)
			itemText.WriteString(text)
		case *ast.List:
			// Render current item first
			if itemText.Len() > len(indent+bullet) {
				label := widget.NewLabel(itemText.String())
				label.Wrapping = fyne.TextWrapWord
				r.widgets = append(r.widgets, label)
			}

			// Render nested list
			r.renderNestedList(child, depth+1)
			return
		}
	}

	// Render the list item
	if itemText.Len() > len(indent+bullet) {
		label := widget.NewLabel(itemText.String())
		label.Wrapping = fyne.TextWrapWord
		r.widgets = append(r.widgets, label)
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

