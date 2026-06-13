package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Snippet represents a text snippet template
type Snippet struct {
	Name    string
	Content string
}

// DefaultSnippets returns the built-in snippets
func DefaultSnippets() []Snippet {
	return []Snippet{
		{
			Name:    "Code Block (Go)",
			Content: "```go\npackage main\n\nfunc main() {\n    \n}\n```",
		},
		{
			Name:    "Code Block (Python)",
			Content: "```python\ndef main():\n    pass\n\nif __name__ == \"__main__\":\n    main()\n```",
		},
		{
			Name:    "Code Block (JavaScript)",
			Content: "```javascript\nfunction main() {\n    \n}\n\nmain();\n```",
		},
		{
			Name:    "Code Block (Bash)",
			Content: "```bash\n#!/bin/bash\n\n```",
		},
		{
			Name:    "Table (3x3)",
			Content: "| Column 1 | Column 2 | Column 3 |\n| --- | --- | --- |\n| Row 1 | Data | Data |\n| Row 2 | Data | Data |\n| Row 3 | Data | Data |",
		},
		{
			Name:    "Task List",
			Content: "- [ ] Task 1\n- [ ] Task 2\n- [ ] Task 3",
		},
		{
			Name:    "Collapsible Section",
			Content: "<details>\n<summary>Click to expand</summary>\n\nContent here...\n\n</details>",
		},
		{
			Name:    "Blockquote",
			Content: "> Quote text here\n>\n> — Author Name",
		},
		{
			Name:    "Image with Caption",
			Content: "![Alt text](image.png)\n*Image caption*",
		},
		{
			Name:    "Link Reference",
			Content: "[link text][ref]\n\n[ref]: https://example.com \"Title\"",
		},
		{
			Name:    "Footnote",
			Content: "Text with footnote[^1].\n\n[^1]: Footnote content here.",
		},
		{
			Name:    "Definition List",
			Content: "Term 1\n: Definition 1\n\nTerm 2\n: Definition 2",
		},
		{
			Name:    "Mermaid Flowchart",
			Content: "```mermaid\nflowchart TD\n    A[Start] --> B{Decision}\n    B -->|Yes| C[Result 1]\n    B -->|No| D[Result 2]\n```",
		},
		{
			Name:    "Mermaid Sequence",
			Content: "```mermaid\nsequenceDiagram\n    participant A\n    participant B\n    A->>B: Request\n    B-->>A: Response\n```",
		},
		{
			Name:    "Math Block",
			Content: "$$\n\\sum_{i=1}^{n} x_i = x_1 + x_2 + \\cdots + x_n\n$$",
		},
		{
			Name:    "YAML Front Matter",
			Content: "---\ntitle: Document Title\nauthor: Author Name\ndate: 2024-01-01\ntags: [tag1, tag2]\n---\n",
		},
	}
}

// ShowSnippetsDialog shows the snippets selection dialog
func ShowSnippetsDialog(window fyne.Window, onInsert func(content string)) {
	snippets := DefaultSnippets()

	list := widget.NewList(
		func() int { return len(snippets) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("Snippet Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Preview..."),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(snippets) {
				return
			}
			box := obj.(*fyne.Container)
			title := box.Objects[0].(*widget.Label)
			preview := box.Objects[1].(*widget.Label)

			title.SetText(snippets[id].Name)
			previewText := snippets[id].Content
			if len(previewText) > 50 {
				previewText = previewText[:50] + "..."
			}
			preview.SetText(previewText)
		},
	)

	var d dialog.Dialog
	list.OnSelected = func(id widget.ListItemID) {
		if id < len(snippets) && onInsert != nil {
			onInsert(snippets[id].Content)
			d.Hide()
		}
	}

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(400, 350))

	d = dialog.NewCustom("Insert Snippet", "Close", scroll, window)
	d.Resize(dialogSizeList)
	d.Show()
}
