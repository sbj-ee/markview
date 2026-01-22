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

### Phase 1: Custom Spacer Widget (Low Risk) ✅ COMPLETED

1. ✅ Created `spacer.go` - Simple `Spacer` widget with configurable height
2. ✅ Used between major elements (headings, paragraphs, code blocks)
3. ✅ Tests added in `spacer_test.go`

### Phase 2: Code Block Styling (Medium Risk) ✅ COMPLETED

1. ✅ Created `codeblock.go` - `CodeBlock` widget with background, padding, rounded corners
2. ✅ Handles sizing explicitly with custom renderer
3. ✅ Tests added in `codeblock_test.go`
4. ✅ Blockquote styling also implemented in `blockquote.go` with left border

### Phase 3: Single RichText Approach (Higher Risk) - SKIPPED

1. Skipped - current widget-based approach works well
2. Widget approach provides better control over visual elements
3. No need for the complexity of a single RichText approach

## Visual Improvements Priority

1. ✅ **Spacing between elements** - Implemented with custom Spacer widget
2. ✅ **Code block backgrounds** - Implemented with CodeBlock widget (background, padding, rounded corners)
3. ✅ **Blockquote styling** - Implemented with Blockquote widget (left border and subtle background)
4. ✅ **Heading sizes** - Already implemented with theme customization
5. ✅ **Horizontal rules** - Implemented with Fyne Separator widget
6. 🔄 **Link styling** - Partially implemented (displayed with URL, not clickable yet)

## Testing Strategy

1. Create test markdown files covering all element types
2. Test each change in isolation
3. Verify no regression in text rendering
4. Test with real-world markdown files (README.md)
5. Test both light and dark themes

## Files Modified/Created

### New Files Created
- `internal/markdown/spacer.go` - Custom Spacer widget for vertical spacing
- `internal/markdown/spacer_test.go` - Tests for Spacer widget
- `internal/markdown/codeblock.go` - Custom CodeBlock widget with background styling
- `internal/markdown/codeblock_test.go` - Tests for CodeBlock widget
- `internal/markdown/blockquote.go` - Custom Blockquote widget with border styling
- `internal/markdown/blockquote_test.go` - Tests for Blockquote widget

### Files Modified
- `internal/markdown/renderer.go` - Updated to use custom widgets
- `internal/markdown/renderer_test.go` - Added normalizeText and additional rendering tests
- `internal/themes/themes.go` - Theme settings for sizing
- `README.md` - Documentation updates

## Success Criteria

1. ✅ All markdown elements render correctly (no missing text)
2. ✅ Visual spacing between elements (custom Spacer widget)
3. ✅ Code blocks visually distinct (background, padding, rounded corners)
4. ✅ Blockquotes visually distinct (left border, subtle background)
5. ✅ No performance regression (all tests pass)
6. ✅ Works in both light and dark themes (theme-aware colors)
