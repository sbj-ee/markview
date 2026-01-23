# MarkView

<p align="center">
  <img src="assets/logo-128.png" alt="MarkView Logo" width="128" height="128">
</p>

<p align="center">
  <strong>A powerful, cross-platform markdown viewer and editor with advanced features</strong>
</p>

<p align="center">
  Built with Go and Fyne • Supports macOS and Linux • Lightweight and Native
</p>

---

## Features

### Viewing & Editing
- **Full Markdown Support** - CommonMark and GitHub Flavored Markdown (GFM) with styled tables
- **Syntax Highlighting** - 250+ languages via Chroma
- **Live Reload** - Auto-refresh on file changes (300ms debounce)
- **Edit Mode** - Toggle between viewing and editing with formatting toolbar
- **Link Autocomplete** - File suggestions when typing `[text](`, `[[`, or `![alt](`
- **Split View** - Side-by-side editor with live preview (Cmd+\\)
- **Focus Mode** - Hide sidebars for distraction-free reading (Cmd+Shift+F)
- **Zen Mode** - Fullscreen distraction-free writing (F11)
- **Typewriter Mode** - Keep cursor centered while typing
- **Auto-Save** - Automatically saves changes every 30 seconds

### Navigation & Search
- **File Browser** - Built-in file tree with filtering (Cmd+F) and action bar
  - Action bar buttons: New File, New Folder, Rename, Delete
  - Select a file or directory, then use the action buttons at the bottom
- **Table of Contents** - Hierarchical navigation with click-to-scroll
- **Quick Switcher** - Fuzzy file finder (Ctrl+P)
- **Full-Text Search** - Search across all markdown files (Cmd+Shift+G)
- **Backlinks Panel** - See documents that link to the current file (Cmd+B)
- **Recent Files** - Quick access to recently opened files (Cmd+Shift+O)

### Organization
- **Library Mode** - Browse and manage all documents in a folder
- **Tags Support** - Organize documents with #hashtags (Cmd+T)
- **Starred Documents** - Mark favorite documents for quick access
- **Document Templates** - Create new files from templates (Cmd+Shift+N)
  - Meeting Notes, Blog Post, README, Technical Spec, Journal, and more

### Themes & Appearance
- **8 Color Themes** - Light, Dark, Nord, Solarized Light/Dark, Monokai, Gruvbox Dark, One Dark
- **Font Families** - System Default, Monospace, Serif, Sans Serif
- **Font Sizes** - Small, Normal, Large, Extra Large
- **Instant Theme Switching** - Click the palette icon to customize

### Export & Print
- **Export Formats** - HTML, PDF (via Chrome or wkhtmltopdf), DOCX, RTF (via pandoc)
- **6 Export Themes** - Default, GitHub, Academic, Dark, Minimal, Print-Friendly
- **Custom CSS** - Add your own styles for exports
- **Print via Browser** - Open in browser for system print dialog
- **Table of Contents in Export** - Automatic TOC generation in HTML exports

### Productivity
- **Word Count & Reading Time** - Live statistics in status bar
- **Word Count Goals** - Set and track writing goals
- **Find & Replace** - Search and replace in edit mode (Cmd+F)
- **Link Validation** - Check for broken links (Cmd+L)
- **Spell Checking** - Highlight misspelled words (requires aspell)
- **Snippets** - Quick-insert common markdown patterns
- **Presentation Mode** - View markdown as slides

### Additional Features
- **Drag & Drop** - Drop markdown files onto the window to open
- **Image Support** - Local image rendering with upload dialog
- **Window Size Persistence** - Remembers your window size
- **Keyboard Shortcuts** - Comprehensive shortcuts for all operations
- **Auto-Update Check** - Notifies when new versions are available

---

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| **File Operations** | |
| Cmd+N | New file |
| Cmd+Shift+N | New from template |
| Cmd+O | Open file |
| Cmd+S | Save file |
| Cmd+Shift+S | Save as |
| Cmd+P | Print/Export |
| Cmd+R | Refresh |
| **Edit Mode** | |
| Cmd+E | Toggle edit mode |
| Cmd+F | Find/Replace (edit) or Filter (view) |
| Cmd+\\ | Toggle split view |
| Escape | Exit edit mode |
| **Link Autocomplete** | |
| Down Arrow | Select next file suggestion |
| Up Arrow | Select previous file suggestion |
| Enter/Tab | Accept selected file |
| Escape | Dismiss suggestions |
| **File Browser** | |
| Action Bar | New File, New Folder, Rename, Delete buttons |
| **Navigation** | |
| Ctrl+P | Quick file switcher |
| Cmd+Shift+G | Search in all files |
| Cmd+B | Show backlinks |
| Cmd+T | Browse by tags |
| Cmd+Shift+O | Recent files |
| Alt+Up | Navigate to parent directory |
| **View** | |
| Cmd+Shift+F | Toggle focus mode |
| F11 | Toggle zen mode (fullscreen) |
| **Tools** | |
| Cmd+L | Validate links |
| Cmd+? | Show shortcuts |

---

## Installation

### Prerequisites

- **Go 1.21+** (for building from source)
- **GCC** or compatible C compiler (required by Fyne for GUI)

### Optional Dependencies

- **aspell** - For spell checking (`brew install aspell` or `apt install aspell`)
- **Google Chrome** - For PDF export (recommended, uses Chrome headless)
- **wkhtmltopdf** - Alternative for PDF export (https://wkhtmltopdf.org/downloads.html)
- **pandoc** - For DOCX/RTF export (`brew install pandoc`)

### Quick Install

```bash
# Clone the repository
git clone https://github.com/sbj-ee/markview.git
cd markview

# Install to $GOPATH/bin
make install
```

### Build from Source

```bash
# Build binary to bin/markview
make build

# The binary will be at: bin/markview
```

---

## Usage

### Command Line Interface

```bash
# Open a specific file
markview -file README.md

# Open file picker dialog
markview

# Show version
markview -version
```

### Edit Mode Toolbar

When in edit mode, the formatting toolbar provides quick access to:

| Button | Action |
|--------|--------|
| **B** | Bold - wrap with `**` |
| *I* | Italic - wrap with `*` |
| H1/H2/H3 | Insert heading |
| Link | Insert `[text](url)` |
| Image | Insert image dialog |
| Table | Visual table editor |
| Code | Inline code |
| Code Block | Fenced code block |
| Quote | Blockquote |
| List | Bullet list |
| HR | Horizontal rule |
| Snippet | Insert template snippets |
| Typewriter | Toggle typewriter mode |
| Goal | Set word count goal |

---

## Themes

### Color Themes

| Theme | Description |
|-------|-------------|
| **Light** | Clean, professional design |
| **Dark** | Dracula-inspired dark theme (default) |
| **Nord** | Arctic, north-bluish color palette |
| **Solarized Light** | Warm, yellowish light theme |
| **Solarized Dark** | Blue-tinted dark theme |
| **Monokai** | Classic code editor colors |
| **Gruvbox Dark** | Retro groove colors |
| **One Dark** | Atom-inspired dark theme |

### Export Themes

| Theme | Best For |
|-------|----------|
| **Default** | General purpose |
| **GitHub** | GitHub-style documentation |
| **Academic** | Papers and formal documents |
| **Dark** | Screen viewing |
| **Minimal** | Clean, serif typography |
| **Print-Friendly** | Physical printing |

---

## Document Templates

Create new documents from templates with Cmd+Shift+N:

- **Blank Document** - Start fresh
- **Meeting Notes** - Attendees, agenda, action items
- **Blog Post** - YAML frontmatter, sections
- **Project README** - Features, installation, usage
- **Technical Specification** - Requirements, architecture, timeline
- **Daily Journal** - Gratitude, goals, reflection
- **Weekly Review** - Accomplishments, lessons, planning
- **Code Review Notes** - Findings, checklist, approval

---

## Architecture

### Project Structure

```
markview/
├── cmd/markview/           # Application entry point
├── internal/
│   ├── gui/                # User interface
│   │   ├── window.go       # Main window and toolbar
│   │   ├── editor.go       # Markdown editor widget
│   │   ├── filetree.go     # File browser
│   │   ├── quickswitcher.go # Fuzzy file finder
│   │   ├── search.go       # Full-text search
│   │   ├── backlinks.go    # Backlinks panel
│   │   ├── tags.go         # Tag management
│   │   ├── templates.go    # Document templates
│   │   ├── export_themes.go # Export styling
│   │   ├── spellcheck.go   # Spell checking
│   │   └── ...
│   ├── markdown/           # Markdown processing
│   │   ├── parser.go       # Goldmark wrapper
│   │   ├── renderer.go     # AST to Fyne conversion
│   │   └── syntax.go       # Chroma integration
│   ├── toc/                # Table of contents
│   ├── watcher/            # File monitoring
│   ├── library/            # Document library
│   └── themes/             # Theme definitions
├── assets/                 # Logo and icons
├── testdata/               # Sample files
└── Makefile                # Build automation
```

### Key Technologies

| Component | Technology |
|-----------|-----------|
| GUI Framework | [Fyne](https://fyne.io) v2.7.2 |
| Markdown Parser | [Goldmark](https://github.com/yuin/goldmark) v1.7.16 |
| Syntax Highlighter | [Chroma](https://github.com/alecthomas/chroma) v2.23.0 |
| File Watcher | [fsnotify](https://github.com/fsnotify/fsnotify) v1.9.0 |
| Logger | [zap](https://github.com/uber-go/zap) v1.27.1 |

---

## Development

### Build Commands

```bash
make build          # Build the application
make run            # Run directly
make test           # Run all tests
make test-coverage  # Run tests with coverage
make fmt            # Format code
make vet            # Run go vet
make lint           # Run linter
make clean          # Clean build artifacts
```

### Cross-Platform Building

```bash
# Requires fyne-cross and Docker
make build-linux    # Build for Linux
make build-macos    # Build for macOS
make build-all      # Build for all platforms
```

### Packaging

```bash
# Build .deb package for Debian/Ubuntu (run on Linux)
make package-deb

# Build .dmg package for macOS (run on macOS)
make package-dmg

# Specify version
VERSION=1.1.0 make package-deb
VERSION=1.1.0 make package-dmg
```

The packaging scripts create:
- **Linux (.deb)**: Full package with desktop entry, icons, MIME types, and AppStream metadata
- **macOS (.dmg)**: Universal binary (Intel + Apple Silicon) in a single DMG with app bundle, Info.plist and icons

---

## Troubleshooting

### Build Issues

**"command not found: go"**
- Install Go 1.21+ from https://go.dev/dl/

**"gcc: command not found"**
- macOS: `xcode-select --install`
- Linux: `sudo apt install build-essential`

### Runtime Issues

**Spell checking not working**
- Install aspell: `brew install aspell` or `apt install aspell`

**PDF export not working**
- Install Google Chrome (recommended) - PDF export will use Chrome headless
- Or install wkhtmltopdf from https://wkhtmltopdf.org/downloads.html
- Or use "Print via Browser" and save as PDF from the print dialog

**DOCX/RTF export not working**
- Install pandoc: `brew install pandoc`

---

## License

MIT License - see [LICENSE](LICENSE) file for details

---

## Acknowledgments

Built with excellent open-source projects:

- [Fyne](https://fyne.io) - Modern GUI toolkit for Go
- [Goldmark](https://github.com/yuin/goldmark) - Fast, extensible markdown parser
- [Chroma](https://github.com/alecthomas/chroma) - Pure Go syntax highlighter
- [fsnotify](https://github.com/fsnotify/fsnotify) - Cross-platform file system notifications
- [zap](https://github.com/uber-go/zap) - Blazing fast, structured logging

---

<p align="center">
  <img src="assets/logo-64.png" alt="MarkView" width="32" height="32">
  <br>
  Made with Go and Fyne
</p>
