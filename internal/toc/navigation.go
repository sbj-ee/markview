package toc

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Navigator handles TOC navigation and scrolling
type Navigator struct {
	tree          *widget.Tree
	scrollContent *container.Scroll
	entries       []*TOCEntry
	entryMap      map[string]*TOCEntry
}

// NewNavigator creates a new TOC navigator
func NewNavigator(entries []*TOCEntry, scrollContent *container.Scroll) *Navigator {
	nav := &Navigator{
		entries:       entries,
		scrollContent: scrollContent,
		entryMap:      make(map[string]*TOCEntry),
	}

	nav.buildEntryMap()
	nav.tree = nav.createTree()

	// Open all tree nodes by default
	nav.openAllNodes()

	return nav
}

// openAllNodes recursively opens all nodes in the tree
func (n *Navigator) openAllNodes() {
	var openNode func(string)
	openNode = func(uid string) {
		if uid == "" {
			// Open root level nodes
			for i := range n.entries {
				childUID := n.getNodeID("", i)
				n.tree.OpenBranch(childUID)
				openNode(childUID)
			}
		} else if entry, ok := n.entryMap[uid]; ok {
			// Open this node's children
			for i := range entry.Children {
				childUID := n.getNodeID(uid, i)
				n.tree.OpenBranch(childUID)
				openNode(childUID)
			}
		}
	}
	openNode("")
}

// GetTree returns the TOC tree widget
func (n *Navigator) GetTree() *widget.Tree {
	return n.tree
}

// buildEntryMap builds a map of UIDs to entries
func (n *Navigator) buildEntryMap() {
	var buildMap func([]*TOCEntry, string)
	buildMap = func(entries []*TOCEntry, parentID string) {
		for i, entry := range entries {
			uid := n.getNodeID(parentID, i)
			n.entryMap[uid] = entry
			if len(entry.Children) > 0 {
				buildMap(entry.Children, uid)
			}
		}
	}
	buildMap(n.entries, "")
}

// createTree creates a Fyne tree widget from TOC entries
func (n *Navigator) createTree() *widget.Tree {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			if uid == "" {
				// Root level
				uids := make([]string, len(n.entries))
				for i := range n.entries {
					uids[i] = n.getNodeID("", i)
				}
				return uids
			}

			// Find entry and return children UIDs
			if entry, ok := n.entryMap[uid]; ok {
				uids := make([]string, len(entry.Children))
				for i := range entry.Children {
					uids[i] = n.getNodeID(uid, i)
				}
				return uids
			}

			return []string{}
		},
		IsBranch: func(uid string) bool {
			if uid == "" {
				return true
			}

			if entry, ok := n.entryMap[uid]; ok {
				return len(entry.Children) > 0
			}

			return false
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			// Use RichText with smaller font for TOC items
			rt := widget.NewRichText(&widget.TextSegment{
				Text: "",
				Style: widget.RichTextStyle{
					SizeName: theme.SizeNameCaptionText,
				},
			})
			rt.Truncation = fyne.TextTruncateEllipsis
			return rt
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			rt := node.(*widget.RichText)

			if uid == "" {
				rt.Segments = []widget.RichTextSegment{
					&widget.TextSegment{
						Text: "",
						Style: widget.RichTextStyle{
							SizeName: theme.SizeNameCaptionText,
						},
					},
				}
				rt.Refresh()
				return
			}

			if entry, ok := n.entryMap[uid]; ok {
				// Show the heading text with smaller font
				rt.Segments = []widget.RichTextSegment{
					&widget.TextSegment{
						Text: entry.Text,
						Style: widget.RichTextStyle{
							SizeName: theme.SizeNameCaptionText,
						},
					},
				}
				rt.Refresh()
			}
		},
		OnSelected: func(uid string) {
			// Handle TOC item selection - scroll to heading
			if entry, ok := n.entryMap[uid]; ok {
				n.scrollToHeading(entry)
			}
		},
	}

	return tree
}

// getNodeID generates a unique node ID
func (n *Navigator) getNodeID(parentID string, index int) string {
	if parentID == "" {
		return fmt.Sprintf("node-%d", index)
	}
	return fmt.Sprintf("%s-%d", parentID, index)
}

// scrollToHeading scrolls the content to show the specified heading
func (n *Navigator) scrollToHeading(entry *TOCEntry) {
	if n.scrollContent == nil || n.scrollContent.Content == nil {
		return
	}

	// Get the content container - may be wrapped in padding/border containers
	vbox := n.findContentVBox(n.scrollContent.Content)
	if vbox == nil {
		return
	}

	// Search for the widget with matching heading text
	var targetOffset float32 = 0
	found := false

	n.searchForHeading(vbox.Objects, entry.Text, &targetOffset, &found)

	if found {
		// Scroll to the calculated offset
		n.scrollContent.Offset = fyne.NewPos(0, targetOffset)
		n.scrollContent.Refresh()
	}
}

// findContentVBox recursively finds the VBox with the actual content
func (n *Navigator) findContentVBox(obj fyne.CanvasObject) *fyne.Container {
	container, ok := obj.(*fyne.Container)
	if !ok {
		return nil
	}

	// Check if this container has RichText children (actual content)
	for _, child := range container.Objects {
		if _, ok := child.(*widget.RichText); ok {
			return container
		}
	}

	// Otherwise, search nested containers
	for _, child := range container.Objects {
		if found := n.findContentVBox(child); found != nil {
			return found
		}
	}

	return nil
}

// searchForHeading recursively searches for a heading in the widget tree
func (n *Navigator) searchForHeading(objects []fyne.CanvasObject, targetText string, offset *float32, found *bool) {
	for _, obj := range objects {
		if *found {
			return
		}

		// Check if this is a RichText widget (headings are RichText)
		if rt, ok := obj.(*widget.RichText); ok {
			// Check if the text matches our heading
			if n.richTextContains(rt, targetText) {
				*found = true
				return
			}
		}

		// Check nested containers
		if container, ok := obj.(*fyne.Container); ok {
			n.searchForHeading(container.Objects, targetText, offset, found)
			if *found {
				return
			}
		}

		if !*found {
			// Add this widget's height to the offset
			*offset += obj.MinSize().Height
		}
	}
}

// richTextContains checks if a RichText widget contains the target text
func (n *Navigator) richTextContains(rt *widget.RichText, targetText string) bool {
	for _, seg := range rt.Segments {
		if para, ok := seg.(*widget.ParagraphSegment); ok {
			for _, text := range para.Texts {
				if textSeg, ok := text.(*widget.TextSegment); ok {
					if textSeg.Text == targetText {
						return true
					}
				}
			}
		}
		// Also check direct TextSegments
		if textSeg, ok := seg.(*widget.TextSegment); ok {
			if textSeg.Text == targetText {
				return true
			}
		}
	}
	return false
}

// Update updates the navigator with new TOC entries
func (n *Navigator) Update(entries []*TOCEntry) {
	n.entries = entries
	n.entryMap = make(map[string]*TOCEntry)
	n.buildEntryMap()
	n.tree.Refresh()
}
