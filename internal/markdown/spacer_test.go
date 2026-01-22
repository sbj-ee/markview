package markdown

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func init() {
	// Initialize a test app for Fyne widget tests
	test.NewApp()
}

func TestNewSpacer(t *testing.T) {
	spacer := NewSpacer(16)

	if spacer == nil {
		t.Fatal("NewSpacer() returned nil")
	}

	if spacer.height != 16 {
		t.Errorf("Spacer.height = %f, want 16", spacer.height)
	}
}

func TestSpacer_MinSize(t *testing.T) {
	tests := []struct {
		name       string
		height     float32
		wantWidth  float32
		wantHeight float32
	}{
		{
			name:       "small spacer",
			height:     8,
			wantWidth:  1,
			wantHeight: 8,
		},
		{
			name:       "medium spacer",
			height:     16,
			wantWidth:  1,
			wantHeight: 16,
		},
		{
			name:       "large spacer",
			height:     32,
			wantWidth:  1,
			wantHeight: 32,
		},
		{
			name:       "zero spacer",
			height:     0,
			wantWidth:  1,
			wantHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spacer := NewSpacer(tt.height)
			size := spacer.MinSize()

			if size.Width != tt.wantWidth {
				t.Errorf("MinSize().Width = %f, want %f", size.Width, tt.wantWidth)
			}
			if size.Height != tt.wantHeight {
				t.Errorf("MinSize().Height = %f, want %f", size.Height, tt.wantHeight)
			}
		})
	}
}

func TestSpacer_CreateRenderer(t *testing.T) {
	spacer := NewSpacer(16)
	renderer := spacer.CreateRenderer()

	if renderer == nil {
		t.Fatal("CreateRenderer() returned nil")
	}
}

func TestSpacerRenderer_MinSize(t *testing.T) {
	spacer := NewSpacer(24)
	renderer := spacer.CreateRenderer()

	size := renderer.MinSize()
	if size.Height != 24 {
		t.Errorf("Renderer MinSize().Height = %f, want 24", size.Height)
	}
}

func TestSpacerRenderer_Objects(t *testing.T) {
	spacer := NewSpacer(16)
	renderer := spacer.CreateRenderer()

	objects := renderer.Objects()
	if len(objects) != 1 {
		t.Errorf("Objects() returned %d objects, want 1", len(objects))
	}

	if objects[0] == nil {
		t.Error("Objects()[0] is nil")
	}
}

func TestSpacerRenderer_Layout(t *testing.T) {
	spacer := NewSpacer(16)
	renderer := spacer.CreateRenderer()

	// Layout should not panic
	renderer.Layout(fyne.NewSize(100, 100))
}

func TestSpacerRenderer_Refresh(t *testing.T) {
	spacer := NewSpacer(16)
	renderer := spacer.CreateRenderer()

	// Refresh should not panic
	renderer.Refresh()
}

func TestSpacerRenderer_Destroy(t *testing.T) {
	spacer := NewSpacer(16)
	renderer := spacer.CreateRenderer()

	// Destroy should not panic
	renderer.Destroy()
}
