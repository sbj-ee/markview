package gui

import (
	"bytes"
	"fmt"
	"image/color"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	fyneWindow      fyne.Window
	app             fyne.App
	parser          *markdown.Parser
	logger          *zap.Logger
	fileTree        *FileTree
	fileTreeScroll  *container.Scroll
	tocTree         *widget.Tree
	tocScroll       *container.Scroll
	scrollContent   *container.Scroll
	leftSplit       *container.Split // File Tree | TOC
	mainSplit       *container.Split // (File Tree | TOC) | Content
	currentFile     string
	currentDir      string
	fileWatcher     *watcher.FileWatcher
	currentTheme    themes.ThemeType
	currentFont     themes.FontFamily
	currentFontSize themes.FontSize

	// Edit mode
	editMode          bool
	splitViewMode     bool // Side-by-side editor and preview
	isDirty           bool
	editor            *MarkdownEditor
	editorScroll      *container.Scroll
	contentStack      *fyne.Container
	splitView         *container.Split  // For side-by-side editing
	splitEditor       *MarkdownEditor   // Editor for split view
	splitEditorScroll *container.Scroll // Scroll container for split editor
	splitPreview      *container.Scroll // Preview scroll for split view
	contentBuffer     string            // Original content for dirty checking
	focusMode         bool              // Hide all UI except content
	typewriterMode    bool              // Keep cursor centered
	autoSaveTicker    *time.Ticker      // Auto-save timer

	// Recent files
	recentFiles []string

	// Word count goal
	wordCountGoal int

	// Custom CSS for exports
	customCSS string

	// Toolbar actions
	editAction      *toolbarAction
	saveAction      *toolbarAction
	discardAction   *toolbarAction
	splitViewAction *toolbarAction

	// Edit toolbar (shown in edit mode)
	editToolbar *widget.Toolbar

	// Outline view (shown in edit mode)
	outline       *Outline
	outlineScroll *container.Scroll

	// Status bar
	statusBar *widget.Label
	wordCount *widget.Label
	cursorPos *widget.Label

	// Library mode
	libraryMode   bool
	libraryView   *LibraryView
	docLibrary    *library.DocumentLibrary
	libraryScroll *container.Scroll

	// Tags and features
	tagManager         *TagManager
	spellChecker       *SpellChecker
	zenMode            bool
	readingTime        *widget.Label
	currentExportTheme string

	// Link autocomplete
	linkAutocomplete      *LinkAutocomplete
	splitLinkAutocomplete *LinkAutocomplete

	// Version and updates
	version       string
	updateChecker *UpdateChecker
}

// NewWindow creates a new application window
func NewWindow(app fyne.App, logger *zap.Logger, version string) *Window {
	// Create file watcher with 300ms debounce
	fw, err := watcher.NewFileWatcher(logger, 300*time.Millisecond)
	if err != nil {
		logger.Error("Failed to create file watcher", zap.Error(err))
		fw = nil
	}

	w := &Window{
		fyneWindow:      app.NewWindow("MarkView"),
		app:             app,
		parser:          markdown.NewParser(logger),
		logger:          logger,
		fileWatcher:     fw,
		currentTheme:    themes.ThemeDark,      // Start with dark theme
		currentFont:     themes.FontDefault,    // Start with default font
		currentFontSize: themes.FontSizeNormal, // Start with normal font size
		version:         version,
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

	// Load saved font size preference
	savedFontSize := app.Preferences().String("fontSize")
	if savedFontSize != "" {
		w.currentFontSize = themes.FontSize(savedFontSize)
	}

	// Apply custom theme with all options
	app.Settings().SetTheme(themes.NewMarkViewThemeWithOptions(w.currentTheme, w.currentFont, w.currentFontSize))

	// Set window icon
	w.fyneWindow.SetIcon(themes.AppLogo())

	w.setupUI()
	w.fyneWindow.Resize(fyne.NewSize(1200, 800))

	// Load last used directory
	w.loadLastDirectory()

	// Load recent files
	w.loadRecentFiles()

	// Load word count goal
	w.wordCountGoal = w.app.Preferences().Int("wordCountGoal")

	// Load custom CSS
	w.customCSS = w.app.Preferences().String("customCSS")

	// Initialize tag manager and spell checker
	w.tagManager = NewTagManager()
	w.spellChecker = NewSpellChecker()

	// Load export theme preference
	w.currentExportTheme = w.app.Preferences().String("exportTheme")
	if w.currentExportTheme == "" {
		w.currentExportTheme = "Default"
	}

	// Load saved window size
	w.loadWindowSize()

	// Start auto-save
	w.startAutoSave()

	// Handle close with unsaved changes check
	w.fyneWindow.SetCloseIntercept(func() {
		w.handleClose()
	})

	// Clean up on close
	w.fyneWindow.SetOnClosed(func() {
		if w.fileWatcher != nil {
			w.fileWatcher.Close()
		}
		w.stopAutoSave()
		w.saveWindowSize()
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

	// Initialize link autocomplete
	w.linkAutocomplete = NewLinkAutocomplete(w.editor, w.fyneWindow, w.currentDir)

	// Initialize update checker
	w.updateChecker = NewUpdateChecker(w.version, w.fyneWindow, w.app)

	// Hook autocomplete key handler to editor
	w.editor.OnKeyEvent = func(key *fyne.KeyEvent) bool {
		return w.linkAutocomplete.HandleKeyEvent(key)
	}

	// Create split view with its own editor and preview containers
	w.splitEditor = NewMarkdownEditor(func(content string) {
		w.onSplitEditorChanged(content)
	})
	w.splitEditorScroll = container.NewScroll(w.splitEditor)
	w.splitPreview = container.NewScroll(container.NewVBox())
	w.splitView = container.NewHSplit(w.splitEditorScroll, w.splitPreview)
	w.splitView.Offset = 0.5
	w.splitView.Hide()

	// Initialize link autocomplete for split editor
	w.splitLinkAutocomplete = NewLinkAutocomplete(w.splitEditor, w.fyneWindow, w.currentDir)
	w.splitEditor.OnKeyEvent = func(key *fyne.KeyEvent) bool {
		return w.splitLinkAutocomplete.HandleKeyEvent(key)
	}

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
	// Set up starred documents persistence
	w.libraryView.SetOnStarredChanged(func(paths []string) {
		w.app.Preferences().SetString("starredDocs", strings.Join(paths, "\n"))
	})
	// Load starred documents
	starredStr := w.app.Preferences().String("starredDocs")
	if starredStr != "" {
		w.libraryView.SetStarredPaths(strings.Split(starredStr, "\n"))
	}
	w.libraryScroll = container.NewScroll(w.libraryView.GetContainer())
	w.libraryScroll.Hide()

	// Create file tree with filter
	w.fileTree = NewFileTree(func(path string) {
		w.checkUnsavedChanges(func() {
			w.loadFile(path)
		})
	})
	w.fileTree.SetWindow(w.fyneWindow)
	w.fileTree.SetOnFileDelete(func(deletedPath string) {
		// If the deleted file is the currently open file, clear the view
		if w.currentFile == deletedPath {
			// Exit edit/split mode first if active
			if w.splitViewMode {
				w.splitViewMode = false
				w.editMode = false
				w.splitView.Hide()
				w.scrollContent.Show()
				w.editorScroll.Hide()
				w.editToolbar.Hide()
				w.tocScroll.Show()
				w.outlineScroll.Hide()
				w.splitViewAction.SetIcon(themes.IconSplitView())
				w.editAction.SetIcon(themes.IconEdit())
				if w.fileWatcher != nil {
					w.fileWatcher.Resume()
				}
			} else if w.editMode {
				w.editMode = false
				w.editorScroll.Hide()
				w.scrollContent.Show()
				w.editToolbar.Hide()
				w.tocScroll.Show()
				w.outlineScroll.Hide()
				w.editAction.SetIcon(themes.IconEdit())
				if w.fileWatcher != nil {
					w.fileWatcher.Resume()
				}
			}
			// Clear the state
			w.currentFile = ""
			w.contentBuffer = ""
			w.isDirty = false
			w.editor.SetText("")
			w.splitEditor.SetText("")
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
			w.updateWindowTitle()
		}
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

	// Create toolbar (file operations - stays at top of entire window)
	toolbar := w.createToolbar()

	// Create edit toolbar (hidden by default, positioned above TOC/content only)
	w.editToolbar = w.createEditToolbar()
	w.editToolbar.Hide()

	// Create three-pane layout: File Tree | TOC/Outline | Content
	// Use a stack for TOC and Outline - only one visible at a time
	tocOutlineStack := container.NewStack(w.tocScroll, w.outlineScroll)

	// Content area: stack of rendered view, editor, split view, and library (only one visible at a time)
	w.contentStack = container.NewStack(w.scrollContent, w.editorScroll, w.splitView, w.libraryScroll)

	// TOC/Content split: TOC | Content
	tocContentSplit := container.NewHSplit(
		tocOutlineStack,
		w.contentStack,
	)
	tocContentSplit.Offset = 0.25 // TOC takes 25% of right pane width

	// Right pane: Edit toolbar at top (hidden by default), TOC/Content split below
	rightPane := container.NewBorder(w.editToolbar, nil, nil, nil, tocContentSplit)

	// Main split: File Tree | Right pane (edit toolbar + TOC + content)
	w.mainSplit = container.NewHSplit(
		w.fileTreeScroll,
		rightPane,
	)
	w.mainSplit.Offset = 0.20 // File tree takes 20% of width

	// Store reference to the left split for TOC toggling (now it's tocContentSplit)
	w.leftSplit = tocContentSplit

	// Create status bar
	w.wordCount = widget.NewLabel("")
	w.cursorPos = widget.NewLabel("")
	w.statusBar = widget.NewLabel("")
	w.readingTime = widget.NewLabel("")
	statusInfo := container.NewHBox(
		w.statusBar,
		widget.NewSeparator(),
		w.wordCount,
		widget.NewSeparator(),
		w.readingTime,
		widget.NewSeparator(),
		w.cursorPos,
	)

	// Create footer with logo on left, status in middle, version on right
	logoIcon := canvas.NewImageFromResource(themes.AppLogo())
	logoIcon.SetMinSize(fyne.NewSize(20, 20))
	logoIcon.FillMode = canvas.ImageFillContain

	versionLabel := widget.NewLabel("v" + w.version)
	versionLabel.TextStyle = fyne.TextStyle{Italic: true}

	footer := container.NewBorder(
		nil, nil,
		container.NewHBox(logoIcon, widget.NewLabel("MarkView")), // Left: logo + name
		versionLabel, // Right: version
		statusInfo,   // Center: status info
	)

	// Create main layout (main toolbar at top over entire window, footer at bottom)
	mainContent := container.NewBorder(toolbar, footer, nil, nil, w.mainSplit)

	// Set up keyboard shortcuts
	w.setupShortcuts()

	// Native menu bar gives every feature a discoverable home
	w.fyneWindow.SetMainMenu(w.createMainMenu())

	w.fyneWindow.SetContent(mainContent)
}

// toolbarAction is a custom toolbar item
type toolbarAction struct {
	button *tooltipButton
}

// newToolbarAction creates a toolbar action with an icon and a hover tooltip
func newToolbarAction(icon fyne.Resource, tooltip string, onTap func()) *toolbarAction {
	btn := newTooltipButton(icon, tooltip, onTap)
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
// menuItem builds a menu item with an optional Super-modifier keyboard
// accelerator shown alongside the label.
func menuItem(label string, key fyne.KeyName, mod fyne.KeyModifier, action func()) *fyne.MenuItem {
	item := fyne.NewMenuItem(label, action)
	if key != "" {
		item.Shortcut = &desktop.CustomShortcut{KeyName: key, Modifier: mod}
	}
	return item
}

// createMainMenu builds the native application menu bar so that every feature
// has a discoverable home with its keyboard shortcut displayed.
func (w *Window) createMainMenu() *fyne.MainMenu {
	const super = fyne.KeyModifierSuper
	const superShift = fyne.KeyModifierSuper | fyne.KeyModifierShift

	fileMenu := fyne.NewMenu("File",
		menuItem("New File", fyne.KeyN, super, w.newFile),
		menuItem("New from Template…", fyne.KeyN, superShift, w.newFileFromTemplate),
		menuItem("Open File…", fyne.KeyO, super, w.showOpenDialog),
		menuItem("Open Folder…", "", 0, w.showFolderDialog),
		menuItem("Open Recent…", fyne.KeyO, superShift, w.showRecentFiles),
		fyne.NewMenuItemSeparator(),
		menuItem("Save", fyne.KeyS, super, func() {
			if w.editMode {
				w.saveFile()
			}
		}),
		menuItem("Save As…", fyne.KeyS, superShift, func() {
			if w.editMode {
				w.saveFileAs()
			}
		}),
		fyne.NewMenuItemSeparator(),
		menuItem("Print", fyne.KeyP, super, w.printDocument),
		menuItem("Export…", "", 0, w.showPrintDialog),
	)

	editMenu := fyne.NewMenu("Edit",
		menuItem("Toggle Edit Mode", fyne.KeyE, super, w.toggleEditMode),
		menuItem("Find & Replace…", fyne.KeyF, super, func() {
			// Matches the Cmd+F canvas handler: find in edit mode,
			// otherwise focus the file browser filter.
			if w.editMode {
				ShowFindReplaceDialog(w.fyneWindow, w.editor)
			} else {
				w.fileTree.FocusFilter(w.fyneWindow.Canvas())
			}
		}),
	)

	viewMenu := fyne.NewMenu("View",
		menuItem("Split View", fyne.KeyBackslash, super, w.toggleSplitView),
		menuItem("Reload", fyne.KeyR, super, func() {
			if w.currentFile != "" {
				w.loadFile(w.currentFile)
			}
		}),
		fyne.NewMenuItemSeparator(),
		menuItem("Toggle File Browser", "", 0, w.toggleFileTree),
		menuItem("Toggle Table of Contents", "", 0, w.toggleTOC),
		menuItem("Toggle Library", "", 0, w.toggleLibraryMode),
		fyne.NewMenuItemSeparator(),
		menuItem("Focus Mode", fyne.KeyF, superShift, w.toggleFocusMode),
		menuItem("Zen Mode", fyne.KeyF11, 0, w.toggleZenMode),
		menuItem("Presentation Mode", "", 0, w.showPresentationMode),
		fyne.NewMenuItemSeparator(),
		menuItem("Appearance Settings…", "", 0, w.toggleTheme),
	)

	goMenu := fyne.NewMenu("Go",
		menuItem("Command Palette", fyne.KeyP, superShift, w.showCommandPalette),
		menuItem("Quick Switcher", fyne.KeyP, fyne.KeyModifierControl, w.showQuickSwitcher),
		menuItem("Search in Files", fyne.KeyG, superShift, w.showFullTextSearch),
		menuItem("Backlinks", fyne.KeyB, super, w.showBacklinks),
		menuItem("Browse Tags", fyne.KeyT, super, w.showTagsBrowser),
	)

	toolsMenu := fyne.NewMenu("Tools",
		menuItem("Validate Links", fyne.KeyL, super, w.validateLinks),
		menuItem("Word Count Goal…", "", 0, w.showWordCountGoalDialog),
		menuItem("Export Theme…", "", 0, w.showExportThemeDialog),
		menuItem("Custom CSS…", "", 0, w.showCustomCSSDialog),
	)

	helpMenu := fyne.NewMenu("Help",
		menuItem("Keyboard Shortcuts", fyne.KeySlash, super, w.showKeyboardShortcuts),
		menuItem("Help", "", 0, w.showHelpMenu),
		menuItem("Check for Updates…", "", 0, func() {
			if w.updateChecker != nil {
				w.updateChecker.CheckForUpdates(false)
			}
		}),
		menuItem("About MarkView", "", 0, w.showAboutDialog),
	)

	return fyne.NewMainMenu(fileMenu, editMenu, viewMenu, goMenu, toolsMenu, helpMenu)
}

func (w *Window) createToolbar() *widget.Toolbar {
	newFileAction := newToolbarAction(themes.IconNewFile(), "New File ("+cmdKey+"+N)", func() {
		w.newFile()
	})

	openFileAction := newToolbarAction(themes.IconDocument(), "Open File ("+cmdKey+"+O)", func() {
		w.showOpenDialog()
	})

	openFolderAction := newToolbarAction(themes.IconFolder(), "Open Folder", func() {
		w.showFolderDialog()
	})

	w.saveAction = newToolbarAction(themes.IconSave(), "Save ("+cmdKey+"+S)", func() {
		w.saveFile()
	})

	w.discardAction = newToolbarAction(themes.IconUndo(), "Discard Changes", func() {
		w.discardChanges()
	})

	w.editAction = newToolbarAction(themes.IconEdit(), "Toggle Edit Mode ("+cmdKey+"+E)", func() {
		w.toggleEditMode()
	})

	w.splitViewAction = newToolbarAction(themes.IconSplitView(), "Toggle Split View ("+cmdKey+"+\\)", func() {
		w.toggleSplitView()
	})

	refreshAction := newToolbarAction(themes.IconRefresh(), "Reload ("+cmdKey+"+R)", func() {
		if w.currentFile != "" {
			w.loadFile(w.currentFile)
		}
	})

	toggleFileTreeAction := newToolbarAction(themes.IconFileTree(), "Toggle File Browser", func() {
		w.toggleFileTree()
	})

	toggleTOCAction := newToolbarAction(themes.IconTOC(), "Toggle Table of Contents", func() {
		w.toggleTOC()
	})

	toggleThemeAction := newToolbarAction(themes.IconTheme(), "Appearance Settings", func() {
		w.toggleTheme()
	})

	toggleLibraryAction := newToolbarAction(themes.IconLibrary(), "Toggle Library", func() {
		w.toggleLibraryMode()
	})

	presentationAction := newToolbarAction(themes.IconPresentation(), "Presentation Mode", func() {
		w.showPresentationMode()
	})

	printAction := newToolbarAction(themes.IconPrint(), "Print ("+cmdKey+"+P)", func() {
		w.printDocument()
	})

	exportAction := newToolbarAction(themes.IconExport(), "Export…", func() {
		w.showPrintDialog()
	})

	helpAction := newToolbarAction(themes.IconHelp(), "Help", func() {
		w.showHelpMenu()
	})

	toolbar := widget.NewToolbar(
		newFileAction,
		openFileAction,
		openFolderAction,
		widget.NewToolbarSeparator(),
		w.editAction,
		w.splitViewAction,
		w.saveAction,
		w.discardAction,
		widget.NewToolbarSeparator(),
		refreshAction,
		presentationAction,
		printAction,
		exportAction,
		widget.NewToolbarSpacer(),
		toggleLibraryAction,
		toggleFileTreeAction,
		toggleTOCAction,
		toggleThemeAction,
		helpAction,
	)

	return toolbar
}

// toggleTheme shows settings dialog with theme and font selection
func (w *Window) toggleTheme() {
	themeNames := themes.ThemeNames()
	fontNames := themes.FontFamilyNames()
	fontSizeNames := themes.FontSizeNames()

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

	// Create radio group for font size selection
	currentFontSizeName := w.currentFontSize.Name()
	fontSizeRadio := widget.NewRadioGroup(fontSizeNames, func(selected string) {
		newFontSize := themes.FontSizeFromName(selected)
		w.setFontSize(newFontSize)
	})
	fontSizeRadio.SetSelected(currentFontSizeName)

	content := container.NewVBox(
		widget.NewLabel("Theme:"),
		themeRadio,
		widget.NewSeparator(),
		widget.NewLabel("Font:"),
		fontRadio,
		widget.NewSeparator(),
		widget.NewLabel("Font Size:"),
		fontSizeRadio,
	)

	d := dialog.NewCustom("Appearance Settings", "Close", content, w.fyneWindow)
	d.Resize(fyne.NewSize(300, 600))
	d.Show()
}

// createEditToolbar creates the markdown editing toolbar
func (w *Window) createEditToolbar() *widget.Toolbar {
	// Text formatting group
	boldAction := newToolbarAction(themes.IconBold(), "Bold", func() {
		w.getActiveEditor().WrapSelection("**", "**")
	})

	italicAction := newToolbarAction(themes.IconItalic(), "Italic", func() {
		w.getActiveEditor().WrapSelection("*", "*")
	})

	strikethroughAction := newToolbarAction(themes.IconStrikethrough(), "Strikethrough", func() {
		w.getActiveEditor().WrapSelection("~~", "~~")
	})

	underlineAction := newToolbarAction(themes.IconUnderline(), "Underline", func() {
		w.getActiveEditor().WrapSelection("<u>", "</u>")
	})

	highlightAction := newToolbarAction(themes.IconHighlight(), "Highlight", func() {
		w.getActiveEditor().WrapSelection("==", "==")
	})

	// Math/Science group - wrap selected text with HTML tags (renders as Unicode sub/superscript)
	subscriptAction := newToolbarAction(themes.IconSubscript(), "Subscript", func() {
		w.getActiveEditor().WrapSelection("<sub>", "</sub>")
	})

	superscriptAction := newToolbarAction(themes.IconSuperscript(), "Superscript", func() {
		w.getActiveEditor().WrapSelection("<sup>", "</sup>")
	})

	symbolAction := newToolbarAction(themes.IconSymbol(), "Insert Symbol…", func() {
		// Sync cursor position before dialog opens
		w.getActiveEditor().SyncLastInsertPos()
		ShowSymbolPickerDialog(w.fyneWindow, func(symbol string) {
			editor := w.getActiveEditor()
			editor.InsertAtCursor(symbol)
		})
	})

	// Headings group
	h1Action := newToolbarAction(themes.IconHeading1(), "Heading 1", func() {
		w.getActiveEditor().InsertAtLineStart("# ")
	})

	h2Action := newToolbarAction(themes.IconHeading2(), "Heading 2", func() {
		w.getActiveEditor().InsertAtLineStart("## ")
	})

	h3Action := newToolbarAction(themes.IconHeading3(), "Heading 3", func() {
		w.getActiveEditor().InsertAtLineStart("### ")
	})

	// Links and media group
	linkAction := newToolbarAction(themes.IconLink(), "Insert Link", func() {
		w.getActiveEditor().WrapSelection("[", "](url)")
	})

	imageAction := newToolbarAction(themes.IconImage(), "Insert Image…", func() {
		// Get base directory for relative paths
		baseDir := ""
		if w.currentFile != "" {
			baseDir = filepath.Dir(w.currentFile)
		}
		ShowImageInsertDialog(w.fyneWindow, baseDir, func(markdown string) {
			w.getActiveEditor().InsertAtCursor(markdown)
		})
	})

	tableAction := newToolbarAction(themes.IconTable(), "Insert Table…", func() {
		ShowTableEditorDialog(w.fyneWindow, func(markdown string) {
			w.getActiveEditor().InsertAtCursor("\n" + markdown)
		})
	})

	footnoteAction := newToolbarAction(themes.IconFootnote(), "Insert Footnote", func() {
		w.getActiveEditor().InsertAtCursor("[^1]")
	})

	// Code group
	codeAction := newToolbarAction(themes.IconCode(), "Inline Code", func() {
		w.getActiveEditor().WrapSelection("`", "`")
	})

	codeBlockAction := newToolbarAction(themes.IconCodeBlock(), "Code Block", func() {
		w.getActiveEditor().InsertAtCursor("\n```\n\n```\n")
	})

	// Lists and structure group
	quoteAction := newToolbarAction(themes.IconQuote(), "Blockquote", func() {
		w.getActiveEditor().InsertAtLineStart("> ")
	})

	listAction := newToolbarAction(themes.IconList(), "Bullet List", func() {
		w.getActiveEditor().InsertAtLineStart("- ")
	})

	numberedListAction := newToolbarAction(themes.IconNumberedList(), "Numbered List", func() {
		w.getActiveEditor().InsertAtLineStart("1. ")
	})

	checkboxAction := newToolbarAction(themes.IconCheckbox(), "Checkbox", func() {
		w.getActiveEditor().InsertAtLineStart("- [ ] ")
	})

	hrAction := newToolbarAction(themes.IconHorizontalRule(), "Horizontal Rule", func() {
		w.getActiveEditor().InsertAtCursor("\n---\n")
	})

	// Productivity group
	snippetAction := newToolbarAction(themes.IconSnippet(), "Snippets…", func() {
		ShowSnippetsDialog(w.fyneWindow, func(content string) {
			w.getActiveEditor().InsertAtCursor(content)
		})
	})

	typewriterAction := newToolbarAction(themes.IconTypewriter(), "Typewriter Mode", func() {
		w.toggleTypewriterMode()
	})

	goalAction := newToolbarAction(themes.IconGoal(), "Word Count Goal…", func() {
		w.showWordCountGoalDialog()
	})

	return widget.NewToolbar(
		// Text formatting
		boldAction,
		italicAction,
		strikethroughAction,
		underlineAction,
		highlightAction,
		widget.NewToolbarSeparator(),
		// Math/Science
		subscriptAction,
		superscriptAction,
		symbolAction,
		widget.NewToolbarSeparator(),
		// Headings
		h1Action,
		h2Action,
		h3Action,
		widget.NewToolbarSeparator(),
		// Links and media
		linkAction,
		imageAction,
		tableAction,
		footnoteAction,
		widget.NewToolbarSeparator(),
		// Code
		codeAction,
		codeBlockAction,
		widget.NewToolbarSeparator(),
		// Lists and structure
		quoteAction,
		listAction,
		numberedListAction,
		checkboxAction,
		hrAction,
		widget.NewToolbarSeparator(),
		// Productivity
		snippetAction,
		typewriterAction,
		goalAction,
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

	// Cmd/Ctrl+F - Find (in edit mode) or Focus file filter (in view mode)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		if w.editMode {
			ShowFindReplaceDialog(w.fyneWindow, w.editor)
		} else {
			w.fileTree.FocusFilter(w.fyneWindow.Canvas())
		}
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

	// Escape - Dismiss autocomplete, clear filter, or exit edit mode (if not dirty)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyEscape,
	}, func(shortcut fyne.Shortcut) {
		// First, check if autocomplete is active
		if w.linkAutocomplete != nil && w.linkAutocomplete.IsActive() {
			w.linkAutocomplete.Dismiss()
			return
		}
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

	// Cmd/Ctrl+\ - Toggle split view
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyBackslash,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.toggleSplitView()
	})

	// Cmd/Ctrl+Shift+F - Toggle focus mode
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.toggleFocusMode()
	})

	// Cmd/Ctrl+? or Cmd+/ - Show keyboard shortcuts
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeySlash,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.showKeyboardShortcuts()
	})

	// Cmd/Ctrl+Shift+O - Show recent files
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.showRecentFiles()
	})

	// Ctrl+P - Quick switcher (not Cmd+P which is print)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		w.showQuickSwitcher()
	})

	// Ctrl+Shift+P - Quick switcher (alternative)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.showQuickSwitcher()
	})

	// Ctrl+Shift+F - Full-text search across files
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyG,
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.showFullTextSearch()
	})

	// Cmd/Ctrl+B - Show backlinks
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyB,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.showBacklinks()
	})

	// Cmd/Ctrl+L - Validate links
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyL,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.validateLinks()
	})

	// Cmd/Ctrl+T - Browse by tags
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyT,
		Modifier: fyne.KeyModifierSuper,
	}, func(shortcut fyne.Shortcut) {
		w.showTagsBrowser()
	})

	// Cmd/Ctrl+Shift+N - New from template
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.newFileFromTemplate()
	})

	// F11 or Cmd+Shift+Enter - Zen mode (fullscreen)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyF11,
	}, func(shortcut fyne.Shortcut) {
		w.toggleZenMode()
	})
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.toggleZenMode()
	})

	// Link autocomplete shortcuts
	// Ctrl+N - Select next autocomplete suggestion
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		if w.linkAutocomplete != nil && w.linkAutocomplete.IsActive() {
			w.linkAutocomplete.SelectNext()
		}
	})

	// Ctrl+P is already used for quick switcher, so use Ctrl+J/K for navigation
	// Ctrl+J - Select next autocomplete suggestion (vim-style)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyJ,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		if w.linkAutocomplete != nil && w.linkAutocomplete.IsActive() {
			w.linkAutocomplete.SelectNext()
		}
	})

	// Ctrl+K - Select previous autocomplete suggestion (vim-style)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyK,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		if w.linkAutocomplete != nil && w.linkAutocomplete.IsActive() {
			w.linkAutocomplete.SelectPrevious()
		}
	})

	// Tab - Accept autocomplete selection (when active)
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyTab,
	}, func(shortcut fyne.Shortcut) {
		if w.linkAutocomplete != nil && w.linkAutocomplete.IsActive() {
			w.linkAutocomplete.AcceptSelection()
		}
		// If autocomplete is not active, Tab is handled by the entry for indentation
	})

	// Ctrl+Space - Accept autocomplete selection or trigger autocomplete
	w.fyneWindow.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeySpace,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		if w.linkAutocomplete != nil && w.linkAutocomplete.IsActive() {
			w.linkAutocomplete.AcceptSelection()
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

// getActiveEditor returns the currently active editor based on view mode
func (w *Window) getActiveEditor() *MarkdownEditor {
	if w.splitViewMode {
		return w.splitEditor
	}
	return w.editor
}

// toggleSplitView toggles split view mode (side-by-side editor and preview)
func (w *Window) toggleSplitView() {
	if w.currentFile == "" {
		dialog.ShowInformation("Split View", "No file is currently open.", w.fyneWindow)
		return
	}

	if w.splitViewMode {
		// Exit split view, back to normal view mode
		w.splitViewMode = false
		w.editMode = false

		// Copy content from split editor back to main editor and buffer
		w.contentBuffer = w.splitEditor.GetText()
		w.editor.SetText(w.contentBuffer)

		w.splitView.Hide()
		w.splitViewAction.SetIcon(themes.IconSplitView())
		w.scrollContent.Show()
		w.editorScroll.Hide()
		w.editToolbar.Hide()

		// Update the main preview with the edited content
		w.updateMainPreview(w.contentBuffer)

		// Show TOC, hide outline
		w.tocScroll.Show()
		w.outlineScroll.Hide()

		w.editAction.SetIcon(themes.IconEdit())

		// Resume file watching
		if w.fileWatcher != nil {
			w.fileWatcher.Resume()
		}
	} else {
		// Enter split view mode
		w.splitViewMode = true
		w.editMode = true

		// Pause file watching
		if w.fileWatcher != nil {
			w.fileWatcher.Pause()
		}

		// Set split editor content
		w.splitEditor.SetText(w.contentBuffer)

		// Update outline
		w.outline.UpdateFromText(w.contentBuffer)

		// Update split preview with current content
		w.updateSplitViewPreview(w.contentBuffer)

		// Hide single views, show split
		w.scrollContent.Hide()
		w.editorScroll.Hide()
		w.splitView.Show()
		w.editToolbar.Show()
		w.splitViewAction.SetIcon(themes.IconSingleView())

		// Show outline in TOC area
		w.tocScroll.Hide()
		w.outlineScroll.Show()

		w.editAction.SetIcon(themes.IconView())

		// Focus the split editor
		w.splitEditor.Focus(w.fyneWindow.Canvas())
	}
	w.updateWindowTitle()
}

// toggleTypewriterMode toggles typewriter mode (keep cursor centered)
func (w *Window) toggleTypewriterMode() {
	w.typewriterMode = !w.typewriterMode
	if w.typewriterMode {
		// Immediately center the cursor
		w.centerCursor()
	}
}

// centerCursor scrolls the editor to keep the cursor line centered
func (w *Window) centerCursor() {
	if !w.editMode || w.editor == nil || w.editorScroll == nil {
		return
	}

	// Get cursor row
	row, _ := w.editor.GetCursorPosition()

	// Estimate line height (approximately 20 pixels per line)
	lineHeight := float32(20)

	// Calculate scroll position to center the cursor row
	scrollHeight := w.editorScroll.Size().Height
	targetY := float32(row)*lineHeight - scrollHeight/2 + lineHeight/2

	if targetY < 0 {
		targetY = 0
	}

	w.editorScroll.Offset = fyne.NewPos(0, targetY)
	w.editorScroll.Refresh()
}

// toggleFocusMode toggles focus mode (hide all UI except content)
func (w *Window) toggleFocusMode() {
	w.focusMode = !w.focusMode

	if w.focusMode {
		// Hide sidebars and toolbars
		w.fileTreeScroll.Hide()
		w.tocScroll.Hide()
		w.outlineScroll.Hide()
		w.editToolbar.Hide()
		w.mainSplit.Offset = 0 // Hide left panel entirely
	} else {
		// Restore UI
		w.fileTreeScroll.Show()
		if w.editMode {
			w.outlineScroll.Show()
			w.editToolbar.Show()
		} else {
			w.tocScroll.Show()
		}
		w.mainSplit.Offset = 0.30
	}
	w.mainSplit.Refresh()
}

// showKeyboardShortcuts shows a dialog with keyboard shortcuts
func (w *Window) showKeyboardShortcuts() {
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"Cmd+N", "New file"},
		{"Cmd+Shift+N", "New from template"},
		{"Cmd+O", "Open file"},
		{"Cmd+S", "Save file"},
		{"Cmd+Shift+S", "Save as"},
		{"Cmd+E", "Toggle edit mode"},
		{"Cmd+F", "Find/Replace (edit) / Filter (view)"},
		{"Cmd+P", "Print/Export"},
		{"Ctrl+P", "Quick file switcher"},
		{"Cmd+Shift+G", "Search in all files"},
		{"Cmd+B", "Show backlinks"},
		{"Cmd+L", "Validate links"},
		{"Cmd+T", "Browse by tags"},
		{"Cmd+R", "Refresh"},
		{"Cmd+\\", "Toggle split view"},
		{"Cmd+Shift+F", "Toggle focus mode"},
		{"F11 / Cmd+Shift+Enter", "Zen mode (fullscreen)"},
		{"Cmd+?", "Show shortcuts"},
		{"Escape", "Exit edit mode / Clear filter / Dismiss autocomplete"},
		{"Alt+Up", "Navigate to parent directory"},
		{"", ""},
		{"Link Autocomplete (in edit mode):", ""},
		{"Tab / Ctrl+Space", "Accept autocomplete suggestion"},
		{"Ctrl+J / Ctrl+N", "Select next suggestion"},
		{"Ctrl+K", "Select previous suggestion"},
	}

	var items []fyne.CanvasObject
	for _, s := range shortcuts {
		row := container.NewHBox(
			widget.NewLabelWithStyle(s.key, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
			widget.NewLabel("-"),
			widget.NewLabel(s.desc),
		)
		items = append(items, row)
	}

	content := container.NewVBox(items...)
	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(350, 400))

	d := dialog.NewCustom("Keyboard Shortcuts", "Close", scroll, w.fyneWindow)
	d.Resize(fyne.NewSize(400, 450))
	d.Show()
}

// showHelpMenu shows the help menu with various options
func (w *Window) showHelpMenu() {
	// Create menu items
	shortcutsBtn := widget.NewButton("Keyboard Shortcuts", func() {
		w.showKeyboardShortcuts()
	})
	shortcutsBtn.Importance = widget.HighImportance

	checkUpdatesBtn := widget.NewButton("Check for Updates...", func() {
		if w.updateChecker != nil {
			w.updateChecker.CheckForUpdates(false)
		}
	})

	aboutBtn := widget.NewButton("About MarkView", func() {
		w.showAboutDialog()
	})

	content := container.NewVBox(
		shortcutsBtn,
		checkUpdatesBtn,
		widget.NewSeparator(),
		aboutBtn,
	)

	popup := widget.NewPopUp(content, w.fyneWindow.Canvas())
	popup.ShowAtPosition(fyne.NewPos(
		w.fyneWindow.Canvas().Size().Width-200,
		50,
	))
}

// showAboutDialog shows information about the application
func (w *Window) showAboutDialog() {
	title := widget.NewLabelWithStyle("MarkView", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	version := widget.NewLabel(fmt.Sprintf("Version %s", w.version))
	version.Alignment = fyne.TextAlignCenter

	desc := widget.NewLabel("A powerful markdown viewer and editor")
	desc.Alignment = fyne.TextAlignCenter
	desc.Wrapping = fyne.TextWrapWord

	copyright := widget.NewLabel("© 2024-2025 sbj-ee • MIT License")
	copyright.Alignment = fyne.TextAlignCenter

	githubLink := widget.NewHyperlink("GitHub Repository", nil)
	if u, err := url.Parse("https://github.com/sbj-ee/markview"); err == nil {
		githubLink.SetURL(u)
	}

	content := container.NewVBox(
		title,
		version,
		widget.NewSeparator(),
		desc,
		copyright,
		githubLink,
	)

	d := dialog.NewCustom("About MarkView", "Close", content, w.fyneWindow)
	d.Resize(fyne.NewSize(300, 200))
	d.Show()
}

// addRecentFile adds a file to the recent files list
func (w *Window) addRecentFile(path string) {
	// Remove if already in list
	for i, p := range w.recentFiles {
		if p == path {
			w.recentFiles = append(w.recentFiles[:i], w.recentFiles[i+1:]...)
			break
		}
	}

	// Add to front
	w.recentFiles = append([]string{path}, w.recentFiles...)

	// Keep only last 10
	if len(w.recentFiles) > 10 {
		w.recentFiles = w.recentFiles[:10]
	}

	// Save to preferences
	w.app.Preferences().SetString("recentFiles", strings.Join(w.recentFiles, "\n"))
}

// loadRecentFiles loads recent files from preferences
func (w *Window) loadRecentFiles() {
	recentStr := w.app.Preferences().String("recentFiles")
	if recentStr != "" {
		w.recentFiles = strings.Split(recentStr, "\n")
	}
}

// showRecentFiles shows a dialog with recent files
func (w *Window) showRecentFiles() {
	if len(w.recentFiles) == 0 {
		dialog.ShowInformation("Recent Files", "No recent files.", w.fyneWindow)
		return
	}

	list := widget.NewList(
		func() int { return len(w.recentFiles) },
		func() fyne.CanvasObject {
			return widget.NewLabel("filename.md")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(w.recentFiles) {
				obj.(*widget.Label).SetText(filepath.Base(w.recentFiles[id]))
			}
		},
	)

	var d dialog.Dialog
	list.OnSelected = func(id widget.ListItemID) {
		if id < len(w.recentFiles) {
			path := w.recentFiles[id]
			d.Hide()
			w.checkUnsavedChanges(func() {
				w.loadFile(path)
			})
		}
	}

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	d = dialog.NewCustom("Recent Files", "Cancel", scroll, w.fyneWindow)
	d.Resize(fyne.NewSize(450, 400))
	d.Show()
}

// startAutoSave starts the auto-save timer
func (w *Window) startAutoSave() {
	if w.autoSaveTicker != nil {
		w.autoSaveTicker.Stop()
	}
	w.autoSaveTicker = time.NewTicker(30 * time.Second)
	go func() {
		for range w.autoSaveTicker.C {
			if w.editMode && w.isDirty && w.currentFile != "" {
				w.saveFile()
				w.logger.Info("Auto-saved file")
			}
		}
	}()
}

// stopAutoSave stops the auto-save timer
func (w *Window) stopAutoSave() {
	if w.autoSaveTicker != nil {
		w.autoSaveTicker.Stop()
		w.autoSaveTicker = nil
	}
}

// showCustomCSSDialog shows a dialog to edit custom CSS for exports
func (w *Window) showCustomCSSDialog() {
	cssEntry := widget.NewMultiLineEntry()
	cssEntry.SetPlaceHolder("/* Custom CSS for exports */\n\nbody {\n    /* your styles here */\n}")
	cssEntry.SetText(w.customCSS)
	cssEntry.Wrapping = fyne.TextWrapOff

	scroll := container.NewScroll(cssEntry)
	scroll.SetMinSize(fyne.NewSize(500, 300))

	resetBtn := widget.NewButton("Reset to Default", func() {
		cssEntry.SetText("")
	})

	content := container.NewBorder(
		widget.NewLabel("Add custom CSS that will be included in HTML/PDF exports:"),
		resetBtn,
		nil, nil,
		scroll,
	)

	dialog.ShowCustomConfirm("Custom Export CSS", "Save", "Cancel", content, func(confirm bool) {
		if confirm {
			w.customCSS = cssEntry.Text
			w.app.Preferences().SetString("customCSS", w.customCSS)
		}
	}, w.fyneWindow)
}

// showWordCountGoalDialog shows a dialog to set the word count goal
func (w *Window) showWordCountGoalDialog() {
	goalEntry := widget.NewEntry()
	goalEntry.SetPlaceHolder("Enter goal (e.g., 500)")
	if w.wordCountGoal > 0 {
		goalEntry.SetText(intToString(w.wordCountGoal))
	}

	clearBtn := widget.NewButton("Clear Goal", func() {
		w.wordCountGoal = 0
		w.app.Preferences().SetInt("wordCountGoal", 0)
		w.updateStatusBar()
	})

	content := container.NewVBox(
		widget.NewLabel("Set a word count goal for this session:"),
		goalEntry,
		clearBtn,
	)

	dialog.ShowCustomConfirm("Word Count Goal", "Set", "Cancel", content, func(confirm bool) {
		if confirm {
			goal := parseIntFromString(goalEntry.Text)
			if goal > 0 {
				w.wordCountGoal = goal
				w.app.Preferences().SetInt("wordCountGoal", goal)
				w.updateStatusBar()
			}
		}
	}, w.fyneWindow)
}

// parseIntFromString parses an integer from a string
func parseIntFromString(s string) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

// showPresentationMode shows the markdown as a presentation
func (w *Window) showPresentationMode() {
	if w.currentFile == "" {
		dialog.ShowInformation("Presentation Mode", "No file is currently open.", w.fyneWindow)
		return
	}
	ShowPresentationMode(w.fyneWindow, w.contentBuffer)
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
	// Update autocomplete with new root path
	if w.linkAutocomplete != nil {
		w.linkAutocomplete.SetRootPath(path)
	}
	if w.splitLinkAutocomplete != nil {
		w.splitLinkAutocomplete.SetRootPath(path)
	}
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

// loadWindowSize loads the saved window size from preferences
func (w *Window) loadWindowSize() {
	width := w.app.Preferences().Float("windowWidth")
	height := w.app.Preferences().Float("windowHeight")
	if width > 0 && height > 0 {
		w.fyneWindow.Resize(fyne.NewSize(float32(width), float32(height)))
	}
}

// saveWindowSize saves the current window size to preferences
func (w *Window) saveWindowSize() {
	size := w.fyneWindow.Canvas().Size()
	w.app.Preferences().SetFloat("windowWidth", float64(size.Width))
	w.app.Preferences().SetFloat("windowHeight", float64(size.Height))
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
	w.app.Settings().SetTheme(themes.NewMarkViewThemeWithOptions(themeType, w.currentFont, w.currentFontSize))

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
	w.app.Settings().SetTheme(themes.NewMarkViewThemeWithOptions(w.currentTheme, fontFamily, w.currentFontSize))

	// Save font preference
	w.app.Preferences().SetString("font", string(fontFamily))

	// Refresh content
	if w.currentFile != "" {
		w.loadFile(w.currentFile)
	}

	w.logger.Info("Font changed", zap.String("font", string(fontFamily)))
}

// setFontSize changes the application font size
func (w *Window) setFontSize(fontSize themes.FontSize) {
	w.currentFontSize = fontSize
	w.app.Settings().SetTheme(themes.NewMarkViewThemeWithOptions(w.currentTheme, w.currentFont, fontSize))

	// Save font size preference
	w.app.Preferences().SetString("fontSize", string(fontSize))

	// Refresh content
	if w.currentFile != "" {
		w.loadFile(w.currentFile)
	}

	w.logger.Info("Font size changed", zap.String("fontSize", string(fontSize)))
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

	// Add to recent files
	w.addRecentFile(filePath)

	// Update window title
	fileName := filepath.Base(filePath)
	w.fyneWindow.SetTitle(fmt.Sprintf("MarkView - %s", fileName))

	// Set file tree root to the file's directory if not already set
	fileDir := filepath.Dir(filePath)
	if w.currentDir == "" || w.currentDir != fileDir {
		w.currentDir = fileDir
		w.fileTree.SetRootPath(fileDir)
		// Update autocomplete with new root path
		if w.linkAutocomplete != nil {
			w.linkAutocomplete.SetRootPath(fileDir)
		}
		if w.splitLinkAutocomplete != nil {
			w.splitLinkAutocomplete.SetRootPath(fileDir)
		}
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

	// Update content with left padding
	w.scrollContent.Content = withLeftPadding(content, 16)
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

	// Update content with left padding
	w.scrollContent.Content = withLeftPadding(content, 16)
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

	// Center cursor in typewriter mode
	if w.typewriterMode {
		w.centerCursor()
	}

	// Check for link autocomplete
	if w.linkAutocomplete != nil {
		w.linkAutocomplete.OnTextChanged(content)
	}

	// Update status bar
	w.updateStatusBar()
}

// updateSplitViewPreview updates the preview pane in split view
func (w *Window) updateSplitViewPreview(content string) {
	fileDir := ""
	if w.currentFile != "" {
		fileDir = filepath.Dir(w.currentFile)
	}

	parsedContent, err := w.parser.ParseWithBasePath([]byte(content), fileDir)
	if err != nil {
		return // Silently fail for preview updates
	}

	w.splitPreview.Content = withLeftPadding(parsedContent, 16)
	w.splitPreview.Refresh()
}

// updateMainPreview updates the main preview pane with the given content
func (w *Window) updateMainPreview(content string) {
	fileDir := ""
	if w.currentFile != "" {
		fileDir = filepath.Dir(w.currentFile)
	}

	parsedContent, err := w.parser.ParseWithBasePath([]byte(content), fileDir)
	if err != nil {
		return
	}

	w.scrollContent.Content = withLeftPadding(parsedContent, 16)
	w.scrollContent.Refresh()
}

// onSplitEditorChanged handles changes in the split view editor
func (w *Window) onSplitEditorChanged(content string) {
	if !w.splitViewMode {
		return
	}

	// Mark as dirty
	if !w.isDirty {
		w.isDirty = true
		w.updateWindowTitle()
	}

	// Update outline
	w.outline.UpdateFromText(content)

	// Update split view preview with live changes
	w.updateSplitViewPreview(content)
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

		// Show goal progress if goal is set
		if w.wordCountGoal > 0 {
			percentage := (words * 100) / w.wordCountGoal
			if percentage > 100 {
				percentage = 100
			}
			w.wordCount.SetText(fmt.Sprintf("%d/%d words (%d%%), %d chars", words, w.wordCountGoal, percentage, chars))
		} else {
			w.wordCount.SetText(fmt.Sprintf("%d words, %d chars", words, chars))
		}

		// Reading time (average 200 words per minute)
		readMins := words / 200
		if readMins < 1 {
			w.readingTime.SetText("< 1 min read")
		} else {
			w.readingTime.SetText(fmt.Sprintf("%d min read", readMins))
		}

		// Cursor position (1-indexed for display)
		row, col := w.editor.GetCursorPosition()
		w.cursorPos.SetText(fmt.Sprintf("Ln %d, Col %d", row+1, col+1))

		w.statusBar.SetText("Edit Mode")
	} else if w.currentFile != "" {
		// Show file path in view mode
		w.statusBar.SetText(w.currentFile)

		// Calculate word count and reading time for view mode too
		words := len(strings.Fields(w.contentBuffer))
		w.wordCount.SetText(fmt.Sprintf("%d words", words))

		readMins := words / 200
		if readMins < 1 {
			w.readingTime.SetText("< 1 min read")
		} else {
			w.readingTime.SetText(fmt.Sprintf("%d min read", readMins))
		}

		w.cursorPos.SetText("")
	} else {
		w.statusBar.SetText("No file open")
		w.wordCount.SetText("")
		w.readingTime.SetText("")
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

	// Get content from the appropriate editor
	var content string
	if w.splitViewMode {
		content = w.splitEditor.GetText()
	} else {
		content = w.editor.GetText()
	}

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

		// Get content from the appropriate editor
		var content string
		if w.splitViewMode {
			content = w.splitEditor.GetText()
		} else {
			content = w.editor.GetText()
		}
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
	// Check for updates silently on startup
	if w.updateChecker != nil {
		w.updateChecker.CheckForUpdates(true)
	}
	w.fyneWindow.ShowAndRun()
}

// GetMarkdown returns the markdown parser for external access
func (w *Window) GetMarkdown() *markdown.Parser {
	return w.parser
}

// showPrintDialog shows a print/export dialog
func (w *Window) showPrintDialog() {
	var d dialog.Dialog

	// Helper to create styled export buttons
	createExportButton := func(label, description string, icon fyne.Resource, onTap func()) fyne.CanvasObject {
		btn := widget.NewButtonWithIcon(label, icon, func() {
			d.Hide()
			onTap()
		})
		btn.Importance = widget.HighImportance

		descLabel := widget.NewLabel(description)
		descLabel.Wrapping = fyne.TextWrapWord
		descLabel.TextStyle = fyne.TextStyle{Italic: true}

		return container.NewVBox(btn, descLabel)
	}

	// File info header
	fileIcon := widget.NewIcon(themes.IconDocument())
	fileName := widget.NewLabelWithStyle(filepath.Base(w.currentFile), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	fileInfo := container.NewHBox(fileIcon, fileName)

	// Section: Export Formats
	sectionExport := widget.NewLabelWithStyle("Export Formats", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	htmlBtn := createExportButton("HTML", "Web-ready format with styling", themes.IconCode(), func() {
		w.exportToHTML()
	})

	pdfBtn := createExportButton("PDF", "Print-ready document (requires wkhtmltopdf)", themes.IconExport(), func() {
		w.exportToPDF()
	})

	docxBtn := createExportButton("DOCX", "Microsoft Word format (requires pandoc)", themes.IconDocument(), func() {
		w.exportWithPandoc("docx")
	})

	rtfBtn := createExportButton("RTF", "Rich Text Format (requires pandoc)", themes.IconDocument(), func() {
		w.exportWithPandoc("rtf")
	})

	// Export format grid - 2x2
	exportGrid := container.NewGridWithColumns(2,
		htmlBtn,
		pdfBtn,
		docxBtn,
		rtfBtn,
	)

	// Section: Print & Settings
	sectionSettings := widget.NewLabelWithStyle("Print & Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	printBtn := widget.NewButtonWithIcon("Print via Browser", themes.IconPrint(), func() {
		d.Hide()
		w.printViaBrowser()
	})

	themeBtn := widget.NewButtonWithIcon("Export Theme", themes.IconTheme(), func() {
		w.showExportThemeDialog()
	})

	cssBtn := widget.NewButtonWithIcon("Custom CSS", themes.IconCode(), func() {
		w.showCustomCSSDialog()
	})

	settingsRow := container.NewGridWithColumns(3, printBtn, themeBtn, cssBtn)

	// Current theme indicator
	themeLabel := widget.NewLabel(fmt.Sprintf("Current theme: %s", w.currentExportTheme))
	themeLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Build dialog content
	dialogContent := container.NewVBox(
		fileInfo,
		widget.NewSeparator(),
		sectionExport,
		exportGrid,
		widget.NewSeparator(),
		sectionSettings,
		settingsRow,
		themeLabel,
	)

	// Add padding
	paddedContent := container.NewPadded(dialogContent)

	d = dialog.NewCustom("Export Document", "Close", paddedContent, w.fyneWindow)
	d.Resize(fyne.NewSize(500, 420))
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

// findChromeBrowser looks for Chrome/Chromium browsers on the system
func findChromeBrowser() string {
	// macOS application paths
	macPaths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
	}

	// Check macOS paths first
	for _, path := range macPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Check PATH for Linux/other systems
	linuxBrowsers := []string{"google-chrome", "chromium", "chromium-browser", "brave-browser", "microsoft-edge"}
	for _, browser := range linuxBrowsers {
		if path, err := exec.LookPath(browser); err == nil {
			return path
		}
	}

	return ""
}

// exportToPDF exports the current markdown to PDF
func (w *Window) exportToPDF() {
	if w.currentFile == "" {
		return
	}

	// Find a PDF converter: prefer Chrome, fall back to wkhtmltopdf
	chromePath := findChromeBrowser()
	wkhtmltopdfPath, _ := exec.LookPath("wkhtmltopdf")

	if chromePath == "" && wkhtmltopdfPath == "" {
		dialog.ShowError(fmt.Errorf("PDF export requires Chrome/Chromium or wkhtmltopdf.\n\n"+
			"Option 1: Install Google Chrome (recommended)\n"+
			"Option 2: Download wkhtmltopdf from https://wkhtmltopdf.org/downloads.html\n"+
			"Option 3: Use 'Print via Browser' and save as PDF"), w.fyneWindow)
		return
	}

	// Read the markdown file
	data, err := os.ReadFile(w.currentFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read file: %w", err), w.fyneWindow)
		return
	}

	// Convert to HTML (without external scripts for offline PDF generation)
	html := w.markdownToHTMLForPDF(data)

	// Show save dialog for PDF
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w.fyneWindow)
			return
		}
		if writer == nil {
			return
		}

		// Get the path and close/remove the empty file Fyne created
		pdfPath := writer.URI().Path()
		writer.Close()

		// URL decode the path in case of encoded characters
		if decoded, err := url.PathUnescape(pdfPath); err == nil {
			pdfPath = decoded
		}

		os.Remove(pdfPath) // Remove empty file so PDF converter can create it

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

		var cmd *exec.Cmd
		if chromePath != "" {
			// Use Chrome headless for PDF generation
			cmd = exec.Command(chromePath,
				"--headless",
				"--disable-gpu",
				"--no-sandbox",
				"--print-to-pdf="+pdfPath,
				"--print-to-pdf-no-header",
				"file://"+tmpFile.Name())
		} else {
			// Fall back to wkhtmltopdf
			cmd = exec.Command(wkhtmltopdfPath,
				"--enable-local-file-access",
				"--load-error-handling", "ignore",
				"--load-media-error-handling", "ignore",
				"--no-stop-slow-scripts",
				"--quiet",
				tmpFile.Name(), pdfPath)
		}

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

// generateTOCHTML generates an HTML table of contents from markdown headings
func generateTOCHTML(data []byte) string {
	lines := strings.Split(string(data), "\n")
	var tocItems []string
	headingID := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			// Count the number of # to determine heading level
			level := 0
			for _, c := range trimmed {
				if c == '#' {
					level++
				} else {
					break
				}
			}
			if level > 0 && level <= 6 {
				text := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
				if text != "" {
					headingID++
					indent := strings.Repeat("  ", level-1)
					anchorID := fmt.Sprintf("heading-%d", headingID)
					tocItems = append(tocItems, fmt.Sprintf("%s<li><a href=\"#%s\">%s</a></li>", indent, anchorID, text))
				}
			}
		}
	}

	if len(tocItems) == 0 {
		return ""
	}

	return "<nav class=\"toc\">\n<h2>Table of Contents</h2>\n<ul>\n" + strings.Join(tocItems, "\n") + "\n</ul>\n</nav>\n"
}

// addHeadingIDs adds id attributes to headings in HTML content
func addHeadingIDs(html string) string {
	headingID := 0
	result := html

	// Simple regex-like replacement for h1-h6 tags
	for level := 1; level <= 6; level++ {
		openTag := fmt.Sprintf("<h%d>", level)
		for strings.Contains(result, openTag) {
			headingID++
			replacement := fmt.Sprintf("<h%d id=\"heading-%d\">", level, headingID)
			result = strings.Replace(result, openTag, replacement, 1)
		}
	}

	return result
}

// markdownToHTML converts markdown to a styled HTML document
func (w *Window) markdownToHTML(data []byte) string {
	return w.markdownToHTMLWithOptions(data, true)
}

// markdownToHTMLForPDF generates HTML without external scripts (for offline PDF export)
func (w *Window) markdownToHTMLForPDF(data []byte) string {
	return w.markdownToHTMLWithOptions(data, false)
}

func (w *Window) markdownToHTMLWithOptions(data []byte, includeScripts bool) string {
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

	// Generate TOC
	tocHTML := generateTOCHTML(data)

	// Add heading IDs to content
	contentHTML := addHeadingIDs(buf.String())

	// Wrap in styled HTML document
	title := filepath.Base(w.currentFile)

	// Get the export theme CSS
	themeCSS := w.getExportThemeCSS()

	// External scripts for Mermaid and MathJax (only when online/browser viewing)
	externalScripts := ""
	if includeScripts {
		externalScripts = `
    <!-- Mermaid for diagram rendering -->
    <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
    <script>mermaid.initialize({startOnLoad:true});</script>
    <!-- MathJax for math rendering -->
    <script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
    <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>`
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <style>
%s
        .mermaid { text-align: center; }
        .toc {
            background-color: #f7fafc;
            border: 1px solid #e2e8f0;
            border-radius: 6px;
            padding: 16px;
            margin-bottom: 24px;
        }
        .toc h2 {
            margin-top: 0;
            border-bottom: none;
            font-size: 1.1em;
            color: #4a5568;
        }
        .toc ul {
            list-style-type: none;
            padding-left: 0;
            margin-bottom: 0;
        }
        .toc li {
            margin: 4px 0;
        }
        .toc a {
            text-decoration: none;
        }
        .toc a:hover {
            text-decoration: underline;
        }
    </style>
    %s
    %s
</head>
<body>
%s
%s
</body>
</html>`, title, themeCSS, externalScripts, w.getCustomCSSStyle(), tocHTML, contentHTML)
}

// getCustomCSSStyle returns the custom CSS wrapped in a style tag, or empty string if no custom CSS
func (w *Window) getCustomCSSStyle() string {
	if w.customCSS == "" {
		return ""
	}
	return "<style>\n/* Custom CSS */\n" + w.customCSS + "\n</style>"
}

// withLeftPadding wraps content with left padding
func withLeftPadding(content fyne.CanvasObject, padding float32) fyne.CanvasObject {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(padding, 0))
	return container.NewBorder(nil, nil, spacer, nil, content)
}

// showQuickSwitcher shows the quick file switcher
func (w *Window) showQuickSwitcher() {
	ShowQuickSwitcher(w.fyneWindow, w.currentDir, func(path string) {
		w.checkUnsavedChanges(func() {
			w.loadFile(path)
		})
	})
}

// shortcutLabel renders a Super-modifier custom shortcut for display, e.g.
// "Cmd+Shift+P". It returns "" for shortcuts it cannot describe.
func shortcutLabel(s fyne.Shortcut) string {
	cs, ok := s.(*desktop.CustomShortcut)
	if !ok || cs == nil {
		return ""
	}
	var parts []string
	if cs.Modifier&fyne.KeyModifierSuper != 0 {
		parts = append(parts, cmdKey)
	}
	if cs.Modifier&fyne.KeyModifierControl != 0 {
		parts = append(parts, "Ctrl")
	}
	if cs.Modifier&fyne.KeyModifierShift != 0 {
		parts = append(parts, "Shift")
	}
	if cs.Modifier&fyne.KeyModifierAlt != 0 {
		parts = append(parts, "Alt")
	}
	parts = append(parts, string(cs.KeyName))
	return strings.Join(parts, "+")
}

// showCommandPalette opens a searchable list of every command in the menu bar,
// keeping the palette automatically in sync with the menus.
func (w *Window) showCommandPalette() {
	menu := w.createMainMenu()
	var cmds []command
	for _, m := range menu.Items {
		for _, item := range m.Items {
			if item.IsSeparator || item.Action == nil {
				continue
			}
			// Skip the palette's own entry to avoid a self-referential command.
			if item.Label == "Command Palette" {
				continue
			}
			cmds = append(cmds, command{
				label:    m.Label + ": " + item.Label,
				shortcut: shortcutLabel(item.Shortcut),
				action:   item.Action,
			})
		}
	}
	NewCommandPalette(cmds).Show(w.fyneWindow)
}

// showFullTextSearch shows the full-text search dialog
func (w *Window) showFullTextSearch() {
	ShowFullTextSearch(w.fyneWindow, w.currentDir, func(path string, line int) {
		w.checkUnsavedChanges(func() {
			w.loadFile(path)
			// Navigate to line if in edit mode
			if w.editMode {
				w.navigateToLine(line - 1)
			}
		})
	})
}

// showBacklinks shows the backlinks panel
func (w *Window) showBacklinks() {
	ShowBacklinks(w.fyneWindow, w.currentDir, w.currentFile, func(path string, line int) {
		w.checkUnsavedChanges(func() {
			w.loadFile(path)
		})
	})
}

// validateLinks validates links in the current document
func (w *Window) validateLinks() {
	if w.currentFile == "" {
		dialog.ShowInformation("Link Validation", "No file is currently open.", w.fyneWindow)
		return
	}

	basePath := filepath.Dir(w.currentFile)
	ShowLinkValidationDialog(w.fyneWindow, basePath, w.contentBuffer, func(line int) {
		if w.editMode {
			w.navigateToLine(line - 1)
		}
	})
}

// showTagsBrowser shows the tags browser dialog
func (w *Window) showTagsBrowser() {
	// Reindex current file if open
	if w.currentFile != "" {
		w.tagManager.IndexFile(w.currentFile, w.contentBuffer)
	}

	ShowTagsDialog(w.fyneWindow, w.tagManager, func(path string) {
		w.checkUnsavedChanges(func() {
			w.loadFile(path)
		})
	})
}

// newFileFromTemplate creates a new file from a template
func (w *Window) newFileFromTemplate() {
	ShowTemplatesDialog(w.fyneWindow, func(content string) {
		w.checkUnsavedChanges(func() {
			// Clear current file
			w.currentFile = ""
			w.contentBuffer = content
			w.isDirty = true

			// Clear rendered content
			w.scrollContent.Content = container.NewVBox()
			w.scrollContent.Refresh()

			// Switch to edit mode with template content
			w.editMode = true
			w.editor.SetText(w.contentBuffer)
			w.scrollContent.Hide()
			w.editorScroll.Show()
			w.editToolbar.Show()
			w.editAction.SetIcon(themes.IconView())

			w.updateWindowTitle()
			w.editor.Focus(w.fyneWindow.Canvas())

			w.logger.Info("Created new file from template")
		})
	})
}

// toggleZenMode toggles zen mode (fullscreen distraction-free)
func (w *Window) toggleZenMode() {
	w.zenMode = !w.zenMode

	if w.zenMode {
		// Enter zen mode - fullscreen and hide all UI
		w.fyneWindow.SetFullScreen(true)

		// Hide all sidebars
		w.fileTreeScroll.Hide()
		w.tocScroll.Hide()
		w.outlineScroll.Hide()
		w.editToolbar.Hide()

		// Set main split to show only content
		w.mainSplit.Offset = 0
	} else {
		// Exit zen mode
		w.fyneWindow.SetFullScreen(false)

		// Restore UI based on current mode
		w.fileTreeScroll.Show()
		if w.editMode {
			w.outlineScroll.Show()
			w.editToolbar.Show()
		} else {
			w.tocScroll.Show()
		}
		w.mainSplit.Offset = 0.30
	}
	w.mainSplit.Refresh()
}

// showExportThemeDialog shows the export theme selection dialog
func (w *Window) showExportThemeDialog() {
	themes := GetExportThemes()
	themeNames := make([]string, len(themes))
	for i, t := range themes {
		themeNames[i] = t.Name
	}

	themeRadio := widget.NewRadioGroup(themeNames, func(selected string) {
		w.currentExportTheme = selected
		w.app.Preferences().SetString("exportTheme", selected)
	})
	themeRadio.SetSelected(w.currentExportTheme)

	// Show theme descriptions
	var descLabels []*widget.Label
	for _, t := range themes {
		descLabels = append(descLabels, widget.NewLabel(t.Description))
	}

	list := widget.NewList(
		func() int { return len(themes) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("Theme Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Description"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(themes) {
				box := obj.(*fyne.Container)
				nameLabel := box.Objects[0].(*widget.Label)
				descLabel := box.Objects[1].(*widget.Label)

				t := themes[id]
				if t.Name == w.currentExportTheme {
					nameLabel.SetText(t.Name + " (selected)")
				} else {
					nameLabel.SetText(t.Name)
				}
				descLabel.SetText(t.Description)
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(themes) {
			w.currentExportTheme = themes[id].Name
			w.app.Preferences().SetString("exportTheme", w.currentExportTheme)
			list.Refresh()
		}
	}

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	content := container.NewBorder(
		widget.NewLabel("Select an export theme style:"),
		nil, nil, nil,
		scroll,
	)

	d := dialog.NewCustom("Export Theme", "Close", content, w.fyneWindow)
	d.Resize(fyne.NewSize(450, 400))
	d.Show()
}

// getExportThemeCSS returns the CSS for the current export theme
func (w *Window) getExportThemeCSS() string {
	themes := GetExportThemes()
	for _, t := range themes {
		if t.Name == w.currentExportTheme {
			return t.CSS
		}
	}
	// Default to first theme
	if len(themes) > 0 {
		return themes[0].CSS
	}
	return ""
}
