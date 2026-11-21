package markdown

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func TestFitColumnWidths_FitsUnchanged(t *testing.T) {
	got := fitColumnWidths([]float32{100, 200}, 500)
	if got[0] != 100 || got[1] != 200 {
		t.Errorf("fitColumnWidths = %v, want [100 200]", got)
	}
}

func TestFitColumnWidths_ShrinksToAvailable(t *testing.T) {
	got := fitColumnWidths([]float32{400, 600}, 500)
	var total float32
	for _, w := range got {
		total += w
	}
	if total > 500.01 {
		t.Errorf("total width = %f, want <= 500", total)
	}
	if got[0] >= got[1] {
		t.Errorf("proportions not kept: %v", got)
	}
}

func TestFitColumnWidths_RespectsMinimum(t *testing.T) {
	got := fitColumnWidths([]float32{1000, 80, 40}, 300)
	if got[1] < minTableColWidth {
		t.Errorf("column 1 = %f, want >= %f", got[1], minTableColWidth)
	}
	if got[2] != 40 {
		t.Errorf("column narrower than minimum should keep its width, got %f", got[2])
	}
}

func TestFitWidth_ScalesDownKeepingAspect(t *testing.T) {
	f := NewFitWidth(canvas.NewRectangle(nil), fyne.NewSize(800, 400))

	if s := f.scaledSize(1000); s != fyne.NewSize(800, 400) {
		t.Errorf("wide viewport: got %v, want natural size", s)
	}
	if s := f.scaledSize(400); s != fyne.NewSize(400, 200) {
		t.Errorf("narrow viewport: got %v, want 400x200", s)
	}
}

func TestFitWidth_MinWidthIsSmall(t *testing.T) {
	f := NewFitWidth(canvas.NewRectangle(nil), fyne.NewSize(800, 400))
	if w := f.MinSize().Width; w > minImageWidth {
		t.Errorf("MinSize().Width = %f, want <= %f", w, minImageWidth)
	}
}
