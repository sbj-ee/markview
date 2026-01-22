package gui

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/sbj-ee/markview/internal/library"
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
	fyneWindow     fyne.Window
	app            fyne.App
	parser         *markdown.Parser
	logger         *zap.Logger
	fileTree       *FileTree
	fileTreeScroll *container.Scroll
	tocTree        *widget.Tree
	tocScroll      *container.Scroll
	scrollContent  *container.Scroll
	leftSplit      *container.Split // File Tree | TOC
	mainSplit      *container.Split // (File Tree | TOC) | Content
	currentFile    string
	currentDir     string
	fileWatcher    *watcher.FileWatcher
	currentTheme   themes.ThemeType
	currentFont    themes.FontFamily

	// Edit mode
	editMode      bool
	isDirty       bool
	editor        *MarkdownEditor
	editorScroll  *container.Scroll
	contentStack  *fyne.Container
	contentBuffer string // Original content for dirty checking

	// Toolbar actions
	editAction    *toolbarAction
	saveAction    *toolbarAction
	discardAction *toolbarAction

	// Edit toolbar (shown in edit mode)
	editToolbar *widget.Toolbar

	// Outline view (shown in edit mode)
	outline       *Outline
	outlineScroll *container.Scroll

	// Status bar
	statusBar   *widget.Label
	wordCount   *widget.Label
	cursorPos   *widget.Label

	// Library mode
	libraryMode   bool
	libraryView   *LibraryView
	docLibrary    *library.DocumentLibrary
	libraryScroll *container.Scroll
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
		currentTheme: themes.ThemeDark,   // Start with dark theme
		currentFont:  themes.FontDefault, // Start with default font
	}

	// Load saved theme preference
	savedTheme := app.Preferences().String("theme")
	if savedTheme != "" {
		w.currentTheme = themes.ThemeFromName(savedTheme)
	}

	// Load saved font preference
	savedFont := app.Preferences().String("font")
	if savedFont != "" {
		w.currentFont = themes.FontFamily(savedFont)
	}

	// Apply custom theme with font
	app.Settings().SetTheme(themes.NewMarkViewThemeWithFont(w.currentTheme, w.currentFont))

	w.setupUI()
	w.fyneWindow.Resize(fyne.NewSize(1200, 800))

	// Load last used directory
	w.loadLastDirectory()

	// Handle close with unsaved changes check
	w.fyneWindow.SetCloseIntercept(func() {
		w.handleClose()
	})

	// Clean up file watcher on close
	w.fyneWindow.SetOnClosed(func() {
		if w.fileWatcher != nil {
			w.fileWatcher.Close()
		}
	})

	// Set up drag and drop
	w.fyneWindow.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			path := uri.Path()
			if isMarkdownFile(path) {
				w.checkUnsavedChanges(func() {
					w.loadFile(path)
				})
				break // Only load the first markdown file
			}
		}
	})

	return w
}

// setupUI sets up the user interface
func (w *Window) setupUI() {
	// Create content display area (will be updated with parsed content)
	w.scrollContent = container.NewScroll(container.NewVBox())

	// Create editor (hidden by default)
	w.editor = NewMarkdownEditor(func(content string) {
		w.onEditorChanged(content)
	})
	w.editorScroll = container.NewScroll(w.editor)
	w.editorScroll.Hide()

	// Create outline view (hidden by default)
	w.outline = NewOutline(func(line int) {
		w.navigateToLine(line)
	})
	w.outlineScroll = container.NewScroll(w.outline.GetContainer())
	w.outlineScroll.Hide()

	// Create library view (hidden by default)
	w.libraryView = NewLibraryView(func(path string) {
		w.checkUnsavedChanges(func() {
			w.toggleLibraryMode() // Exit library mode
			w.loadFile(path)
		})
	})
	w.libraryScroll = container.NewScroll(w.libraryView.GetContainer())
	w.libraryScroll.Hide()

	// Create file tree with filter
	w.fileTree = NewFileTree(func(path string) {
		w.checkUnsavedChanges(func() {
			w.loadFile(path)
		})
	})
	w.fileTreeScroll = container.NewScroll(w.fileTree.GetContainer())

	// Create placeholder TOC
	w.tocTree = widget.NewTree(
		func(uid string) []string { return []string{} },
		func(uid string) bool { return false },
		func(branch bool) fyne.CanvasObject { return widget.NewLabel("") },
		func(uid string, branch bool, node fyne.CanvasObject) {},
	)
	w.tocScroll = container.NewScroll(w.tocTree)

	// Create three-pane layout: File Tree | TOC/Outline | Content
	// Left split: File Tree | TOC (or Outline in edit mode)
	// Use a stack for TOC and Outline - only one visible at a time
	tocOutlineStack := container.NewStack(w.tocScroll, w.outlineScroll)
	w.leftSplit = container.NewHSplit(
		w.fileTreeScroll,
		tocOutlineStack,
	)
	w.leftSplit.Offset = 0.5 // Equal split between file tree and TOC/Outline

	// Content area: stack of rendered view, editor, and library (only one visible at a time)
	w.contentStack = container.NewStack(w.scrollContent, w.editorScroll, w.libraryScroll)

	// Main split: (File Tree | TOC) | Content
	w.mainSplit = container.NewHSplit(
		w.leftSplit,
		w.contentStack,
	)
	w.mainSplit.Offset = 0.30 // Left panes take 30% of width

	// Create toolbar
	toolbar := w.createToolbar()

	// Create edit toolbar (hidden by default)
	w.editToolbar = w.createEditToolbar()
	w.editToolbar.Hide()

	// Create toolbar container with both toolbars
	toolbarContainer := container.NewVBox(toolbar, w.editToolbar)

	// Create status bar
	w.wordCount = widget.NewLabel("")
	w.cursorPos = widget.NewLabel("")
	w.statusBar = widget.NewLabel("")
	statusContainer := container.NewHBox(
		w.statusBar,
		widget.NewSeparator(),
		w.wordCount,
		widget.NewSeparator(),
		w.cursorPos,
	)

	// Create main layout
	mainContent := container.NewBorder(toolbarContainer, statusContainer, nil, nil, w.mainSplit)

	// Set up keyboard shortcuts
	w.setupShortcuts()

	w.fyneWindow.SetContent(mainContent)
}

// toolbarAction is a custom toolbar item
type toolbarAction struct {
	button *widget.Button
}

// newToolbarAction creates a toolbar action with an icon
func newToolbarAction(icon fyne.Resource, onTap func()) *toolbarAction {
	btn := widget.NewButtonWithIcon("", icon, onTap)
	btn.Importance = widget.LowImportance
	return &toolbarAction{button: btn}
}

// ToolbarObject implements widget.ToolbarItem
func (t *toolbarAction) ToolbarObject() fyne.CanvasObject {
	return container.NewPadded(t.button)
}

// SetIcon updates the icon of the toolbar action
func (t *toolbarAction) SetIcon(icon fyne.Resource) {
	t.button.SetIcon(icon)
}

// createToolbar creates the application toolbar
func (w *Window) createToolbar() *widget.Toolbar {
	newFileAction := newToolbarAction(themes.IconNewFile(), func() {
		w.newFile()
	})

	openFileAction := newToolbarAction(themes.IconDocument(), func() {
		w.showOpenDialog()
	})

	openFolderAction := newToolbarAction(themes.IconFolder(), func() {
		w.showFolderDialog()
	})

	w.saveAction = newToolbarAction(themes.IconSave(), func() {
		w.saveFile()
	})

	w.discardAction = newToolbarAction(themes.IconUndo(), func() {
		w.discardChanges()
	})

	w.editAction = newToolbarAction(themes.IconEdit(), func() {
		w.toggleEditMode()
	})

	refreshAction := newToolbarAction(themes.IconRefresh(), func() {
		if w.currentFile != "" {
			w.loadFile(w.currentFile)
		}
	})

	toggleFileTreeAction := newToolbarAction(themes.IconFileTree(), func() {
		w.toggleFileTree()
	})

	toggleTOCAction := newToolbarAction(themes.IconTOC(), func() {
		w.toggleTOC()
	})

	toggleThemeAction := newToolbarAction(themes.IconTheme(), func() {
		w.toggleTheme()
	})

	toggleLibraryAction := newToolbarAction(themes.IconLibrary(), func() {
		w.toggleLibraryMode()
	})

	toolbar := widget.NewToolbar(
		newFileAction,
		openFileAction,
		openFolderAction,
		widget.NewToolbarSeparator(),
		w.editAction,
		w.saveAction,
		w.discardAction,
		widget.NewToolbarSeparator(),
		refreshAction,
		widget.NewToolbarSpacer(),
		toggleLibraryAction,
		toggleFileTreeAction,
		toggleTOCAction,
		toggleThemeAction,
	)

	return toolbar
}

// toggleTheme shows settings dialog with theme and font selection
func (w *Window) toggleTheme() {
	themeNames := themes.ThemeNames()
	fontNames := themes.FontFamilyNames()

	// Create radio group for theme selection
	currentThemeName := w.currentTheme.Name()
	themeRadio := widget.NewRadioGroup(themeNames, func(selected string) {
		newTheme := themes.ThemeFromName(selected)
		w.setTheme(newTheme)
	})
	themeRadio.SetSelected(currentThemeName)

	// Create radio group for font selection
	var currentFontName string
	switch w.currentFont {
	case themes.FontMonospace:
		currentFontName = "Monospace"
	case themes.FontSerif:
		currentFontName = "Serif"
	case themes.FontSansSerif:
		currentFontName = "Sans Serif"
	default:
		currentFontName = "System Default"
	}
	fontRadio := widget.NewRadioGroup(fontNames, func(selected string) {
		newFont := themes.FontFamilyFromName(selected)
		w.setFont(newFont)
	})
	fontRadio.SetSelected(currentFontName)

	content := container.NewVBox(
		widget.NewLabel("Theme:"),
		themeRadio,
		widget.NewSeparator(),
		widget.NewLabel("Font:"),
		fontRadio,
	)

	d := dialog.NewCustom("Appearance Settings", "Close", content, w.fyneWindow)
	d.Resize(fyne.NewSize(300, 500))
	d.Show()
}

// createEditToolbar creates the markdown editing toolbar
func (w *Window) createEditToolbar() *widget.Toolbar {
	boldAction := newToolbarAction(themes.IconBold(), func() {
		w.editor.WrapSelection("**", "**")
	})

	italicAction := newToolbarAction(themes.IconItalic(), func() {
		w.editor.WrapSelection("*", "*")
	})

	h1Action := newToolbarAction(themes.IconHeading1(), func() {
		w.editor.InsertAtLineStart("# ")
	})

	h2Action := newToolbarAction(themes.IconHeading2(), func() {
		w.editor.InsertAtLineStart("## ")
	})

	h3Action := newToolbarAction(themes.IconHeading3(), func() {
		w.editor.InsertAtLineStart("### ")
	})

	linkAction := newToolbarAction(themes.IconLink(), func() {
		w.editor.WrapSelection("[", "](url)")
	})

	imageAction := newToolbarAction(themes.IconImage(), func() {
		// Get base directory for relative paths
		baseDir := ""
		if w.currentFile != "" {
			baseDir = filepath.Dir(w.currentFile)
		}
		ShowImageInsertDialog(w.fyneWindow, baseDir, func(markdown string) {
			w.editor.InsertAtCursor(markdown)
		})
	})

	codeAction := newToolbarAction(themes.IconCode(), func() {
		w.editor.WrapSelection("`", "`")
	})

	codeBlockAction := newToolbarAction(themes.IconCodeBlock(), func() {
		w.editor.InsertAtCursor("\n```\n\n```\n")
	})

	quoteAction := newToolbarAction(themes.IconQuote(), func() {
		w.editor.InsertAtLineStart("> ")
	})

	listAction := newToolbarAction(themes.IconList(), func() {
		w.editor.InsertAtLineStart("- ")
	})

	hrAction := newToolbarAction(themes.IconHorizontalRule(), func() {
		w.editor.InsertAtCursor("\n---\n")
	})

	tableAction := newToolbarAction(themes.IconTable(), func() {
		ShowTableEditorDialog(w.fyneWindow, func(markdown string) {
			w.editor.InsertAtCursor("\n" + markdown)
		})
	})

	return widget.NewToolbar(
		boldAction,
		italicAction,
		widget.NewToolbarSeparator(),
		h1Action,
		h2Action,
		h3Action,
		widget.NewToolbarSeparator(),
		linkAction,
		imageAction,
		tableAction,
		widget.NewToolbarSeparator(),
		codeAction,
		codeBlockAction,
		widget.NewToolbarSeparator(),
		quoteAction,
		listAction,
		hrAction,
	)
}

// setupShortcuts sets up keyboard shortcuts
func (w *Window) setupShortcuts() {
	// Cmd/Ctrl+N - New file
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.newFile()
	})

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

	// Cmd/Ctrl+F - Focus file filter
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.fileTree.FocusFilter(w.fyneWindow.Canvas())
	})

	// Cmd/Ctrl+E - Toggle edit mode
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.toggleEditMode()
	})

	// Cmd/Ctrl+S - Save file
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		if w.editMode {
			w.saveFile()
		}
	})

	// Cmd/Ctrl+Shift+S - Save As
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		if w.editMode {
			w.saveFileAs()
		}
	})

	// Escape - Clear filter or exit edit mode (if not dirty)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyEscape,
	}, func(shortcut fyne.Shortcut) {
		if w.editMode && !w.isDirty {
			w.switchToViewMode()
		} else {
			w.fileTree.ClearFilter()
		}
	})

	// Alt+Up - Navigate to parent directory
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyUp,
		Modifier: fyne.KeyModifierAlt,
	}, func(shortcut fyne.Shortcut) {
		w.fileTree.NavigateUp()
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

// toggleLibraryMode toggles between library mode and normal mode
func (w *Window) toggleLibraryMode() {
	if w.libraryMode {
		// Exit library mode
		w.libraryMode = false
		w.libraryScroll.Hide()
		w.scrollContent.Show()
		w.fyneWindow.SetTitle("MarkView")
		w.updateWindowTitle()
	} else {
		// Enter library mode
		if w.currentDir == "" {
			dialog.ShowInformation("Library Mode", "Please open a folder first to use library mode.", w.fyneWindow)
			return
		}

		w.libraryMode = true

		// Initialize or refresh library
		if w.docLibrary == nil || w.docLibrary.RootPath != w.currentDir {
			w.docLibrary = library.NewDocumentLibrary(w.currentDir)
		}
		w.docLibrary.Scan()
		w.libraryView.SetLibrary(w.docLibrary)

		// Update UI
		w.scrollContent.Hide()
		w.editorScroll.Hide()
		w.libraryScroll.Show()
		w.fyneWindow.SetTitle("MarkView - Library")
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
	// Save to preferences
	w.app.Preferences().SetString("lastDirectory", path)
	w.logger.Info("Set root folder", zap.String("path", path))
}

// loadLastDirectory loads the last used directory from preferences
func (w *Window) loadLastDirectory() {
	lastDir := w.app.Preferences().String("lastDirectory")
	if lastDir != "" {
		// Verify directory still exists
		if info, err := os.Stat(lastDir); err == nil && info.IsDir() {
			w.setRootFolder(lastDir)
		}
	}
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
	w.app.Settings().SetTheme(themes.NewMarkViewThemeWithFont(themeType, w.currentFont))

	// Save theme preference
	w.app.Preferences().SetString("theme", themeType.Name())

	// Reload current file to apply new theme to syntax highlighting
	if w.currentFile != "" {
		w.loadFile(w.currentFile)
	}

	w.logger.Info("Theme changed", zap.String("theme", w.getThemeName()))
}

// setFont changes the application font
func (w *Window) setFont(fontFamily themes.FontFamily) {
	w.currentFont = fontFamily
	w.app.Settings().SetTheme(themes.NewMarkViewThemeWithFont(w.currentTheme, fontFamily))

	// Save font preference
	w.app.Preferences().SetString("font", string(fontFamily))

	// Refresh content
	if w.currentFile != "" {
		w.loadFile(w.currentFile)
	}

	w.logger.Info("Font changed", zap.String("font", string(fontFamily)))
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
	w.checkUnsavedChanges(func() {
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
	})
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

	// Highlight current file in file tree
	w.fileTree.SetCurrentFile(filePath)

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		w.logger.Error("Failed to read file", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to read file: %w", err), w.fyneWindow)
		return
	}

	// Store content buffer for edit mode
	w.contentBuffer = string(data)
	w.isDirty = false

	// Parse markdown with base path for relative images
	content, err := w.parser.ParseWithBasePath(data, fileDir)
	if err != nil {
		w.logger.Error("Failed to parse markdown", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to parse markdown: %w", err), w.fyneWindow)
		return
	}

	// Update content
	w.scrollContent.Content = content
	w.scrollContent.Refresh()

	// Update editor content if in edit mode
	if w.editMode {
		w.editor.SetText(w.contentBuffer)
	}

	// Generate TOC
	w.updateTOC(data)

	// Update status bar
	w.updateStatusBar()

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

	// Start watching the file for changes (only if not in edit mode)
	if w.fileWatcher != nil && !w.editMode {
		err := w.fileWatcher.Watch(w.currentFile, func() {
			w.logger.Info("File changed, reloading", zap.String("path", w.currentFile))
			w.loadFile(w.currentFile)
		})
		if err != nil {
			w.logger.Error("Failed to watch file", zap.Error(err))
		}
	}
}

// toggleEditMode toggles between edit and view mode
func (w *Window) toggleEditMode() {
	if w.editMode {
		w.switchToViewMode()
	} else {
		w.switchToEditMode()
	}
}

// switchToEditMode switches to edit mode
func (w *Window) switchToEditMode() {
	if w.currentFile == "" {
		dialog.ShowInformation("Edit Mode", "No file is currently open.", w.fyneWindow)
		return
	}

	w.editMode = true
	w.logger.Info("Switching to edit mode")

	// Pause file watching
	if w.fileWatcher != nil {
		w.fileWatcher.Pause()
	}

	// Set editor content
	w.editor.SetText(w.contentBuffer)

	// Update outline
	w.outline.UpdateFromText(w.contentBuffer)

	// Update UI
	w.scrollContent.Hide()
	w.editorScroll.Show()
	w.editToolbar.Show()

	// Show outline in TOC area
	w.tocScroll.Hide()
	w.outlineScroll.Show()

	// Update toolbar icon to show "View" action
	w.editAction.SetIcon(themes.IconView())

	// Focus the editor
	w.editor.Focus(w.fyneWindow.Canvas())

	w.updateWindowTitle()
	w.updateStatusBar()
}

// switchToViewMode switches to view mode
func (w *Window) switchToViewMode() {
	w.editMode = false
	w.logger.Info("Switching to view mode")

	// Resume file watching
	if w.fileWatcher != nil {
		w.fileWatcher.Resume()
	}

	// Update content buffer from editor
	w.contentBuffer = w.editor.GetText()

	// Re-render the markdown
	fileDir := filepath.Dir(w.currentFile)
	content, err := w.parser.ParseWithBasePath([]byte(w.contentBuffer), fileDir)
	if err != nil {
		w.logger.Error("Failed to parse markdown", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to parse markdown: %w", err), w.fyneWindow)
		return
	}

	// Update content
	w.scrollContent.Content = content
	w.scrollContent.Refresh()

	// Update TOC
	w.updateTOCOnly([]byte(w.contentBuffer))

	// Update UI
	w.editorScroll.Hide()
	w.scrollContent.Show()
	w.editToolbar.Hide()

	// Hide outline, show TOC
	w.outlineScroll.Hide()
	w.tocScroll.Show()

	// Update toolbar icon to show "Edit" action
	w.editAction.SetIcon(themes.IconEdit())

	w.updateWindowTitle()
	w.updateStatusBar()
}

// updateTOCOnly updates the TOC without restarting the file watcher
func (w *Window) updateTOCOnly(data []byte) {
	reader := text.NewReader(data)
	w.logger.Debug("Updating TOC")
	doc := w.parser.GetMarkdown().Parser().Parse(reader)

	tocGen := toc.NewGenerator(data)
	entries := tocGen.Generate(doc)

	navigator := toc.NewNavigator(entries, w.scrollContent)
	w.tocTree = navigator.GetTree()

	w.tocScroll.Content = w.tocTree
	w.tocScroll.Refresh()
}

// onEditorChanged handles content changes in the editor
func (w *Window) onEditorChanged(content string) {
	// Skip dirty checking if we're just loading content
	if !w.editMode {
		return
	}

	// Mark as dirty on any change (don't compare full content - too slow)
	if !w.isDirty {
		w.isDirty = true
		w.updateWindowTitle()
	}

	// Update outline (debounced - only update occasionally)
	w.outline.UpdateFromText(content)

	// Update status bar
	w.updateStatusBar()
}

// navigateToLine navigates the editor to a specific line
func (w *Window) navigateToLine(line int) {
	if !w.editMode || w.editor == nil {
		return
	}

	// Get the text content
	text := w.editor.GetText()
	lines := strings.Split(text, "\n")

	// Calculate position at start of the target line
	pos := 0
	for i := 0; i < line && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}

	// Set cursor position
	w.editor.setCursorPosition(pos)

	// Focus the editor
	w.editor.Focus(w.fyneWindow.Canvas())
}

// updateStatusBar updates the status bar with current info
func (w *Window) updateStatusBar() {
	if w.editMode {
		content := w.editor.GetText()

		// Word count
		words := len(strings.Fields(content))
		chars := len(content)
		w.wordCount.SetText(fmt.Sprintf("%d words, %d chars", words, chars))

		// Cursor position (1-indexed for display)
		row, col := w.editor.GetCursorPosition()
		w.cursorPos.SetText(fmt.Sprintf("Ln %d, Col %d", row+1, col+1))

		w.statusBar.SetText("Edit Mode")
	} else if w.currentFile != "" {
		// Show file path in view mode
		w.statusBar.SetText(w.currentFile)
		w.wordCount.SetText("")
		w.cursorPos.SetText("")
	} else {
		w.statusBar.SetText("No file open")
		w.wordCount.SetText("")
		w.cursorPos.SetText("")
	}
}

// updateWindowTitle updates the window title with dirty indicator
func (w *Window) updateWindowTitle() {
	if w.currentFile == "" {
		w.fyneWindow.SetTitle("MarkView")
		return
	}

	fileName := filepath.Base(w.currentFile)
	modeIndicator := ""
	if w.editMode {
		modeIndicator = " [Edit]"
	}
	dirtyIndicator := ""
	if w.isDirty {
		dirtyIndicator = " *"
	}
	w.fyneWindow.SetTitle(fmt.Sprintf("MarkView - %s%s%s", fileName, modeIndicator, dirtyIndicator))
}

// saveFile saves the current editor content to the file
func (w *Window) saveFile() {
	if w.currentFile == "" {
		// No file yet, use Save As
		w.saveFileAs()
		return
	}

	content := w.editor.GetText()

	err := os.WriteFile(w.currentFile, []byte(content), 0644)
	if err != nil {
		w.logger.Error("Failed to save file", zap.Error(err))
		dialog.ShowError(fmt.Errorf("failed to save file: %w", err), w.fyneWindow)
		return
	}

	w.contentBuffer = content
	w.isDirty = false
	w.updateWindowTitle()

	w.logger.Info("File saved", zap.String("path", w.currentFile))
}

// saveFileAs saves the current content to a new file
func (w *Window) saveFileAs() {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		content := w.editor.GetText()
		_, err = writer.Write([]byte(content))
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to write file: %w", err), w.fyneWindow)
			return
		}

		// Update current file to the new path
		w.currentFile = writer.URI().Path()
		w.contentBuffer = content
		w.isDirty = false

		// Update file tree
		fileDir := filepath.Dir(w.currentFile)
		if w.currentDir != fileDir {
			w.currentDir = fileDir
			w.fileTree.SetRootPath(fileDir)
		}
		w.fileTree.SetCurrentFile(w.currentFile)

		w.updateWindowTitle()
		w.logger.Info("File saved as", zap.String("path", w.currentFile))
	}, w.fyneWindow)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".md", ".markdown"}))
	if w.currentFile != "" {
		fd.SetFileName(filepath.Base(w.currentFile))
	} else {
		fd.SetFileName("untitled.md")
	}
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
}

// newFile creates a new empty markdown file
func (w *Window) newFile() {
	w.checkUnsavedChanges(func() {
		// Clear current file
		w.currentFile = ""
		w.contentBuffer = "# New Document\n\n"
		w.isDirty = false

		// Clear rendered content
		w.scrollContent.Content = container.NewVBox()
		w.scrollContent.Refresh()

		// Clear TOC
		w.tocTree = widget.NewTree(
			func(uid string) []string { return []string{} },
			func(uid string) bool { return false },
			func(branch bool) fyne.CanvasObject { return widget.NewLabel("") },
			func(uid string, branch bool, node fyne.CanvasObject) {},
		)
		w.tocScroll.Content = w.tocTree
		w.tocScroll.Refresh()

		// Switch to edit mode with default content
		w.editMode = true
		w.editor.SetText(w.contentBuffer)
		w.scrollContent.Hide()
		w.editorScroll.Show()
		w.editToolbar.Show()
		w.editAction.SetIcon(themes.IconView())

		w.updateWindowTitle()
		w.editor.Focus(w.fyneWindow.Canvas())

		w.logger.Info("Created new file")
	})
}

// discardChanges discards unsaved changes and reverts to original content
func (w *Window) discardChanges() {
	if !w.isDirty {
		w.switchToViewMode()
		return
	}

	dialog.ShowConfirm("Discard Changes",
		"Are you sure you want to discard your changes?",
		func(discard bool) {
			if discard {
				w.editor.SetText(w.contentBuffer)
				w.isDirty = false
				w.switchToViewMode()
			}
		}, w.fyneWindow)
}

// handleClose handles window close with unsaved changes check
func (w *Window) handleClose() {
	if !w.isDirty {
		w.fyneWindow.Close()
		return
	}

	dialog.ShowConfirm("Unsaved Changes",
		"You have unsaved changes. Do you want to close without saving?",
		func(close bool) {
			if close {
				w.isDirty = false // Prevent re-prompt
				w.fyneWindow.Close()
			}
		}, w.fyneWindow)
}

// checkUnsavedChanges prompts the user about unsaved changes before a file switch
// Returns true if the caller should proceed, false if they should wait for user input
func (w *Window) checkUnsavedChanges(onProceed func()) {
	if !w.isDirty {
		onProceed()
		return
	}

	dialog.ShowConfirm("Unsaved Changes",
		"You have unsaved changes. Do you want to switch files without saving?",
		func(proceed bool) {
			if proceed {
				w.isDirty = false
				if w.editMode {
					w.switchToViewMode()
				}
				onProceed()
			}
		}, w.fyneWindow)
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
		widget.NewLabel("Export Options"),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("File: %s", filepath.Base(w.currentFile))),
	)

	// Export to HTML button
	exportHTMLBtn := widget.NewButton("Export to HTML", func() {
		w.exportToHTML()
	})

	// Export to PDF button
	exportPDFBtn := widget.NewButton("Export to PDF", func() {
		w.exportToPDF()
	})

	// Export to DOCX button (via pandoc)
	exportDOCXBtn := widget.NewButton("Export to DOCX", func() {
		w.exportWithPandoc("docx")
	})

	// Export to RTF button (via pandoc)
	exportRTFBtn := widget.NewButton("Export to RTF", func() {
		w.exportWithPandoc("rtf")
	})

	// Print button (opens system print dialog via HTML export)
	printBtn := widget.NewButton("Print (via Browser)", func() {
		w.printViaBrowser()
	})

	row1 := container.NewHBox(
		exportHTMLBtn,
		exportPDFBtn,
	)
	row2 := container.NewHBox(
		exportDOCXBtn,
		exportRTFBtn,
	)
	row3 := container.NewHBox(
		printBtn,
	)

	dialogContent := container.NewVBox(
		content,
		widget.NewSeparator(),
		row1,
		row2,
		row3,
	)

	d := dialog.NewCustom("Export Document", "Cancel", dialogContent, w.fyneWindow)
	d.Resize(fyne.NewSize(450, 280))
	d.Show()
}

// exportWithPandoc exports using pandoc to various formats (docx, rtf, etc.)
func (w *Window) exportWithPandoc(format string) {
	if w.currentFile == "" {
		return
	}

	// Check if pandoc is available
	_, err := exec.LookPath("pandoc")
	if err != nil {
		dialog.ShowError(fmt.Errorf("Export to %s requires pandoc.\n\nInstall with:\n  macOS: brew install pandoc\n  Linux: sudo apt install pandoc\n  Windows: choco install pandoc", strings.ToUpper(format)), w.fyneWindow)
		return
	}

	// Show save dialog
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if writer == nil {
			return
		}
		writer.Close() // Close immediately, we'll write via pandoc

		outputPath := writer.URI().Path()

		// Run pandoc
		cmd := exec.Command("pandoc", "-f", "markdown", "-t", format, "-o", outputPath, w.currentFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Pandoc conversion failed: %s\n%s", err, string(output)), w.fyneWindow)
			return
		}

		dialog.ShowInformation("Export Complete", fmt.Sprintf("%s file saved successfully.", strings.ToUpper(format)), w.fyneWindow)
	}, w.fyneWindow)

	// Suggest filename
	baseName := filepath.Base(w.currentFile)
	outputName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + "." + format
	fd.SetFileName(outputName)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{"." + format}))
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
}

// exportToPDF exports the current markdown to PDF
func (w *Window) exportToPDF() {
	if w.currentFile == "" {
		return
	}

	// Check if wkhtmltopdf is available
	_, err := exec.LookPath("wkhtmltopdf")
	if err != nil {
		dialog.ShowError(fmt.Errorf("PDF export requires wkhtmltopdf.\n\nInstall with:\n  macOS: brew install wkhtmltopdf\n  Linux: sudo apt install wkhtmltopdf"), w.fyneWindow)
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

	// Show save dialog for PDF
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if writer == nil {
			return
		}
		writer.Close() // Close immediately, we'll write via wkhtmltopdf

		pdfPath := writer.URI().Path()

		// Write HTML to temp file
		tmpFile, err := os.CreateTemp("", "markview-export-*.html")
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to create temp file: %w", err), w.fyneWindow)
			return
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(html)
		tmpFile.Close()
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to write temp file: %w", err), w.fyneWindow)
			return
		}

		// Run wkhtmltopdf
		cmd := exec.Command("wkhtmltopdf", "--quiet", tmpFile.Name(), pdfPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			dialog.ShowError(fmt.Errorf("PDF conversion failed: %s\n%s", err, string(output)), w.fyneWindow)
			return
		}

		dialog.ShowInformation("Export Complete", "PDF file saved successfully.", w.fyneWindow)
	}, w.fyneWindow)

	// Suggest filename
	baseName := filepath.Base(w.currentFile)
	pdfName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".pdf"
	fd.SetFileName(pdfName)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
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
        .mermaid { text-align: center; }
        @media print {
            body { max-width: none; }
            pre { white-space: pre-wrap; }
        }
    </style>
    <!-- Mermaid for diagram rendering -->
    <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
    <script>mermaid.initialize({startOnLoad:true});</script>
    <!-- MathJax for math rendering -->
    <script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
    <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
</head>
<body>
%s
</body>
</html>`, title, buf.String())
}
