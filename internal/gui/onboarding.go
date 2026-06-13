package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sbj-ee/markview/internal/themes"
)

// onboardingPrefKey marks that the first-run onboarding has been shown.
const onboardingPrefKey = "onboardingCompleted"

// shouldShowOnboarding reports whether the first-run onboarding should be
// displayed, recording that it has been shown so it only appears once.
func (w *Window) shouldShowOnboarding() bool {
	if w.app.Preferences().Bool(onboardingPrefKey) {
		return false
	}
	w.app.Preferences().SetBool(onboardingPrefKey, true)
	return true
}

// maybeShowOnboarding shows the first-run onboarding dialog the first time the
// app is launched.
func (w *Window) maybeShowOnboarding() {
	if w.shouldShowOnboarding() {
		w.showOnboarding()
	}
}

// showOnboarding presents a one-time welcome dialog introducing the key ways to
// get around the app.
func (w *Window) showOnboarding() {
	logo := canvas.NewImageFromResource(themes.AppLogo())
	logo.SetMinSize(fyne.NewSize(64, 64))
	logo.FillMode = canvas.ImageFillContain

	heading := widget.NewLabelWithStyle("Welcome to MarkView", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	intro := widget.NewLabelWithStyle(
		"A fast, native Markdown viewer and editor. A few ways to get around:",
		fyne.TextAlignCenter, fyne.TextStyle{})
	intro.Wrapping = fyne.TextWrapWord

	tips := widget.NewRichTextFromMarkdown(
		"- **Command Palette** — press **" + cmdKey + "+Shift+P** to run any command by name\n" +
			"- **Menus** — every feature lives in the menu bar\n" +
			"- **Edit & preview** — **" + cmdKey + "+E** to edit, **" + cmdKey + "+\\\\** for side-by-side preview\n" +
			"- **Jump around** — **Ctrl+P** to open a file, **" + cmdKey + "+Shift+G** to search all files\n" +
			"- **Appearance** — change theme and font from the palette icon")
	tips.Wrapping = fyne.TextWrapWord

	var d dialog.Dialog
	openFolderBtn := widget.NewButtonWithIcon("Open a Folder", theme.FolderOpenIcon(), func() {
		d.Hide()
		w.showFolderDialog()
	})
	openFolderBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		container.NewCenter(logo),
		heading,
		intro,
		widget.NewSeparator(),
		tips,
		widget.NewSeparator(),
		container.NewCenter(openFolderBtn),
	)

	d = dialog.NewCustom("Welcome to MarkView", "Get Started", content, w.fyneWindow)
	d.Resize(fyne.NewSize(540, 520))
	d.Show()
}
