package markdown

import (
	"testing"

	"fyne.io/fyne/v2"
	"go.uber.org/zap"
)

func TestNewParser(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	parser := NewParser(logger)

	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}

	if parser.md == nil {
		t.Error("Parser.md is nil")
	}

	if parser.logger == nil {
		t.Error("Parser.logger is nil")
	}
}

func TestParse_SimpleMarkdown(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	parser := NewParser(logger)

	tests := []struct {
		name       string
		input      []byte
		wantErr    bool
		minWidgets int // minimum expected widgets
	}{
		{
			name:       "empty input",
			input:      []byte(""),
			wantErr:    false,
			minWidgets: 0,
		},
		{
			name:       "simple heading",
			input:      []byte("# Hello World"),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "paragraph",
			input:      []byte("This is a paragraph."),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "heading and paragraph",
			input:      []byte("# Title\n\nParagraph text."),
			wantErr:    false,
			minWidgets: 2,
		},
		{
			name: "code block",
			input: []byte("```go\n" +
				"func main() {\n" +
				"    fmt.Println(\"Hello\")\n" +
				"}\n" +
				"```"),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "list",
			input:      []byte("- Item 1\n- Item 2\n- Item 3"),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "bold text",
			input:      []byte("This is **bold** text."),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "italic text",
			input:      []byte("This is *italic* text."),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "link",
			input:      []byte("[Fyne](https://fyne.io)"),
			wantErr:    false,
			minWidgets: 1,
		},
		{
			name:       "blockquote",
			input:      []byte("> This is a quote."),
			wantErr:    false,
			minWidgets: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := parser.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if content is a Container
			if cont, ok := content.(*fyne.Container); ok {
				if len(cont.Objects) < tt.minWidgets {
					t.Errorf("Parse() returned %d widgets, want at least %d", len(cont.Objects), tt.minWidgets)
				}
			} else if tt.minWidgets > 0 {
				t.Error("Parse() did not return a Container")
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

> This is a quote.

[Link to Fyne](https://fyne.io)
`

	content, err := parser.Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check if content is a Container with widgets
	if cont, ok := content.(*fyne.Container); ok {
		if len(cont.Objects) < 5 {
			t.Errorf("Parse() returned %d widgets, expected at least 5 for complex document", len(cont.Objects))
		}
	} else {
		t.Error("Parse() did not return a Container")
	}
}
