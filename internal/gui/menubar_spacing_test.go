package gui

import (
	"runtime"
	"strings"
	"testing"
)

// TestMenuTitle verifies titles are padded for menu-bar spacing off macOS (and
// left plain on macOS), and that trimming recovers the original name — which
// the command palette relies on.
func TestMenuTitle(t *testing.T) {
	got := menuTitle("File")

	if runtime.GOOS == "darwin" {
		if got != "File" {
			t.Errorf("on macOS the title must be plain, got %q", got)
		}
		return
	}

	if got == "File" {
		t.Errorf("expected padding around the title off macOS, got %q", got)
	}
	if strings.TrimSpace(got) != "File" {
		t.Errorf("trimming the padded title should recover the name, got %q -> %q", got, strings.TrimSpace(got))
	}
}
