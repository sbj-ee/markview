# CLAUDE.md

This file provides guidance to Claude Code when working with the MarkView codebase.

## Project Overview

MarkView is a cross-platform markdown viewer built with Go and the Fyne GUI toolkit. It features live reload, syntax highlighting, a file browser, table of contents navigation, and markdown editing capabilities.

## Build Commands

```bash
# Build the application
make build          # Creates bin/markview

# Run directly
make run            # Builds and runs
make run-sample     # Runs with testdata/sample.md

# Testing
make test           # Run all tests
go test ./internal/... -v  # Verbose test output

# Code quality
make fmt            # Format code
make vet            # Run go vet
```

## Project Structure

```
markview/
├── cmd/markview/main.go     # Entry point, CLI parsing
├── internal/
│   ├── gui/                 # User interface
│   │   ├── window.go        # Main window, toolbar, edit mode
│   │   ├── editor.go        # Markdown editor widget
│   │   └── filetree.go      # File browser tree
│   ├── markdown/            # Markdown processing
│   │   ├── parser.go        # Goldmark wrapper
│   │   ├── renderer.go      # AST to Fyne widgets
│   │   ├── syntax.go        # Chroma syntax highlighting
│   │   ├── codeblock.go     # Code block widget
│   │   ├── blockquote.go    # Blockquote widget
│   │   └── spacer.go        # Vertical spacing widget
│   ├── toc/                 # Table of contents
│   │   ├── generator.go     # Heading extraction
│   │   └── navigation.go    # TOC tree widget
│   ├── watcher/             # File monitoring
│   │   ├── watcher.go       # fsnotify wrapper with pause/resume
│   │   └── debounce.go      # Change debouncing
│   └── themes/              # Custom themes
│       ├── themes.go        # Light/dark theme definitions
│       └── icons.go         # SVG icons for toolbar
└── testdata/                # Sample markdown files
```

## Key Dependencies

- **Fyne v2** - GUI toolkit (`fyne.io/fyne/v2`)
- **Goldmark** - Markdown parser (`github.com/yuin/goldmark`)
- **Chroma** - Syntax highlighter (`github.com/alecthomas/chroma/v2`)
- **fsnotify** - File system watcher (`github.com/fsnotify/fsnotify`)
- **zap** - Structured logging (`go.uber.org/zap`)

## Architecture Notes

### Edit Mode
The application supports toggling between view and edit modes:
- `Window.editMode` tracks the current mode
- `Window.isDirty` tracks unsaved changes
- `Window.contentStack` is a Stack container holding both the rendered view and editor
- File watcher is paused during edit mode to prevent conflicts

### Toolbar
Uses icon-based toolbar (menu bar was removed):
- Custom `toolbarAction` type implements `widget.ToolbarItem`
- Icons defined as SVG resources in `themes/icons.go`
- Theme functions in `themes/themes.go` expose icons

### Rendering Pipeline
```
Markdown → Goldmark Parser → AST → Renderer → Fyne Widgets
                                     ↳ Chroma (code highlighting)
```

### Custom Widgets
- `CodeBlock` - Renders code with syntax highlighting and background
- `Blockquote` - Styled quotes with left border
- `Spacer` - Vertical spacing between elements
- `MarkdownEditor` - Multi-line entry with text manipulation helpers

## Testing

Tests are located alongside source files with `_test.go` suffix:
- `internal/gui/editor_test.go` - Editor text manipulation
- `internal/watcher/watcher_test.go` - File watcher pause/resume
- `internal/watcher/debounce_test.go` - Debouncer logic
- `internal/markdown/*_test.go` - Parser, renderer, widgets
- `internal/toc/generator_test.go` - TOC generation

Run specific package tests:
```bash
go test ./internal/gui/... -v
go test ./internal/watcher/... -v
```

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
- `lightColor()` - Light theme colors
- `darkColor()` - Dark theme colors
- Custom markdown color names: "heading1", "heading2", "bold", "link", etc.

## Fyne-Specific Notes

- Use `container.NewStack()` for overlapping widgets (only one visible)
- Call `widget.Refresh()` after modifying widget state
- `widget.Entry` has built-in undo/redo support
- File dialogs: `dialog.NewFileOpen()`, `dialog.NewFolderOpen()`
- Theme changes: `app.Settings().SetTheme()`
