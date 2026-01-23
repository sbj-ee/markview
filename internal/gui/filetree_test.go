package gui

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestNewFileTree(t *testing.T) {
	var selectedPath string
	ft := NewFileTree(func(path string) {
		selectedPath = path
	})

	if ft == nil {
		t.Fatal("NewFileTree returned nil")
	}
	if ft.tree == nil {
		t.Error("FileTree.tree is nil")
	}
	if ft.pathMap == nil {
		t.Error("FileTree.pathMap is nil")
	}
	if ft.filterEntry == nil {
		t.Error("FileTree.filterEntry is nil")
	}
	if ft.container == nil {
		t.Error("FileTree.container is nil")
	}

	// Verify callback is stored (indirectly by checking it's not nil behavior)
	_ = selectedPath
}

func TestFileTree_SetRootPath(t *testing.T) {
	ft := NewFileTree(nil)

	// Create a temp directory
	tmpDir := t.TempDir()

	ft.SetRootPath(tmpDir)
	if ft.rootPath != tmpDir {
		t.Errorf("Expected rootPath %q, got %q", tmpDir, ft.rootPath)
	}
}

func TestFileTree_SetCurrentFile(t *testing.T) {
	ft := NewFileTree(nil)
	testPath := "/some/test/file.md"

	ft.SetCurrentFile(testPath)
	if ft.currentFile != testPath {
		t.Errorf("Expected currentFile %q, got %q", testPath, ft.currentFile)
	}
}

func TestFileTree_NavigateUp(t *testing.T) {
	ft := NewFileTree(nil)

	// Test with empty root path
	ft.NavigateUp()
	if ft.rootPath != "" {
		t.Errorf("Expected empty rootPath after NavigateUp with no root, got %q", ft.rootPath)
	}

	// Test with actual path
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)

	ft.SetRootPath(subDir)
	ft.NavigateUp()
	if ft.rootPath != tmpDir {
		t.Errorf("Expected rootPath %q after NavigateUp, got %q", tmpDir, ft.rootPath)
	}
}

func TestFileTree_GetTree(t *testing.T) {
	ft := NewFileTree(nil)
	tree := ft.GetTree()
	if tree == nil {
		t.Error("GetTree returned nil")
	}
	if tree != ft.tree {
		t.Error("GetTree returned different tree instance")
	}
}

func TestFileTree_GetContainer(t *testing.T) {
	ft := NewFileTree(nil)
	container := ft.GetContainer()
	if container == nil {
		t.Error("GetContainer returned nil")
	}
}

func TestFileTree_Filter(t *testing.T) {
	ft := NewFileTree(nil)

	// Set filter text
	ft.filterEntry.SetText("test")
	if ft.filter != "test" {
		t.Errorf("Expected filter %q, got %q", "test", ft.filter)
	}

	// Clear filter
	ft.ClearFilter()
	if ft.filter != "" {
		t.Errorf("Expected empty filter after ClearFilter, got %q", ft.filter)
	}
	if ft.filterEntry.Text != "" {
		t.Errorf("Expected empty filterEntry text after ClearFilter, got %q", ft.filterEntry.Text)
	}
}

func TestFileTree_SetWindow(t *testing.T) {
	ft := NewFileTree(nil)
	app := test.NewApp()
	w := app.NewWindow("Test")

	ft.SetWindow(w)
	if ft.fyneWindow != w {
		t.Error("SetWindow did not set fyneWindow")
	}
	if ft.canvas == nil {
		t.Error("SetWindow did not set canvas")
	}
}

func TestFileTree_SetOnFileChange(t *testing.T) {
	ft := NewFileTree(nil)
	called := false

	ft.SetOnFileChange(func() {
		called = true
	})

	if ft.onFileChange == nil {
		t.Error("SetOnFileChange did not set callback")
	}

	// Call it to verify
	ft.onFileChange()
	if !called {
		t.Error("onFileChange callback was not called")
	}
}

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"test.md", true},
		{"test.MD", true},
		{"test.markdown", true},
		{"test.MARKDOWN", true},
		{"test.txt", false},
		{"test.html", false},
		{"readme", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isMarkdownFile(tc.name)
			if result != tc.expected {
				t.Errorf("isMarkdownFile(%q) = %v, expected %v", tc.name, result, tc.expected)
			}
		})
	}
}

func TestFileTreeNode_SetContent(t *testing.T) {
	node := newFileTreeNode()

	node.SetContent(nil, "test.md", false)
	if node.label.Text != "test.md" {
		t.Errorf("Expected label text %q, got %q", "test.md", node.label.Text)
	}
	if node.label.TextStyle.Bold {
		t.Error("Expected non-bold text style")
	}

	node.SetContent(nil, "current.md", true)
	if !node.label.TextStyle.Bold {
		t.Error("Expected bold text style for current file")
	}
}

func TestFileTreeNode_SetPath(t *testing.T) {
	node := newFileTreeNode()

	node.SetPath("/test/path", true, false)
	if node.path != "/test/path" {
		t.Errorf("Expected path %q, got %q", "/test/path", node.path)
	}
	if !node.isDirectory {
		t.Error("Expected isDirectory to be true")
	}
	if node.isParentDir {
		t.Error("Expected isParentDir to be false")
	}

	node.SetPath("", false, true)
	if node.path != "" {
		t.Errorf("Expected empty path, got %q", node.path)
	}
	if node.isDirectory {
		t.Error("Expected isDirectory to be false")
	}
	if !node.isParentDir {
		t.Error("Expected isParentDir to be true")
	}
}

func TestFileTreeNode_SetContextCallback(t *testing.T) {
	node := newFileTreeNode()
	var calledPath string
	var calledIsDir bool

	node.SetContextCallback(func(path string, isDir bool, pos fyne.Position) {
		calledPath = path
		calledIsDir = isDir
	})

	if node.onContext == nil {
		t.Error("SetContextCallback did not set callback")
	}

	// Simulate callback
	node.onContext("/test/path", true, fyne.Position{})
	if calledPath != "/test/path" {
		t.Errorf("Expected path %q, got %q", "/test/path", calledPath)
	}
	if !calledIsDir {
		t.Error("Expected isDir to be true")
	}
}

func TestFileTreeNode_TappedSecondary(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		isParentDir  bool
		expectCalled bool
	}{
		{"normal file", "/test/file.md", false, true},
		{"normal directory", "/test/dir", false, true},
		{"parent dir entry", "", true, false},
		{"empty path", "", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := newFileTreeNode()
			called := false

			node.SetPath(tc.path, false, tc.isParentDir)
			node.SetContextCallback(func(path string, isDir bool, pos fyne.Position) {
				called = true
			})

			node.TappedSecondary(&fyne.PointEvent{
				AbsolutePosition: fyne.Position{X: 100, Y: 100},
			})

			if called != tc.expectCalled {
				t.Errorf("Expected callback called=%v, got called=%v", tc.expectCalled, called)
			}
		})
	}
}

func TestFileTree_HasMarkdownFiles(t *testing.T) {
	ft := NewFileTree(nil)

	// Create temp directory structure
	tmpDir := t.TempDir()

	// Empty directory should return false
	if ft.hasMarkdownFiles(tmpDir) {
		t.Error("Expected hasMarkdownFiles to return false for empty directory")
	}

	// Create a markdown file
	mdFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(mdFile, []byte("# Test"), 0644)

	if !ft.hasMarkdownFiles(tmpDir) {
		t.Error("Expected hasMarkdownFiles to return true when markdown file exists")
	}

	// Create subdirectory with markdown
	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	subMdFile := filepath.Join(subDir, "nested.md")
	os.WriteFile(subMdFile, []byte("# Nested"), 0644)

	if !ft.hasMarkdownFiles(tmpDir) {
		t.Error("Expected hasMarkdownFiles to return true for nested markdown files")
	}
}

func TestFileTree_HasMatchingFiles(t *testing.T) {
	ft := NewFileTree(nil)

	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create markdown files
	os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# Readme"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "notes.md"), []byte("# Notes"), 0644)

	// Test matching
	if !ft.hasMatchingFiles(tmpDir, "read") {
		t.Error("Expected hasMatchingFiles to return true for 'read' filter")
	}

	if !ft.hasMatchingFiles(tmpDir, "note") {
		t.Error("Expected hasMatchingFiles to return true for 'note' filter")
	}

	if ft.hasMatchingFiles(tmpDir, "xyz") {
		t.Error("Expected hasMatchingFiles to return false for 'xyz' filter")
	}
}

func TestFileTree_MakeUID(t *testing.T) {
	ft := NewFileTree(nil)

	// Root level UID
	uid := ft.makeUID("", "file.md")
	if uid != "file.md" {
		t.Errorf("Expected UID %q, got %q", "file.md", uid)
	}

	// Nested UID
	uid = ft.makeUID("parent", "child.md")
	if uid != "parent/child.md" {
		t.Errorf("Expected UID %q, got %q", "parent/child.md", uid)
	}
}
