package markdown

import (
	"bytes"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// SyntaxHighlighter handles syntax highlighting using Chroma
type SyntaxHighlighter struct {
	style *chroma.Style
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter() *SyntaxHighlighter {
	return &SyntaxHighlighter{
		style: styles.Get("monokai"),
	}
}

// SetStyle sets the color scheme for syntax highlighting
func (h *SyntaxHighlighter) SetStyle(styleName string) {
	if style := styles.Get(styleName); style != nil {
		h.style = style
	}
}

// Highlight highlights code and returns Fyne RichText segments
func (h *SyntaxHighlighter) Highlight(code, language string) []widget.RichTextSegment {
	// Get lexer for the language
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		// On error, return code as monospace text
		return []widget.RichTextSegment{
			&widget.TextSegment{
				Text: code,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Monospace: true},
				},
			},
		}
	}

	// Convert tokens to Fyne segments
	var segments []widget.RichTextSegment
	tokens := iterator.Tokens()

	for _, token := range tokens {
		if token.Value == "" {
			continue
		}

		style := h.getTokenStyle(token.Type)
		segments = append(segments, &widget.TextSegment{
			Text:  token.Value,
			Style: style,
		})
	}

	return segments
}

// HighlightAsPlainText returns highlighted code as plain text (for fallback)
func (h *SyntaxHighlighter) HighlightAsPlainText(code, language string) string {
	// Get lexer for the language
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	// Create a string formatter (simple text output)
	formatter := formatters.Get("terminal")
	if formatter == nil {
		return code
	}

	// Tokenize and format
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, h.style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}

// getTokenStyle maps Chroma token types to Fyne text styles
func (h *SyntaxHighlighter) getTokenStyle(tokenType chroma.TokenType) widget.RichTextStyle {
	style := widget.RichTextStyle{
		TextStyle: fyne.TextStyle{Monospace: true},
	}

	// Get color from Chroma style
	entry := h.style.Get(tokenType)
	if entry.Colour.IsSet() {
		// Fyne doesn't support arbitrary colors in RichText easily
		// Map to semantic theme colors based on token type
		style.ColorName = h.mapTokenTypeToColorName(tokenType)
		style.Inline = true
	}

	if entry.Bold == chroma.Yes {
		style.TextStyle.Bold = true
	}
	if entry.Italic == chroma.Yes {
		style.TextStyle.Italic = true
	}

	return style
}

// mapTokenTypeToColorName maps Chroma token types to Fyne theme color names
func (h *SyntaxHighlighter) mapTokenTypeToColorName(tokenType chroma.TokenType) fyne.ThemeColorName {
	// Map common token types to Fyne theme colors
	typeStr := tokenType.String()

	switch {
	case tokenType == chroma.Keyword || tokenType == chroma.KeywordNamespace ||
		tokenType == chroma.KeywordType || tokenType == chroma.KeywordDeclaration ||
		tokenType == chroma.KeywordConstant || tokenType == chroma.KeywordReserved:
		return theme.ColorNameError // Pink/Red for keywords

	case tokenType == chroma.String || tokenType == chroma.LiteralString ||
		strings.HasPrefix(typeStr, "LiteralString") || tokenType == chroma.LiteralStringDouble ||
		tokenType == chroma.LiteralStringSingle:
		return theme.ColorNameWarning // Yellow/Orange for strings

	case tokenType == chroma.Comment || tokenType == chroma.CommentSingle ||
		tokenType == chroma.CommentMultiline || tokenType == chroma.CommentPreproc:
		return theme.ColorNameDisabled // Gray for comments

	case tokenType == chroma.LiteralNumber || strings.HasPrefix(typeStr, "LiteralNumber") ||
		tokenType == chroma.LiteralNumberInteger || tokenType == chroma.LiteralNumberFloat:
		return theme.ColorNameSuccess // Green for numbers

	case tokenType == chroma.NameFunction || tokenType == chroma.NameBuiltin ||
		tokenType == chroma.NameBuiltinPseudo:
		return theme.ColorNameSuccess // Green for functions

	case tokenType == chroma.NameClass || tokenType == chroma.NameException:
		return theme.ColorNamePrimary // Cyan for types/classes

	case tokenType == chroma.Operator || tokenType == chroma.OperatorWord:
		return theme.ColorNameError // Pink/Red for operators

	case tokenType == chroma.NameVariable || tokenType == chroma.NameAttribute:
		return theme.ColorNameForeground // Default text color

	default:
		return theme.ColorNameForeground
	}
}

// GetAvailableStyles returns a list of available Chroma styles
func GetAvailableStyles() []string {
	return styles.Names()
}

// CreateHighlightedLabel creates a Fyne label with syntax-highlighted code
func CreateHighlightedLabel(code, language string, highlighter *SyntaxHighlighter) fyne.CanvasObject {
	segments := highlighter.Highlight(code, language)

	richText := widget.NewRichText(segments...)
	richText.Wrapping = fyne.TextWrapOff

	// Add background for code block
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))

	return fyne.NewContainerWithLayout(
		&codeBlockLayout{padding: fyne.NewSize(10, 10)},
		bg,
		richText,
	)
}

// codeBlockLayout is a simple layout for code blocks with background
type codeBlockLayout struct {
	padding fyne.Size
}

func (l *codeBlockLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(0, 0)
	}

	contentSize := objects[1].MinSize()
	return fyne.NewSize(
		contentSize.Width+l.padding.Width*2,
		contentSize.Height+l.padding.Height*2,
	)
}

func (l *codeBlockLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}

	// Background fills entire area
	objects[0].Resize(size)
	objects[0].Move(fyne.NewPos(0, 0))

	// Content is padded
	objects[1].Resize(fyne.NewSize(
		size.Width-l.padding.Width*2,
		size.Height-l.padding.Height*2,
	))
	objects[1].Move(fyne.NewPos(l.padding.Width, l.padding.Height))
}
