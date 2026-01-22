package toc

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestNewGenerator(t *testing.T) {
	source := []byte("# Test")
	gen := NewGenerator(source)

	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}

	if gen.source == nil {
		t.Error("Generator source is nil")
	}
}

func TestGenerate_SingleHeading(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	source := []byte("# Main Title")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	gen := NewGenerator(source)
	entries := gen.Generate(doc)

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Level != 1 {
		t.Errorf("Expected level 1, got %d", entries[0].Level)
	}

	if entries[0].Text != "Main Title" {
		t.Errorf("Expected text 'Main Title', got '%s'", entries[0].Text)
	}
}

func TestGenerate_MultipleHeadings(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	source := []byte(`# Title 1
## Title 2
### Title 3`)

	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	gen := NewGenerator(source)
	entries := gen.Generate(doc)

	if len(entries) != 1 {
		t.Fatalf("Expected 1 top-level entry, got %d", len(entries))
	}

	// Check hierarchy
	if entries[0].Level != 1 {
		t.Errorf("Expected level 1 for first entry, got %d", entries[0].Level)
	}

	if len(entries[0].Children) != 1 {
		t.Fatalf("Expected 1 child for H1, got %d", len(entries[0].Children))
	}

	if entries[0].Children[0].Level != 2 {
		t.Errorf("Expected level 2 for child, got %d", entries[0].Children[0].Level)
	}

	if len(entries[0].Children[0].Children) != 1 {
		t.Fatalf("Expected 1 child for H2, got %d", len(entries[0].Children[0].Children))
	}

	if entries[0].Children[0].Children[0].Level != 3 {
		t.Errorf("Expected level 3 for grandchild, got %d", entries[0].Children[0].Children[0].Level)
	}
}

func TestGenerate_MultipleSections(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	source := []byte(`# Section 1

Content here.

## Subsection 1.1

More content.

# Section 2

More content.

## Subsection 2.1`)

	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	gen := NewGenerator(source)
	entries := gen.Generate(doc)

	if len(entries) != 2 {
		t.Fatalf("Expected 2 top-level entries, got %d", len(entries))
	}

	if entries[0].Text != "Section 1" {
		t.Errorf("Expected 'Section 1', got '%s'", entries[0].Text)
	}

	if entries[1].Text != "Section 2" {
		t.Errorf("Expected 'Section 2', got '%s'", entries[1].Text)
	}

	if len(entries[0].Children) != 1 {
		t.Errorf("Expected 1 child for Section 1, got %d", len(entries[0].Children))
	}

	if len(entries[1].Children) != 1 {
		t.Errorf("Expected 1 child for Section 2, got %d", len(entries[1].Children))
	}
}

func TestGenerate_SkipNonHeadings(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	source := []byte(`# Title

This is a paragraph.

**Bold text** and *italic text*.

- List item 1
- List item 2

## Subtitle`)

	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	gen := NewGenerator(source)
	entries := gen.Generate(doc)

	if len(entries) != 1 {
		t.Fatalf("Expected 1 top-level entry, got %d", len(entries))
	}

	if entries[0].Text != "Title" {
		t.Errorf("Expected 'Title', got '%s'", entries[0].Text)
	}

	if len(entries[0].Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(entries[0].Children))
	}

	if entries[0].Children[0].Text != "Subtitle" {
		t.Errorf("Expected 'Subtitle', got '%s'", entries[0].Children[0].Text)
	}
}

func TestGenerate_EmptyDocument(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	source := []byte("")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	gen := NewGenerator(source)
	entries := gen.Generate(doc)

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for empty document, got %d", len(entries))
	}
}

func TestGenerate_NoHeadings(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	source := []byte("This is just a paragraph.\n\nAnother paragraph.")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	gen := NewGenerator(source)
	entries := gen.Generate(doc)

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for document with no headings, got %d", len(entries))
	}
}

func TestFlattenEntries(t *testing.T) {
	entries := []*TOCEntry{
		{
			Level: 1,
			Text:  "Section 1",
			Children: []*TOCEntry{
				{
					Level: 2,
					Text:  "Subsection 1.1",
					Children: []*TOCEntry{
						{
							Level:    3,
							Text:     "Subsubsection 1.1.1",
							Children: []*TOCEntry{},
						},
					},
				},
				{
					Level:    2,
					Text:     "Subsection 1.2",
					Children: []*TOCEntry{},
				},
			},
		},
		{
			Level:    1,
			Text:     "Section 2",
			Children: []*TOCEntry{},
		},
	}

	flat := FlattenEntries(entries)

	expected := []string{
		"Section 1",
		"Subsection 1.1",
		"Subsubsection 1.1.1",
		"Subsection 1.2",
		"Section 2",
	}

	if len(flat) != len(expected) {
		t.Fatalf("Expected %d entries, got %d", len(expected), len(flat))
	}

	for i, exp := range expected {
		if flat[i].Text != exp {
			t.Errorf("Entry %d: expected '%s', got '%s'", i, exp, flat[i].Text)
		}
	}
}

func TestFlattenEntries_Empty(t *testing.T) {
	entries := []*TOCEntry{}
	flat := FlattenEntries(entries)

	if len(flat) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(flat))
	}
}
