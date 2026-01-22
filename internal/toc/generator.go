package toc

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

// TOCEntry represents a single entry in the table of contents
type TOCEntry struct {
	Level    int
	Text     string
	ID       string
	Children []*TOCEntry
}

// Generator extracts table of contents from markdown AST
type Generator struct {
	source []byte
}

// NewGenerator creates a new TOC generator
func NewGenerator(source []byte) *Generator {
	return &Generator{
		source: source,
	}
}

// Generate extracts TOC entries from the AST
func (g *Generator) Generate(root ast.Node) []*TOCEntry {
	var entries []*TOCEntry
	var stack []*TOCEntry

	// Walk the AST to find all headings
	ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := node.(*ast.Heading); ok {
			entry := &TOCEntry{
				Level:    heading.Level,
				Text:     g.extractText(heading),
				ID:       g.generateID(heading),
				Children: make([]*TOCEntry, 0),
			}

			// Build hierarchy based on heading levels
			for len(stack) > 0 && stack[len(stack)-1].Level >= entry.Level {
				stack = stack[:len(stack)-1]
			}

			if len(stack) == 0 {
				// Top-level entry
				entries = append(entries, entry)
			} else {
				// Child entry
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, entry)
			}

			stack = append(stack, entry)
		}

		return ast.WalkContinue, nil
	})

	return entries
}

// extractText extracts text content from a heading node
func (g *Generator) extractText(node ast.Node) string {
	var text string

	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			text += string(textNode.Segment.Value(g.source))
		} else if str, ok := child.(*ast.String); ok {
			text += string(str.Value)
		} else {
			// Recursively extract text from nested elements
			text += g.extractText(child)
		}
	}

	return text
}

// generateID generates an ID for a heading node
func (g *Generator) generateID(heading *ast.Heading) string {
	// Check if heading has an explicit ID attribute
	if id, ok := heading.AttributeString("id"); ok {
		if idStr, ok := id.([]byte); ok {
			return string(idStr)
		}
	}

	// Generate ID from text (simple version - just use level)
	// In a more complete implementation, this would slugify the heading text
	return fmt.Sprintf("heading-%d", heading.Level)
}

// FlattenEntries returns a flat list of all TOC entries
func FlattenEntries(entries []*TOCEntry) []*TOCEntry {
	var flat []*TOCEntry

	var flatten func([]*TOCEntry)
	flatten = func(list []*TOCEntry) {
		for _, entry := range list {
			flat = append(flat, entry)
			if len(entry.Children) > 0 {
				flatten(entry.Children)
			}
		}
	}

	flatten(entries)
	return flat
}
