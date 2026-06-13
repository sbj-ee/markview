package gui

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Backlink represents a link from another document
type Backlink struct {
	SourcePath  string
	RelPath     string
	LineNumber  int
	LineContent string
	LinkText    string
}

// BacklinksPanel shows documents that link to the current document
type BacklinksPanel struct {
	rootPath    string
	currentFile string
	backlinks   []Backlink
	onSelect    func(path string, line int)
}

// NewBacklinksPanel creates a new backlinks panel
func NewBacklinksPanel(rootPath string, onSelect func(path string, line int)) *BacklinksPanel {
	return &BacklinksPanel{
		rootPath: rootPath,
		onSelect: onSelect,
	}
}

// SetCurrentFile sets the current file and finds backlinks
func (bp *BacklinksPanel) SetCurrentFile(filePath string) {
	bp.currentFile = filePath
	bp.findBacklinks()
}

// findBacklinks searches for links to the current file
func (bp *BacklinksPanel) findBacklinks() {
	bp.backlinks = []Backlink{}
	if bp.rootPath == "" || bp.currentFile == "" {
		return
	}

	// Get the current file's name and relative path
	currentName := filepath.Base(bp.currentFile)
	currentRel, _ := filepath.Rel(bp.rootPath, bp.currentFile)

	// Patterns to match markdown links
	// [text](file.md) or [text](./path/file.md) or [[file]] (wiki-style)
	linkPatterns := []*regexp.Regexp{
		regexp.MustCompile(`\[([^\]]*)\]\(([^)]+\.md)\)`), // Standard markdown
		regexp.MustCompile(`\[\[([^\]]+)\]\]`),            // Wiki-style
	}

	filepath.Walk(bp.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if !isMarkdownFile(path) || path == bp.currentFile {
			return nil
		}

		// Read file content
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(bp.rootPath, path)
		lines := strings.Split(string(data), "\n")

		for lineNum, line := range lines {
			for _, pattern := range linkPatterns {
				matches := pattern.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					var linkTarget string
					var linkText string

					if len(match) >= 3 {
						// Standard markdown: [text](target)
						linkText = match[1]
						linkTarget = match[2]
					} else if len(match) >= 2 {
						// Wiki-style: [[target]]
						linkTarget = match[1]
						linkText = match[1]
					}

					// Check if link points to current file
					linkTarget = strings.TrimPrefix(linkTarget, "./")
					linkBase := filepath.Base(linkTarget)

					if linkBase == currentName || linkTarget == currentRel || linkTarget == currentName {
						bp.backlinks = append(bp.backlinks, Backlink{
							SourcePath:  path,
							RelPath:     relPath,
							LineNumber:  lineNum + 1,
							LineContent: strings.TrimSpace(line),
							LinkText:    linkText,
						})
					}
				}
			}
		}
		return nil
	})
}

// GetBacklinks returns the current backlinks
func (bp *BacklinksPanel) GetBacklinks() []Backlink {
	return bp.backlinks
}

// Show displays the backlinks panel
func (bp *BacklinksPanel) Show(window fyne.Window) {
	if bp.currentFile == "" {
		dialog.ShowInformation("Backlinks", "No file is currently open.", window)
		return
	}

	bp.findBacklinks()

	if len(bp.backlinks) == 0 {
		dialog.ShowInformation("Backlinks",
			"No other documents link to this file.",
			window)
		return
	}

	var d dialog.Dialog

	list := widget.NewList(
		func() int { return len(bp.backlinks) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("filename.md:1", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("link context..."),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(bp.backlinks) {
				box := obj.(*fyne.Container)
				fileLabel := box.Objects[0].(*widget.Label)
				contentLabel := box.Objects[1].(*widget.Label)

				bl := bp.backlinks[id]
				fileLabel.SetText(bl.RelPath)

				// Truncate long lines
				content := bl.LineContent
				if len(content) > 80 {
					content = content[:80] + "..."
				}
				contentLabel.SetText(content)
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(bp.backlinks) {
			bl := bp.backlinks[id]
			d.Hide()
			if bp.onSelect != nil {
				bp.onSelect(bl.SourcePath, bl.LineNumber)
			}
		}
	}

	currentFileName := filepath.Base(bp.currentFile)
	title := widget.NewLabelWithStyle(
		currentFileName,
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	countLabel := widget.NewLabel(intToString(len(bp.backlinks)) + " documents link here")

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(500, 350))

	content := container.NewBorder(
		container.NewVBox(title, countLabel, widget.NewSeparator()),
		nil, nil, nil,
		scroll,
	)

	d = dialog.NewCustom("Backlinks", "Close", content, window)
	d.Resize(dialogSizeList)
	d.Show()
}

// ShowBacklinks shows the backlinks dialog
func ShowBacklinks(window fyne.Window, rootPath, currentFile string, onSelect func(path string, line int)) {
	bp := NewBacklinksPanel(rootPath, onSelect)
	bp.SetCurrentFile(currentFile)
	bp.Show(window)
}
