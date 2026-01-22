package gui

import (
	"regexp"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// TagManager manages document tags
type TagManager struct {
	tags     map[string][]string // tag -> list of file paths
	fileTags map[string][]string // file path -> list of tags
}

// NewTagManager creates a new tag manager
func NewTagManager() *TagManager {
	return &TagManager{
		tags:     make(map[string][]string),
		fileTags: make(map[string][]string),
	}
}

// ExtractTagsFromContent extracts tags from document content
// Supports formats: #tag, tags: [tag1, tag2], tags: tag1, tag2
func (tm *TagManager) ExtractTagsFromContent(content string) []string {
	var tags []string
	seen := make(map[string]bool)

	// Match #hashtag style tags (not in code blocks)
	hashtagPattern := regexp.MustCompile(`(?m)^[^` + "`" + `]*#([a-zA-Z][a-zA-Z0-9_-]*)`)
	matches := hashtagPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			tag := strings.ToLower(match[1])
			if !seen[tag] && !isCommonWord(tag) {
				tags = append(tags, tag)
				seen[tag] = true
			}
		}
	}

	// Match YAML frontmatter tags: [tag1, tag2] or tags: tag1, tag2
	frontmatterPattern := regexp.MustCompile(`(?m)^tags:\s*\[([^\]]+)\]|^tags:\s*(.+)$`)
	fmMatches := frontmatterPattern.FindAllStringSubmatch(content, -1)
	for _, match := range fmMatches {
		var tagStr string
		if match[1] != "" {
			tagStr = match[1]
		} else if match[2] != "" {
			tagStr = match[2]
		}
		if tagStr != "" {
			parts := strings.Split(tagStr, ",")
			for _, part := range parts {
				tag := strings.ToLower(strings.TrimSpace(part))
				tag = strings.Trim(tag, "\"'")
				if tag != "" && !seen[tag] {
					tags = append(tags, tag)
					seen[tag] = true
				}
			}
		}
	}

	sort.Strings(tags)
	return tags
}

// isCommonWord returns true if the word is too common to be a tag
func isCommonWord(word string) bool {
	common := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"it": true, "this": true, "that": true, "these": true, "those": true,
		"todo": true, "fixme": true, "note": true, "warning": true,
	}
	return common[word]
}

// AddFileTag adds a tag to a file
func (tm *TagManager) AddFileTag(filePath, tag string) {
	tag = strings.ToLower(tag)

	// Add to tags map
	if !contains(tm.tags[tag], filePath) {
		tm.tags[tag] = append(tm.tags[tag], filePath)
	}

	// Add to fileTags map
	if !contains(tm.fileTags[filePath], tag) {
		tm.fileTags[filePath] = append(tm.fileTags[filePath], tag)
	}
}

// RemoveFileTag removes a tag from a file
func (tm *TagManager) RemoveFileTag(filePath, tag string) {
	tag = strings.ToLower(tag)

	// Remove from tags map
	if files, ok := tm.tags[tag]; ok {
		tm.tags[tag] = removeFromSlice(files, filePath)
		if len(tm.tags[tag]) == 0 {
			delete(tm.tags, tag)
		}
	}

	// Remove from fileTags map
	if tags, ok := tm.fileTags[filePath]; ok {
		tm.fileTags[filePath] = removeFromSlice(tags, tag)
	}
}

// GetFileTags returns all tags for a file
func (tm *TagManager) GetFileTags(filePath string) []string {
	return tm.fileTags[filePath]
}

// GetFilesByTag returns all files with a given tag
func (tm *TagManager) GetFilesByTag(tag string) []string {
	return tm.tags[strings.ToLower(tag)]
}

// GetAllTags returns all unique tags
func (tm *TagManager) GetAllTags() []string {
	var tags []string
	for tag := range tm.tags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// IndexFile indexes tags from a file's content
func (tm *TagManager) IndexFile(filePath, content string) {
	// Clear existing tags for this file
	if existingTags, ok := tm.fileTags[filePath]; ok {
		for _, tag := range existingTags {
			tm.tags[tag] = removeFromSlice(tm.tags[tag], filePath)
		}
		delete(tm.fileTags, filePath)
	}

	// Extract and add new tags
	tags := tm.ExtractTagsFromContent(content)
	for _, tag := range tags {
		tm.AddFileTag(filePath, tag)
	}
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeFromSlice(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// ShowTagsDialog shows the tags browser dialog
func ShowTagsDialog(window fyne.Window, tagManager *TagManager, onSelectFile func(path string)) {
	allTags := tagManager.GetAllTags()

	if len(allTags) == 0 {
		dialog.ShowInformation("Tags", "No tags found in documents.\n\nUse #hashtags or YAML frontmatter to add tags.", window)
		return
	}

	var d dialog.Dialog
	var fileList *widget.List
	var selectedTag string
	var filesForTag []string

	tagList := widget.NewList(
		func() int { return len(allTags) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("#tag"),
				widget.NewLabel("(0)"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(allTags) {
				box := obj.(*fyne.Container)
				tagLabel := box.Objects[0].(*widget.Label)
				countLabel := box.Objects[1].(*widget.Label)

				tag := allTags[id]
				tagLabel.SetText("#" + tag)
				countLabel.SetText("(" + intToString(len(tagManager.GetFilesByTag(tag))) + ")")
			}
		},
	)

	fileList = widget.NewList(
		func() int { return len(filesForTag) },
		func() fyne.CanvasObject {
			return widget.NewLabel("filename.md")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(filesForTag) {
				obj.(*widget.Label).SetText(filesForTag[id])
			}
		},
	)

	tagList.OnSelected = func(id widget.ListItemID) {
		if id < len(allTags) {
			selectedTag = allTags[id]
			filesForTag = tagManager.GetFilesByTag(selectedTag)
			fileList.Refresh()
		}
	}

	fileList.OnSelected = func(id widget.ListItemID) {
		if id < len(filesForTag) {
			path := filesForTag[id]
			d.Hide()
			if onSelectFile != nil {
				onSelectFile(path)
			}
		}
	}

	tagScroll := container.NewScroll(tagList)
	tagScroll.SetMinSize(fyne.NewSize(200, 350))

	fileScroll := container.NewScroll(fileList)
	fileScroll.SetMinSize(fyne.NewSize(300, 350))

	split := container.NewHSplit(
		container.NewBorder(widget.NewLabel("Tags"), nil, nil, nil, tagScroll),
		container.NewBorder(widget.NewLabel("Documents"), nil, nil, nil, fileScroll),
	)
	split.Offset = 0.35

	d = dialog.NewCustom("Browse by Tag", "Close", split, window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
