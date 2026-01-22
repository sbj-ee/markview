package gui

import (
	"image/color"

	"github.com/sbj-ee/markview/internal/themes"
)

// MarkdownSyntaxColors defines colors for markdown syntax highlighting
type MarkdownSyntaxColors struct {
	Header     color.Color
	Bold       color.Color
	Italic     color.Color
	Code       color.Color
	CodeBg     color.Color
	Link       color.Color
	List       color.Color
	Blockquote color.Color
	Default    color.Color
}

// GetMarkdownSyntaxColors returns syntax colors based on theme
func GetMarkdownSyntaxColors(themeType themes.ThemeType) MarkdownSyntaxColors {
	if themeType == themes.ThemeDark {
		return MarkdownSyntaxColors{
			Header:     color.RGBA{R: 86, G: 182, B: 194, A: 255},  // Cyan/teal
			Bold:       color.RGBA{R: 229, G: 181, B: 103, A: 255}, // Orange
			Italic:     color.RGBA{R: 189, G: 147, B: 249, A: 255}, // Purple
			Code:       color.RGBA{R: 80, G: 250, B: 123, A: 255},  // Green
			CodeBg:     color.RGBA{R: 40, G: 42, B: 48, A: 255},    // Dark bg
			Link:       color.RGBA{R: 86, G: 182, B: 194, A: 255},  // Cyan
			List:       color.RGBA{R: 255, G: 121, B: 198, A: 255}, // Pink
			Blockquote: color.RGBA{R: 98, G: 114, B: 164, A: 255},  // Purple-gray
			Default:    color.RGBA{R: 200, G: 200, B: 200, A: 255}, // Light gray
		}
	}

	return MarkdownSyntaxColors{
		Header:     color.RGBA{R: 0, G: 92, B: 197, A: 255},    // Blue
		Bold:       color.RGBA{R: 215, G: 58, B: 73, A: 255},   // Red
		Italic:     color.RGBA{R: 111, G: 66, B: 193, A: 255},  // Purple
		Code:       color.RGBA{R: 3, G: 47, B: 98, A: 255},     // Dark blue
		CodeBg:     color.RGBA{R: 246, G: 248, B: 250, A: 255}, // Light bg
		Link:       color.RGBA{R: 0, G: 122, B: 204, A: 255},   // Blue
		List:       color.RGBA{R: 215, G: 58, B: 73, A: 255},   // Red
		Blockquote: color.RGBA{R: 106, G: 115, B: 125, A: 255}, // Gray
		Default:    color.RGBA{R: 36, G: 41, B: 47, A: 255},    // Dark gray
	}
}
