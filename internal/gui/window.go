package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/sbj-ee/markview/internal/markdown"
	"github.com/sbj-ee/markview/internal/themes"
	"github.com/sbj-ee/markview/internal/toc"
	"github.com/sbj-ee/markview/internal/watcher"
	"github.com/yuin/goldmark/text"
	"go.uber.org/zap"
)

// Window represents the main application window
type Window struct {
	fyneWindow    fyne.Window
	app           fyne.App
	parser        *markdown.Parser
	logger        *zap.Logger
	content       *widget.RichText
	tocTree       *widget.Tree
	tocScroll     *container.Scroll
	scrollContent *container.Scroll
	split         *container.Split
	currentFile   string
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
	// Create content display
	w.content = widget.NewRichText()
	w.content.Wrapping = fyne.TextWrapWord
	w.scrollContent = container.NewScroll(w.content)

	// Create placeholder TOC
	w.tocTree = widget.NewTree(
		func(uid string) []string { return []string{} },
		func(uid string) bool { return false },
		func(branch bool) fyne.CanvasObject { return widget.NewLabel("") },
		func(uid string, branch bool, node fyne.CanvasObject) {},
	)
	w.tocScroll = container.NewScroll(w.tocTree)

	// Create split layout
	w.split = container.NewHSplit(
		w.tocScroll,
		w.scrollContent,
	)
	w.split.Offset = 0.25 // TOC takes 25% of width

	// Create toolbar
	toolbar := w.createToolbar()

	// Create main layout
	mainContent := container.NewBorder(toolbar, nil, nil, nil, w.split)

	// Set up menu
	w.setupMenu()

	w.fyneWindow.SetContent(mainContent)
}

// createToolbar creates the application toolbar
func (w *Window) createToolbar() *fyne.Container {
	openButton := widget.NewButton("Open File", func() {
		w.showOpenDialog()
	})

	refreshButton := widget.NewButton("Refresh", func() {
		if w.currentFile != "" {
			w.loadFile(w.currentFile)
		}
	})

	toolbar := container.NewHBox(
		openButton,
		refreshButton,
	)

	return toolbar
}

// setupMenu sets up the application menu
func (w *Window) setupMenu() {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", func() {
			w.showOpenDialog()
		}),
		fyne.NewMenuItem("Refresh", func() {
			if w.currentFile != "" {
				w.loadFile(w.currentFile)
			}
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			w.fyneWindow.Close()
		}),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Toggle TOC", func() {
			// Toggle TOC visibility
			if w.split.Leading.Visible() {
				w.split.Leading.Hide()
			} else {
				w.split.Leading.Show()
			}
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

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		w.logger.Error("Failed to read file", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to read file: %w", err), w.fyneWindow)
		return
	}

	// Parse markdown
	segments, err := w.parser.Parse(data)
	if err != nil {
		w.logger.Error("Failed to parse markdown", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to parse markdown: %w", err), w.fyneWindow)
		return
	}

	// Update content
	w.content.Segments = segments
	w.content.Refresh()

	// Generate TOC
	w.updateTOC(data)

	w.logger.Info("File loaded successfully")
}

// updateTOC updates the table of contents
func (w *Window) updateTOC(data []byte) {
	// Parse markdown to get AST
	reader := text.NewReader(data)
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
