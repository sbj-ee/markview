package gui

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version string
		want    []int
	}{
		{"1.0.0", []int{1, 0, 0}},
		{"v1.0.0", []int{1, 0, 0}},
		{"1.2.3", []int{1, 2, 3}},
		{"2.0", []int{2, 0}},
		{"1.0.0-beta", []int{1, 0, 0}},
		{"1.0.0-rc1", []int{1, 0, 0}}, // "rc1" is not parsed as a number
		{"", []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := parseVersion(tt.version)
			if len(got) != len(tt.want) {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.version, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseVersion(%q) = %v, want %v", tt.version, got, tt.want)
					break
				}
			}
		})
	}
}

func TestIsNewerVersion(t *testing.T) {
	uc := &UpdateChecker{}

	tests := []struct {
		latest  string
		current string
		want    bool
	}{
		{"1.0.1", "1.0.0", true},
		{"1.1.0", "1.0.0", true},
		{"2.0.0", "1.9.9", true},
		{"1.0.0", "1.0.0", false},
		{"1.0.0", "1.0.1", false},
		{"1.0.0", "2.0.0", false},
		{"v1.0.1", "v1.0.0", true},
		{"1.0.4", "1.0.3", true},
		{"1.0.13", "1.0.12", true},
		{"1.0.0", "1.0.0-beta", false}, // Same major.minor.patch
	}

	for _, tt := range tests {
		name := tt.latest + "_vs_" + tt.current
		t.Run(name, func(t *testing.T) {
			got := uc.isNewerVersion(tt.latest, tt.current)
			if got != tt.want {
				t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}

func TestNewUpdateChecker(t *testing.T) {
	uc := NewUpdateChecker("1.0.0", nil, nil)
	if uc == nil {
		t.Error("NewUpdateChecker returned nil")
	}
	if uc.currentVersion != "1.0.0" {
		t.Errorf("currentVersion = %q, want %q", uc.currentVersion, "1.0.0")
	}
}
