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
			return widget.NewLabel("")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			label := node.(*widget.Label)

			if uid == "" {
				label.SetText("")
				return
			}

			if entry, ok := n.entryMap[uid]; ok {
				// Just show the heading text - tree widget handles hierarchy
				label.SetText(entry.Text)
				label.Wrapping = fyne.TextTruncate
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
	n.entryMap = make(map[string]*TOCEntry)
	n.buildEntryMap()
	n.tree.Refresh()
}
