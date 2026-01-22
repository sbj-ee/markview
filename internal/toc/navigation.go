package toc

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Navigator handles TOC navigation and scrolling
type Navigator struct {
	tree          *widget.Tree
	scrollContent *container.Scroll
	entries       []*TOCEntry
}

// NewNavigator creates a new TOC navigator
func NewNavigator(entries []*TOCEntry, scrollContent *container.Scroll) *Navigator {
	nav := &Navigator{
		entries:       entries,
		scrollContent: scrollContent,
	}

	nav.tree = nav.createTree()
	return nav
}

// GetTree returns the TOC tree widget
func (n *Navigator) GetTree() *widget.Tree {
	return n.tree
}

// createTree creates a Fyne tree widget from TOC entries
func (n *Navigator) createTree() *widget.Tree {
	// Build a map for quick lookup
	entryMap := make(map[string]*TOCEntry)
	var buildMap func([]*TOCEntry, string)
	buildMap = func(entries []*TOCEntry, parentID string) {
		for i, entry := range entries {
			id := n.getNodeID(parentID, i)
			entryMap[id] = entry
			if len(entry.Children) > 0 {
				buildMap(entry.Children, id)
			}
		}
	}
	buildMap(n.entries, "")

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
			if entry, ok := entryMap[uid]; ok {
				uids := make([]string, len(entry.Children))
				for i := range entry.Children {
					uids[i] = n.getNodeID(uid, i)
				}
				return uids
			}

			return nil
		},
		IsBranch: func(uid string) bool {
			if uid == "" {
				return true
			}

			if entry, ok := entryMap[uid]; ok {
				return len(entry.Children) > 0
			}

			return false
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			label := node.(*widget.Label)

			if uid == "" {
				label.SetText("")
				return
			}

			if entry, ok := entryMap[uid]; ok {
				// Indent based on level
				indent := ""
				for i := 1; i < entry.Level; i++ {
					indent += "  "
				}
				label.SetText(indent + entry.Text)
			}
		},
		OnSelected: func(uid string) {
			// Handle TOC item selection - scroll to heading
			if entry, ok := entryMap[uid]; ok {
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
	// This is a placeholder implementation
	// In a real implementation, we would:
	// 1. Track the position of each heading in the rendered content
	// 2. Calculate the scroll offset to bring the heading into view
	// 3. Scroll to that offset

	// For now, we'll just scroll to the top as a basic implementation
	n.scrollContent.ScrollToTop()

	// TODO: Implement proper scroll-to-heading functionality
	// This requires tracking heading positions during rendering
}

// Update updates the navigator with new TOC entries
func (n *Navigator) Update(entries []*TOCEntry) {
	n.entries = entries
	n.tree = n.createTree()
}
