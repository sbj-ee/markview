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

// LinkType represents the type of link being autocompleted
type LinkType int

const (
	LinkTypeStandard LinkType = iota // [text](path)
	LinkTypeWiki                     // [[path]]
	LinkTypeImage                    // ![alt](path)
)

// LinkAutocomplete provides inline autocomplete for markdown links
type LinkAutocomplete struct {
	editor      *MarkdownEditor
	window      fyne.Window
	rootPath    string
	popup       *widget.PopUp
	list        *widget.List
	files       []string // All markdown files (relative paths)
	images      []string // All image files (relative paths)
	filtered    []string // Filtered matches
	selectedIdx int
	active      bool
	linkStart   int      // Position where link path starts
	partialPath string   // Current partial path being typed
	linkType    LinkType // Type of link being completed
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

// scanFiles scans the root path for markdown and image files
func (la *LinkAutocomplete) scanFiles() {
	la.files = []string{}
	la.images = []string{}
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
		relPath, _ := filepath.Rel(la.rootPath, path)
		if isMarkdownFile(path) {
			la.files = append(la.files, relPath)
		} else if isImageFile(path) {
			la.images = append(la.images, relPath)
		}
		return nil
	})

	// Sort by filename
	sort.Slice(la.files, func(i, j int) bool {
		return strings.ToLower(filepath.Base(la.files[i])) < strings.ToLower(filepath.Base(la.files[j]))
	})
	sort.Slice(la.images, func(i, j int) bool {
		return strings.ToLower(filepath.Base(la.images[i])) < strings.ToLower(filepath.Base(la.images[j]))
	})
}

// isImageFile checks if a file is an image based on extension
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".bmp", ".ico":
		return true
	}
	return false
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
				icon := box.Objects[0].(*widget.Icon)
				label := box.Objects[1].(*widget.Label)

				// Set icon based on link type
				if la.linkType == LinkTypeImage {
					icon.SetResource(theme.FileImageIcon())
				} else {
					icon.SetResource(theme.DocumentIcon())
				}
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

	// Check if we're in a link destination (standard, wiki, or image)
	inLink, partialPath, linkStart, linkType := la.detectLinkContext(content, cursorPos)

	if inLink {
		la.partialPath = partialPath
		la.linkStart = linkStart
		la.linkType = linkType
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

// detectLinkContext detects if cursor is in a link and returns the type
// Supports: [text](path), [[path]], ![alt](path)
// Returns: inLink, partialPath, startPosition, linkType
func (la *LinkAutocomplete) detectLinkContext(content string, cursorPos int) (bool, string, int, LinkType) {
	if cursorPos > len(content) {
		cursorPos = len(content)
	}

	textBefore := content[:cursorPos]

	// Check for wiki-style link [[partial
	if inLink, partial, start := la.isInWikiLink(textBefore); inLink {
		return true, partial, start, LinkTypeWiki
	}

	// Check for image link ![alt](partial
	if inLink, partial, start := la.isInImageLink(textBefore); inLink {
		return true, partial, start, LinkTypeImage
	}

	// Check for standard link [text](partial
	if inLink, partial, start := la.isInStandardLink(textBefore); inLink {
		return true, partial, start, LinkTypeStandard
	}

	return false, "", 0, LinkTypeStandard
}

// isInWikiLink checks for [[partial pattern
func (la *LinkAutocomplete) isInWikiLink(textBefore string) (bool, string, int) {
	// Find last [[ that isn't closed
	wikiPattern := regexp.MustCompile(`\[\[[^\]]*$`)
	match := wikiPattern.FindStringIndex(textBefore)

	if match == nil {
		return false, "", 0
	}

	linkStart := match[0] + 2 // Position after "[["
	partialPath := textBefore[linkStart:]

	// Don't trigger if there's a | (for [[file|display]] syntax)
	if strings.Contains(partialPath, "|") {
		return false, "", 0
	}

	return true, partialPath, linkStart
}

// isInImageLink checks for ![alt](partial pattern
func (la *LinkAutocomplete) isInImageLink(textBefore string) (bool, string, int) {
	// Find ![...]( pattern
	imgPattern := regexp.MustCompile(`!\[[^\]]*\]\([^)]*$`)
	match := imgPattern.FindStringIndex(textBefore)

	if match == nil {
		return false, "", 0
	}

	// Find the position after ](
	bracketPos := strings.LastIndex(textBefore[match[0]:], "](")
	if bracketPos == -1 {
		return false, "", 0
	}

	linkStart := match[0] + bracketPos + 2
	partialPath := textBefore[linkStart:]

	// Don't trigger for URLs
	if strings.HasPrefix(partialPath, "http://") || strings.HasPrefix(partialPath, "https://") {
		return false, "", 0
	}

	return true, partialPath, linkStart
}

// isInStandardLink checks for [text](partial pattern (but not image)
func (la *LinkAutocomplete) isInStandardLink(textBefore string) (bool, string, int) {
	// Find ](  but not !]( which would be an image
	linkPattern := regexp.MustCompile(`\]\([^)]*$`)
	match := linkPattern.FindStringIndex(textBefore)

	if match == nil {
		return false, "", 0
	}

	// Check it's not an image link (no ! before [)
	if match[0] > 0 {
		// Look back for the opening [
		prefix := textBefore[:match[0]]
		lastBracket := strings.LastIndex(prefix, "[")
		if lastBracket > 0 && prefix[lastBracket-1] == '!' {
			return false, "", 0 // It's an image link
		}
	}

	linkStart := match[0] + 2 // Position after "]("
	partialPath := textBefore[linkStart:]

	// Don't trigger for URLs or anchors
	if strings.HasPrefix(partialPath, "http://") || strings.HasPrefix(partialPath, "https://") ||
		strings.HasPrefix(partialPath, "mailto:") || strings.HasPrefix(partialPath, "#") {
		return false, "", 0
	}

	return true, partialPath, linkStart
}

// isInLinkDestination checks if cursor is inside a markdown link destination (legacy)
// Pattern: [any text](partial_path|cursor_here
// Returns: inLink, partialPath, startPosition
func (la *LinkAutocomplete) isInLinkDestination(content string, cursorPos int) (bool, string, int) {
	inLink, partial, start, _ := la.detectLinkContext(content, cursorPos)
	return inLink, partial, start
}

// filterFiles filters files matching the partial path
func (la *LinkAutocomplete) filterFiles(partialPath string) {
	partialPath = strings.ToLower(partialPath)
	la.filtered = []string{}
	la.selectedIdx = 0

	// Choose source based on link type
	var source []string
	if la.linkType == LinkTypeImage {
		source = la.images
	} else {
		source = la.files
	}

	for _, f := range source {
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

	// Format insertion based on link type
	var insertion string
	switch la.linkType {
	case LinkTypeWiki:
		// For wiki links, strip .md extension and add closing ]]
		name := selectedFile
		if strings.HasSuffix(strings.ToLower(name), ".md") {
			name = name[:len(name)-3]
		}
		insertion = name + "]]"
	default:
		insertion = selectedFile
	}

	// Build new content
	newContent := content[:replaceStart] + insertion + content[replaceEnd:]

	// Set the new content
	la.editor.SetText(newContent)

	// Position cursor after the inserted path
	newCursorPos := replaceStart + len(insertion)
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
