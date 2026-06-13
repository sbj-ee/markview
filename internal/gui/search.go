package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SearchResult represents a search result
type SearchResult struct {
	FilePath    string
	RelPath     string
	LineNumber  int
	LineContent string
	MatchStart  int
	MatchEnd    int
}

// FullTextSearch searches for text across all markdown files
type FullTextSearch struct {
	rootPath string
	results  []SearchResult
	onSelect func(path string, line int)
}

// NewFullTextSearch creates a new full-text search
func NewFullTextSearch(rootPath string, onSelect func(path string, line int)) *FullTextSearch {
	return &FullTextSearch{
		rootPath: rootPath,
		onSelect: onSelect,
	}
}

// Search searches for the query in all markdown files
func (fts *FullTextSearch) Search(query string) []SearchResult {
	fts.results = []SearchResult{}
	if fts.rootPath == "" || query == "" {
		return fts.results
	}

	queryLower := strings.ToLower(query)

	filepath.Walk(fts.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if !isMarkdownFile(path) {
			return nil
		}

		// Read file content
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(fts.rootPath, path)
		lines := strings.Split(string(data), "\n")

		for lineNum, line := range lines {
			lineLower := strings.ToLower(line)
			if idx := strings.Index(lineLower, queryLower); idx != -1 {
				// Trim long lines for display
				displayLine := line
				if len(displayLine) > 100 {
					start := idx - 30
					if start < 0 {
						start = 0
					}
					end := idx + len(query) + 50
					if end > len(displayLine) {
						end = len(displayLine)
					}
					displayLine = "..." + displayLine[start:end]
					if end < len(line) {
						displayLine += "..."
					}
				}

				fts.results = append(fts.results, SearchResult{
					FilePath:    path,
					RelPath:     relPath,
					LineNumber:  lineNum + 1,
					LineContent: strings.TrimSpace(displayLine),
					MatchStart:  idx,
					MatchEnd:    idx + len(query),
				})
			}
		}
		return nil
	})

	return fts.results
}

// Show displays the search dialog
func (fts *FullTextSearch) Show(window fyne.Window) {
	if fts.rootPath == "" {
		dialog.ShowInformation("Search", "Please open a folder first.", window)
		return
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search in all files...")

	var d dialog.Dialog
	var list *widget.List
	var results []SearchResult

	statusLabel := widget.NewLabel("")

	list = widget.NewList(
		func() int { return len(results) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("filename.md:1", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("matching line content..."),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(results) {
				box := obj.(*fyne.Container)
				fileLabel := box.Objects[0].(*widget.Label)
				contentLabel := box.Objects[1].(*widget.Label)

				r := results[id]
				fileLabel.SetText(fmt.Sprintf("%s:%d", r.RelPath, r.LineNumber))
				contentLabel.SetText(r.LineContent)
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(results) {
			r := results[id]
			d.Hide()
			if fts.onSelect != nil {
				fts.onSelect(r.FilePath, r.LineNumber)
			}
		}
	}

	// Debounce search
	var searchTimer *fyne.Animation
	searchEntry.OnChanged = func(query string) {
		if searchTimer != nil {
			searchTimer.Stop()
		}
		if query == "" {
			results = []SearchResult{}
			statusLabel.SetText("")
			list.Refresh()
			return
		}
		// Simple delay before searching
		results = fts.Search(query)
		statusLabel.SetText(fmt.Sprintf("%d results found", len(results)))
		list.Refresh()
	}

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	content := container.NewBorder(
		container.NewVBox(searchEntry, statusLabel),
		nil, nil, nil,
		scroll,
	)

	d = dialog.NewCustom("Search in Files (Ctrl+Shift+F)", "Close", content, window)
	d.Resize(dialogSizeList)
	d.Show()

	window.Canvas().Focus(searchEntry)
}

// ShowFullTextSearch shows the full-text search dialog
func ShowFullTextSearch(window fyne.Window, rootPath string, onSelect func(path string, line int)) {
	fts := NewFullTextSearch(rootPath, onSelect)
	fts.Show(window)
}
