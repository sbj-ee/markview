package markdown

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewParser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	parser := NewParser(logger)

	if parser == nil {
		t.Fatal("NewParser returned nil")
	}

	if parser.md == nil {
		t.Error("Parser markdown instance is nil")
	}

	if parser.logger == nil {
		t.Error("Parser logger is nil")
	}
}

func TestParse_BasicMarkdown(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	parser := NewParser(logger)

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		minSegs  int // Minimum number of segments expected
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
			minSegs: 0,
		},
		{
			name:    "simple heading",
			input:   "# Hello World",
			wantErr: false,
			minSegs: 1,
		},
		{
			name:    "heading with paragraph",
			input:   "# Title\n\nThis is a paragraph.",
			wantErr: false,
			minSegs: 2,
		},
		{
			name:    "multiple headings",
			input:   "# H1\n## H2\n### H3",
			wantErr: false,
			minSegs: 3,
		},
		{
			name:    "bold text",
			input:   "This is **bold** text.",
			wantErr: false,
			minSegs: 1,
		},
		{
			name:    "italic text",
			input:   "This is *italic* text.",
			wantErr: false,
			minSegs: 1,
		},
		{
			name:    "code block",
			input:   "```go\nfunc main() {}\n```",
			wantErr: false,
			minSegs: 1,
		},
		{
			name:    "unordered list",
			input:   "- Item 1\n- Item 2\n- Item 3",
			wantErr: false,
			minSegs: 3,
		},
		{
			name:    "ordered list",
			input:   "1. First\n2. Second\n3. Third",
			wantErr: false,
			minSegs: 3,
		},
		{
			name:    "blockquote",
			input:   "> This is a quote",
			wantErr: false,
			minSegs: 1,
		},
		{
			name:    "link",
			input:   "[Google](https://google.com)",
			wantErr: false,
			minSegs: 1,
		},
		{
			name:    "inline code",
			input:   "Use `var x = 10` for variables.",
			wantErr: false,
			minSegs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, err := parser.Parse([]byte(tt.input))

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(segments) < tt.minSegs {
				t.Errorf("Parse() returned %d segments, want at least %d", len(segments), tt.minSegs)
			}
		})
	}
}

func TestParse_ComplexMarkdown(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	parser := NewParser(logger)

	input := `# Main Title

This is a paragraph with **bold** and *italic* text.

## Section 1

Here's a code block:

` + "```go\nfunc hello() {\n    fmt.Println(\"Hello\")\n}\n```" + `

### Subsection

- Item 1
- Item 2
- Item 3

> A quote here

[Link](https://example.com)
`

	segments, err := parser.Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(segments) == 0 {
		t.Error("Parse() returned no segments for complex markdown")
	}

	// Should have at least headings, paragraphs, code block, list items, quote
	if len(segments) < 10 {
		t.Errorf("Parse() returned %d segments, expected more for complex markdown", len(segments))
	}
}

func TestGetMarkdown(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	parser := NewParser(logger)

	md := parser.GetMarkdown()
	if md == nil {
		t.Error("GetMarkdown() returned nil")
	}
}
