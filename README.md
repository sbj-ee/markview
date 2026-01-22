# MarkView

<p align="center">
  <strong>A fast, cross-platform markdown viewer with live reload and syntax highlighting</strong>
</p>

<p align="center">
  Built with Go and Fyne • Supports macOS and Linux • Lightweight and Native
</p>

---

## Features

### Core Functionality
- ✅ **Full Markdown Support** - CommonMark and GitHub Flavored Markdown (GFM)
- ✅ **Syntax Highlighting** - 250+ languages via Chroma with beautiful color schemes
- ✅ **Live Reload** - Auto-refresh on file changes (300ms debounce)
- ✅ **Table of Contents** - Hierarchical navigation with click-to-scroll
- ✅ **Split View** - Resizable TOC sidebar with content area
- ✅ **Custom Themes** - Beautiful light and dark themes with instant switching
- ✅ **Smart Typography** - Automatic smart quotes, em-dashes, and proper character rendering
- ✅ **Cross-Platform** - Native performance on macOS and Linux

### Markdown Features Supported
- **Headings** (H1-H6) with bold styling
- **Text Formatting** - *italic*, **bold**, ***bold italic***
- **Code Blocks** - Fenced (` ```language `) and indented, with syntax highlighting
- **Inline Code** - Monospace rendering with `backticks`
- **Lists** - Ordered and unordered, nested support
- **Blockquotes** - Italic rendering
- **Links** - Displayed with URL
- **Images** - Alt text display (full image rendering coming soon)
- **Tables** - GFM table support
- **Task Lists** - GFM task list support
- **Strikethrough** - GFM strikethrough support
- **Smart Typography** - Smart quotes, em dashes, etc.

---

## Installation

### Prerequisites

- **Go 1.21+** (for building from source)
- **GCC** or compatible C compiler (required by Fyne for GUI)

### Quick Install

```bash
# Clone the repository
git clone <repository-url>
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

#### Open a specific file:
```bash
markview -file README.md
```

#### Open file picker dialog:
```bash
markview
```

#### Show version:
```bash
markview -version
```

### GUI Operations

| Action | Method |
|--------|--------|
| **Open File** | Click "Open File" button or File → Open |
| **Refresh** | Click "Refresh" button or File → Refresh |
| **Toggle TOC** | View → Toggle TOC |
| **Switch Theme** | View → Light Theme / Dark Theme |
| **Navigate** | Click TOC entries to jump to sections |
| **Resize TOC** | Drag the split divider |
| **Quit** | File → Quit or Cmd+Q (macOS) / Ctrl+Q (Linux) |

### Live Reload

MarkView automatically watches your file for changes. Perfect for:
- 📝 Real-time markdown previewing while writing
- 📖 Viewing auto-generated documentation
- 🔄 Live editing workflows

The file watcher uses a 300ms debounce to prevent rapid re-renders during saves.

### Beautiful Themes

MarkView includes two carefully crafted themes with enhanced typography:

#### 🌙 Dark Theme (Default)
- Dracula-inspired color palette
- Deep purple-gray background (#282A36)
- Vibrant syntax highlighting
- **Large, blue-colored headings** for excellent hierarchy
- Easy on the eyes for extended reading sessions
- Perfect for low-light environments

#### ☀️ Light Theme
- GitHub-inspired clean design
- Pure white background with professional colors
- **Blue headings with increased font sizes** (H1: 28pt, H2: 22pt)
- High contrast for readability
- Ideal for daytime use and printing
- Easy on battery life

**Switch themes instantly** via View → Light Theme / Dark Theme menu.

### Typography & Rendering

- **Larger headings** with visual hierarchy (H1 @ 28pt, H2 @ 22pt, H3 @ 21pt)
- **Colored headings** in professional blue shades
- **Bold heading styling** for improved scannability
- **Smart typography** - Automatic conversion of quotes and dashes (e.g., "smart quotes", em-dashes)
- **Clean rendering** without markdown syntax symbols (no fence markers shown)
- **Proper code formatting** - Syntax highlighted code blocks with no line wrapping
- **Proper spacing** between elements
- **Widget-based rendering** for precise control over appearance

---

## Development

### Build Commands

```bash
# Build the application
make build

# Run directly
make run

# Run with sample file
make run-sample

# Run with theme showcase
make run-theme

# Clean build artifacts
make clean
```

### Testing

The project includes comprehensive unit tests:

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./internal/... -v

# Run tests with coverage
make test-coverage
```

#### Test Coverage

- ✅ **Parser Tests** - Markdown parsing and AST generation
- ✅ **Renderer Tests** - AST to Fyne widget conversion
- ✅ **TOC Tests** - Heading extraction and hierarchy
- ✅ **Syntax Highlighter Tests** - Code highlighting for multiple languages
- ✅ **Debouncer Tests** - File change debouncing logic
- ✅ **Watcher Tests** - File system monitoring

**Test Results:**
```
internal/markdown   - 22 tests (21 passed, 1 skipped)
internal/toc        - 9 tests (all passed)
internal/watcher    - 8 tests (all passed)
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint
```

---

## Cross-Platform Building

MarkView supports cross-compilation for multiple platforms using `fyne-cross`.

### Setup

```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Requires Docker for cross-compilation
```

### Build Commands

```bash
# Build for Linux (amd64 & arm64)
make build-linux

# Build for macOS (amd64 & arm64)
make build-macos

# Build for all platforms
make build-all
```

Binaries will be output to `fyne-cross/dist/`.

### Packaging

```bash
# Create macOS .app bundle
make package-macos

# Create Linux package
make package-linux
```

---

## Architecture

### Project Structure

```
markview/
├── cmd/markview/           # Application entry point
│   └── main.go             # CLI parsing and app initialization
├── internal/
│   ├── gui/                # User interface
│   │   └── window.go       # Main window, menus, toolbar
│   ├── markdown/           # Markdown processing
│   │   ├── parser.go       # Goldmark wrapper
│   │   ├── renderer.go     # AST to Fyne conversion
│   │   └── syntax.go       # Chroma integration
│   ├── toc/                # Table of contents
│   │   ├── generator.go    # Heading extraction
│   │   └── navigation.go   # TOC tree widget
│   ├── watcher/            # File monitoring
│   │   ├── watcher.go      # fsnotify wrapper
│   │   └── debounce.go     # Change debouncing
│   └── images/             # Image handling (future)
├── testdata/               # Sample markdown files
├── assets/                 # Application assets
├── Makefile                # Build automation
├── README.md               # This file
└── IMPLEMENTATION.md       # Technical details
```

### Rendering Pipeline

```
Markdown File
    ↓
Goldmark Parser (CommonMark + GFM)
    ↓
Abstract Syntax Tree (AST)
    ├→ TOC Generator → Heading Tree
    └→ Renderer
        ├→ Chroma (syntax highlighting)
        └→ Fyne Widgets (RichText, Tree, etc.)
            ↓
        Display in Window
```

### Key Technologies

| Component | Technology | Version |
|-----------|-----------|---------|
| GUI Framework | [Fyne](https://fyne.io) | v2.7.2 |
| Markdown Parser | [Goldmark](https://github.com/yuin/goldmark) | v1.7.16 |
| Syntax Highlighter | [Chroma](https://github.com/alecthomas/chroma) | v2.23.0 |
| File Watcher | [fsnotify](https://github.com/fsnotify/fsnotify) | v1.9.0 |
| Logger | [zap](https://github.com/uber-go/zap) | v1.27.1 |

---

## Performance

- **Fast Parsing** - Goldmark is optimized for speed
- **Efficient Rendering** - Only re-renders on file changes
- **Low Memory** - ~50MB typical memory usage
- **Native UI** - Fyne provides native performance
- **Debounced Reload** - Prevents excessive re-renders

### Benchmark (Apple M1, macOS)

| File Size | Parse Time | Render Time | Total |
|-----------|-----------|-------------|-------|
| 10 KB | < 1ms | 2-3ms | ~3ms |
| 100 KB | 3-5ms | 10-15ms | ~18ms |
| 1 MB | 30-40ms | 100-150ms | ~180ms |

---

## Roadmap

### ✅ Implemented (v0.1.0)
- Basic markdown rendering with smart typography
- Syntax highlighting with beautiful colors (no line wrapping in code blocks)
- Live file reload with auto-watch
- Table of contents with click-to-scroll navigation
- Split view layout with resizable panels
- File operations (open, refresh)
- Custom light and dark themes
- Theme switching with instant reload
- HTML entity rendering (smart quotes, em-dashes, etc.)
- Clean code block display (no fence markers)

### 🚧 In Progress
- Image rendering (local files)
- Additional theme options (Nord, Solarized, etc.)

### 📋 Planned
- **v0.2.0**
  - Full image support (local + remote)
  - Keyboard shortcuts (Ctrl+O, Ctrl+R, etc.)
  - Additional themes (Nord, Solarized, Monokai, etc.)
  - Enhanced TOC with current position tracking

- **v0.3.0**
  - Export to HTML/PDF
  - Search functionality
  - Recent files list
  - Preferences dialog

- **v0.4.0**
  - Custom CSS styling
  - Plugin system
  - Multi-file viewing (tabs)
  - Markdown editing (optional)

---

## Contributing

Contributions are welcome! Here's how to get started:

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/yourusername/markview.git
   cd markview
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make your changes**
   - Write code
   - Add tests
   - Update documentation

5. **Test your changes**
   ```bash
   make test
   make fmt
   make vet
   ```

6. **Submit a pull request**

### Code Guidelines

- Follow Go best practices
- Add tests for new features
- Update documentation
- Run `make fmt` before committing
- Keep commits atomic and well-described

---

## Troubleshooting

### Build Issues

**"command not found: go"**
- Install Go 1.21+ from https://go.dev/dl/

**"gcc: command not found"**
- macOS: Install Xcode Command Line Tools: `xcode-select --install`
- Linux: Install build-essential: `sudo apt install build-essential`

### Runtime Issues

**"Failed to initialize logger"**
- Check file permissions in the application directory

**"File watcher not working"**
- Verify the file exists and is readable
- Check system file descriptor limits (Linux: `ulimit -n`)

**"Syntax highlighting not working"**
- Verify language name in fence (e.g., ` ```go ` not ` ```golang `)
- Check Chroma supports the language: https://github.com/alecthomas/chroma

---

## FAQ

**Q: What markdown flavor is supported?**
A: CommonMark with GitHub Flavored Markdown (GFM) extensions.

**Q: Can I use MarkView on Windows?**
A: Not currently, but Fyne supports Windows. Contributions welcome!

**Q: Does MarkView support custom CSS?**
A: Not yet, but it's planned for v0.4.0.

**Q: Can I edit markdown in MarkView?**
A: Currently read-only. Editing support is being considered for future versions.

**Q: How do I change the theme?**
A: Go to View → Light Theme or View → Dark Theme. The entire interface updates instantly!

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

Special thanks to the Go community for creating amazing tools!

---

## Support & Contact

- **Issues**: [GitHub Issues](https://github.com/yourusername/markview/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/markview/discussions)
- **Email**: your.email@example.com

---

<p align="center">
  Made with ❤️ using Go and Fyne
</p>
