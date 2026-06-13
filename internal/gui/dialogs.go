package gui

import "fyne.io/fyne/v2"

// Shared dialog sizes, used so that similar dialogs open at a consistent size.
//
// Labeling convention for the dismiss button:
//   - single-button dialogs (navigators, pickers, viewers, settings) use "Close"
//   - two-button forms use an action label ("Create", "Save", …) plus "Cancel"
//
// Presentation mode is an intentional exception and uses "Exit".
var (
	// dialogSizeList is for search/list/result dialogs: a search field or tabs
	// above a scrolling list (quick switcher, command palette, search, etc.).
	dialogSizeList = fyne.NewSize(600, 520)

	// dialogSizeForm is for small input forms with a few fields.
	dialogSizeForm = fyne.NewSize(420, 170)

	// dialogSizeFileChooser is for native file/folder open and save dialogs.
	dialogSizeFileChooser = fyne.NewSize(900, 650)
)
