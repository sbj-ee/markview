package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ThemeType represents the type of theme
type ThemeType int

const (
	ThemeLight ThemeType = iota
	ThemeDark
	ThemeNord
	ThemeSolarizedLight
	ThemeSolarizedDark
	ThemeMonokai
	ThemeGruvboxDark
	ThemeOneDark
	ThemeAgua
	ThemeGroovy
	ThemeChiliPepper
)

// ThemeNames returns all available theme names
func ThemeNames() []string {
	return []string{
		"Light",
		"Dark",
		"Nord",
		"Solarized Light",
		"Solarized Dark",
		"Monokai",
		"Gruvbox Dark",
		"One Dark",
		"Agua",
		"Groovy",
		"Chili Pepper",
	}
}

// ThemeFromName returns the ThemeType for a given name
func ThemeFromName(name string) ThemeType {
	switch name {
	case "Light":
		return ThemeLight
	case "Dark":
		return ThemeDark
	case "Nord":
		return ThemeNord
	case "Solarized Light":
		return ThemeSolarizedLight
	case "Solarized Dark":
		return ThemeSolarizedDark
	case "Monokai":
		return ThemeMonokai
	case "Gruvbox Dark":
		return ThemeGruvboxDark
	case "One Dark":
		return ThemeOneDark
	case "Agua":
		return ThemeAgua
	case "Groovy":
		return ThemeGroovy
	case "Chili Pepper":
		return ThemeChiliPepper
	default:
		return ThemeDark
	}
}

// ThemeName returns the name for a ThemeType
func (t ThemeType) Name() string {
	return ThemeNames()[int(t)]
}

// IsDark returns true if the theme is a dark theme
func (t ThemeType) IsDark() bool {
	switch t {
	case ThemeLight, ThemeSolarizedLight:
		return false
	default:
		return true
	}
}

// FontFamily represents a font family option
type FontFamily string

const (
	FontDefault   FontFamily = "default"
	FontMonospace FontFamily = "monospace"
	FontSerif     FontFamily = "serif"
	FontSansSerif FontFamily = "sans-serif"
)

// FontSize represents a font size option (scale factor)
type FontSize string

const (
	FontSizeSmall      FontSize = "small"  // 0.85x
	FontSizeNormal     FontSize = "normal" // 1.0x
	FontSizeLarge      FontSize = "large"  // 1.15x
	FontSizeExtraLarge FontSize = "xlarge" // 1.30x
)

// FontSizeNames returns all available font size names
func FontSizeNames() []string {
	return []string{
		"Small",
		"Normal",
		"Large",
		"Extra Large",
	}
}

// FontSizeFromName returns the FontSize for a given name
func FontSizeFromName(name string) FontSize {
	switch name {
	case "Small":
		return FontSizeSmall
	case "Large":
		return FontSizeLarge
	case "Extra Large":
		return FontSizeExtraLarge
	default:
		return FontSizeNormal
	}
}

// FontSizeName returns the display name for a FontSize
func (f FontSize) Name() string {
	switch f {
	case FontSizeSmall:
		return "Small"
	case FontSizeLarge:
		return "Large"
	case FontSizeExtraLarge:
		return "Extra Large"
	default:
		return "Normal"
	}
}

// Scale returns the scale factor for the font size
func (f FontSize) Scale() float32 {
	switch f {
	case FontSizeSmall:
		return 0.85
	case FontSizeLarge:
		return 1.15
	case FontSizeExtraLarge:
		return 1.30
	default:
		return 1.0
	}
}

// FontFamilyNames returns all available font family names
func FontFamilyNames() []string {
	return []string{
		"System Default",
		"Monospace",
		"Serif",
		"Sans Serif",
	}
}

// FontFamilyFromName returns the FontFamily for a given name
func FontFamilyFromName(name string) FontFamily {
	switch name {
	case "Monospace":
		return FontMonospace
	case "Serif":
		return FontSerif
	case "Sans Serif":
		return FontSansSerif
	default:
		return FontDefault
	}
}

// MarkViewTheme is a custom theme for MarkView
type MarkViewTheme struct {
	themeType  ThemeType
	fontFamily FontFamily
	fontSize   FontSize
}

// NewMarkViewTheme creates a new MarkView theme
func NewMarkViewTheme(themeType ThemeType) fyne.Theme {
	return &MarkViewTheme{themeType: themeType, fontFamily: FontDefault, fontSize: FontSizeNormal}
}

// NewMarkViewThemeWithFont creates a new MarkView theme with custom font
func NewMarkViewThemeWithFont(themeType ThemeType, fontFamily FontFamily) fyne.Theme {
	return &MarkViewTheme{themeType: themeType, fontFamily: fontFamily, fontSize: FontSizeNormal}
}

// NewMarkViewThemeWithOptions creates a new MarkView theme with all options
func NewMarkViewThemeWithOptions(themeType ThemeType, fontFamily FontFamily, fontSize FontSize) fyne.Theme {
	return &MarkViewTheme{themeType: themeType, fontFamily: fontFamily, fontSize: fontSize}
}

// Color returns the color for the given theme color name
func (m *MarkViewTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch m.themeType {
	case ThemeLight:
		return m.lightColor(name)
	case ThemeDark:
		return m.darkColor(name)
	case ThemeNord:
		return m.nordColor(name)
	case ThemeSolarizedLight:
		return m.solarizedLightColor(name)
	case ThemeSolarizedDark:
		return m.solarizedDarkColor(name)
	case ThemeMonokai:
		return m.monokaiColor(name)
	case ThemeGruvboxDark:
		return m.gruvboxDarkColor(name)
	case ThemeOneDark:
		return m.oneDarkColor(name)
	case ThemeAgua:
		return m.aguaColor(name)
	case ThemeGroovy:
		return m.groovyColor(name)
	case ThemeChiliPepper:
		return m.chiliPepperColor(name)
	default:
		return m.darkColor(name)
	}
}

// Font returns the font resource for the given theme text style
func (m *MarkViewTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Apply custom font family preference
	switch m.fontFamily {
	case FontMonospace:
		style.Monospace = true
	case FontSerif:
		// Fyne doesn't have built-in serif support, use default
	case FontSansSerif:
		// Default is sans-serif in most systems
	}
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given theme icon name
func (m *MarkViewTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the given theme size name
func (m *MarkViewTheme) Size(name fyne.ThemeSizeName) float32 {
	scale := m.fontSize.Scale()

	// Increase text sizes for better readability
	switch name {
	case theme.SizeNameText:
		return 14 * scale // Compact text size for dialogs and lists
	case theme.SizeNameHeadingText:
		return 28 * scale // H1 size - large and prominent
	case theme.SizeNameSubHeadingText:
		return 22 * scale // H2 size
	case theme.SizeNameCaptionText:
		return 12 * scale // Smaller text for TOC and captions
	case theme.SizeNameInlineIcon:
		return 20 // Icons - don't scale
	case theme.SizeNamePadding:
		return 2 // Minimal padding - don't scale
	case theme.SizeNameInnerPadding:
		return 2 // Minimal inner padding - don't scale
	case theme.SizeNameLineSpacing:
		return 0 // No extra line spacing
	case "heading1":
		return 32 * scale // H1: Extra large
	case "heading2":
		return 26 * scale // H2: Large
	case "heading3":
		return 21 * scale // H3: Medium-large
	case "heading4":
		return 18 * scale // H4: Medium
	case "heading5":
		return 16 * scale // H5: Small
	case "heading6":
		return 15 * scale // H6: Base size
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// lightColor returns colors for light theme
func (m *MarkViewTheme) lightColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // Pure white
	case theme.ColorNameForeground:
		return color.RGBA{R: 36, G: 41, B: 47, A: 255} // Dark gray text

	case theme.ColorNameButton:
		return color.RGBA{R: 240, G: 244, B: 248, A: 255} // Light gray button
	case theme.ColorNameHover:
		return color.RGBA{R: 225, G: 231, B: 237, A: 255} // Slightly darker on hover
	case theme.ColorNamePressed:
		return color.RGBA{R: 210, G: 218, B: 226, A: 255} // Even darker when pressed

	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 122, B: 204, A: 255} // Nice blue
	case theme.ColorNameSuccess:
		return color.RGBA{R: 40, G: 167, B: 69, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 255, G: 152, B: 0, A: 255} // Orange
	case theme.ColorNameError:
		return color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red

	case theme.ColorNameSelection:
		return color.RGBA{R: 0, G: 122, B: 204, A: 76} // Blue with transparency
	case theme.ColorNameFocus:
		return color.RGBA{R: 0, G: 122, B: 204, A: 128} // Blue focus

	case theme.ColorNameInputBackground:
		return color.RGBA{R: 250, G: 251, B: 252, A: 255} // Very light gray
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 209, G: 213, B: 219, A: 255} // Light border

	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 30} // Subtle shadow

	case theme.ColorNameDisabled:
		return color.RGBA{R: 149, G: 157, B: 165, A: 255} // Gray for disabled
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 149, G: 157, B: 165, A: 255} // Gray for placeholders

	case theme.ColorNameScrollBar:
		return color.RGBA{R: 209, G: 213, B: 219, A: 255} // Light scrollbar

	// Custom markdown colors
	case theme.ColorNameHeaderBackground:
		return color.RGBA{R: 246, G: 248, B: 250, A: 255} // Very light blue-gray
	case "math":
		return color.RGBA{R: 128, G: 90, B: 213, A: 255} // Purple for math
	case "tableHeader":
		return color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange for table headers

	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

// darkColor returns colors for dark theme (inspired by the reference markdown viewer)
func (m *MarkViewTheme) darkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 30, G: 32, B: 36, A: 255} // Dark charcoal background
	case theme.ColorNameForeground:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Light gray text

	case theme.ColorNameButton:
		return color.RGBA{R: 50, G: 52, B: 58, A: 255} // Slightly lighter than bg
	case theme.ColorNameHover:
		return color.RGBA{R: 60, G: 62, B: 70, A: 255} // Lighter on hover
	case theme.ColorNamePressed:
		return color.RGBA{R: 70, G: 72, B: 80, A: 255} // Lightest when pressed

	case theme.ColorNamePrimary:
		return color.RGBA{R: 86, G: 182, B: 194, A: 255} // Cyan/teal for headings
	case theme.ColorNameSuccess:
		return color.RGBA{R: 80, G: 250, B: 123, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 229, G: 181, B: 103, A: 255} // Orange/gold for sub-headings
	case theme.ColorNameError:
		return color.RGBA{R: 255, G: 85, B: 85, A: 255} // Red

	case theme.ColorNameSelection:
		return color.RGBA{R: 68, G: 71, B: 90, A: 200} // Selection with transparency
	case theme.ColorNameFocus:
		return color.RGBA{R: 86, G: 182, B: 194, A: 128} // Cyan focus

	case theme.ColorNameHyperlink:
		return color.RGBA{R: 86, G: 182, B: 194, A: 255} // Cyan for links

	case theme.ColorNameInputBackground:
		return color.RGBA{R: 40, G: 42, B: 48, A: 255} // Darker background
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 60, G: 62, B: 70, A: 255} // Border

	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 50} // Shadow

	case theme.ColorNameDisabled:
		return color.RGBA{R: 100, G: 100, B: 110, A: 255} // Muted gray
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 140, G: 148, B: 160, A: 255} // Brighter gray for comments

	case theme.ColorNameScrollBar:
		return color.RGBA{R: 60, G: 62, B: 70, A: 255} // Scrollbar

	// Custom markdown colors
	case theme.ColorNameHeaderBackground:
		return color.RGBA{R: 40, G: 42, B: 48, A: 255} // Darker section bg

	// Custom colors for markdown elements
	case "heading1", "heading2":
		return color.RGBA{R: 86, G: 182, B: 194, A: 255} // Cyan/teal for H1, H2
	case "heading3", "heading4":
		return color.RGBA{R: 229, G: 181, B: 103, A: 255} // Orange/gold for H3, H4
	case "heading5", "heading6":
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Light gray for H5, H6
	case "bold":
		return color.RGBA{R: 229, G: 181, B: 103, A: 255} // Orange/gold for bold text
	case "link":
		return color.RGBA{R: 86, G: 182, B: 194, A: 255} // Cyan for links
	case "separator":
		return color.RGBA{R: 70, G: 72, B: 80, A: 255} // Subtle gray for horizontal rules
	case "math":
		return color.RGBA{R: 189, G: 147, B: 249, A: 255} // Purple for math
	case "tableHeader":
		return color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange for table headers

	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// nordColor returns colors for Nord theme
func (m *MarkViewTheme) nordColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 46, G: 52, B: 64, A: 255} // Nord0 - Polar Night
	case theme.ColorNameForeground:
		return color.RGBA{R: 216, G: 222, B: 233, A: 255} // Nord4 - Snow Storm
	case theme.ColorNameButton:
		return color.RGBA{R: 59, G: 66, B: 82, A: 255} // Nord1
	case theme.ColorNameHover:
		return color.RGBA{R: 67, G: 76, B: 94, A: 255} // Nord2
	case theme.ColorNamePressed:
		return color.RGBA{R: 76, G: 86, B: 106, A: 255} // Nord3
	case theme.ColorNamePrimary:
		return color.RGBA{R: 136, G: 192, B: 208, A: 255} // Nord8 - Frost
	case theme.ColorNameSuccess:
		return color.RGBA{R: 163, G: 190, B: 140, A: 255} // Nord14 - Aurora Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 235, G: 203, B: 139, A: 255} // Nord13 - Aurora Yellow
	case theme.ColorNameError:
		return color.RGBA{R: 191, G: 97, B: 106, A: 255} // Nord11 - Aurora Red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 129, G: 161, B: 193, A: 255} // Nord9
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 59, G: 66, B: 82, A: 255} // Nord1
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 76, G: 86, B: 106, A: 255} // Nord3
	case theme.ColorNameSelection:
		return color.RGBA{R: 76, G: 86, B: 106, A: 200} // Nord3
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 76, G: 86, B: 106, A: 255} // Nord3
	case "heading1", "heading2":
		return color.RGBA{R: 136, G: 192, B: 208, A: 255} // Nord8
	case "heading3", "heading4":
		return color.RGBA{R: 235, G: 203, B: 139, A: 255} // Nord13
	case "math":
		return color.RGBA{R: 180, G: 142, B: 173, A: 255} // Nord15 - Purple
	case "tableHeader":
		return color.RGBA{R: 235, G: 203, B: 139, A: 255} // Nord13 - Orange/gold for table headers
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 136, G: 148, B: 166, A: 255} // Visible gray for comments
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// solarizedLightColor returns colors for Solarized Light theme
func (m *MarkViewTheme) solarizedLightColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 253, G: 246, B: 227, A: 255} // Base3
	case theme.ColorNameForeground:
		return color.RGBA{R: 101, G: 123, B: 131, A: 255} // Base00
	case theme.ColorNameButton:
		return color.RGBA{R: 238, G: 232, B: 213, A: 255} // Base2
	case theme.ColorNameHover:
		return color.RGBA{R: 227, G: 221, B: 202, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 216, G: 210, B: 191, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 38, G: 139, B: 210, A: 255} // Blue
	case theme.ColorNameSuccess:
		return color.RGBA{R: 133, G: 153, B: 0, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 181, G: 137, B: 0, A: 255} // Yellow
	case theme.ColorNameError:
		return color.RGBA{R: 220, G: 50, B: 47, A: 255} // Red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 38, G: 139, B: 210, A: 255} // Blue
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 238, G: 232, B: 213, A: 255} // Base2
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 147, G: 161, B: 161, A: 255} // Base1
	case "heading1", "heading2":
		return color.RGBA{R: 38, G: 139, B: 210, A: 255} // Blue
	case "heading3", "heading4":
		return color.RGBA{R: 181, G: 137, B: 0, A: 255} // Yellow
	case "math":
		return color.RGBA{R: 108, G: 113, B: 196, A: 255} // Violet
	case "tableHeader":
		return color.RGBA{R: 203, G: 75, B: 22, A: 255} // Orange for table headers
	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

// solarizedDarkColor returns colors for Solarized Dark theme
func (m *MarkViewTheme) solarizedDarkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 0, G: 43, B: 54, A: 255} // Base03
	case theme.ColorNameForeground:
		return color.RGBA{R: 131, G: 148, B: 150, A: 255} // Base0
	case theme.ColorNameButton:
		return color.RGBA{R: 7, G: 54, B: 66, A: 255} // Base02
	case theme.ColorNameHover:
		return color.RGBA{R: 27, G: 74, B: 86, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 47, G: 94, B: 106, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 38, G: 139, B: 210, A: 255} // Blue
	case theme.ColorNameSuccess:
		return color.RGBA{R: 133, G: 153, B: 0, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 181, G: 137, B: 0, A: 255} // Yellow
	case theme.ColorNameError:
		return color.RGBA{R: 220, G: 50, B: 47, A: 255} // Red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 38, G: 139, B: 210, A: 255} // Blue
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 7, G: 54, B: 66, A: 255} // Base02
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 88, G: 110, B: 117, A: 255} // Base01
	case theme.ColorNameSelection:
		return color.RGBA{R: 7, G: 54, B: 66, A: 200} // Base02
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 88, G: 110, B: 117, A: 255} // Base01
	case "heading1", "heading2":
		return color.RGBA{R: 38, G: 139, B: 210, A: 255} // Blue
	case "heading3", "heading4":
		return color.RGBA{R: 181, G: 137, B: 0, A: 255} // Yellow
	case "math":
		return color.RGBA{R: 108, G: 113, B: 196, A: 255} // Violet
	case "tableHeader":
		return color.RGBA{R: 203, G: 75, B: 22, A: 255} // Orange for table headers
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 101, G: 123, B: 131, A: 255} // Base00 - visible gray for comments
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// monokaiColor returns colors for Monokai theme
func (m *MarkViewTheme) monokaiColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 39, G: 40, B: 34, A: 255} // Background
	case theme.ColorNameForeground:
		return color.RGBA{R: 248, G: 248, B: 242, A: 255} // Foreground
	case theme.ColorNameButton:
		return color.RGBA{R: 55, G: 56, B: 50, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 70, G: 71, B: 65, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 85, G: 86, B: 80, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 102, G: 217, B: 239, A: 255} // Cyan
	case theme.ColorNameSuccess:
		return color.RGBA{R: 166, G: 226, B: 46, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 253, G: 151, B: 31, A: 255} // Orange
	case theme.ColorNameError:
		return color.RGBA{R: 249, G: 38, B: 114, A: 255} // Pink/Red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 102, G: 217, B: 239, A: 255} // Cyan
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 55, G: 56, B: 50, A: 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 117, G: 113, B: 94, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 73, G: 72, B: 62, A: 200}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 117, G: 113, B: 94, A: 255}
	case "heading1", "heading2":
		return color.RGBA{R: 249, G: 38, B: 114, A: 255} // Pink
	case "heading3", "heading4":
		return color.RGBA{R: 253, G: 151, B: 31, A: 255} // Orange
	case "math":
		return color.RGBA{R: 174, G: 129, B: 255, A: 255} // Purple
	case "tableHeader":
		return color.RGBA{R: 253, G: 151, B: 31, A: 255} // Orange for table headers
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 117, G: 113, B: 94, A: 255} // Visible gray for comments
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// gruvboxDarkColor returns colors for Gruvbox Dark theme
func (m *MarkViewTheme) gruvboxDarkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 40, G: 40, B: 40, A: 255} // bg0
	case theme.ColorNameForeground:
		return color.RGBA{R: 235, G: 219, B: 178, A: 255} // fg
	case theme.ColorNameButton:
		return color.RGBA{R: 60, G: 56, B: 54, A: 255} // bg1
	case theme.ColorNameHover:
		return color.RGBA{R: 80, G: 73, B: 69, A: 255} // bg2
	case theme.ColorNamePressed:
		return color.RGBA{R: 102, G: 92, B: 84, A: 255} // bg3
	case theme.ColorNamePrimary:
		return color.RGBA{R: 131, G: 165, B: 152, A: 255} // Aqua
	case theme.ColorNameSuccess:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 250, G: 189, B: 47, A: 255} // Yellow
	case theme.ColorNameError:
		return color.RGBA{R: 251, G: 73, B: 52, A: 255} // Red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 131, G: 165, B: 152, A: 255} // Aqua
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 60, G: 56, B: 54, A: 255} // bg1
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 102, G: 92, B: 84, A: 255} // bg3
	case theme.ColorNameSelection:
		return color.RGBA{R: 80, G: 73, B: 69, A: 200} // bg2
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 102, G: 92, B: 84, A: 255} // bg3
	case "heading1", "heading2":
		return color.RGBA{R: 131, G: 165, B: 152, A: 255} // Aqua
	case "heading3", "heading4":
		return color.RGBA{R: 250, G: 189, B: 47, A: 255} // Yellow
	case "math":
		return color.RGBA{R: 211, G: 134, B: 155, A: 255} // Purple
	case "tableHeader":
		return color.RGBA{R: 254, G: 128, B: 25, A: 255} // Orange for table headers
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 146, G: 131, B: 116, A: 255} // fg3 - visible gray for comments
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// oneDarkColor returns colors for One Dark theme
func (m *MarkViewTheme) oneDarkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 40, G: 44, B: 52, A: 255} // Background
	case theme.ColorNameForeground:
		return color.RGBA{R: 171, G: 178, B: 191, A: 255} // Foreground
	case theme.ColorNameButton:
		return color.RGBA{R: 50, G: 54, B: 62, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 60, G: 64, B: 72, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 70, G: 74, B: 82, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 97, G: 175, B: 239, A: 255} // Blue
	case theme.ColorNameSuccess:
		return color.RGBA{R: 152, G: 195, B: 121, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 229, G: 192, B: 123, A: 255} // Yellow
	case theme.ColorNameError:
		return color.RGBA{R: 224, G: 108, B: 117, A: 255} // Red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 97, G: 175, B: 239, A: 255} // Blue
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 50, G: 54, B: 62, A: 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 76, G: 82, B: 99, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 62, G: 68, B: 81, A: 200}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 76, G: 82, B: 99, A: 255}
	case "heading1", "heading2":
		return color.RGBA{R: 97, G: 175, B: 239, A: 255} // Blue
	case "heading3", "heading4":
		return color.RGBA{R: 229, G: 192, B: 123, A: 255} // Yellow
	case "math":
		return color.RGBA{R: 198, G: 120, B: 221, A: 255} // Purple
	case "tableHeader":
		return color.RGBA{R: 209, G: 154, B: 102, A: 255} // Orange for table headers
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 128, G: 135, B: 150, A: 255} // Comment gray - visible
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// aguaColor returns colors for the Agua theme: a deep-ocean palette of blues.
func (m *MarkViewTheme) aguaColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 13, G: 27, B: 42, A: 255} // Deep ocean navy
	case theme.ColorNameForeground:
		return color.RGBA{R: 202, G: 227, B: 240, A: 255} // Pale sea foam text
	case theme.ColorNameButton:
		return color.RGBA{R: 27, G: 49, B: 77, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 35, G: 62, B: 94, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 45, G: 76, B: 112, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 180, B: 216, A: 255} // Bright azure
	case theme.ColorNameSuccess:
		return color.RGBA{R: 64, G: 192, B: 160, A: 255} // Teal
	case theme.ColorNameWarning:
		return color.RGBA{R: 240, G: 196, B: 100, A: 255} // Sandy gold
	case theme.ColorNameError:
		return color.RGBA{R: 235, G: 110, B: 110, A: 255} // Coral
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 72, G: 202, B: 228, A: 255} // Cyan
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 22, G: 41, B: 64, A: 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 60, G: 100, B: 140, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 40, G: 90, B: 130, A: 200}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 60, G: 100, B: 140, A: 255}
	case "heading1", "heading2":
		return color.RGBA{R: 0, G: 180, B: 216, A: 255} // Azure
	case "heading3", "heading4":
		return color.RGBA{R: 144, G: 224, B: 239, A: 255} // Light aqua
	case "math":
		return color.RGBA{R: 150, G: 170, B: 230, A: 255} // Periwinkle
	case "tableHeader":
		return color.RGBA{R: 72, G: 202, B: 228, A: 255} // Cyan
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 110, G: 140, B: 170, A: 255} // Muted slate blue
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// groovyColor returns colors for the Groovy theme: 1970s burnt orange, brown,
// and mustard yellow.
func (m *MarkViewTheme) groovyColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 42, G: 30, B: 22, A: 255} // Dark chocolate brown
	case theme.ColorNameForeground:
		return color.RGBA{R: 235, G: 220, B: 195, A: 255} // Cream
	case theme.ColorNameButton:
		return color.RGBA{R: 66, G: 46, B: 30, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 84, G: 58, B: 38, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 102, G: 70, B: 46, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 204, G: 102, B: 0, A: 255} // Burnt orange
	case theme.ColorNameSuccess:
		return color.RGBA{R: 124, G: 138, B: 50, A: 255} // Avocado green
	case theme.ColorNameWarning:
		return color.RGBA{R: 224, G: 168, B: 0, A: 255} // Mustard yellow
	case theme.ColorNameError:
		return color.RGBA{R: 170, G: 57, B: 30, A: 255} // Rust red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 230, G: 150, B: 40, A: 255} // Amber
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 60, G: 42, B: 28, A: 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 120, G: 84, B: 50, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 140, G: 90, B: 40, A: 200}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 120, G: 84, B: 50, A: 255}
	case "heading1", "heading2":
		return color.RGBA{R: 204, G: 102, B: 0, A: 255} // Burnt orange
	case "heading3", "heading4":
		return color.RGBA{R: 224, G: 168, B: 0, A: 255} // Mustard yellow
	case "math":
		return color.RGBA{R: 124, G: 138, B: 50, A: 255} // Avocado green
	case "tableHeader":
		return color.RGBA{R: 204, G: 102, B: 0, A: 255} // Burnt orange
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 150, G: 125, B: 95, A: 255} // Muted tan
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// chiliPepperColor returns colors for the Chili Pepper theme: pepper reds,
// greens, and oranges on a warm charcoal.
func (m *MarkViewTheme) chiliPepperColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 28, G: 24, B: 22, A: 255} // Warm charcoal
	case theme.ColorNameForeground:
		return color.RGBA{R: 240, G: 230, B: 220, A: 255} // Warm off-white
	case theme.ColorNameButton:
		return color.RGBA{R: 44, G: 38, B: 34, A: 255}
	case theme.ColorNameHover:
		return color.RGBA{R: 58, G: 50, B: 44, A: 255}
	case theme.ColorNamePressed:
		return color.RGBA{R: 72, G: 62, B: 54, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 213, G: 43, B: 30, A: 255} // Chili red
	case theme.ColorNameSuccess:
		return color.RGBA{R: 76, G: 175, B: 80, A: 255} // Jalapeño green
	case theme.ColorNameWarning:
		return color.RGBA{R: 240, G: 140, B: 0, A: 255} // Habanero orange
	case theme.ColorNameError:
		return color.RGBA{R: 200, G: 30, B: 30, A: 255} // Deep red
	case theme.ColorNameHyperlink:
		return color.RGBA{R: 240, G: 140, B: 0, A: 255} // Orange
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 40, G: 34, B: 30, A: 255}
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 90, G: 70, B: 60, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 160, G: 50, B: 40, A: 200}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 90, G: 70, B: 60, A: 255}
	case "heading1", "heading2":
		return color.RGBA{R: 213, G: 43, B: 30, A: 255} // Chili red
	case "heading3", "heading4":
		return color.RGBA{R: 76, G: 175, B: 80, A: 255} // Pepper green
	case "math":
		return color.RGBA{R: 240, G: 140, B: 0, A: 255} // Orange
	case "tableHeader":
		return color.RGBA{R: 240, G: 140, B: 0, A: 255} // Orange
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 150, G: 130, B: 120, A: 255} // Muted warm gray
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// CodeColors returns color scheme for syntax highlighting
type CodeColors struct {
	Keyword    color.Color
	String     color.Color
	Comment    color.Color
	Number     color.Color
	Function   color.Color
	Operator   color.Color
	Type       color.Color
	Variable   color.Color
	Background color.Color
}

// AppLogo returns the MarkView application logo
func AppLogo() fyne.Resource {
	return resourceAppLogoSvg
}

// IconDocument returns a document/file icon
func IconDocument() fyne.Resource {
	return theme.DocumentIcon()
}

// IconFolder returns a folder icon
func IconFolder() fyne.Resource {
	return theme.FolderOpenIcon()
}

// IconRefresh returns a refresh icon
func IconRefresh() fyne.Resource {
	return theme.ViewRefreshIcon()
}

// IconEdit returns an edit/content icon
func IconEdit() fyne.Resource {
	return theme.DocumentCreateIcon()
}

// IconView returns a view/visibility icon
func IconView() fyne.Resource {
	return theme.VisibilityIcon()
}

// IconSave returns a save icon (using storage/download as proxy)
func IconSave() fyne.Resource {
	return theme.DocumentSaveIcon()
}

// IconUndo returns an undo/cancel icon
func IconUndo() fyne.Resource {
	return theme.ContentUndoIcon()
}

// IconFileTree returns a list/tree icon
func IconFileTree() fyne.Resource {
	return theme.ListIcon()
}

// IconTOC returns a menu/TOC icon
func IconTOC() fyne.Resource {
	return theme.MenuIcon()
}

// IconTheme returns a color/theme icon
func IconTheme() fyne.Resource {
	return theme.ColorPaletteIcon()
}

// IconBold returns a bold text icon
func IconBold() fyne.Resource {
	return theme.NewThemedResource(resourceBoldSvg)
}

// IconItalic returns an italic text icon
func IconItalic() fyne.Resource {
	return theme.NewThemedResource(resourceItalicSvg)
}

// IconHeading returns a heading icon
func IconHeading() fyne.Resource {
	return theme.NewThemedResource(resourceHeadingSvg)
}

// IconHeading1 returns a H1 heading icon
func IconHeading1() fyne.Resource {
	return theme.NewThemedResource(resourceHeading1Svg)
}

// IconHeading2 returns a H2 heading icon
func IconHeading2() fyne.Resource {
	return theme.NewThemedResource(resourceHeading2Svg)
}

// IconHeading3 returns a H3 heading icon
func IconHeading3() fyne.Resource {
	return theme.NewThemedResource(resourceHeading3Svg)
}

// IconLink returns a link icon
func IconLink() fyne.Resource {
	return theme.NewThemedResource(resourceLinkSvg)
}

// IconCode returns an inline code icon
func IconCode() fyne.Resource {
	return theme.NewThemedResource(resourceCodeSvg)
}

// IconCodeBlock returns a code block icon
func IconCodeBlock() fyne.Resource {
	return theme.NewThemedResource(resourceCodeBlockSvg)
}

// IconQuote returns a blockquote icon
func IconQuote() fyne.Resource {
	return theme.NewThemedResource(resourceQuoteSvg)
}

// IconList returns a list icon
func IconList() fyne.Resource {
	return theme.ListIcon()
}

// IconHorizontalRule returns a horizontal rule icon
func IconHorizontalRule() fyne.Resource {
	return theme.NewThemedResource(resourceHrSvg)
}

// IconImage returns an image icon
func IconImage() fyne.Resource {
	return theme.NewThemedResource(resourceImageSvg)
}

// IconNewFile returns a new file icon
func IconNewFile() fyne.Resource {
	return theme.NewThemedResource(resourceNewFileSvg)
}

// IconSaveAs returns a save as icon
func IconSaveAs() fyne.Resource {
	return theme.NewThemedResource(resourceSaveAsSvg)
}

// IconTable returns a table icon
func IconTable() fyne.Resource {
	return theme.NewThemedResource(resourceTableSvg)
}

// IconLibrary returns a library/book icon
func IconLibrary() fyne.Resource {
	return theme.NewThemedResource(resourceLibrarySvg)
}

// IconSplitView returns a split view icon
func IconSplitView() fyne.Resource {
	return theme.NewThemedResource(resourceSplitViewSvg)
}

// IconSingleView returns a single view icon (for exiting split view)
func IconSingleView() fyne.Resource {
	return theme.NewThemedResource(resourceSingleViewSvg)
}

// IconFocus returns a focus mode icon
func IconFocus() fyne.Resource {
	return theme.NewThemedResource(resourceFocusSvg)
}

// IconHelp returns a help icon
func IconHelp() fyne.Resource {
	return theme.NewThemedResource(resourceHelpSvg)
}

// IconSearch returns a search icon
func IconSearch() fyne.Resource {
	return theme.NewThemedResource(resourceSearchSvg)
}

// IconStar returns a star/favorite icon
func IconStar() fyne.Resource {
	return theme.NewThemedResource(resourceStarSvg)
}

// IconPresentation returns a presentation mode icon
func IconPresentation() fyne.Resource {
	return theme.NewThemedResource(resourcePresentationSvg)
}

// IconSnippet returns a snippet icon
func IconSnippet() fyne.Resource {
	return theme.NewThemedResource(resourceSnippetSvg)
}

// IconSort returns a sort icon
func IconSort() fyne.Resource {
	return theme.NewThemedResource(resourceSortSvg)
}

// IconGoal returns a goal/target icon
func IconGoal() fyne.Resource {
	return theme.NewThemedResource(resourceGoalSvg)
}

// IconTypewriter returns a typewriter mode icon
func IconTypewriter() fyne.Resource {
	return theme.NewThemedResource(resourceTypewriterSvg)
}

// IconPrint returns a print icon
func IconPrint() fyne.Resource {
	return theme.NewThemedResource(resourcePrintSvg)
}

// IconExport returns an export icon
func IconExport() fyne.Resource {
	return theme.NewThemedResource(resourceExportSvg)
}

// IconBacklinks returns a backlinks icon
func IconBacklinks() fyne.Resource {
	return theme.NewThemedResource(resourceBacklinksSvg)
}

// IconTag returns a tag icon
func IconTag() fyne.Resource {
	return theme.NewThemedResource(resourceTagSvg)
}

// IconTemplate returns a template icon
func IconTemplate() fyne.Resource {
	return theme.NewThemedResource(resourceTemplateSvg)
}

// IconZen returns a zen/fullscreen icon
func IconZen() fyne.Resource {
	return theme.NewThemedResource(resourceZenSvg)
}

// IconQuickSwitch returns a quick switch/search icon
func IconQuickSwitch() fyne.Resource {
	return theme.NewThemedResource(resourceQuickSwitchSvg)
}

// IconLinkCheck returns a link validation icon
func IconLinkCheck() fyne.Resource {
	return theme.NewThemedResource(resourceLinkCheckSvg)
}

// IconSubscript returns a subscript formatting icon
func IconSubscript() fyne.Resource {
	return theme.NewThemedResource(resourceSubscriptSvg)
}

// IconSuperscript returns a superscript formatting icon
func IconSuperscript() fyne.Resource {
	return theme.NewThemedResource(resourceSuperscriptSvg)
}

// IconStrikethrough returns a strikethrough formatting icon
func IconStrikethrough() fyne.Resource {
	return theme.NewThemedResource(resourceStrikethroughSvg)
}

// IconNumberedList returns a numbered list icon
func IconNumberedList() fyne.Resource {
	return theme.NewThemedResource(resourceNumberedListSvg)
}

// IconCheckbox returns a checkbox/task list icon
func IconCheckbox() fyne.Resource {
	return theme.NewThemedResource(resourceCheckboxSvg)
}

// IconSymbol returns a symbol picker icon
func IconSymbol() fyne.Resource {
	return theme.NewThemedResource(resourceSymbolSvg)
}

// IconUnderline returns an underline formatting icon
func IconUnderline() fyne.Resource {
	return theme.NewThemedResource(resourceUnderlineSvg)
}

// IconHighlight returns a highlight/mark formatting icon
func IconHighlight() fyne.Resource {
	return theme.NewThemedResource(resourceHighlightSvg)
}

// IconFootnote returns a footnote icon
func IconFootnote() fyne.Resource {
	return theme.NewThemedResource(resourceFootnoteSvg)
}

// GetCodeColors returns syntax highlighting colors for the current theme
func (m *MarkViewTheme) GetCodeColors() CodeColors {
	if m.themeType == ThemeDark {
		return CodeColors{
			Keyword:    color.RGBA{R: 255, G: 121, B: 198, A: 255}, // Pink
			String:     color.RGBA{R: 241, G: 250, B: 140, A: 255}, // Yellow
			Comment:    color.RGBA{R: 98, G: 114, B: 164, A: 255},  // Purple-gray
			Number:     color.RGBA{R: 189, G: 147, B: 249, A: 255}, // Purple
			Function:   color.RGBA{R: 80, G: 250, B: 123, A: 255},  // Green
			Operator:   color.RGBA{R: 255, G: 121, B: 198, A: 255}, // Pink
			Type:       color.RGBA{R: 139, G: 233, B: 253, A: 255}, // Cyan
			Variable:   color.RGBA{R: 248, G: 248, B: 242, A: 255}, // Off-white
			Background: color.RGBA{R: 30, G: 31, B: 41, A: 255},    // Dark bg
		}
	}

	return CodeColors{
		Keyword:    color.RGBA{R: 215, G: 58, B: 73, A: 255},   // Red
		String:     color.RGBA{R: 3, G: 47, B: 98, A: 255},     // Dark blue
		Comment:    color.RGBA{R: 106, G: 115, B: 125, A: 255}, // Gray
		Number:     color.RGBA{R: 0, G: 92, B: 197, A: 255},    // Blue
		Function:   color.RGBA{R: 111, G: 66, B: 193, A: 255},  // Purple
		Operator:   color.RGBA{R: 215, G: 58, B: 73, A: 255},   // Red
		Type:       color.RGBA{R: 0, G: 92, B: 197, A: 255},    // Blue
		Variable:   color.RGBA{R: 36, G: 41, B: 47, A: 255},    // Dark gray
		Background: color.RGBA{R: 246, G: 248, B: 250, A: 255}, // Light bg
	}
}
