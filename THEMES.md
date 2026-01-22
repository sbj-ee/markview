# MarkView Theme Documentation

## Overview

MarkView includes two beautiful, carefully crafted themes designed for optimal markdown reading experience.

## Theme Switching

Switch themes via:
- **Menu**: View → Light Theme / Dark Theme
- **Effect**: Instant theme application with automatic content re-render

## Dark Theme (Default)

### Color Palette

Inspired by the popular **Dracula** theme, optimized for extended reading sessions.

#### Background Colors
- **Main Background**: `#282A36` (Deep purple-gray)
- **Input Background**: `#1E1F29` (Darker)
- **Header Background**: `#1E1F29` (Darker section bg)
- **Button**: `#44475A` (Purple-gray)
- **Hover**: `#505466` (Lighter on hover)
- **Pressed**: `#626780` (Lightest when pressed)

#### Text Colors
- **Foreground**: `#F8F8F2` (Off-white)
- **Disabled**: `#6272A4` (Muted purple)

#### Accent Colors
- **Primary (Cyan)**: `#8BE9FD` - Keywords, types
- **Success (Green)**: `#50FA7B` - Strings, functions
- **Warning (Orange)**: `#FFB86C` - Numbers
- **Error (Red)**: `#FF5555` - Keywords, operators

#### Syntax Highlighting Colors
- **Keyword**: `#FF79C6` (Pink)
- **String**: `#F1FA8C` (Yellow)
- **Comment**: `#6272A4` (Purple-gray)
- **Number**: `#BD93F9` (Purple)
- **Function**: `#50FA7B` (Green)
- **Operator**: `#FF79C6` (Pink)
- **Type**: `#8BE9FD` (Cyan)
- **Variable**: `#F8F8F2` (Off-white)
- **Code Background**: `#1E1F29` (Dark)

### Visual Characteristics
- **Eye Strain**: Minimal - designed for low-light environments
- **Contrast Ratio**: Medium (easier on eyes than pure black/white)
- **Best For**: Night reading, extended sessions, low-light environments
- **Mood**: Modern, vibrant, focused

---

## Light Theme

### Color Palette

Inspired by **GitHub's** clean design, professional and easy to read.

#### Background Colors
- **Main Background**: `#FFFFFF` (Pure white)
- **Input Background**: `#FAFBFC` (Very light gray)
- **Header Background**: `#F6F8FA` (Very light blue-gray)
- **Button**: `#F0F4F8` (Light gray)
- **Hover**: `#E1E7ED` (Slightly darker)
- **Pressed**: `#D2DAE2` (Even darker)

#### Text Colors
- **Foreground**: `#24292F` (Dark gray)
- **Disabled**: `#959DA5` (Gray)

#### Accent Colors
- **Primary (Blue)**: `#007ACC` - Links, keywords
- **Success (Green)**: `#28A745` - Strings, success states
- **Warning (Orange)**: `#FF9800` - Numbers, warnings
- **Error (Red)**: `#DC3545` - Errors, operators

#### Syntax Highlighting Colors
- **Keyword**: `#D73A49` (Red)
- **String**: `#032F62` (Dark blue)
- **Comment**: `#6A737D` (Gray)
- **Number**: `#005CC5` (Blue)
- **Function**: `#6F42C1` (Purple)
- **Operator**: `#D73A49` (Red)
- **Type**: `#005CC5` (Blue)
- **Variable**: `#24292F` (Dark gray)
- **Code Background**: `#F6F8FA` (Light gray)

### Visual Characteristics
- **Eye Strain**: Low - professional contrast
- **Contrast Ratio**: High (better for readability)
- **Best For**: Daytime use, printing, bright environments
- **Mood**: Clean, professional, minimal

---

## Theme Comparison

| Aspect | Dark Theme | Light Theme |
|--------|-----------|-------------|
| **Background** | Deep Purple-Gray | Pure White |
| **Text** | Off-White | Dark Gray |
| **Best Time** | Night / Low-light | Day / Bright light |
| **Eye Strain** | Minimal (dark) | Low (professional) |
| **Inspiration** | Dracula | GitHub |
| **Code Highlight** | Vibrant | Professional |
| **Battery Impact** | Lower (OLED) | Higher |
| **Print Quality** | Poor | Excellent |

---

## Syntax Highlighting Features

### Supported Token Types

Both themes provide rich syntax highlighting for:

- **Keywords** - Language keywords (if, for, while, etc.)
- **Strings** - String literals
- **Comments** - Single and multi-line comments
- **Numbers** - Integer and floating-point literals
- **Functions** - Function names and calls
- **Operators** - Arithmetic, logical, comparison operators
- **Types** - Class names, type declarations
- **Variables** - Variable names and identifiers

### Languages Supported

250+ languages including:
- Go, Python, JavaScript, TypeScript
- Rust, C, C++, C#
- Java, Kotlin, Swift
- Ruby, PHP, Perl
- SQL, HTML, CSS, SCSS
- Bash, PowerShell, Shell
- Markdown, YAML, JSON, TOML
- And many more...

---

## UI Element Styling

### Buttons
- Distinct colors for normal, hover, and pressed states
- Smooth visual feedback on interaction
- Theme-appropriate shadows

### Inputs & Text Fields
- Contrasting background from main content
- Clear borders in theme colors
- Focused state with accent color

### Scrollbars
- Themed to match overall color scheme
- Non-intrusive but visible
- Smooth scrolling experience

### TOC (Table of Contents)
- Tree structure with proper indentation
- Hover effects for navigation
- Selected item highlighting

### Menus
- Native Fyne menu styling
- Theme-aware colors
- Clear visual hierarchy

---

## Design Philosophy

### Dark Theme
- **Goal**: Reduce eye strain during extended reading
- **Inspiration**: Dracula color scheme
- **Approach**: Vibrant but not overwhelming colors
- **Use Case**: Primary theme for developers and night owls

### Light Theme
- **Goal**: Professional, clean markdown viewer
- **Inspiration**: GitHub's markdown rendering
- **Approach**: High contrast for clarity
- **Use Case**: Documentation, daytime use, presentations

---

## Future Theme Options

Planned for future releases:

- **Nord** - Arctic, north-bluish color palette
- **Solarized** - Precision colors for readability
- **Monokai** - Classic editor theme
- **One Dark** - Atom's popular theme
- **Gruvbox** - Retro groove color scheme
- **Custom Themes** - User-defined color schemes

---

## Technical Implementation

### Theme Structure
```go
type MarkViewTheme struct {
    themeType ThemeType
}
```

### Color Methods
- `Color()` - Returns themed colors for UI elements
- `GetCodeColors()` - Returns syntax highlighting colors
- Dynamic theme switching without restart

### Integration
- Fyne's theme system
- Custom color mapping
- Real-time theme updates

---

## Tips for Best Experience

### Dark Theme
- Use in dimly lit environments
- Ideal for 2-3+ hour reading sessions
- Great for OLED displays (battery saving)
- Reduce screen brightness for comfort

### Light Theme
- Best in well-lit rooms
- Perfect for printing or exporting
- Higher perceived sharpness
- Better for photography/image-heavy content

---

*Themes designed with ❤️ for the markdown community*
