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
)

// MarkViewTheme is a custom theme for MarkView
type MarkViewTheme struct {
	themeType ThemeType
}

// NewMarkViewTheme creates a new MarkView theme
func NewMarkViewTheme(themeType ThemeType) fyne.Theme {
	return &MarkViewTheme{themeType: themeType}
}

// Color returns the color for the given theme color name
func (m *MarkViewTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if m.themeType == ThemeDark {
		return m.darkColor(name)
	}
	return m.lightColor(name)
}

// Font returns the font resource for the given theme text style
func (m *MarkViewTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given theme icon name
func (m *MarkViewTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the given theme size name
func (m *MarkViewTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
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

	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

// darkColor returns colors for dark theme (inspired by Dracula/Nord)
func (m *MarkViewTheme) darkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 40, G: 42, B: 54, A: 255} // Dark purple-gray (Dracula bg)
	case theme.ColorNameForeground:
		return color.RGBA{R: 248, G: 248, B: 242, A: 255} // Off-white text

	case theme.ColorNameButton:
		return color.RGBA{R: 68, G: 71, B: 90, A: 255} // Lighter purple-gray
	case theme.ColorNameHover:
		return color.RGBA{R: 80, G: 84, B: 106, A: 255} // Even lighter on hover
	case theme.ColorNamePressed:
		return color.RGBA{R: 98, G: 103, B: 128, A: 255} // Lightest when pressed

	case theme.ColorNamePrimary:
		return color.RGBA{R: 139, G: 233, B: 253, A: 255} // Cyan
	case theme.ColorNameSuccess:
		return color.RGBA{R: 80, G: 250, B: 123, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 255, G: 184, B: 108, A: 255} // Orange
	case theme.ColorNameError:
		return color.RGBA{R: 255, G: 85, B: 85, A: 255} // Red

	case theme.ColorNameSelection:
		return color.RGBA{R: 68, G: 71, B: 90, A: 200} // Purple-gray with transparency
	case theme.ColorNameFocus:
		return color.RGBA{R: 139, G: 233, B: 253, A: 128} // Cyan focus

	case theme.ColorNameInputBackground:
		return color.RGBA{R: 30, G: 31, B: 41, A: 255} // Darker background
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 68, G: 71, B: 90, A: 255} // Purple-gray border

	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 50} // Darker shadow

	case theme.ColorNameDisabled:
		return color.RGBA{R: 98, G: 114, B: 164, A: 255} // Muted purple
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 98, G: 114, B: 164, A: 255} // Muted purple

	case theme.ColorNameScrollBar:
		return color.RGBA{R: 68, G: 71, B: 90, A: 255} // Purple-gray scrollbar

	// Custom markdown colors
	case theme.ColorNameHeaderBackground:
		return color.RGBA{R: 30, G: 31, B: 41, A: 255} // Darker section bg

	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

// CodeColors returns color scheme for syntax highlighting
type CodeColors struct {
	Keyword   color.Color
	String    color.Color
	Comment   color.Color
	Number    color.Color
	Function  color.Color
	Operator  color.Color
	Type      color.Color
	Variable  color.Color
	Background color.Color
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
