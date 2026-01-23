package gui

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LinkAutocomplete provides inline autocomplete for markdown links
type LinkAutocomplete struct {
	editor      *MarkdownEditor
	window      fyne.Window
	rootPath    string
	popup       *widget.PopUp
	list        *widget.List
	files       []string // All markdown files (relative paths)
	filtered    []string // Filtered matches
	selectedIdx int
	active      bool
	linkStart   int    // Position where link path starts
	partialPath string // Current partial path being typed
}

// NewLinkAutocomplete creates a new link autocomplete helper
func NewLinkAutocomplete(editor *MarkdownEditor, window fyne.Window, rootPath string) *LinkAutocomplete {
	la := &LinkAutocomplete{
		editor:   editor,
		window:   window,
		rootPath: rootPath,
		filtered: []string{},
	}
	la.scanFiles()
	la.createPopup()
	return la
}

// SetRootPath updates the root path and rescans files
func (la *LinkAutocomplete) SetRootPath(rootPath string) {
	la.rootPath = rootPath
	la.scanFiles()
}

// scanFiles scans the root path for markdown files
func (la *LinkAutocomplete) scanFiles() {
	la.files = []string{}
	if la.rootPath == "" {
		return
	}

	filepath.Walk(la.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			// Skip hidden directories
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if isMarkdownFile(path) {
			// Store relative path for display
			relPath, _ := filepath.Rel(la.rootPath, path)
			la.files = append(la.files, relPath)
		}
		return nil
	})

	// Sort by filename
	sort.Slice(la.files, func(i, j int) bool {
		return strings.ToLower(filepath.Base(la.files[i])) < strings.ToLower(filepath.Base(la.files[j]))
	})
}

// createPopup creates the autocomplete popup widget
func (la *LinkAutocomplete) createPopup() {
	la.list = widget.NewList(
		func() int { return len(la.filtered) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("filename.md"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(la.filtered) {
				box := obj.(*fyne.Container)
				label := box.Objects[1].(*widget.Label)
				label.SetText(la.filtered[id])
			}
		},
	)

	la.list.OnSelected = func(id widget.ListItemID) {
		la.selectedIdx = id
		la.AcceptSelection()
	}
}

// OnTextChanged is called when the editor content changes
func (la *LinkAutocomplete) OnTextChanged(content string) {
	if la.editor == nil {
		return
	}

	row, col := la.editor.GetCursorPosition()
	cursorPos := la.getCursorPositionInText(content, row, col)

	// Check if we're in a link destination
	inLink, partialPath, linkStart := la.isInLinkDestination(content, cursorPos)

	if inLink && len(la.files) > 0 {
		la.partialPath = partialPath
		la.linkStart = linkStart
		la.filterFiles(partialPath)

		if len(la.filtered) > 0 {
			la.showPopup()
		} else {
			la.hidePopup()
		}
	} else {
		la.hidePopup()
	}
}

// getCursorPositionInText converts row,col to absolute position
func (la *LinkAutocomplete) getCursorPositionInText(content string, row, col int) int {
	lines := strings.Split(content, "\n")
	pos := 0
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	if row < len(lines) && col <= len(lines[row]) {
		pos += col
	}
	return pos
}

// isInLinkDestination checks if cursor is inside a markdown link destination
// Pattern: [any text](partial_path|cursor_here
// Returns: inLink, partialPath, startPosition
func (la *LinkAutocomplete) isInLinkDestination(content string, cursorPos int) (bool, string, int) {
	if cursorPos > len(content) {
		cursorPos = len(content)
	}

	// Get text before cursor
	textBefore := content[:cursorPos]

	// Find the last occurrence of "](" before cursor
	linkPattern := regexp.MustCompile(`\]\([^)]*$`)
	match := linkPattern.FindStringIndex(textBefore)

	if match == nil {
		return false, "", 0
	}

	// Check that there's no closing ")" between "](" and cursor
	afterBracket := textBefore[match[0]+2:] // Skip "]("
	if strings.Contains(afterBracket, ")") {
		return false, "", 0
	}

	// Extract the partial path
	linkStart := match[0] + 2 // Position after "]("
	partialPath := textBefore[linkStart:]

	// Don't trigger for URLs
	if strings.HasPrefix(partialPath, "http://") || strings.HasPrefix(partialPath, "https://") ||
		strings.HasPrefix(partialPath, "mailto:") || strings.HasPrefix(partialPath, "#") {
		return false, "", 0
	}

	return true, partialPath, linkStart
}

// filterFiles filters files matching the partial path
func (la *LinkAutocomplete) filterFiles(partialPath string) {
	partialPath = strings.ToLower(partialPath)
	la.filtered = []string{}
	la.selectedIdx = 0

	for _, f := range la.files {
		fLower := strings.ToLower(f)
		// Match if the file path contains the partial path
		if partialPath == "" || strings.Contains(fLower, partialPath) || fuzzyMatch(fLower, partialPath) {
			la.filtered = append(la.filtered, f)
		}
	}

	// Limit to 10 suggestions
	if len(la.filtered) > 10 {
		la.filtered = la.filtered[:10]
	}

	if la.list != nil {
		la.list.Refresh()
	}
}

// showPopup displays the autocomplete popup
func (la *LinkAutocomplete) showPopup() {
	if la.popup != nil {
		la.popup.Hide()
	}

	// Create a background for the popup
	bg := canvas.NewRectangle(theme.OverlayBackgroundColor())
	bg.CornerRadius = 4

	// Create the list with minimum size
	la.list.Refresh()
	listContainer := container.NewStack(bg, la.list)

	// Calculate popup size based on number of items
	itemHeight := float32(30)
	listHeight := float32(len(la.filtered)) * itemHeight
	if listHeight > 300 {
		listHeight = 300
	}
	if listHeight < itemHeight {
		listHeight = itemHeight
	}

	listContainer.Resize(fyne.NewSize(300, listHeight))

	// Create and position the popup
	la.popup = widget.NewPopUp(listContainer, la.window.Canvas())

	// Position near the editor cursor
	// For now, position at a fixed location since getting exact cursor position is complex
	editorPos := la.editor.Position()
	row, _ := la.editor.GetCursorPosition()
	lineHeight := float32(20)
	popupPos := fyne.NewPos(editorPos.X+50, editorPos.Y+float32(row+1)*lineHeight+30)

	la.popup.ShowAtPosition(popupPos)
	la.popup.Resize(fyne.NewSize(300, listHeight))
	la.active = true
}

// hidePopup hides the autocomplete popup
func (la *LinkAutocomplete) hidePopup() {
	if la.popup != nil {
		la.popup.Hide()
	}
	la.active = false
	la.selectedIdx = 0
}

// IsActive returns whether the autocomplete is currently showing
func (la *LinkAutocomplete) IsActive() bool {
	return la.active
}

// SelectNext moves selection to the next item
func (la *LinkAutocomplete) SelectNext() {
	if !la.active || len(la.filtered) == 0 {
		return
	}
	la.selectedIdx++
	if la.selectedIdx >= len(la.filtered) {
		la.selectedIdx = 0
	}
	if la.list != nil {
		la.list.Select(la.selectedIdx)
	}
}

// SelectPrevious moves selection to the previous item
func (la *LinkAutocomplete) SelectPrevious() {
	if !la.active || len(la.filtered) == 0 {
		return
	}
	la.selectedIdx--
	if la.selectedIdx < 0 {
		la.selectedIdx = len(la.filtered) - 1
	}
	if la.list != nil {
		la.list.Select(la.selectedIdx)
	}
}

// AcceptSelection inserts the selected file path
func (la *LinkAutocomplete) AcceptSelection() {
	if !la.active || len(la.filtered) == 0 || la.selectedIdx >= len(la.filtered) {
		la.hidePopup()
		return
	}

	// Editor is required to insert selection
	if la.editor == nil {
		la.hidePopup()
		return
	}

	selectedFile := la.filtered[la.selectedIdx]

	// Get current editor content
	content := la.editor.GetText()

	// Replace the partial path with the selected file
	// We need to find where the partial path starts and replace from there to cursor
	row, col := la.editor.GetCursorPosition()
	cursorPos := la.getCursorPositionInText(content, row, col)

	// Calculate replacement range
	replaceStart := la.linkStart
	replaceEnd := cursorPos

	if replaceStart < 0 || replaceStart > len(content) {
		la.hidePopup()
		return
	}
	if replaceEnd > len(content) {
		replaceEnd = len(content)
	}

	// Build new content
	newContent := content[:replaceStart] + selectedFile + content[replaceEnd:]

	// Set the new content
	la.editor.SetText(newContent)

	// Position cursor after the inserted path
	newCursorPos := replaceStart + len(selectedFile)
	la.editor.setCursorPosition(newCursorPos)

	la.hidePopup()
}

// Dismiss hides the popup without selecting
func (la *LinkAutocomplete) Dismiss() {
	la.hidePopup()
}

// HandleKeyEvent handles key events for the autocomplete
// Returns true if the event was handled
func (la *LinkAutocomplete) HandleKeyEvent(key *fyne.KeyEvent) bool {
	if !la.active {
		return false
	}

	switch key.Name {
	case fyne.KeyDown:
		la.SelectNext()
		return true
	case fyne.KeyUp:
		la.SelectPrevious()
		return true
	case fyne.KeyReturn, fyne.KeyEnter:
		la.AcceptSelection()
		return true
	case fyne.KeyTab:
		la.AcceptSelection()
		return true
	case fyne.KeyEscape:
		la.Dismiss()
		return true
	}

	return false
}
