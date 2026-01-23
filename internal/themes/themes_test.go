package themes

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestThemeNames(t *testing.T) {
	names := ThemeNames()
	if len(names) != 8 {
		t.Errorf("ThemeNames() returned %d names, want 8", len(names))
	}

	expected := []string{"Light", "Dark", "Nord", "Solarized Light", "Solarized Dark", "Monokai", "Gruvbox Dark", "One Dark"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("ThemeNames()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestThemeFromName(t *testing.T) {
	tests := []struct {
		name string
		want ThemeType
	}{
		{"Light", ThemeLight},
		{"Dark", ThemeDark},
		{"Nord", ThemeNord},
		{"Solarized Light", ThemeSolarizedLight},
		{"Solarized Dark", ThemeSolarizedDark},
		{"Monokai", ThemeMonokai},
		{"Gruvbox Dark", ThemeGruvboxDark},
		{"One Dark", ThemeOneDark},
		{"Unknown", ThemeDark}, // Default
		{"", ThemeDark},        // Empty string defaults to Dark
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThemeFromName(tt.name)
			if got != tt.want {
				t.Errorf("ThemeFromName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestThemeType_Name(t *testing.T) {
	tests := []struct {
		theme ThemeType
		want  string
	}{
		{ThemeLight, "Light"},
		{ThemeDark, "Dark"},
		{ThemeNord, "Nord"},
		{ThemeSolarizedLight, "Solarized Light"},
		{ThemeSolarizedDark, "Solarized Dark"},
		{ThemeMonokai, "Monokai"},
		{ThemeGruvboxDark, "Gruvbox Dark"},
		{ThemeOneDark, "One Dark"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.theme.Name()
			if got != tt.want {
				t.Errorf("ThemeType.Name() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestThemeType_IsDark(t *testing.T) {
	tests := []struct {
		theme  ThemeType
		isDark bool
	}{
		{ThemeLight, false},
		{ThemeDark, true},
		{ThemeNord, true},
		{ThemeSolarizedLight, false},
		{ThemeSolarizedDark, true},
		{ThemeMonokai, true},
		{ThemeGruvboxDark, true},
		{ThemeOneDark, true},
	}

	for _, tt := range tests {
		t.Run(tt.theme.Name(), func(t *testing.T) {
			got := tt.theme.IsDark()
			if got != tt.isDark {
				t.Errorf("ThemeType.IsDark() = %v, want %v", got, tt.isDark)
			}
		})
	}
}

func TestFontSizeNames(t *testing.T) {
	names := FontSizeNames()
	expected := []string{"Small", "Normal", "Large", "Extra Large"}
	if len(names) != len(expected) {
		t.Errorf("FontSizeNames() returned %d names, want %d", len(names), len(expected))
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("FontSizeNames()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestFontSizeFromName(t *testing.T) {
	tests := []struct {
		name string
		want FontSize
	}{
		{"Small", FontSizeSmall},
		{"Normal", FontSizeNormal},
		{"Large", FontSizeLarge},
		{"Extra Large", FontSizeExtraLarge},
		{"Unknown", FontSizeNormal}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FontSizeFromName(tt.name)
			if got != tt.want {
				t.Errorf("FontSizeFromName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestFontSize_Name(t *testing.T) {
	tests := []struct {
		size FontSize
		want string
	}{
		{FontSizeSmall, "Small"},
		{FontSizeNormal, "Normal"},
		{FontSizeLarge, "Large"},
		{FontSizeExtraLarge, "Extra Large"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.size.Name()
			if got != tt.want {
				t.Errorf("FontSize.Name() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFontSize_Scale(t *testing.T) {
	tests := []struct {
		size  FontSize
		scale float32
	}{
		{FontSizeSmall, 0.85},
		{FontSizeNormal, 1.0},
		{FontSizeLarge, 1.15},
		{FontSizeExtraLarge, 1.30},
	}

	for _, tt := range tests {
		t.Run(tt.size.Name(), func(t *testing.T) {
			got := tt.size.Scale()
			if got != tt.scale {
				t.Errorf("FontSize.Scale() = %v, want %v", got, tt.scale)
			}
		})
	}
}

func TestFontFamilyNames(t *testing.T) {
	names := FontFamilyNames()
	expected := []string{"System Default", "Monospace", "Serif", "Sans Serif"}
	if len(names) != len(expected) {
		t.Errorf("FontFamilyNames() returned %d names, want %d", len(names), len(expected))
	}
}

func TestFontFamilyFromName(t *testing.T) {
	tests := []struct {
		name string
		want FontFamily
	}{
		{"System Default", FontDefault},
		{"Monospace", FontMonospace},
		{"Serif", FontSerif},
		{"Sans Serif", FontSansSerif},
		{"Unknown", FontDefault}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FontFamilyFromName(tt.name)
			if got != tt.want {
				t.Errorf("FontFamilyFromName(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestNewMarkViewTheme(t *testing.T) {
	th := NewMarkViewTheme(ThemeDark)
	if th == nil {
		t.Error("NewMarkViewTheme() returned nil")
	}
}

func TestNewMarkViewThemeWithFont(t *testing.T) {
	th := NewMarkViewThemeWithFont(ThemeDark, FontMonospace)
	if th == nil {
		t.Error("NewMarkViewThemeWithFont() returned nil")
	}
}

func TestNewMarkViewThemeWithOptions(t *testing.T) {
	th := NewMarkViewThemeWithOptions(ThemeDark, FontMonospace, FontSizeLarge)
	if th == nil {
		t.Error("NewMarkViewThemeWithOptions() returned nil")
	}
}

func TestMarkViewTheme_Color(t *testing.T) {
	themes := []ThemeType{
		ThemeLight, ThemeDark, ThemeNord,
		ThemeSolarizedLight, ThemeSolarizedDark,
		ThemeMonokai, ThemeGruvboxDark, ThemeOneDark,
	}

	for _, themeType := range themes {
		t.Run(themeType.Name(), func(t *testing.T) {
			th := NewMarkViewTheme(themeType).(*MarkViewTheme)

			// Test that colors are returned without panicking
			_ = th.Color(theme.ColorNameBackground, fyne.ThemeVariant(0))
			_ = th.Color(theme.ColorNameForeground, fyne.ThemeVariant(0))
			_ = th.Color(theme.ColorNamePrimary, fyne.ThemeVariant(0))
			_ = th.Color(theme.ColorNameButton, fyne.ThemeVariant(0))
		})
	}
}

func TestMarkViewTheme_Size(t *testing.T) {
	th := NewMarkViewTheme(ThemeDark).(*MarkViewTheme)

	tests := []struct {
		name     fyne.ThemeSizeName
		minValue float32
	}{
		{theme.SizeNameText, 10},
		{theme.SizeNameHeadingText, 20},
		{theme.SizeNameSubHeadingText, 15},
		{"heading1", 25},
		{"heading2", 20},
		{"heading3", 15},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			got := th.Size(tt.name)
			if got < tt.minValue {
				t.Errorf("Size(%q) = %v, want >= %v", tt.name, got, tt.minValue)
			}
		})
	}
}

func TestMarkViewTheme_Font(t *testing.T) {
	th := NewMarkViewTheme(ThemeDark).(*MarkViewTheme)
	font := th.Font(fyne.TextStyle{})
	if font == nil {
		t.Error("Font() returned nil")
	}
}

func TestMarkViewTheme_Icon(t *testing.T) {
	th := NewMarkViewTheme(ThemeDark).(*MarkViewTheme)
	icon := th.Icon(theme.IconNameDocument)
	if icon == nil {
		t.Error("Icon() returned nil")
	}
}

func TestMarkViewTheme_GetCodeColors(t *testing.T) {
	// Test dark theme code colors
	thDark := NewMarkViewTheme(ThemeDark).(*MarkViewTheme)
	colorsDark := thDark.GetCodeColors()
	if colorsDark.Keyword == nil {
		t.Error("GetCodeColors() dark theme Keyword is nil")
	}

	// Test light theme code colors
	thLight := NewMarkViewTheme(ThemeLight).(*MarkViewTheme)
	colorsLight := thLight.GetCodeColors()
	if colorsLight.Keyword == nil {
		t.Error("GetCodeColors() light theme Keyword is nil")
	}
}

func TestIconFunctions(t *testing.T) {
	// Test that icon functions return non-nil resources
	icons := []struct {
		name string
		fn   func() fyne.Resource
	}{
		{"AppLogo", AppLogo},
		{"IconDocument", IconDocument},
		{"IconFolder", IconFolder},
		{"IconRefresh", IconRefresh},
		{"IconEdit", IconEdit},
		{"IconView", IconView},
		{"IconSave", IconSave},
		{"IconUndo", IconUndo},
		{"IconFileTree", IconFileTree},
		{"IconTOC", IconTOC},
		{"IconTheme", IconTheme},
		{"IconBold", IconBold},
		{"IconItalic", IconItalic},
		{"IconHeading", IconHeading},
		{"IconHeading1", IconHeading1},
		{"IconHeading2", IconHeading2},
		{"IconHeading3", IconHeading3},
		{"IconLink", IconLink},
		{"IconCode", IconCode},
		{"IconCodeBlock", IconCodeBlock},
		{"IconQuote", IconQuote},
		{"IconList", IconList},
		{"IconHorizontalRule", IconHorizontalRule},
		{"IconImage", IconImage},
		{"IconNewFile", IconNewFile},
		{"IconSaveAs", IconSaveAs},
		{"IconTable", IconTable},
		{"IconLibrary", IconLibrary},
		{"IconSplitView", IconSplitView},
		{"IconFocus", IconFocus},
		{"IconHelp", IconHelp},
		{"IconSearch", IconSearch},
		{"IconStar", IconStar},
		{"IconPresentation", IconPresentation},
		{"IconSnippet", IconSnippet},
		{"IconSort", IconSort},
		{"IconGoal", IconGoal},
		{"IconTypewriter", IconTypewriter},
		{"IconPrint", IconPrint},
		{"IconExport", IconExport},
		{"IconBacklinks", IconBacklinks},
		{"IconTag", IconTag},
		{"IconTemplate", IconTemplate},
		{"IconZen", IconZen},
		{"IconQuickSwitch", IconQuickSwitch},
		{"IconLinkCheck", IconLinkCheck},
	}

	for _, tt := range icons {
		t.Run(tt.name, func(t *testing.T) {
			icon := tt.fn()
			if icon == nil {
				t.Errorf("%s() returned nil", tt.name)
			}
		})
	}
}
