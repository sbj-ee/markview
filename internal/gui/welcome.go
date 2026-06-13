package gui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sbj-ee/markview/internal/themes"
)

// maxWelcomeRecentFiles caps how many recent files the welcome panel lists.
const maxWelcomeRecentFiles = 8

// buildWelcomeContent builds the empty-state panel shown when no document is
// open: quick actions, recent files, and a few key shortcuts.
func (w *Window) buildWelcomeContent() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(themes.AppLogo())
	logo.SetMinSize(fyne.NewSize(72, 72))
	logo.FillMode = canvas.ImageFillContain

	title := widget.NewLabelWithStyle("Welcome to MarkView", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabelWithStyle("Open a file or folder to get started", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	openFile := widget.NewButtonWithIcon("Open File…", theme.DocumentIcon(), w.showOpenDialog)
	openFile.Importance = widget.HighImportance
	openFolder := widget.NewButtonWithIcon("Open Folder…", theme.FolderOpenIcon(), w.showFolderDialog)
	actions := container.NewHBox(layout.NewSpacer(), openFile, openFolder, layout.NewSpacer())

	items := []fyne.CanvasObject{
		container.NewCenter(logo),
		title,
		subtitle,
		widget.NewSeparator(),
		actions,
	}

	if len(w.recentFiles) > 0 {
		items = append(items,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Recent", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)
		for i, p := range w.recentFiles {
			if i >= maxWelcomeRecentFiles {
				break
			}
			path := p // capture per iteration
			btn := widget.NewButtonWithIcon(filepath.Base(path), theme.FileTextIcon(), func() {
				w.checkUnsavedChanges(func() {
					w.loadFile(path)
				})
			})
			btn.Importance = widget.LowImportance
			btn.Alignment = widget.ButtonAlignLeading
			items = append(items, btn)
		}
	}

	hint := widget.NewLabelWithStyle(
		cmdKey+"+O  Open    Ctrl+P  Quick switch    "+cmdKey+"+Shift+P  Command palette",
		fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	items = append(items, widget.NewSeparator(), hint)

	return container.NewPadded(container.NewCenter(container.NewVBox(items...)))
}

// showWelcomeContent displays the empty-state panel in the main content area.
func (w *Window) showWelcomeContent() {
	if w.scrollContent == nil {
		return
	}
	w.scrollContent.Content = w.buildWelcomeContent()
	w.scrollContent.Refresh()
}
