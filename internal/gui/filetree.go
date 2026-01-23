package gui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	parentDirUID = ".."
)

// fileTreeNode is a custom widget for tree nodes with icon and label
type fileTreeNode struct {
	widget.BaseWidget
	icon        *widget.Icon
	label       *widget.Label
	box         *fyne.Container
	path        string // full path to the file/directory
	isDirectory bool
	isParentDir bool // true if this is the ".." entry
	onContext   func(path string, isDir bool, pos fyne.Position) // callback for context menu
}

func newFileTreeNode() *fileTreeNode {
	icon := widget.NewIcon(theme.FileIcon())
	label := widget.NewLabel("")
	n := &fileTreeNode{
		icon:  icon,
		label: label,
		box:   container.NewBorder(nil, nil, icon, nil, label),
	}
	n.ExtendBaseWidget(n)
	return n
}

func (n *fileTreeNode) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(n.box)
}

func (n *fileTreeNode) SetContent(icon fyne.Resource, text string, bold bool) {
	n.icon.SetResource(icon)
	n.label.SetText(text)
	if bold {
		n.label.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		n.label.TextStyle = fyne.TextStyle{}
	}
	n.label.Refresh()
}

// SetPath stores the path and directory state for context menu operations
func (n *fileTreeNode) SetPath(path string, isDir bool, isParent bool) {
	n.path = path
	n.isDirectory = isDir
	n.isParentDir = isParent
}

// SetContextCallback sets the callback for right-click context menu
func (n *fileTreeNode) SetContextCallback(callback func(path string, isDir bool, pos fyne.Position)) {
	n.onContext = callback
}

// TappedSecondary handles right-click to show context menu
func (n *fileTreeNode) TappedSecondary(e *fyne.PointEvent) {
	// Don't show context menu for parent directory ".."
	if n.isParentDir || n.path == "" {
		return
	}
	if n.onContext != nil {
		n.onContext(n.path, n.isDirectory, e.AbsolutePosition)
	}
}

// FileTree represents a file browser tree widget
type FileTree struct {
	tree         *widget.Tree
	rootPath     string
	currentFile  string // currently open file path
	onFileSelect func(path string)
	pathMap      map[string]string // uid -> path mapping
	filter       string            // current filter text
	filterEntry  *widget.Entry     // search entry widget
	container    *fyne.Container   // container with search + tree
	fyneWindow   fyne.Window       // reference to window for dialogs
	canvas       fyne.Canvas       // canvas for popup menus
	onFileChange func()            // callback when files/dirs are created/deleted/renamed
}

// NewFileTree creates a new file tree browser
func NewFileTree(onFileSelect func(path string)) *FileTree {
	ft := &FileTree{
		onFileSelect: onFileSelect,
		pathMap:      make(map[string]string),
		filter:       "",
	}
	ft.tree = ft.createTree()
	ft.filterEntry = ft.createFilterEntry()
	ft.container = container.NewBorder(ft.filterEntry, nil, nil, nil, ft.tree)
	return ft
}

// SetWindow sets the window reference for showing dialogs
func (ft *FileTree) SetWindow(w fyne.Window) {
	ft.fyneWindow = w
	ft.canvas = w.Canvas()
}

// SetOnFileChange sets a callback to be called when files/directories are modified
func (ft *FileTree) SetOnFileChange(callback func()) {
	ft.onFileChange = callback
}

// createFilterEntry creates the search/filter entry
func (ft *FileTree) createFilterEntry() *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Filter files...")
	entry.OnChanged = func(text string) {
		ft.filter = strings.ToLower(text)
		ft.tree.Refresh()
	}
	return entry
}

// SetRootPath sets the root directory for the file tree
func (ft *FileTree) SetRootPath(path string) {
	ft.rootPath = path
	ft.pathMap = make(map[string]string)
	ft.tree.Refresh()
}

// SetCurrentFile sets the currently open file for highlighting
func (ft *FileTree) SetCurrentFile(path string) {
	ft.currentFile = path
	ft.tree.Refresh()
}

// GetTree returns the underlying tree widget
func (ft *FileTree) GetTree() *widget.Tree {
	return ft.tree
}

// GetContainer returns the container with search and tree
func (ft *FileTree) GetContainer() fyne.CanvasObject {
	return ft.container
}

// NavigateUp navigates to the parent directory
func (ft *FileTree) NavigateUp() {
	if ft.rootPath == "" {
		return
	}
	parent := filepath.Dir(ft.rootPath)
	if parent != ft.rootPath { // Not at filesystem root
		ft.SetRootPath(parent)
	}
}

// createTree creates the file tree widget
func (ft *FileTree) createTree() *widget.Tree {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			// For virtual root "", use rootPath directly
			var path string
			if uid == "" {
				path = ft.rootPath
			} else if uid == parentDirUID {
				return []string{} // Parent dir has no children
			} else {
				path = ft.getPath(uid)
			}
			if path == "" {
				return []string{}
			}

			entries, err := os.ReadDir(path)
			if err != nil {
				return []string{}
			}

			var result []string

			// Add parent directory navigation at root level
			if uid == "" && ft.rootPath != "" {
				parent := filepath.Dir(ft.rootPath)
				if parent != ft.rootPath { // Not at filesystem root
					result = append(result, parentDirUID)
				}
			}

			var dirs []string
			var files []string

			for _, entry := range entries {
				name := entry.Name()
				// Skip hidden files
				if strings.HasPrefix(name, ".") {
					continue
				}

				// Apply filter
				if ft.filter != "" && !strings.Contains(strings.ToLower(name), ft.filter) {
					// For directories, check if any children match
					if entry.IsDir() {
						fullPath := filepath.Join(path, name)
						if !ft.hasMatchingFiles(fullPath, ft.filter) {
							continue
						}
					} else {
						continue
					}
				}

				fullPath := filepath.Join(path, name)
				childUID := ft.makeUID(uid, name)
				ft.pathMap[childUID] = fullPath

				if entry.IsDir() {
					// Check if directory contains any markdown files
					if ft.hasMarkdownFiles(fullPath) {
						dirs = append(dirs, childUID)
					}
				} else if isMarkdownFile(name) {
					files = append(files, childUID)
				}
			}

			// Sort directories and files alphabetically
			sort.Strings(dirs)
			sort.Strings(files)

			// Return parent dir first, then directories, then files
			result = append(result, dirs...)
			result = append(result, files...)
			return result
		},
		IsBranch: func(uid string) bool {
			// Virtual root "" is a branch but hidden
			if uid == "" {
				return true
			}
			// Parent dir is not a branch (no children shown)
			if uid == parentDirUID {
				return false
			}
			path := ft.getPath(uid)
			if path == "" {
				return false
			}
			info, err := os.Stat(path)
			if err != nil {
				return false
			}
			return info.IsDir()
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return newFileTreeNode()
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			treeNode := node.(*fileTreeNode)

			// Set up context menu callback
			treeNode.SetContextCallback(ft.showContextMenu)

			// Hide the root node - only show its children
			if uid == "" {
				treeNode.SetContent(theme.FolderIcon(), "", false)
				treeNode.SetPath("", false, false)
				return
			}

			// Parent directory navigation
			if uid == parentDirUID {
				treeNode.SetContent(theme.FolderOpenIcon(), "..", false)
				treeNode.SetPath("", false, true)
				return
			}

			path := ft.getPath(uid)
			if path == "" {
				treeNode.SetContent(theme.FileIcon(), "", false)
				treeNode.SetPath("", false, false)
				return
			}

			name := filepath.Base(path)
			isCurrent := path == ft.currentFile
			if branch {
				treeNode.SetContent(theme.FolderIcon(), name, false)
				treeNode.SetPath(path, true, false)
			} else {
				treeNode.SetContent(theme.DocumentIcon(), name, isCurrent)
				treeNode.SetPath(path, false, false)
			}
		},
		OnSelected: func(uid string) {
			// Handle parent directory navigation
			if uid == parentDirUID {
				ft.NavigateUp()
				ft.tree.UnselectAll()
				return
			}

			path := ft.getPath(uid)
			if path == "" {
				return
			}
			info, err := os.Stat(path)
			if err != nil {
				return
			}
			if !info.IsDir() && ft.onFileSelect != nil {
				ft.onFileSelect(path)
			}
		},
	}
	return tree
}

// getPath returns the file path for a given UID
func (ft *FileTree) getPath(uid string) string {
	// "" is virtual root with no path - use rootPath in ChildUIDs instead
	if uid == "" || uid == parentDirUID {
		return ""
	}
	if path, ok := ft.pathMap[uid]; ok {
		return path
	}
	return ""
}

// makeUID creates a unique ID for a tree node
func (ft *FileTree) makeUID(parentUID, name string) string {
	if parentUID == "" {
		return name
	}
	return parentUID + "/" + name
}

// hasMarkdownFiles recursively checks if a directory contains markdown files
func (ft *FileTree) hasMarkdownFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if entry.IsDir() {
			if ft.hasMarkdownFiles(filepath.Join(path, name)) {
				return true
			}
		} else if isMarkdownFile(name) {
			return true
		}
	}
	return false
}

// hasMatchingFiles checks if a directory contains files matching the filter
func (ft *FileTree) hasMatchingFiles(path string, filter string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if entry.IsDir() {
			if ft.hasMatchingFiles(filepath.Join(path, name), filter) {
				return true
			}
		} else if isMarkdownFile(name) && strings.Contains(strings.ToLower(name), filter) {
			return true
		}
	}
	return false
}

// isMarkdownFile checks if a filename has a markdown extension
func isMarkdownFile(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown")
}

// FocusFilter focuses the filter entry for keyboard input
func (ft *FileTree) FocusFilter(canvas fyne.Canvas) {
	canvas.Focus(ft.filterEntry)
}

// ClearFilter clears the filter and refreshes the tree
func (ft *FileTree) ClearFilter() {
	ft.filterEntry.SetText("")
	ft.filter = ""
	ft.tree.Refresh()
}

// showContextMenu displays the right-click context menu
func (ft *FileTree) showContextMenu(path string, isDir bool, pos fyne.Position) {
	if ft.fyneWindow == nil {
		return
	}

	var menuItems []*fyne.MenuItem

	if isDir {
		// Directory context menu
		menuItems = []*fyne.MenuItem{
			fyne.NewMenuItem("New File", func() {
				ft.showNewFileDialog(path)
			}),
			fyne.NewMenuItem("New Directory", func() {
				ft.showNewDirectoryDialog(path)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Rename Directory", func() {
				ft.showRenameDialog(path, true)
			}),
			fyne.NewMenuItem("Delete Directory", func() {
				ft.showDeleteConfirmation(path, true)
			}),
		}
	} else {
		// File context menu
		menuItems = []*fyne.MenuItem{
			fyne.NewMenuItem("Rename File", func() {
				ft.showRenameDialog(path, false)
			}),
			fyne.NewMenuItem("Delete File", func() {
				ft.showDeleteConfirmation(path, false)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("New File Here", func() {
				ft.showNewFileDialog(filepath.Dir(path))
			}),
			fyne.NewMenuItem("New Directory Here", func() {
				ft.showNewDirectoryDialog(filepath.Dir(path))
			}),
		}
	}

	menu := fyne.NewMenu("", menuItems...)
	popUp := widget.NewPopUpMenu(menu, ft.canvas)
	popUp.ShowAtPosition(pos)
}

// showNewDirectoryDialog shows a dialog to create a new directory
func (ft *FileTree) showNewDirectoryDialog(parentPath string) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Directory name")

	dialog.ShowForm("New Directory", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", entry),
		},
		func(confirmed bool) {
			if !confirmed || entry.Text == "" {
				return
			}
			newPath := filepath.Join(parentPath, entry.Text)
			if err := os.MkdirAll(newPath, 0755); err != nil {
				dialog.ShowError(err, ft.fyneWindow)
				return
			}
			ft.tree.Refresh()
			if ft.onFileChange != nil {
				ft.onFileChange()
			}
		},
		ft.fyneWindow,
	)
}

// showNewFileDialog shows a dialog to create a new file
func (ft *FileTree) showNewFileDialog(parentPath string) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("filename.md")

	dialog.ShowForm("New File", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", entry),
		},
		func(confirmed bool) {
			if !confirmed || entry.Text == "" {
				return
			}
			fileName := entry.Text
			// Add .md extension if not present
			if !strings.HasSuffix(strings.ToLower(fileName), ".md") && !strings.HasSuffix(strings.ToLower(fileName), ".markdown") {
				fileName = fileName + ".md"
			}
			newPath := filepath.Join(parentPath, fileName)

			// Check if file already exists
			if _, err := os.Stat(newPath); err == nil {
				dialog.ShowError(os.ErrExist, ft.fyneWindow)
				return
			}

			// Create empty file
			file, err := os.Create(newPath)
			if err != nil {
				dialog.ShowError(err, ft.fyneWindow)
				return
			}
			file.Close()

			ft.tree.Refresh()
			if ft.onFileChange != nil {
				ft.onFileChange()
			}

			// Open the newly created file
			if ft.onFileSelect != nil {
				ft.onFileSelect(newPath)
			}
		},
		ft.fyneWindow,
	)
}

// showRenameDialog shows a dialog to rename a file or directory
func (ft *FileTree) showRenameDialog(path string, isDir bool) {
	oldName := filepath.Base(path)
	entry := widget.NewEntry()
	entry.SetText(oldName)

	title := "Rename File"
	if isDir {
		title = "Rename Directory"
	}

	dialog.ShowForm(title, "Rename", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("New name", entry),
		},
		func(confirmed bool) {
			if !confirmed || entry.Text == "" || entry.Text == oldName {
				return
			}
			newPath := filepath.Join(filepath.Dir(path), entry.Text)

			// Check if target already exists
			if _, err := os.Stat(newPath); err == nil {
				dialog.ShowError(os.ErrExist, ft.fyneWindow)
				return
			}

			if err := os.Rename(path, newPath); err != nil {
				dialog.ShowError(err, ft.fyneWindow)
				return
			}

			// Update current file reference if we renamed the current file
			if ft.currentFile == path {
				ft.currentFile = newPath
			}

			ft.tree.Refresh()
			if ft.onFileChange != nil {
				ft.onFileChange()
			}
		},
		ft.fyneWindow,
	)
}

// showDeleteConfirmation shows a confirmation dialog before deleting
func (ft *FileTree) showDeleteConfirmation(path string, isDir bool) {
	name := filepath.Base(path)
	itemType := "file"
	if isDir {
		itemType = "directory"
	}

	dialog.ShowConfirm(
		"Confirm Delete",
		"Are you sure you want to delete the "+itemType+" \""+name+"\"?\n\nThis action cannot be undone.",
		func(confirmed bool) {
			if !confirmed {
				return
			}

			var err error
			if isDir {
				err = os.RemoveAll(path)
			} else {
				err = os.Remove(path)
			}

			if err != nil {
				dialog.ShowError(err, ft.fyneWindow)
				return
			}

			ft.tree.Refresh()
			if ft.onFileChange != nil {
				ft.onFileChange()
			}
		},
		ft.fyneWindow,
	)
}
