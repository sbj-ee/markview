# Visual Rendering Improvement Plan

## Status: COMPLETED

All major visual improvements have been implemented and the markdown viewer now has a polished, modern appearance.

## Final Architecture

```
Markdown Source
    ↓
Goldmark Parser → AST
    ↓
Renderer.Render() → []fyne.CanvasObject (custom widgets)
    ├── Spacer (vertical spacing)
    ├── RichText (headings with level-based colors)
    ├── Label (paragraphs)
    ├── CodeBlock (syntax highlighted code with background)
    ├── Blockquote (styled quotes with left border)
    └── Separator (horizontal rules)
    ↓
container.NewVBox(widgets...) → displayed in ScrollContainer
```

## Implemented Solutions

### Custom Widgets Created

1. **Spacer Widget** (`spacer.go`)
   - Configurable height for vertical spacing
   - Used between all major elements (headings, paragraphs, code blocks, lists)

2. **CodeBlock Widget** (`codeblock.go`)
   - Dark background (#282A30)
   - 8px padding
   - 4px rounded corners
   - Contains syntax-highlighted RichText
   - No text wrapping (horizontal scroll for long lines)

3. **Blockquote Widget** (`blockquote.go`)
   - Cyan left border (4px wide)
   - Subtle dark background
   - 12px padding
   - Italic text styling

### Theme Customization (`themes.go`)

**Color Scheme (Dark Theme):**
- Background: Dark charcoal (#1E2024)
- Foreground: Light gray (#C8C8C8)
- H1, H2 headings: Cyan/teal (#56B6C2)
- H3, H4 headings: Orange/gold (#E5B567)
- H5, H6 headings: Light gray
- Links: Cyan (#56B6C2)
- Code/blockquote backgrounds: Darker gray (#282A30)

**Font Sizes:**
- Body text: 18pt
- H1: 28pt (HeadingText)
- H2: 22pt (SubHeadingText)
- H3: 21pt
- H4: 18pt
- H5, H6: 16pt, 15pt
- TOC/captions: 13pt

### Renderer Updates (`renderer.go`)

- Level-based heading colors and sizes
- H1 headings include horizontal rule underneath
- Proper spacing between all elements
- Text normalization to prevent unwanted line breaks
- HTML entity decoding for smart quotes and em-dashes

### UI Updates

**TOC Navigation (`navigation.go`):**
- Smaller font (13pt CaptionText)
- Uses RichText instead of Label for size control
- Hierarchical display with tree expansion

**Toolbar (`window.go`):**
- Compact buttons with icons
- LowImportance styling for smaller appearance

## Files Created/Modified

### New Files
- `internal/markdown/spacer.go` - Custom Spacer widget
- `internal/markdown/spacer_test.go` - Spacer tests
- `internal/markdown/codeblock.go` - Custom CodeBlock widget
- `internal/markdown/codeblock_test.go` - CodeBlock tests
- `internal/markdown/blockquote.go` - Custom Blockquote widget
- `internal/markdown/blockquote_test.go` - Blockquote tests

### Modified Files
- `internal/markdown/renderer.go` - Heading colors, spacing, text normalization
- `internal/markdown/renderer_test.go` - Additional tests for normalizeText
- `internal/themes/themes.go` - Color scheme, font sizes, icon helpers
- `internal/toc/navigation.go` - Smaller TOC font using RichText
- `internal/gui/window.go` - Compact toolbar with icons
- `README.md` - Removed emojis for Fyne Label compatibility

## Visual Improvements Completed

| Feature | Status | Implementation |
|---------|--------|----------------|
| Spacing between elements | ✅ | Custom Spacer widget |
| Code block backgrounds | ✅ | CodeBlock widget with dark bg, padding, rounded corners |
| Blockquote styling | ✅ | Blockquote widget with cyan left border |
| Heading hierarchy | ✅ | Level-based colors (cyan H1-H2, orange H3-H4) |
| Horizontal rules | ✅ | Fyne Separator widget, auto after H1 |
| Smaller TOC font | ✅ | RichText with CaptionText size (13pt) |
| Compact toolbar | ✅ | Icon buttons with LowImportance |
| Dark theme colors | ✅ | Charcoal bg, cyan/orange accents |
| Link styling | ✅ | Cyan color (not clickable yet) |

## Known Limitations

1. **Emojis**: Fyne Label widgets may not render certain emojis properly on all systems. Emojis have been removed from README.md.

2. **Clickable links**: Links are displayed with cyan color but are not yet clickable. Would require Hyperlink widget or tap handler.

3. **Inline formatting in paragraphs**: Paragraphs use plain Label widgets to avoid line break issues, which means bold/italic formatting within paragraphs is lost. Only heading text preserves bold styling.

## Future Improvements (Optional)

1. **Clickable links** - Implement tap handlers or use Hyperlink segments
2. **Inline formatting** - Explore RichText improvements to preserve bold/italic in paragraphs
3. **Image rendering** - Display actual images instead of alt text
4. **Light theme polish** - Update light theme colors to match the dark theme aesthetic
5. **Custom fonts** - Bundle Inter or similar font for consistent cross-platform rendering

## Test Results

```
internal/markdown   - 47 tests (46 passed, 1 skipped)
internal/toc        - 9 tests (all passed)
internal/watcher    - 8 tests (all passed)
```

All tests pass. The visual rendering is stable and matches the target aesthetic.
