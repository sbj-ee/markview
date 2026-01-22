package gui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// FileTree represents a file browser tree widget
type FileTree struct {
	tree        *widget.Tree
	rootPath    string
	onFileSelect func(path string)
	pathMap     map[string]string // uid -> path mapping
}

// NewFileTree creates a new file tree browser
func NewFileTree(onFileSelect func(path string)) *FileTree {
	ft := &FileTree{
		onFileSelect: onFileSelect,
		pathMap:      make(map[string]string),
	}
	ft.tree = ft.createTree()
	return ft
}

// SetRootPath sets the root directory for the file tree
func (ft *FileTree) SetRootPath(path string) {
	ft.rootPath = path
	ft.pathMap = make(map[string]string)
	// Don't map "" to path - "" is a virtual hidden root
	ft.tree.Refresh()
}

// GetTree returns the underlying tree widget
func (ft *FileTree) GetTree() *widget.Tree {
	return ft.tree
}

// createTree creates the file tree widget
func (ft *FileTree) createTree() *widget.Tree {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			// For virtual root "", use rootPath directly
			var path string
			if uid == "" {
				path = ft.rootPath
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

			var dirs []string
			var files []string

			for _, entry := range entries {
				name := entry.Name()
				// Skip hidden files
				if strings.HasPrefix(name, ".") {
					continue
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

			// Return directories first, then files
			return append(dirs, files...)
		},
		IsBranch: func(uid string) bool {
			// Virtual root "" is a branch but hidden
			if uid == "" {
				return true
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
			return widget.NewLabel("")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			label := node.(*widget.Label)
			// Hide the root node - only show its children
			if uid == "" {
				label.SetText("")
				return
			}
			path := ft.getPath(uid)
			if path == "" {
				label.SetText("")
				return
			}
			name := filepath.Base(path)
			label.SetText(name)
			label.Truncation = fyne.TextTruncateEllipsis
		},
		OnSelected: func(uid string) {
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
	if uid == "" {
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

// isMarkdownFile checks if a filename has a markdown extension
func isMarkdownFile(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown")
}

// IconFolder returns the folder icon
func IconFolder() fyne.Resource {
	return theme.FolderIcon()
}

// IconFile returns the file icon
func IconFile() fyne.Resource {
	return theme.FileIcon()
}
