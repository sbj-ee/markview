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
	// Use Inter font - for now, use default which is similar to Inter on modern systems
	// TODO: Bundle Inter font files for consistent cross-platform rendering
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given theme icon name
func (m *MarkViewTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the given theme size name
func (m *MarkViewTheme) Size(name fyne.ThemeSizeName) float32 {
	// Increase text sizes for better readability
	switch name {
	case theme.SizeNameText:
		return 14 // Compact text size for dialogs and lists
	case theme.SizeNameHeadingText:
		return 28 // H1 size - large and prominent
	case theme.SizeNameSubHeadingText:
		return 22 // H2 size
	case theme.SizeNameCaptionText:
		return 12 // Smaller text for TOC and captions
	case theme.SizeNameInlineIcon:
		return 20 // Icons in toolbar and lists
	case theme.SizeNamePadding:
		return 2 // Minimal padding for tight spacing
	case theme.SizeNameInnerPadding:
		return 2 // Minimal inner padding
	case theme.SizeNameLineSpacing:
		return 0 // No extra line spacing
	case "heading1":
		return 32 // H1: Extra large
	case "heading2":
		return 26 // H2: Large
	case "heading3":
		return 21 // H3: Medium-large
	case "heading4":
		return 18 // H4: Medium
	case "heading5":
		return 16 // H5: Small
	case "heading6":
		return 15 // H6: Base size
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
		return color.RGBA{R: 100, G: 100, B: 110, A: 255} // Muted gray

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
