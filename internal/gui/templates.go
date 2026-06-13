package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// DocumentTemplate represents a document template
type DocumentTemplate struct {
	Name        string
	Description string
	Content     string
	Icon        fyne.Resource
}

// GetDefaultTemplates returns the default document templates
func GetDefaultTemplates() []DocumentTemplate {
	today := time.Now().Format("2006-01-02")

	return []DocumentTemplate{
		{
			Name:        "Blank Document",
			Description: "Start with a clean slate",
			Content:     "# Untitled\n\n",
		},
		{
			Name:        "Meeting Notes",
			Description: "Template for meeting notes",
			Content: fmt.Sprintf(`# Meeting Notes - %s

## Attendees
-

## Agenda
1.

## Discussion Points


## Action Items
- [ ]

## Next Meeting

`, today),
		},
		{
			Name:        "Blog Post",
			Description: "Template for a blog post",
			Content: fmt.Sprintf(`---
title: "Your Title Here"
date: %s
author:
tags: []
draft: true
---

# Your Title Here

## Introduction

Write your introduction here...

## Main Content

### Section 1

### Section 2

## Conclusion

## References

`, today),
		},
		{
			Name:        "Project README",
			Description: "Template for project documentation",
			Content: `# Project Name

Brief description of the project.

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

` + "```bash\n# Installation commands\n```" + `

## Usage

` + "```bash\n# Usage examples\n```" + `

## Configuration

Describe configuration options here.

## Contributing

Guidelines for contributing to the project.

## License

Specify the license here.
`,
		},
		{
			Name:        "Technical Specification",
			Description: "Template for technical documentation",
			Content: fmt.Sprintf(`# Technical Specification

**Document Version:** 1.0
**Date:** %s
**Author:**

## Overview

Brief overview of what this document covers.

## Requirements

### Functional Requirements

1.

### Non-Functional Requirements

1.

## Architecture

### System Components

### Data Flow

## Implementation Details

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
|          |        |             |

### Database Schema

## Testing Strategy

## Deployment

## Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
|      |            |        |            |

## Timeline

## Appendix

`, today),
		},
		{
			Name:        "Daily Journal",
			Description: "Template for daily journaling",
			Content: fmt.Sprintf(`# Journal - %s

## Gratitude
Three things I'm grateful for today:
1.
2.
3.

## Today's Goals
- [ ]
- [ ]
- [ ]

## Notes & Thoughts


## Evening Reflection

### What went well?


### What could be improved?


### Tomorrow's priorities

`, today),
		},
		{
			Name:        "Weekly Review",
			Description: "Template for weekly review",
			Content: fmt.Sprintf(`# Weekly Review - Week of %s

## Accomplishments
-

## Challenges
-

## Lessons Learned
-

## Next Week's Goals
1.
2.
3.

## Notes

`, today),
		},
		{
			Name:        "Code Review Notes",
			Description: "Template for code review documentation",
			Content: fmt.Sprintf(`# Code Review Notes

**Date:** %s
**Reviewer:**
**PR/MR Link:**

## Summary

Brief description of changes being reviewed.

## Files Reviewed

- [ ] file1.go
- [ ] file2.go

## Findings

### Critical Issues

### Suggestions

### Positive Notes

## Checklist

- [ ] Code follows project style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Error handling is appropriate

## Approval Status

[ ] Approved / [ ] Changes Requested

`, today),
		},
	}
}

// ShowTemplatesDialog shows the document templates dialog
func ShowTemplatesDialog(window fyne.Window, onCreate func(content string)) {
	templates := GetDefaultTemplates()

	var d dialog.Dialog

	list := widget.NewList(
		func() int { return len(templates) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("Template Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Template description..."),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(templates) {
				box := obj.(*fyne.Container)
				nameLabel := box.Objects[0].(*widget.Label)
				descLabel := box.Objects[1].(*widget.Label)

				t := templates[id]
				nameLabel.SetText(t.Name)
				descLabel.SetText(t.Description)
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(templates) {
			t := templates[id]
			d.Hide()
			if onCreate != nil {
				onCreate(t.Content)
			}
		}
	}

	scroll := container.NewScroll(list)
	scroll.SetMinSize(fyne.NewSize(400, 350))

	content := container.NewBorder(
		widget.NewLabel("Select a template to create a new document:"),
		nil, nil, nil,
		scroll,
	)

	d = dialog.NewCustom("Document Templates", "Close", content, window)
	d.Resize(dialogSizeList)
	d.Show()
}
