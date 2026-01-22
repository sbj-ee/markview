package gui

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sbj-ee/markview/internal/library"
)

// LibraryView provides a view of the document library
type LibraryView struct {
	widget.BaseWidget
	library       *library.DocumentLibrary
	container     *fyne.Container
	onFileSelect  func(path string)
	categoryList  *widget.List
	tagList       *widget.List
	documentList  *widget.List
	searchEntry   *widget.Entry
	currentFilter string
	filterType    string // "all", "category", "tag", "search"
	filteredDocs  []*library.Document
}

// NewLibraryView creates a new library view
func NewLibraryView(onFileSelect func(path string)) *LibraryView {
	lv := &LibraryView{
		onFileSelect: onFileSelect,
		filterType:   "all",
	}
	lv.ExtendBaseWidget(lv)
	lv.buildUI()
	return lv
}

// buildUI builds the library view UI
func (lv *LibraryView) buildUI() {
	// Search entry
	lv.searchEntry = widget.NewEntry()
	lv.searchEntry.SetPlaceHolder("Search documents...")
	lv.searchEntry.OnChanged = func(query string) {
		if query == "" {
			lv.filterType = "all"
			lv.currentFilter = ""
		} else {
			lv.filterType = "search"
			lv.currentFilter = query
		}
		lv.updateFilteredDocs()
	}

	// Category list
	lv.categoryList = widget.NewList(
		func() int {
			if lv.library == nil {
				return 0
			}
			return len(lv.library.GetCategories()) + 1 // +1 for "All"
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Category")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id == 0 {
				label.SetText("All Documents")
			} else if lv.library != nil {
				categories := lv.library.GetCategories()
				sort.Strings(categories)
				if id-1 < len(categories) {
					label.SetText(categories[id-1])
				}
			}
		},
	)
	lv.categoryList.OnSelected = func(id widget.ListItemID) {
		if id == 0 {
			lv.filterType = "all"
			lv.currentFilter = ""
		} else if lv.library != nil {
			categories := lv.library.GetCategories()
			sort.Strings(categories)
			if id-1 < len(categories) {
				lv.filterType = "category"
				lv.currentFilter = categories[id-1]
			}
		}
		lv.tagList.UnselectAll()
		lv.updateFilteredDocs()
	}

	// Tag list
	lv.tagList = widget.NewList(
		func() int {
			if lv.library == nil {
				return 0
			}
			return len(lv.library.GetTags())
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Tag")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if lv.library != nil {
				tags := lv.library.GetTags()
				sort.Strings(tags)
				if id < len(tags) {
					label.SetText("#" + tags[id])
				}
			}
		},
	)
	lv.tagList.OnSelected = func(id widget.ListItemID) {
		if lv.library != nil {
			tags := lv.library.GetTags()
			sort.Strings(tags)
			if id < len(tags) {
				lv.filterType = "tag"
				lv.currentFilter = tags[id]
			}
		}
		lv.categoryList.UnselectAll()
		lv.updateFilteredDocs()
	}

	// Document list
	lv.documentList = widget.NewList(
		func() int {
			return len(lv.filteredDocs)
		},
		func() fyne.CanvasObject {
			title := widget.NewLabel("Title")
			title.TextStyle = fyne.TextStyle{Bold: true}
			preview := widget.NewLabel("Preview text here...")
			preview.Wrapping = fyne.TextWrapWord
			meta := widget.NewLabel("Category | Tags | Words")
			meta.TextStyle = fyne.TextStyle{Italic: true}
			return container.NewVBox(title, preview, meta)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(lv.filteredDocs) {
				return
			}
			doc := lv.filteredDocs[id]
			box := obj.(*fyne.Container)

			title := box.Objects[0].(*widget.Label)
			title.SetText(doc.Title)

			preview := box.Objects[1].(*widget.Label)
			previewText := doc.Preview
			if len(previewText) > 100 {
				previewText = previewText[:100] + "..."
			}
			preview.SetText(previewText)

			meta := box.Objects[2].(*widget.Label)
			tagsStr := ""
			if len(doc.Tags) > 0 {
				tagsStr = " | #" + fmt.Sprintf("%v", doc.Tags)
			}
			meta.SetText(fmt.Sprintf("%s%s | %d words", doc.Category, tagsStr, doc.WordCount))
		},
	)
	lv.documentList.OnSelected = func(id widget.ListItemID) {
		if id < len(lv.filteredDocs) && lv.onFileSelect != nil {
			lv.onFileSelect(lv.filteredDocs[id].Path)
		}
	}

	// Create sidebar with categories and tags
	categoryHeader := widget.NewLabel("Categories")
	categoryHeader.TextStyle = fyne.TextStyle{Bold: true}
	tagHeader := widget.NewLabel("Tags")
	tagHeader.TextStyle = fyne.TextStyle{Bold: true}

	sidebar := container.NewVBox(
		categoryHeader,
		container.NewVScroll(lv.categoryList),
		widget.NewSeparator(),
		tagHeader,
		container.NewVScroll(lv.tagList),
	)

	// Create main content area
	docHeader := widget.NewLabel("Documents")
	docHeader.TextStyle = fyne.TextStyle{Bold: true}

	mainContent := container.NewBorder(
		container.NewVBox(lv.searchEntry, docHeader),
		nil, nil, nil,
		lv.documentList,
	)

	// Combine sidebar and main content
	lv.container = container.NewBorder(
		nil, nil,
		container.NewVScroll(sidebar),
		nil,
		mainContent,
	)
}

// SetLibrary sets the document library
func (lv *LibraryView) SetLibrary(lib *library.DocumentLibrary) {
	lv.library = lib
	lv.filterType = "all"
	lv.currentFilter = ""
	lv.updateFilteredDocs()
	lv.categoryList.Refresh()
	lv.tagList.Refresh()
}

// updateFilteredDocs updates the filtered document list
func (lv *LibraryView) updateFilteredDocs() {
	if lv.library == nil {
		lv.filteredDocs = nil
		lv.documentList.Refresh()
		return
	}

	switch lv.filterType {
	case "category":
		lv.filteredDocs = lv.library.FilterByCategory(lv.currentFilter)
	case "tag":
		lv.filteredDocs = lv.library.FilterByTag(lv.currentFilter)
	case "search":
		lv.filteredDocs = lv.library.Search(lv.currentFilter)
	default:
		lv.filteredDocs = lv.library.Documents
	}

	// Sort by modification time (newest first)
	sort.Slice(lv.filteredDocs, func(i, j int) bool {
		return lv.filteredDocs[i].ModTime.After(lv.filteredDocs[j].ModTime)
	})

	lv.documentList.Refresh()
}

// Refresh refreshes the library from disk
func (lv *LibraryView) Refresh() {
	if lv.library != nil {
		lv.library.Scan()
		lv.updateFilteredDocs()
		lv.categoryList.Refresh()
		lv.tagList.Refresh()
	}
}

// CreateRenderer implements fyne.Widget
func (lv *LibraryView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(lv.container)
}

// GetContainer returns the library view container
func (lv *LibraryView) GetContainer() *fyne.Container {
	return lv.container
}
