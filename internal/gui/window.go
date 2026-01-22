package gui

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/sbj-ee/markview/internal/markdown"
	"github.com/sbj-ee/markview/internal/themes"
	"github.com/sbj-ee/markview/internal/toc"
	"github.com/sbj-ee/markview/internal/watcher"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.uber.org/zap"
)

// Window represents the main application window
type Window struct {
	fyneWindow    fyne.Window
	app           fyne.App
	parser        *markdown.Parser
	logger        *zap.Logger
	fileTree      *FileTree
	fileTreeScroll *container.Scroll
	tocTree       *widget.Tree
	tocScroll     *container.Scroll
	scrollContent *container.Scroll
	leftSplit     *container.Split  // File Tree | TOC
	mainSplit     *container.Split  // (File Tree | TOC) | Content
	currentFile   string
	currentDir    string
	fileWatcher   *watcher.FileWatcher
	currentTheme  themes.ThemeType
}

// NewWindow creates a new application window
func NewWindow(app fyne.App, logger *zap.Logger) *Window {
	// Create file watcher with 300ms debounce
	fw, err := watcher.NewFileWatcher(logger, 300*time.Millisecond)
	if err != nil {
		logger.Error("Failed to create file watcher", zap.Error(err))
		fw = nil
	}

	w := &Window{
		fyneWindow:   app.NewWindow("MarkView"),
		app:          app,
		parser:       markdown.NewParser(logger),
		logger:       logger,
		fileWatcher:  fw,
		currentTheme: themes.ThemeDark, // Start with dark theme
	}

	// Apply custom theme
	app.Settings().SetTheme(themes.NewMarkViewTheme(w.currentTheme))

	w.setupUI()
	w.fyneWindow.Resize(fyne.NewSize(1200, 800))

	// Clean up file watcher on close
	w.fyneWindow.SetOnClosed(func() {
		if w.fileWatcher != nil {
			w.fileWatcher.Close()
		}
	})

	return w
}

// setupUI sets up the user interface
func (w *Window) setupUI() {
	// Create content display area (will be updated with parsed content)
	w.scrollContent = container.NewScroll(container.NewVBox())

	// Create file tree
	w.fileTree = NewFileTree(func(path string) {
		w.loadFile(path)
	})
	w.fileTreeScroll = container.NewScroll(w.fileTree.GetTree())

	// Create placeholder TOC
	w.tocTree = widget.NewTree(
		func(uid string) []string { return []string{} },
		func(uid string) bool { return false },
		func(branch bool) fyne.CanvasObject { return widget.NewLabel("") },
		func(uid string, branch bool, node fyne.CanvasObject) {},
	)
	w.tocScroll = container.NewScroll(w.tocTree)

	// Create three-pane layout: File Tree | TOC | Content
	// Left split: File Tree | TOC
	w.leftSplit = container.NewHSplit(
		w.fileTreeScroll,
		w.tocScroll,
	)
	w.leftSplit.Offset = 0.5 // Equal split between file tree and TOC

	// Main split: (File Tree | TOC) | Content
	w.mainSplit = container.NewHSplit(
		w.leftSplit,
		w.scrollContent,
	)
	w.mainSplit.Offset = 0.30 // Left panes take 30% of width

	// Create toolbar
	toolbar := w.createToolbar()

	// Create main layout
	mainContent := container.NewBorder(toolbar, nil, nil, nil, w.mainSplit)

	// Set up menu
	w.setupMenu()

	// Set up keyboard shortcuts
	w.setupShortcuts()

	w.fyneWindow.SetContent(mainContent)
}

// createToolbar creates the application toolbar with smaller buttons
func (w *Window) createToolbar() *fyne.Container {
	// Use icon buttons for a more compact toolbar
	openButton := widget.NewButtonWithIcon("Open", themes.IconDocument(), func() {
		w.showOpenDialog()
	})
	openButton.Importance = widget.LowImportance

	refreshButton := widget.NewButtonWithIcon("Refresh", themes.IconRefresh(), func() {
		if w.currentFile != "" {
			w.loadFile(w.currentFile)
		}
	})
	refreshButton.Importance = widget.LowImportance

	toolbar := container.NewHBox(
		openButton,
		refreshButton,
	)

	return toolbar
}

// setupMenu sets up the application menu
func (w *Window) setupMenu() {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open File...", func() {
			w.showOpenDialog()
		}),
		fyne.NewMenuItem("Open Folder...", func() {
			w.showFolderDialog()
		}),
		fyne.NewMenuItem("Refresh", func() {
			if w.currentFile != "" {
				w.loadFile(w.currentFile)
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Print...", func() {
			w.printDocument()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			w.fyneWindow.Close()
		}),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Toggle File Tree", func() {
			w.toggleFileTree()
		}),
		fyne.NewMenuItem("Toggle TOC", func() {
			w.toggleTOC()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Light Theme", func() {
			w.setTheme(themes.ThemeLight)
		}),
		fyne.NewMenuItem("Dark Theme", func() {
			w.setTheme(themes.ThemeDark)
		}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, viewMenu)
	w.fyneWindow.SetMainMenu(mainMenu)
}

// setupShortcuts sets up keyboard shortcuts
func (w *Window) setupShortcuts() {
	// Cmd/Ctrl+O - Open file
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.showOpenDialog()
	})

	// Cmd/Ctrl+P - Print
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.printDocument()
	})

	// Cmd/Ctrl+R - Refresh
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyR,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		if w.currentFile != "" {
			w.loadFile(w.currentFile)
		}
	})
}

// toggleFileTree toggles the file tree visibility
func (w *Window) toggleFileTree() {
	if w.fileTreeScroll.Visible() {
		w.fileTreeScroll.Hide()
	} else {
		w.fileTreeScroll.Show()
	}
}

// toggleTOC toggles the TOC visibility
func (w *Window) toggleTOC() {
	if w.tocScroll.Visible() {
		w.tocScroll.Hide()
	} else {
		w.tocScroll.Show()
	}
}

// showFolderDialog shows the folder selection dialog
func (w *Window) showFolderDialog() {
	fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if uri == nil {
			return
		}
		w.setRootFolder(uri.Path())
	}, w.fyneWindow)

	fd.Resize(fyne.NewSize(1050, 700))
	fd.Show()
}

// setRootFolder sets the root folder for the file tree
func (w *Window) setRootFolder(path string) {
	w.currentDir = path
	w.fileTree.SetRootPath(path)
	w.logger.Info("Set root folder", zap.String("path", path))
}

// printDocument prints the current document
func (w *Window) printDocument() {
	if w.currentFile == "" {
		dialog.ShowInformation("Print", "No document is currently open.", w.fyneWindow)
		return
	}

	// Show print dialog
	w.showPrintDialog()
}

// setTheme changes the application theme
func (w *Window) setTheme(themeType themes.ThemeType) {
	w.currentTheme = themeType
	w.app.Settings().SetTheme(themes.NewMarkViewTheme(themeType))

	// Reload current file to apply new theme to syntax highlighting
	if w.currentFile != "" {
		w.loadFile(w.currentFile)
	}

	w.logger.Info("Theme changed", zap.String("theme", w.getThemeName()))
}

// getThemeName returns the current theme name
func (w *Window) getThemeName() string {
	if w.currentTheme == themes.ThemeDark {
		return "dark"
	}
	return "light"
}

// showOpenDialog shows the file open dialog
func (w *Window) showOpenDialog() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		w.loadFile(reader.URI().Path())
	}, w.fyneWindow)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".md", ".markdown"}))

	// Make the file dialog 1.75x larger
	fd.Resize(fyne.NewSize(1050, 700))

	fd.Show()
}

// LoadFile loads and displays a markdown file
func (w *Window) LoadFile(filePath string) {
	w.loadFile(filePath)
}

// loadFile loads and displays a markdown file
func (w *Window) loadFile(filePath string) {
	w.logger.Info("Loading file", zap.String("path", filePath))
	w.currentFile = filePath

	// Update window title
	fileName := filepath.Base(filePath)
	w.fyneWindow.SetTitle(fmt.Sprintf("MarkView - %s", fileName))

	// Set file tree root to the file's directory if not already set
	fileDir := filepath.Dir(filePath)
	if w.currentDir == "" || w.currentDir != fileDir {
		w.currentDir = fileDir
		w.fileTree.SetRootPath(fileDir)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		w.logger.Error("Failed to read file", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to read file: %w", err), w.fyneWindow)
		return
	}

	// Parse markdown
	content, err := w.parser.Parse(data)
	if err != nil {
		w.logger.Error("Failed to parse markdown", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to parse markdown: %w", err), w.fyneWindow)
		return
	}

	// Update content
	w.scrollContent.Content = content
	w.scrollContent.Refresh()

	// Generate TOC
	w.updateTOC(data)

	w.logger.Info("File loaded successfully")
}

// updateTOC updates the table of contents
func (w *Window) updateTOC(data []byte) {
	// Parse markdown to get AST
	reader := text.NewReader(data)
	w.logger.Debug("Updating TOC")
	doc := w.parser.GetMarkdown().Parser().Parse(reader)

	// Generate TOC
	tocGen := toc.NewGenerator(data)
	entries := tocGen.Generate(doc)

	// Create TOC navigator
	navigator := toc.NewNavigator(entries, w.scrollContent)
	w.tocTree = navigator.GetTree()

	// Update TOC scroll container
	w.tocScroll.Content = w.tocTree
	w.tocScroll.Refresh()

	// Start watching the file for changes
	if w.fileWatcher != nil {
		err := w.fileWatcher.Watch(w.currentFile, func() {
			w.logger.Info("File changed, reloading", zap.String("path", w.currentFile))
			w.loadFile(w.currentFile)
		})
		if err != nil {
			w.logger.Error("Failed to watch file", zap.Error(err))
		}
	}
}

// Show shows the window
func (w *Window) Show() {
	w.fyneWindow.Show()
}

// ShowAndRun shows the window and runs the application
func (w *Window) ShowAndRun() {
	w.fyneWindow.ShowAndRun()
}

// GetMarkdown returns the markdown parser for external access
func (w *Window) GetMarkdown() *markdown.Parser {
	return w.parser
}

// showPrintDialog shows a print/export dialog
func (w *Window) showPrintDialog() {
	// Create print options dialog
	content := container.NewVBox(
		widget.NewLabel("Print Options"),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("File: %s", filepath.Base(w.currentFile))),
	)

	// Export to HTML button
	exportHTMLBtn := widget.NewButton("Export to HTML", func() {
		w.exportToHTML()
	})

	// Print button (opens system print dialog via HTML export)
	printBtn := widget.NewButton("Print (via Browser)", func() {
		w.printViaBrowser()
	})

	buttons := container.NewHBox(
		exportHTMLBtn,
		printBtn,
	)

	dialogContent := container.NewVBox(
		content,
		widget.NewSeparator(),
		buttons,
	)

	d := dialog.NewCustom("Print Document", "Cancel", dialogContent, w.fyneWindow)
	d.Resize(fyne.NewSize(400, 200))
	d.Show()
}

// exportToHTML exports the current markdown to HTML
func (w *Window) exportToHTML() {
	if w.currentFile == "" {
		return
	}

	// Read the markdown file
	data, err := os.ReadFile(w.currentFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read file: %w", err), w.fyneWindow)
		return
	}

	// Convert to HTML
	html := w.markdownToHTML(data)

	// Show save dialog
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(html))
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to write file: %w", err), w.fyneWindow)
			return
		}

		dialog.ShowInformation("Export Complete", "HTML file saved successfully.", w.fyneWindow)
	}, w.fyneWindow)

	// Suggest filename
	baseName := filepath.Base(w.currentFile)
	htmlName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".html"
	fd.SetFileName(htmlName)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".html"}))
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
}

// printViaBrowser exports to HTML and opens in browser for printing
func (w *Window) printViaBrowser() {
	if w.currentFile == "" {
		return
	}

	// Read the markdown file
	data, err := os.ReadFile(w.currentFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read file: %w", err), w.fyneWindow)
		return
	}

	// Convert to HTML with print-friendly styling
	html := w.markdownToHTML(data)

	// Write to temp file
	tmpFile, err := os.CreateTemp("", "markview-print-*.html")
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to create temp file: %w", err), w.fyneWindow)
		return
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(html)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to write temp file: %w", err), w.fyneWindow)
		return
	}

	// Open in default browser
	err = w.app.OpenURL(parseFileURL(tmpFile.Name()))
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to open browser: %w", err), w.fyneWindow)
		return
	}

	dialog.ShowInformation("Print", "Document opened in browser.\nUse your browser's print function (Cmd+P / Ctrl+P) to print.", w.fyneWindow)
}

// parseFileURL creates a file:// URL from a path
func parseFileURL(path string) *url.URL {
	u, _ := url.Parse("file://" + path)
	return u
}

// markdownToHTML converts markdown to a styled HTML document
func (w *Window) markdownToHTML(data []byte) string {
	// Use goldmark to convert markdown to HTML
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
		),
	)

	err := md.Convert(data, &buf)
	if err != nil {
		return fmt.Sprintf("<html><body><pre>Error: %s</pre></body></html>", err.Error())
	}

	// Wrap in styled HTML document
	title := filepath.Base(w.currentFile)

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            font-size: 16px;
            line-height: 1.6;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1, h2 { color: #2c5282; border-bottom: 1px solid #e2e8f0; padding-bottom: 0.3em; }
        h3, h4 { color: #c05621; }
        code {
            background-color: #f7fafc;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
            font-size: 85%%;
        }
        pre {
            background-color: #2d3748;
            color: #e2e8f0;
            padding: 16px;
            border-radius: 6px;
            overflow-x: auto;
        }
        pre code {
            background-color: transparent;
            padding: 0;
            color: inherit;
        }
        blockquote {
            border-left: 4px solid #4299e1;
            margin: 0;
            padding-left: 16px;
            color: #4a5568;
            font-style: italic;
        }
        table {
            border-collapse: collapse;
            width: 100%%;
        }
        th, td {
            border: 1px solid #e2e8f0;
            padding: 8px 12px;
            text-align: left;
        }
        th { background-color: #f7fafc; }
        a { color: #4299e1; }
        hr { border: none; border-top: 1px solid #e2e8f0; }
        @media print {
            body { max-width: none; }
            pre { white-space: pre-wrap; }
        }
    </style>
</head>
<body>
%s
</body>
</html>`, title, buf.String())
}
