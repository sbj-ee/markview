package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SymbolCategory represents a category of symbols
type SymbolCategory struct {
	Name    string
	Symbols []string
}

// GetSymbolCategories returns all available symbol categories
func GetSymbolCategories() []SymbolCategory {
	return []SymbolCategory{
		{
			Name: "Greek Letters",
			Symbols: []string{
				"α", "β", "γ", "δ", "ε", "ζ", "η", "θ",
				"ι", "κ", "λ", "μ", "ν", "ξ", "ο", "π",
				"ρ", "σ", "τ", "υ", "φ", "χ", "ψ", "ω",
				"Α", "Β", "Γ", "Δ", "Ε", "Ζ", "Η", "Θ",
				"Ι", "Κ", "Λ", "Μ", "Ν", "Ξ", "Ο", "Π",
				"Ρ", "Σ", "Τ", "Υ", "Φ", "Χ", "Ψ", "Ω",
			},
		},
		{
			Name: "Math Operators",
			Symbols: []string{
				"±", "∓", "×", "÷", "∙", "·", "∘", "√",
				"∛", "∜", "∑", "∏", "∫", "∬", "∭", "∮",
				"∂", "∇", "∞", "≈", "≠", "≡", "≢", "≤",
				"≥", "≪", "≫", "∝", "∈", "∉", "⊂", "⊃",
				"⊆", "⊇", "∪", "∩", "∧", "∨", "⊕", "⊗",
			},
		},
		{
			Name: "Subscripts",
			Symbols: []string{
				"₀", "₁", "₂", "₃", "₄", "₅", "₆", "₇", "₈", "₉",
				"₊", "₋", "₌", "₍", "₎",
				"ₐ", "ₑ", "ₕ", "ᵢ", "ⱼ", "ₖ", "ₗ", "ₘ", "ₙ", "ₒ", "ₚ", "ᵣ", "ₛ", "ₜ", "ᵤ", "ᵥ", "ₓ",
			},
		},
		{
			Name: "Superscripts",
			Symbols: []string{
				"⁰", "¹", "²", "³", "⁴", "⁵", "⁶", "⁷", "⁸", "⁹",
				"⁺", "⁻", "⁼", "⁽", "⁾",
				"ᵃ", "ᵇ", "ᶜ", "ᵈ", "ᵉ", "ᶠ", "ᵍ", "ʰ", "ⁱ", "ʲ", "ᵏ", "ˡ", "ᵐ", "ⁿ", "ᵒ", "ᵖ", "ʳ", "ˢ", "ᵗ", "ᵘ", "ᵛ", "ʷ", "ˣ", "ʸ", "ᶻ",
			},
		},
		{
			Name: "Arrows",
			Symbols: []string{
				"←", "→", "↑", "↓", "↔", "↕", "↖", "↗",
				"↘", "↙", "⇐", "⇒", "⇑", "⇓", "⇔", "⇕",
				"↵", "↩", "↪", "↺", "↻", "⟵", "⟶", "⟷",
			},
		},
		{
			Name: "Units & Science",
			Symbols: []string{
				"°", "′", "″", "‰", "‱", "µ", "Å", "℃",
				"℉", "Ω", "℧", "℮", "ℏ", "℞", "℠", "™",
				"©", "®", "№", "℗", "℁", "℅", "℆", "⅍",
			},
		},
		{
			Name: "Fractions",
			Symbols: []string{
				"½", "⅓", "⅔", "¼", "¾", "⅕", "⅖", "⅗",
				"⅘", "⅙", "⅚", "⅛", "⅜", "⅝", "⅞", "⅐",
				"⅑", "⅒",
			},
		},
		{
			Name: "Currency",
			Symbols: []string{
				"$", "€", "£", "¥", "¢", "₹", "₽", "₿",
				"₩", "₪", "₫", "₭", "₮", "₯", "₱", "₲",
				"₳", "₴", "₵", "₸", "₺", "₼", "₾", "৳",
			},
		},
		{
			Name: "Punctuation & Misc",
			Symbols: []string{
				"…", "–", "—", "‹", "›", "«", "»", "•",
				"◦", "‣", "⁃", "※", "†", "‡", "§", "¶",
				"¦", "‖", "·", "⁂", "⁑", "⁎", "⁕", "❧",
			},
		},
		{
			Name: "Geometric Shapes",
			Symbols: []string{
				"■", "□", "▪", "▫", "▬", "▭", "▮", "▯",
				"▲", "△", "▴", "▵", "▶", "▷", "▸", "▹",
				"►", "▻", "▼", "▽", "▾", "▿", "◀", "◁",
				"◂", "◃", "◄", "◅", "◆", "◇", "◈", "◉",
				"○", "◌", "◍", "◎", "●", "◐", "◑", "◒",
			},
		},
		{
			Name: "Checkmarks & Ballots",
			Symbols: []string{
				"✓", "✔", "✕", "✖", "✗", "✘", "☐", "☑",
				"☒", "✙", "✚", "✛", "✜", "✝", "✞", "✟",
			},
		},
	}
}

// ShowSymbolPickerDialog shows a dialog for picking symbols
func ShowSymbolPickerDialog(parent fyne.Window, onSelect func(symbol string)) {
	categories := GetSymbolCategories()

	// Create category tabs
	tabs := container.NewAppTabs()

	for _, category := range categories {
		cat := category // Capture for closure
		grid := container.NewGridWrap(fyne.NewSize(40, 40))

		for _, symbol := range cat.Symbols {
			sym := symbol // Capture for closure
			btn := widget.NewButton(sym, func() {
				onSelect(sym)
			})
			grid.Add(btn)
		}

		scroll := container.NewScroll(grid)
		scroll.SetMinSize(fyne.NewSize(400, 300))

		tabs.Append(container.NewTabItem(cat.Name, scroll))
	}

	// Create a custom dialog
	content := container.NewBorder(
		widget.NewLabel("Click a symbol to insert it:"),
		nil, nil, nil,
		tabs,
	)

	d := dialog.NewCustom("Symbol Picker", "Close", content, parent)
	d.Resize(fyne.NewSize(500, 450))
	d.Show()
}

// ShowSubscriptPickerDialog shows a dialog for picking subscript characters
func ShowSubscriptPickerDialog(parent fyne.Window, onSelect func(symbol string)) {
	// Unicode subscript characters
	subscripts := []string{
		"₀", "₁", "₂", "₃", "₄", "₅", "₆", "₇", "₈", "₉",
		"₊", "₋", "₌", "₍", "₎",
		"ₐ", "ₑ", "ₕ", "ᵢ", "ⱼ", "ₖ", "ₗ", "ₘ", "ₙ", "ₒ", "ₚ", "ᵣ", "ₛ", "ₜ", "ᵤ", "ᵥ", "ₓ",
	}

	grid := container.NewGridWrap(fyne.NewSize(40, 40))
	for _, sym := range subscripts {
		s := sym // Capture for closure
		btn := widget.NewButton(s, func() {
			onSelect(s)
		})
		grid.Add(btn)
	}

	content := container.NewBorder(
		widget.NewLabel("Click a subscript character to insert:"),
		nil, nil, nil,
		container.NewScroll(grid),
	)

	d := dialog.NewCustom("Subscript Characters", "Close", content, parent)
	d.Resize(fyne.NewSize(400, 250))
	d.Show()
}

// ShowSuperscriptPickerDialog shows a dialog for picking superscript characters
func ShowSuperscriptPickerDialog(parent fyne.Window, onSelect func(symbol string)) {
	// Unicode superscript characters
	superscripts := []string{
		"⁰", "¹", "²", "³", "⁴", "⁵", "⁶", "⁷", "⁸", "⁹",
		"⁺", "⁻", "⁼", "⁽", "⁾",
		"ᵃ", "ᵇ", "ᶜ", "ᵈ", "ᵉ", "ᶠ", "ᵍ", "ʰ", "ⁱ", "ʲ", "ᵏ", "ˡ", "ᵐ", "ⁿ", "ᵒ", "ᵖ", "ʳ", "ˢ", "ᵗ", "ᵘ", "ᵛ", "ʷ", "ˣ", "ʸ", "ᶻ",
	}

	grid := container.NewGridWrap(fyne.NewSize(40, 40))
	for _, sym := range superscripts {
		s := sym // Capture for closure
		btn := widget.NewButton(s, func() {
			onSelect(s)
		})
		grid.Add(btn)
	}

	content := container.NewBorder(
		widget.NewLabel("Click a superscript character to insert:"),
		nil, nil, nil,
		container.NewScroll(grid),
	)

	d := dialog.NewCustom("Superscript Characters", "Close", content, parent)
	d.Resize(fyne.NewSize(400, 250))
	d.Show()
}
