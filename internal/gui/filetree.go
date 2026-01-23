package gui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	parentDirUID = ".."
)

// fileTreeNode is a custom widget for tree nodes with icon and label
type fileTreeNode struct {
	widget.BaseWidget
	icon  *widget.Icon
	label *widget.Label
	box   *fyne.Container
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

			// Hide the root node - only show its children
			if uid == "" {
				treeNode.SetContent(theme.FolderIcon(), "", false)
				return
			}

			// Parent directory navigation
			if uid == parentDirUID {
				treeNode.SetContent(theme.FolderOpenIcon(), "..", false)
				return
			}

			path := ft.getPath(uid)
			if path == "" {
				treeNode.SetContent(theme.FileIcon(), "", false)
				return
			}

			name := filepath.Base(path)
			isCurrent := path == ft.currentFile
			if branch {
				treeNode.SetContent(theme.FolderIcon(), name, false)
			} else {
				treeNode.SetContent(theme.DocumentIcon(), name, isCurrent)
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
