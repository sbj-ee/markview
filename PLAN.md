# Visual Rendering Improvement Plan

## Current State

The markdown renderer produces functional output but lacks visual polish. Attempts to add visual improvements (spacing, backgrounds, borders) have broken the rendering.

### Current Architecture

```
Markdown Source
    ↓
Goldmark Parser → AST
    ↓
Renderer.Render() → []fyne.CanvasObject (flat slice of widgets)
    ↓
container.NewVBox(widgets...) → displayed in ScrollContainer
```

### Why Visual Improvements Break Rendering

1. **Flat Widget Slice**: The renderer builds a flat slice of widgets that goes into a VBox
2. **Container Nesting**: When we wrap widgets in containers (NewStack, NewHBox, NewPadded), the VBox layout calculates sizes differently
3. **Empty Labels for Spacing**: Using `widget.NewLabel("")` for spacing adds zero-height elements that don't create visual space
4. **Mixed Widget Types**: Mixing simple widgets with container widgets causes layout inconsistencies

## Proposed Solutions

### Option 1: Custom Spacer Widget

Create a custom widget that has a fixed height for spacing:

```go
type Spacer struct {
    widget.BaseWidget
    Height float32
}

func NewSpacer(height float32) *Spacer {
    s := &Spacer{Height: height}
    s.ExtendBaseWidget(s)
    return s
}

func (s *Spacer) MinSize() fyne.Size {
    return fyne.NewSize(0, s.Height)
}

func (s *Spacer) CreateRenderer() fyne.WidgetRenderer {
    return &spacerRenderer{spacer: s}
}
```

**Pros**: Clean, predictable spacing
**Cons**: Requires custom widget implementation

### Option 2: Padding via Theme

Modify the theme to add padding around text elements:

```go
func (m *MarkViewTheme) Size(name fyne.ThemeSizeName) float32 {
    switch name {
    case theme.SizeNamePadding:
        return 8 // Increase default padding
    case theme.SizeNameInnerPadding:
        return 6
    // ...
    }
}
```

**Pros**: Simple, uses existing Fyne mechanisms
**Cons**: Affects all widgets globally, less control

### Option 3: Use RichText with Newlines for Spacing

Add newline characters within RichText segments for spacing:

```go
// Add spacing by including newlines in text
rt := widget.NewRichText(&widget.ParagraphSegment{
    Texts: []widget.RichTextSegment{
        &widget.TextSegment{Text: "\n"}, // Spacing before
        &widget.TextSegment{Text: heading},
        &widget.TextSegment{Text: "\n"}, // Spacing after
    },
})
```

**Pros**: Works within existing widget types
**Cons**: May affect text selection, copy/paste

### Option 4: Complete Architecture Refactor

Restructure to use a single RichText widget for the entire document:

```go
func (r *Renderer) Render(node ast.Node) fyne.CanvasObject {
    // Build one large RichText with all segments
    var allSegments []widget.RichTextSegment

    r.walkAST(node, func(n ast.Node) {
        segments := r.nodeToSegments(n)
        allSegments = append(allSegments, segments...)
    })

    rt := widget.NewRichText(allSegments...)
    rt.Wrapping = fyne.TextWrapWord
    return rt
}
```

**Pros**:
- Single widget, consistent rendering
- RichText handles internal layout
- Supports inline formatting (bold, italic, links)

**Cons**:
- Significant refactor
- May have performance issues with large documents
- Code blocks need special handling (no wrapping)

### Option 5: Custom Document Widget

Create a custom widget specifically for markdown documents:

```go
type MarkdownDocument struct {
    widget.BaseWidget
    blocks []DocumentBlock
}

type DocumentBlock interface {
    Render() fyne.CanvasObject
    MinSize() fyne.Size
    Type() BlockType
}

type HeadingBlock struct {
    Level int
    Text  string
}

type ParagraphBlock struct {
    Text string
}

type CodeBlock struct {
    Language string
    Code     string
}
```

**Pros**:
- Full control over layout
- Can implement proper spacing logic
- Clean separation of concerns

**Cons**:
- Most complex solution
- Requires implementing custom renderer

## Recommended Approach

### Phase 1: Custom Spacer Widget (Low Risk)

1. Create a simple `Spacer` widget with configurable height
2. Use it between major elements (headings, paragraphs, code blocks)
3. Test thoroughly before proceeding

### Phase 2: Code Block Styling (Medium Risk)

1. Create a `CodeBlockWidget` that wraps RichText with background
2. Handle sizing explicitly to avoid layout issues
3. Test in isolation before integrating

### Phase 3: Single RichText Approach (Higher Risk)

1. Experiment with rendering entire document as one RichText
2. Use different segment types for formatting
3. Handle code blocks as separate widgets (they need different wrapping)

## Visual Improvements Priority

1. **Spacing between elements** - Most impactful for readability
2. **Code block backgrounds** - Visual distinction for code
3. **Blockquote styling** - Left border and background
4. **Heading sizes** - Already implemented, may need tweaking
5. **Horizontal rules** - Better visual separators
6. **Link styling** - Clickable, colored links

## Testing Strategy

1. Create test markdown files covering all element types
2. Test each change in isolation
3. Verify no regression in text rendering
4. Test with real-world markdown files (README.md)
5. Test both light and dark themes

## Files to Modify

- `internal/markdown/renderer.go` - Main rendering logic
- `internal/themes/themes.go` - Theme settings for spacing/padding
- `internal/gui/window.go` - May need to adjust scroll container setup

## Success Criteria

1. All markdown elements render correctly (no missing text)
2. Visual spacing between elements
3. Code blocks visually distinct
4. Blockquotes visually distinct
5. No performance regression
6. Works in both light and dark themes
