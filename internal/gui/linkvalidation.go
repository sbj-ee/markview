package gui

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// LinkStatus represents the status of a link
type LinkStatus int

const (
	LinkValid LinkStatus = iota
	LinkBroken
	LinkUnchecked
	LinkExternal
)

// LinkInfo represents information about a link
type LinkInfo struct {
	Text       string
	Target     string
	LineNumber int
	Status     LinkStatus
	StatusMsg  string
}

// LinkValidator validates links in markdown documents
type LinkValidator struct {
	basePath      string
	links         []LinkInfo
	httpClient    *http.Client
	checkExternal bool
}

// NewLinkValidator creates a new link validator
func NewLinkValidator(basePath string) *LinkValidator {
	return &LinkValidator{
		basePath:      basePath,
		checkExternal: false,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// SetCheckExternal enables or disables checking external URLs
func (lv *LinkValidator) SetCheckExternal(check bool) {
	lv.checkExternal = check
}

// ValidateContent validates all links in the content
func (lv *LinkValidator) ValidateContent(content string) []LinkInfo {
	lv.links = []LinkInfo{}

	// Pattern for markdown links: [text](target)
	linkPattern := regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		matches := linkPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				linkText := match[1]
				linkTarget := match[2]

				// Skip anchor-only links
				if strings.HasPrefix(linkTarget, "#") {
					continue
				}

				info := LinkInfo{
					Text:       linkText,
					Target:     linkTarget,
					LineNumber: lineNum + 1,
					Status:     LinkUnchecked,
				}

				// Check if it's an external link
				if strings.HasPrefix(linkTarget, "http://") || strings.HasPrefix(linkTarget, "https://") {
					if lv.checkExternal {
						lv.validateExternalLink(&info)
					} else {
						info.Status = LinkExternal
						info.StatusMsg = "External link (not checked)"
					}
				} else {
					// Local file link
					lv.validateLocalLink(&info)
				}

				lv.links = append(lv.links, info)
			}
		}
	}

	return lv.links
}

// validateLocalLink validates a local file link
func (lv *LinkValidator) validateLocalLink(info *LinkInfo) {
	target := info.Target

	// Remove anchor from target
	if idx := strings.Index(target, "#"); idx != -1 {
		target = target[:idx]
	}

	// Handle relative paths
	target = strings.TrimPrefix(target, "./")

	// Build full path
	fullPath := filepath.Join(lv.basePath, target)

	// Check if file exists
	if _, err := os.Stat(fullPath); err == nil {
		info.Status = LinkValid
		info.StatusMsg = "File exists"
	} else {
		info.Status = LinkBroken
		info.StatusMsg = "File not found"
	}
}

// validateExternalLink validates an external URL
func (lv *LinkValidator) validateExternalLink(info *LinkInfo) {
	resp, err := lv.httpClient.Head(info.Target)
	if err != nil {
		info.Status = LinkBroken
		info.StatusMsg = "Connection failed"
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		info.Status = LinkValid
		info.StatusMsg = "OK"
	} else {
		info.Status = LinkBroken
		info.StatusMsg = "HTTP " + resp.Status
	}
}

// GetBrokenLinks returns only broken links
func (lv *LinkValidator) GetBrokenLinks() []LinkInfo {
	var broken []LinkInfo
	for _, link := range lv.links {
		if link.Status == LinkBroken {
			broken = append(broken, link)
		}
	}
	return broken
}

// ShowLinkValidationDialog shows the link validation results
func ShowLinkValidationDialog(window fyne.Window, basePath, content string, onNavigate func(line int)) {
	validator := NewLinkValidator(basePath)

	// Create progress dialog
	progress := dialog.NewCustomWithoutButtons("Validating Links...",
		widget.NewProgressBarInfinite(), window)
	progress.Show()

	// Validate in background
	go func() {
		links := validator.ValidateContent(content)
		progress.Hide()

		// Show results on main thread
		fyne.CurrentApp().Driver().CanvasForObject(window.Content()).Content().Refresh()

		if len(links) == 0 {
			dialog.ShowInformation("Link Validation", "No links found in the document.", window)
			return
		}

		brokenCount := 0
		for _, l := range links {
			if l.Status == LinkBroken {
				brokenCount++
			}
		}

		var d dialog.Dialog

		list := widget.NewList(
			func() int { return len(links) },
			func() fyne.CanvasObject {
				return container.NewVBox(
					container.NewHBox(
						widget.NewLabel(""), // Status icon placeholder
						widget.NewLabel("Link text"),
					),
					widget.NewLabel("target → status"),
				)
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				if id < len(links) {
					box := obj.(*fyne.Container)
					header := box.Objects[0].(*fyne.Container)
					statusLabel := header.Objects[0].(*widget.Label)
					textLabel := header.Objects[1].(*widget.Label)
					detailLabel := box.Objects[1].(*widget.Label)

					link := links[id]

					// Status indicator
					switch link.Status {
					case LinkValid:
						statusLabel.SetText("[OK]")
					case LinkBroken:
						statusLabel.SetText("[X]")
					case LinkExternal:
						statusLabel.SetText("[?]")
					default:
						statusLabel.SetText("[ ]")
					}

					textLabel.SetText(link.Text)

					// Truncate long targets
					target := link.Target
					if len(target) > 50 {
						target = target[:50] + "..."
					}
					detailLabel.SetText("Line " + intToString(link.LineNumber) + ": " + target + " → " + link.StatusMsg)
				}
			},
		)

		list.OnSelected = func(id widget.ListItemID) {
			if id < len(links) {
				link := links[id]
				d.Hide()
				if onNavigate != nil {
					onNavigate(link.LineNumber)
				}
			}
		}

		// Summary
		summary := widget.NewLabel("")
		if brokenCount == 0 {
			summary.SetText("All " + intToString(len(links)) + " links are valid")
		} else {
			summary.SetText(intToString(brokenCount) + " broken link(s) found out of " + intToString(len(links)))
		}

		scroll := container.NewScroll(list)
		scroll.SetMinSize(fyne.NewSize(500, 350))

		content := container.NewBorder(
			summary,
			nil, nil, nil,
			scroll,
		)

		d = dialog.NewCustom("Link Validation Results", "Close", content, window)
		d.Resize(fyne.NewSize(550, 450))
		d.Show()
	}()
}
