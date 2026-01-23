# CLAUDE.md

This file provides guidance to Claude Code when working with the MarkView codebase.

## Git Pre-commit Hook

This repository has a pre-commit hook that shows files being committed. In interactive mode (terminal), it prompts for authorization. In non-interactive mode (scripts, Claude Code), it auto-authorizes.

## Project Overview

MarkView is a cross-platform markdown viewer and editor built with Go and the Fyne GUI toolkit. It features live reload, syntax highlighting, multiple themes, file browser, table of contents, edit mode with formatting toolbar, and extensive productivity features.

## Build Commands

```bash
# Build the application
make build          # Creates bin/markview

# Run directly
make run            # Builds and runs
make run-sample     # Runs with testdata/sample.md

# Testing
make test           # Run all tests
make test-coverage  # Tests with coverage report

# Code quality
make fmt            # Format code
make vet            # Run go vet
make lint           # Run golangci-lint

# Packaging
make package-deb    # Build .deb for Debian/Ubuntu (Linux only)
make package-dmg    # Build .dmg for macOS (macOS only)

# Cross-platform (requires fyne-cross + Docker)
make build-linux    # Build for Linux
make build-macos    # Build for macOS
```

## Project Structure

```
markview/
├── cmd/markview/main.go     # Entry point, CLI parsing
├── internal/
│   ├── gui/                 # User interface
│   │   ├── window.go        # Main window, toolbar, menus
│   │   ├── editor.go        # Markdown editor widget
│   │   ├── autocomplete.go  # Link autocomplete for [text](path
│   │   ├── updater.go       # Auto-update checker via GitHub API
│   │   ├── filetree.go      # File browser tree
│   │   ├── quickswitcher.go # Fuzzy file finder (Ctrl+P)
│   │   ├── search.go        # Full-text search
│   │   ├── backlinks.go     # Backlinks panel
│   │   ├── tags.go          # Tag management
│   │   ├── templates.go     # Document templates
│   │   ├── export_themes.go # Export styling (6 themes)
│   │   ├── spellcheck.go    # Spell checking (aspell)
│   │   ├── linkvalidation.go # Link checker
│   │   └── imagepaste.go    # Image upload dialog
│   ├── markdown/            # Markdown processing
│   │   ├── parser.go        # Goldmark wrapper
│   │   ├── renderer.go      # AST to Fyne widgets (includes table rendering)
│   │   ├── syntax.go        # Chroma syntax highlighting
│   │   ├── codeblock.go     # Code block widget
│   │   ├── blockquote.go    # Blockquote widget
│   │   ├── spacer.go        # Vertical spacing widget
│   │   └── tablecontainer.go # Table container widget
│   ├── toc/                 # Table of contents
│   │   ├── generator.go     # Heading extraction
│   │   └── navigation.go    # TOC tree widget
│   ├── watcher/             # File monitoring
│   │   ├── watcher.go       # fsnotify wrapper with pause/resume
│   │   └── debounce.go      # Change debouncing
│   ├── library/             # Document library management
│   └── themes/              # Theme definitions
│       ├── themes.go        # 8 color themes + icon functions
│       └── icons.go         # SVG icons for toolbar
├── scripts/                 # Build scripts
│   ├── build-deb.sh         # Debian/Ubuntu packaging
│   └── build-dmg.sh         # macOS DMG packaging
├── assets/                  # Logo and icons
└── testdata/                # Sample markdown files
```

## Key Dependencies

- **Fyne v2.7** - GUI toolkit (`fyne.io/fyne/v2`)
- **Goldmark v1.7** - Markdown parser (`github.com/yuin/goldmark`)
- **Chroma v2.23** - Syntax highlighter (`github.com/alecthomas/chroma/v2`)
- **fsnotify v1.9** - File system watcher (`github.com/fsnotify/fsnotify`)
- **zap v1.27** - Structured logging (`go.uber.org/zap`)

## Architecture Notes

### Edit Mode
The application supports toggling between view and edit modes:
- `Window.editMode` tracks the current mode
- `Window.isDirty` tracks unsaved changes
- `Window.contentStack` is a Stack container holding both the rendered view and editor
- File watcher is paused during edit mode to prevent conflicts

### Toolbar
Uses icon-based toolbar:
- Custom `toolbarAction` type implements `widget.ToolbarItem`
- Icons defined as SVG resources in `themes/icons.go`
- Theme functions in `themes/themes.go` expose icons

### Rendering Pipeline
```
Markdown → Goldmark Parser → AST → Renderer → Fyne Widgets
                                     ↳ Chroma (code highlighting)
```

### Split View Mode
Split view provides side-by-side editing with live preview:
- `Window.splitEditor` - Separate editor instance for split view
- `Window.splitEditorScroll` - Scroll container for split editor
- `Window.splitPreview` - Live preview pane showing rendered markdown
- `Window.splitLinkAutocomplete` - Autocomplete for split editor
- `onSplitEditorChanged()` - Re-renders preview on each keystroke
- Editor and preview use separate instances to avoid conflicts
- Toggle with Cmd+\\ or toolbar icon (switches between split/single view icons)

### Table Rendering
GFM tables are rendered with styled grid layout:
- Uses `HBox`/`VBox` containers instead of Fyne's Table widget (no scrollbars)
- `FixedWidthContainer` - Custom widget for fixed-width cells with flexible height
- Headers styled with primary background color and orange text (`tableHeader` color)
- Alternating row colors for readability
- Word wrapping with automatic row height calculation
- Column widths based on content with padding

### Custom Widgets
- `CodeBlock` - Renders code with syntax highlighting and background
- `Blockquote` - Styled quotes with left border
- `Spacer` - Vertical spacing between elements
- `MarkdownEditor` - Multi-line entry with text manipulation helpers
- `LinkAutocomplete` - Inline file suggestions for markdown links
- `FixedWidthContainer` - Fixed-width container with flexible height for table cells

### Link Autocomplete
Shows file suggestions when typing markdown links:
- **Standard links** `[text](path...` - suggests `.md` files
- **Wiki links** `[[path...` - suggests `.md` files, auto-closes with `]]`
- **Image links** `![alt](path...` - suggests image files (png, jpg, gif, svg, etc.)

Implementation in `autocomplete.go`:
- `LinkType` enum: `LinkTypeStandard`, `LinkTypeWiki`, `LinkTypeImage`
- `detectLinkContext()` determines link type from cursor position
- Scans directory for markdown and image files separately
- Fuzzy matching filters suggestions as user types
- Keyboard navigation: Arrow keys, Enter/Tab to accept, Escape to dismiss

### File Tree Action Bar
Action bar at the bottom of the file tree for file/directory operations in `filetree.go`:
- **New File** - Creates a new markdown file in selected directory (or parent of selected file)
- **New Folder** - Creates a new directory
- **Rename** - Renames the selected file or directory
- **Delete** - Deletes the selected file or directory with confirmation

Implementation:
- `createActionBar()` creates four icon buttons using Fyne theme icons
- `selectedPath` and `selectedIsDir` track the current selection
- Dialog methods: `showNewFileDialog()`, `showNewDirectoryDialog()`, `showRenameDialog()`, `showDeleteConfirmation()`
- All dialogs are resizable (400x150) using `dialog.NewForm()` with `Resize()`
- Auto-adds `.md` extension to new files if not present
- Empty directories are shown in the tree (no longer filtered out)

### Update Checker
Auto-update checking in `updater.go`:
- Fetches latest release from GitHub API
- Downloads universal DMG (works on both Intel and Apple Silicon)
- Supports "Skip This Version" preference to dismiss updates
- Checks once per day in silent mode (startup), on-demand via Help menu

### Color Themes (8 total)
Light, Dark, Nord, Solarized Light, Solarized Dark, Monokai, Gruvbox Dark, One Dark

### Export Themes (6 total)
Default, GitHub, Academic, Dark, Minimal, Print-Friendly

## Testing

Tests are located alongside source files with `_test.go` suffix:
```bash
go test ./internal/gui/... -v
go test ./internal/watcher/... -v
go test ./internal/markdown/... -v
go test ./internal/toc/... -v
go test ./internal/themes/... -v
```

Key test files:
- `internal/gui/autocomplete_test.go` - Link autocomplete tests
- `internal/gui/filetree_test.go` - File tree and context menu tests
- `internal/themes/themes_test.go` - Theme, font, and icon tests

## Common Patterns

### Adding a Toolbar Icon
1. Add SVG resource in `internal/themes/icons.go`
2. Add icon function in `internal/themes/themes.go`
3. Create `toolbarAction` in `window.go` using `newToolbarAction()`

### Adding Keyboard Shortcuts
Add in `Window.setupShortcuts()` using `desktop.CustomShortcut`:
```go
w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
    KeyName:  fyne.KeyX,
    Modifier: fyne.KeyModifierSuper,
}, func(shortcut fyne.Shortcut) {
    // handler
})
```

### Theme Colors
Custom colors defined in `themes.go`:
- Theme-specific palettes for each of the 8 themes
- Custom markdown color names: "heading1", "heading2", "bold", "link", etc.

### Adding Autocomplete Features
The `LinkAutocomplete` pattern can be extended:
1. Create detection function (like `isInLinkDestination`)
2. Add filtering logic with fuzzy matching
3. Create popup with `widget.NewPopUp()` and `widget.List`
4. Handle keyboard events in `HandleKeyEvent()`
5. Insert selection with cursor position calculation

## Key Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Cmd+E | Toggle edit mode |
| Cmd+S | Save file |
| Cmd+\\ | Toggle split view |
| Ctrl+P | Quick file switcher |
| Cmd+Shift+G | Search in all files |
| Cmd+B | Show backlinks |
| Cmd+T | Browse tags |
| Cmd+L | Validate links |
| F11 | Zen mode |
| **File Browser** | |
| Action Bar | New File, New Folder, Rename, Delete buttons |
| **Link Autocomplete** | |
| Up/Down | Navigate suggestions |
| Enter/Tab | Accept selection |
| Escape | Dismiss popup |

## Fyne-Specific Notes

- Use `container.NewStack()` for overlapping widgets (only one visible)
- Call `widget.Refresh()` after modifying widget state
- `widget.Entry` has built-in undo/redo support
- File dialogs: `dialog.NewFileOpen()`, `dialog.NewFolderOpen()`
- Theme changes: `app.Settings().SetTheme()`
