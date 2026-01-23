package gui

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestIsInLinkDestination(t *testing.T) {
	la := &LinkAutocomplete{}

	tests := []struct {
		name        string
		content     string
		cursorPos   int
		wantInLink  bool
		wantPartial string
	}{
		{
			name:        "cursor after opening bracket",
			content:     "[link text](",
			cursorPos:   12,
			wantInLink:  true,
			wantPartial: "",
		},
		{
			name:        "cursor with partial path",
			content:     "[link text](test",
			cursorPos:   16,
			wantInLink:  true,
			wantPartial: "test",
		},
		{
			name:        "cursor with partial path containing slash",
			content:     "[link text](docs/readme",
			cursorPos:   23,
			wantInLink:  true,
			wantPartial: "docs/readme",
		},
		{
			name:        "not in link - no bracket",
			content:     "just some text",
			cursorPos:   14,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "not in link - closed link",
			content:     "[link text](file.md)",
			cursorPos:   20,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "not in link - cursor before bracket",
			content:     "[link text](file.md)",
			cursorPos:   5,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "not in link - http URL",
			content:     "[link](https://example.com",
			cursorPos:   26,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "not in link - http URL short",
			content:     "[link](http://",
			cursorPos:   14,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "not in link - mailto",
			content:     "[email](mailto:test",
			cursorPos:   19,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "not in link - anchor",
			content:     "[heading](#section",
			cursorPos:   18,
			wantInLink:  false,
			wantPartial: "",
		},
		{
			name:        "multiline - cursor in link on second line",
			content:     "First line\n[link](",
			cursorPos:   18,
			wantInLink:  true,
			wantPartial: "",
		},
		{
			name:        "multiline - cursor in link with partial",
			content:     "First line\n[link](doc",
			cursorPos:   21,
			wantInLink:  true,
			wantPartial: "doc",
		},
		{
			name:        "empty content",
			content:     "",
			cursorPos:   0,
			wantInLink:  false,
			wantPartial: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInLink, gotPartial, _ := la.isInLinkDestination(tt.content, tt.cursorPos)
			if gotInLink != tt.wantInLink {
				t.Errorf("isInLinkDestination() inLink = %v, want %v", gotInLink, tt.wantInLink)
			}
			if gotPartial != tt.wantPartial {
				t.Errorf("isInLinkDestination() partial = %q, want %q", gotPartial, tt.wantPartial)
			}
		})
	}
}

func TestFilterFiles(t *testing.T) {
	la := &LinkAutocomplete{
		files: []string{
			"README.md",
			"docs/guide.md",
			"docs/api.md",
			"notes/meeting.md",
			"archive/old-readme.md",
		},
	}

	tests := []struct {
		name        string
		partialPath string
		wantCount   int
		wantFirst   string
	}{
		{
			name:        "empty filter returns all",
			partialPath: "",
			wantCount:   5,
			wantFirst:   "README.md",
		},
		{
			name:        "filter by filename",
			partialPath: "readme",
			wantCount:   2, // README.md and archive/old-readme.md
			wantFirst:   "README.md",
		},
		{
			name:        "filter by directory",
			partialPath: "docs",
			wantCount:   2, // docs/guide.md and docs/api.md
			wantFirst:   "docs/guide.md",
		},
		{
			name:        "filter by partial name",
			partialPath: "gui",
			wantCount:   1, // docs/guide.md
			wantFirst:   "docs/guide.md",
		},
		{
			name:        "no matches",
			partialPath: "xyz",
			wantCount:   0,
			wantFirst:   "",
		},
		{
			name:        "case insensitive",
			partialPath: "README",
			wantCount:   2,
			wantFirst:   "README.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			la.filterFiles(tt.partialPath)
			if len(la.filtered) != tt.wantCount {
				t.Errorf("filterFiles() count = %d, want %d", len(la.filtered), tt.wantCount)
			}
			if tt.wantCount > 0 && la.filtered[0] != tt.wantFirst {
				t.Errorf("filterFiles() first = %q, want %q", la.filtered[0], tt.wantFirst)
			}
		})
	}
}

func TestGetCursorPositionInText(t *testing.T) {
	la := &LinkAutocomplete{}

	tests := []struct {
		name    string
		content string
		row     int
		col     int
		want    int
	}{
		{
			name:    "first position",
			content: "hello",
			row:     0,
			col:     0,
			want:    0,
		},
		{
			name:    "middle of first line",
			content: "hello",
			row:     0,
			col:     3,
			want:    3,
		},
		{
			name:    "end of first line",
			content: "hello",
			row:     0,
			col:     5,
			want:    5,
		},
		{
			name:    "start of second line",
			content: "hello\nworld",
			row:     1,
			col:     0,
			want:    6,
		},
		{
			name:    "middle of second line",
			content: "hello\nworld",
			row:     1,
			col:     3,
			want:    9,
		},
		{
			name:    "third line",
			content: "a\nb\nc",
			row:     2,
			col:     1,
			want:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := la.getCursorPositionInText(tt.content, tt.row, tt.col)
			if got != tt.want {
				t.Errorf("getCursorPositionInText() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSelectNextPrevious(t *testing.T) {
	la := &LinkAutocomplete{
		filtered:    []string{"a.md", "b.md", "c.md"},
		selectedIdx: 0,
		active:      true,
	}

	// Test SelectNext
	la.SelectNext()
	if la.selectedIdx != 1 {
		t.Errorf("SelectNext() selectedIdx = %d, want 1", la.selectedIdx)
	}

	la.SelectNext()
	if la.selectedIdx != 2 {
		t.Errorf("SelectNext() selectedIdx = %d, want 2", la.selectedIdx)
	}

	// Wrap around
	la.SelectNext()
	if la.selectedIdx != 0 {
		t.Errorf("SelectNext() wrap selectedIdx = %d, want 0", la.selectedIdx)
	}

	// Test SelectPrevious
	la.SelectPrevious()
	if la.selectedIdx != 2 {
		t.Errorf("SelectPrevious() wrap selectedIdx = %d, want 2", la.selectedIdx)
	}

	la.SelectPrevious()
	if la.selectedIdx != 1 {
		t.Errorf("SelectPrevious() selectedIdx = %d, want 1", la.selectedIdx)
	}
}

func TestSelectWhenNotActive(t *testing.T) {
	la := &LinkAutocomplete{
		filtered:    []string{"a.md", "b.md"},
		selectedIdx: 0,
		active:      false,
	}

	la.SelectNext()
	if la.selectedIdx != 0 {
		t.Errorf("SelectNext() when not active should not change selectedIdx")
	}

	la.SelectPrevious()
	if la.selectedIdx != 0 {
		t.Errorf("SelectPrevious() when not active should not change selectedIdx")
	}
}

func TestSelectWithEmptyList(t *testing.T) {
	la := &LinkAutocomplete{
		filtered:    []string{},
		selectedIdx: 0,
		active:      true,
	}

	la.SelectNext()
	if la.selectedIdx != 0 {
		t.Errorf("SelectNext() with empty list should not change selectedIdx")
	}

	la.SelectPrevious()
	if la.selectedIdx != 0 {
		t.Errorf("SelectPrevious() with empty list should not change selectedIdx")
	}
}

func TestIsActive(t *testing.T) {
	la := &LinkAutocomplete{active: false}
	if la.IsActive() {
		t.Error("IsActive() should return false when not active")
	}

	la.active = true
	if !la.IsActive() {
		t.Error("IsActive() should return true when active")
	}
}

func TestDismiss(t *testing.T) {
	la := &LinkAutocomplete{
		active:      true,
		selectedIdx: 5,
	}

	la.Dismiss()
	if la.active {
		t.Error("Dismiss() should set active to false")
	}
	if la.selectedIdx != 0 {
		t.Error("Dismiss() should reset selectedIdx to 0")
	}
}

func TestHandleKeyEvent(t *testing.T) {
	// Test keys that don't require editor access
	tests := []struct {
		name     string
		active   bool
		key      fyne.KeyName
		consumed bool
	}{
		{"down when active", true, fyne.KeyDown, true},
		{"up when active", true, fyne.KeyUp, true},
		{"escape when active", true, fyne.KeyEscape, true},
		{"other key when active", true, fyne.KeyA, false},
		{"down when not active", false, fyne.KeyDown, false},
		{"escape when not active", false, fyne.KeyEscape, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			la := &LinkAutocomplete{
				active:   tt.active,
				filtered: []string{"a.md", "b.md"},
			}
			event := &fyne.KeyEvent{Name: tt.key}
			got := la.HandleKeyEvent(event)
			if got != tt.consumed {
				t.Errorf("HandleKeyEvent() = %v, want %v", got, tt.consumed)
			}
		})
	}
}

func TestHandleKeyEventEnterTab(t *testing.T) {
	// Enter and Tab call AcceptSelection which needs an editor
	// Test that they return true (consumed) when active
	// but will hide popup when no editor is available
	la := &LinkAutocomplete{
		active:   true,
		filtered: []string{"a.md"},
	}

	// These should consume the event even though they can't complete
	enterEvent := &fyne.KeyEvent{Name: fyne.KeyReturn}
	if !la.HandleKeyEvent(enterEvent) {
		t.Error("HandleKeyEvent(Enter) should return true when active")
	}

	la.active = true // Reset since AcceptSelection hides
	tabEvent := &fyne.KeyEvent{Name: fyne.KeyTab}
	if !la.HandleKeyEvent(tabEvent) {
		t.Error("HandleKeyEvent(Tab) should return true when active")
	}
}

func TestFilterFilesLimit(t *testing.T) {
	// Create a list with more than 10 files
	files := make([]string, 15)
	for i := 0; i < 15; i++ {
		files[i] = "file" + string(rune('a'+i)) + ".md"
	}

	la := &LinkAutocomplete{files: files}
	la.filterFiles("")

	if len(la.filtered) > 10 {
		t.Errorf("filterFiles() should limit to 10 results, got %d", len(la.filtered))
	}
}

func TestSetRootPath(t *testing.T) {
	la := &LinkAutocomplete{rootPath: "/old/path"}
	la.SetRootPath("/new/path")
	if la.rootPath != "/new/path" {
		t.Errorf("SetRootPath() did not update rootPath")
	}
}

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		text    string
		pattern string
		want    bool
	}{
		{"hello", "hel", true},
		{"hello", "hlo", true},
		{"hello", "xyz", false},
		{"readme.md", "rdm", true},
		{"README.md", "RDM", true}, // matches R, D (from .mD), M (from .Md)
		{"docs/guide.md", "dg", true},
		{"", "", true},
		{"hello", "", true},
		{"", "a", false},
		{"abcdef", "ace", true},
		{"abcdef", "aec", false}, // out of order
	}

	for _, tt := range tests {
		t.Run(tt.text+"/"+tt.pattern, func(t *testing.T) {
			got := fuzzyMatch(tt.text, tt.pattern)
			if got != tt.want {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tt.text, tt.pattern, got, tt.want)
			}
		})
	}
}
