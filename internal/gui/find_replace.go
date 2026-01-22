package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// FindReplaceDialog provides find and replace functionality
type FindReplaceDialog struct {
	window      fyne.Window
	editor      *MarkdownEditor
	findEntry   *widget.Entry
	replaceEntry *widget.Entry
	matchCase   *widget.Check
	resultLabel *widget.Label
	currentPos  int
	matches     []int
}

// NewFindReplaceDialog creates a new find and replace dialog
func NewFindReplaceDialog(window fyne.Window, editor *MarkdownEditor) *FindReplaceDialog {
	return &FindReplaceDialog{
		window:     window,
		editor:     editor,
		currentPos: -1,
	}
}

// Show displays the find and replace dialog
func (f *FindReplaceDialog) Show() {
	f.findEntry = widget.NewEntry()
	f.findEntry.SetPlaceHolder("Find...")
	f.findEntry.OnChanged = func(text string) {
		f.updateMatches()
	}

	f.replaceEntry = widget.NewEntry()
	f.replaceEntry.SetPlaceHolder("Replace with...")

	f.matchCase = widget.NewCheck("Match case", func(bool) {
		f.updateMatches()
	})

	f.resultLabel = widget.NewLabel("")

	findNextBtn := widget.NewButton("Find Next", func() {
		f.findNext()
	})

	findPrevBtn := widget.NewButton("Find Previous", func() {
		f.findPrevious()
	})

	replaceBtn := widget.NewButton("Replace", func() {
		f.replace()
	})

	replaceAllBtn := widget.NewButton("Replace All", func() {
		f.replaceAll()
	})

	findRow := container.NewBorder(nil, nil, widget.NewLabel("Find:"), nil, f.findEntry)
	replaceRow := container.NewBorder(nil, nil, widget.NewLabel("Replace:"), nil, f.replaceEntry)

	buttonRow := container.NewHBox(
		findPrevBtn,
		findNextBtn,
		replaceBtn,
		replaceAllBtn,
	)

	content := container.NewVBox(
		findRow,
		replaceRow,
		f.matchCase,
		buttonRow,
		f.resultLabel,
	)

	d := dialog.NewCustom("Find and Replace", "Close", content, f.window)
	d.Resize(fyne.NewSize(450, 250))
	d.Show()

	// Focus the find entry
	f.window.Canvas().Focus(f.findEntry)
}

// updateMatches finds all matches in the text
func (f *FindReplaceDialog) updateMatches() {
	searchText := f.findEntry.Text
	if searchText == "" {
		f.matches = nil
		f.currentPos = -1
		f.resultLabel.SetText("")
		return
	}

	content := f.editor.GetText()
	if !f.matchCase.Checked {
		content = strings.ToLower(content)
		searchText = strings.ToLower(searchText)
	}

	f.matches = nil
	pos := 0
	for {
		idx := strings.Index(content[pos:], searchText)
		if idx == -1 {
			break
		}
		f.matches = append(f.matches, pos+idx)
		pos += idx + 1
	}

	if len(f.matches) == 0 {
		f.resultLabel.SetText("No matches found")
		f.currentPos = -1
	} else {
		f.resultLabel.SetText(strings.ReplaceAll("X matches found", "X", string(rune('0'+len(f.matches)%10))))
		if len(f.matches) > 9 {
			f.resultLabel.SetText(strings.Replace(f.resultLabel.Text, string(rune('0'+len(f.matches)%10)),
				string(rune('0'+len(f.matches)/10))+string(rune('0'+len(f.matches)%10)), 1))
		}
		f.resultLabel.SetText(formatMatchCount(len(f.matches)))
	}
}

// formatMatchCount formats the match count message
func formatMatchCount(count int) string {
	if count == 1 {
		return "1 match found"
	}
	return strings.Replace("N matches found", "N", intToString(count), 1)
}

// intToString converts an int to string without fmt
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// findNext moves to the next match
func (f *FindReplaceDialog) findNext() {
	if len(f.matches) == 0 {
		return
	}

	f.currentPos++
	if f.currentPos >= len(f.matches) {
		f.currentPos = 0
	}

	pos := f.matches[f.currentPos]
	f.editor.setCursorPosition(pos)
	f.resultLabel.SetText(formatCurrentMatch(f.currentPos+1, len(f.matches)))
}

// findPrevious moves to the previous match
func (f *FindReplaceDialog) findPrevious() {
	if len(f.matches) == 0 {
		return
	}

	f.currentPos--
	if f.currentPos < 0 {
		f.currentPos = len(f.matches) - 1
	}

	pos := f.matches[f.currentPos]
	f.editor.setCursorPosition(pos)
	f.resultLabel.SetText(formatCurrentMatch(f.currentPos+1, len(f.matches)))
}

// formatCurrentMatch formats the current match position message
func formatCurrentMatch(current, total int) string {
	return intToString(current) + " of " + intToString(total)
}

// replace replaces the current match
func (f *FindReplaceDialog) replace() {
	if len(f.matches) == 0 || f.currentPos < 0 {
		return
	}

	searchText := f.findEntry.Text
	replaceText := f.replaceEntry.Text
	content := f.editor.GetText()

	pos := f.matches[f.currentPos]
	newContent := content[:pos] + replaceText + content[pos+len(searchText):]
	f.editor.SetText(newContent)

	f.updateMatches()
	if len(f.matches) > 0 && f.currentPos >= len(f.matches) {
		f.currentPos = len(f.matches) - 1
	}
}

// replaceAll replaces all matches
func (f *FindReplaceDialog) replaceAll() {
	searchText := f.findEntry.Text
	if searchText == "" {
		return
	}

	replaceText := f.replaceEntry.Text
	content := f.editor.GetText()

	var newContent string
	if f.matchCase.Checked {
		newContent = strings.ReplaceAll(content, searchText, replaceText)
	} else {
		// Case-insensitive replace
		newContent = caseInsensitiveReplaceAll(content, searchText, replaceText)
	}

	f.editor.SetText(newContent)
	f.updateMatches()
}

// caseInsensitiveReplaceAll performs case-insensitive replacement
func caseInsensitiveReplaceAll(content, search, replace string) string {
	lowerContent := strings.ToLower(content)
	lowerSearch := strings.ToLower(search)

	var result strings.Builder
	pos := 0
	for {
		idx := strings.Index(lowerContent[pos:], lowerSearch)
		if idx == -1 {
			result.WriteString(content[pos:])
			break
		}
		result.WriteString(content[pos : pos+idx])
		result.WriteString(replace)
		pos += idx + len(search)
	}
	return result.String()
}

// ShowFindReplaceDialog shows the find and replace dialog
func ShowFindReplaceDialog(window fyne.Window, editor *MarkdownEditor) {
	d := NewFindReplaceDialog(window, editor)
	d.Show()
}
