package gui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// QuickSwitcher provides a fuzzy file finder dialog
type QuickSwitcher struct {
	rootPath string
	files    []string
	onSelect func(path string)
}

// NewQuickSwitcher creates a new quick switcher
func NewQuickSwitcher(rootPath string, onSelect func(path string)) *QuickSwitcher {
	qs := &QuickSwitcher{
		rootPath: rootPath,
		onSelect: onSelect,
	}
	qs.scanFiles()
	return qs
}

// scanFiles scans the root path for markdown files
func (qs *QuickSwitcher) scanFiles() {
	qs.files = []string{}
	if qs.rootPath == "" {
		return
	}

	filepath.Walk(qs.rootPath, func(path string, info os.FileInfo, err error) error {
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
			relPath, _ := filepath.Rel(qs.rootPath, path)
			qs.files = append(qs.files, relPath)
		}
		return nil
	})

	// Sort by filename
	sort.Slice(qs.files, func(i, j int) bool {
		return strings.ToLower(filepath.Base(qs.files[i])) < strings.ToLower(filepath.Base(qs.files[j]))
	})
}

// Show displays the quick switcher dialog
func (qs *QuickSwitcher) Show(window fyne.Window) {
	if qs.rootPath == "" {
		dialog.ShowInformation("Quick Switcher", "Please open a folder first.", window)
		return
	}

	// Rescan files in case they changed
	qs.scanFiles()

	if len(qs.files) == 0 {
		dialog.ShowInformation("Quick Switcher", "No markdown files found in the current folder.", window)
		return
	}

	// Create search entry
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Type to search files...")

	// Create results list
	filteredFiles := make([]string, len(qs.files))
	copy(filteredFiles, qs.files)

	var d dialog.Dialog
	var list *widget.List

	list = widget.NewList(
		func() int { return len(filteredFiles) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(nil),
				widget.NewLabel("filename.md"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(filteredFiles) {
				box := obj.(*fyne.Container)
				label := box.Objects[1].(*widget.Label)
				label.SetText(filteredFiles[id])
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(filteredFiles) {
			fullPath := filepath.Join(qs.rootPath, filteredFiles[id])
			d.Hide()
			if qs.onSelect != nil {
				qs.onSelect(fullPath)
			}
		}
	}

	// Filter function
	filterFiles := func(query string) {
		query = strings.ToLower(query)
		if query == "" {
			filteredFiles = make([]string, len(qs.files))
			copy(filteredFiles, qs.files)
		} else {
			filteredFiles = []string{}
			for _, f := range qs.files {
				if fuzzyMatch(strings.ToLower(f), query) {
					filteredFiles = append(filteredFiles, f)
				}
			}
		}
		list.Refresh()
	}

	searchEntry.OnChanged = filterFiles

	// Layout
	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	content := container.NewBorder(
		searchEntry,
		nil, nil, nil,
		scroll,
	)

	d = dialog.NewCustom("Quick Switcher (Ctrl+P)", "Cancel", content, window)
	d.Resize(fyne.NewSize(550, 500))
	d.Show()

	// Focus the search entry
	window.Canvas().Focus(searchEntry)
}

// fuzzyMatch performs a simple fuzzy match
func fuzzyMatch(text, pattern string) bool {
	// Simple contains match for now
	if strings.Contains(text, pattern) {
		return true
	}

	// Also try matching each character in sequence
	patternIdx := 0
	for _, c := range text {
		if patternIdx < len(pattern) && byte(c) == pattern[patternIdx] {
			patternIdx++
		}
	}
	return patternIdx == len(pattern)
}

// ShowQuickSwitcher shows the quick switcher dialog
func ShowQuickSwitcher(window fyne.Window, rootPath string, onSelect func(path string)) {
	qs := NewQuickSwitcher(rootPath, onSelect)
	qs.Show(window)
}
