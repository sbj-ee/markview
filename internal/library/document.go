package library

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Document represents a markdown document in the library
type Document struct {
	Path      string
	Title     string
	Tags      []string
	Category  string
	ModTime   time.Time
	Preview   string
	WordCount int
	Starred   bool
}

// DocumentLibrary manages a collection of documents
type DocumentLibrary struct {
	RootPath  string
	Documents []*Document
}

// NewDocumentLibrary creates a new document library
func NewDocumentLibrary(rootPath string) *DocumentLibrary {
	return &DocumentLibrary{
		RootPath:  rootPath,
		Documents: make([]*Document, 0),
	}
}

// Scan scans the root path for markdown files
func (lib *DocumentLibrary) Scan() error {
	lib.Documents = make([]*Document, 0)

	err := filepath.Walk(lib.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Only process markdown files
		if !info.IsDir() && isMarkdownFile(path) {
			doc, err := lib.parseDocument(path)
			if err == nil {
				lib.Documents = append(lib.Documents, doc)
			}
		}

		return nil
	})

	return err
}

// parseDocument parses a markdown file and extracts metadata
func (lib *DocumentLibrary) parseDocument(path string) (*Document, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := &Document{
		Path:     path,
		ModTime:  info.ModTime(),
		Category: lib.getCategory(path),
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	inFrontMatter := false
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Check for YAML front matter
		if lineCount == 1 && line == "---" {
			inFrontMatter = true
			continue
		}

		if inFrontMatter {
			if line == "---" {
				inFrontMatter = false
				continue
			}
			// Parse front matter
			if strings.HasPrefix(line, "title:") {
				doc.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
				doc.Title = strings.Trim(doc.Title, "\"'")
			}
			if strings.HasPrefix(line, "tags:") {
				tagStr := strings.TrimSpace(strings.TrimPrefix(line, "tags:"))
				tagStr = strings.Trim(tagStr, "[]")
				for _, tag := range strings.Split(tagStr, ",") {
					tag = strings.TrimSpace(tag)
					tag = strings.Trim(tag, "\"'")
					if tag != "" {
						doc.Tags = append(doc.Tags, tag)
					}
				}
			}
			if strings.HasPrefix(line, "category:") {
				doc.Category = strings.TrimSpace(strings.TrimPrefix(line, "category:"))
				doc.Category = strings.Trim(doc.Category, "\"'")
			}
			continue
		}

		// Count words
		words := strings.Fields(line)
		wordCount += len(words)

		// Collect preview lines (first 5 non-empty lines after front matter)
		if len(lines) < 5 && line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	doc.WordCount = wordCount
	doc.Preview = strings.Join(lines, " ")
	if len(doc.Preview) > 200 {
		doc.Preview = doc.Preview[:200] + "..."
	}

	// If no title was found in front matter, try to extract from first heading
	if doc.Title == "" {
		doc.Title = lib.extractTitleFromFile(path)
	}

	// Fallback to filename
	if doc.Title == "" {
		doc.Title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	return doc, scanner.Err()
}

// extractTitleFromFile extracts the first heading from a markdown file
func (lib *DocumentLibrary) extractTitleFromFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	headingRegex := regexp.MustCompile(`^#+\s+(.+)$`)
	inFrontMatter := false

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Skip front matter
		if lineCount == 1 && line == "---" {
			inFrontMatter = true
			continue
		}
		if inFrontMatter {
			if line == "---" {
				inFrontMatter = false
			}
			continue
		}

		// Look for heading
		matches := headingRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// getCategory gets the category (parent folder name) for a document
func (lib *DocumentLibrary) getCategory(path string) string {
	relPath, err := filepath.Rel(lib.RootPath, path)
	if err != nil {
		return ""
	}

	dir := filepath.Dir(relPath)
	if dir == "." {
		return "Uncategorized"
	}

	return dir
}

// GetCategories returns all unique categories
func (lib *DocumentLibrary) GetCategories() []string {
	categoryMap := make(map[string]bool)
	for _, doc := range lib.Documents {
		categoryMap[doc.Category] = true
	}

	categories := make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	return categories
}

// GetTags returns all unique tags
func (lib *DocumentLibrary) GetTags() []string {
	tagMap := make(map[string]bool)
	for _, doc := range lib.Documents {
		for _, tag := range doc.Tags {
			tagMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}

// FilterByCategory returns documents in a specific category
func (lib *DocumentLibrary) FilterByCategory(category string) []*Document {
	var docs []*Document
	for _, doc := range lib.Documents {
		if doc.Category == category {
			docs = append(docs, doc)
		}
	}
	return docs
}

// FilterByTag returns documents with a specific tag
func (lib *DocumentLibrary) FilterByTag(tag string) []*Document {
	var docs []*Document
	for _, doc := range lib.Documents {
		for _, t := range doc.Tags {
			if t == tag {
				docs = append(docs, doc)
				break
			}
		}
	}
	return docs
}

// Search searches documents by title or content
func (lib *DocumentLibrary) Search(query string) []*Document {
	query = strings.ToLower(query)
	var docs []*Document
	for _, doc := range lib.Documents {
		if strings.Contains(strings.ToLower(doc.Title), query) ||
			strings.Contains(strings.ToLower(doc.Preview), query) {
			docs = append(docs, doc)
		}
	}
	return docs
}

// isMarkdownFile checks if a file is a markdown file
func isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mkd"
}

// FilterStarred returns only starred documents
func (lib *DocumentLibrary) FilterStarred() []*Document {
	var docs []*Document
	for _, doc := range lib.Documents {
		if doc.Starred {
			docs = append(docs, doc)
		}
	}
	return docs
}

// ToggleStarred toggles the starred status of a document
func (lib *DocumentLibrary) ToggleStarred(path string) bool {
	for _, doc := range lib.Documents {
		if doc.Path == path {
			doc.Starred = !doc.Starred
			return doc.Starred
		}
	}
	return false
}

// SetStarred sets the starred status for a document
func (lib *DocumentLibrary) SetStarred(path string, starred bool) {
	for _, doc := range lib.Documents {
		if doc.Path == path {
			doc.Starred = starred
			return
		}
	}
}

// GetStarredPaths returns paths of all starred documents
func (lib *DocumentLibrary) GetStarredPaths() []string {
	var paths []string
	for _, doc := range lib.Documents {
		if doc.Starred {
			paths = append(paths, doc.Path)
		}
	}
	return paths
}

// LoadStarredFromPaths marks documents as starred based on provided paths
func (lib *DocumentLibrary) LoadStarredFromPaths(paths []string) {
	pathSet := make(map[string]bool)
	for _, p := range paths {
		pathSet[p] = true
	}
	for _, doc := range lib.Documents {
		doc.Starred = pathSet[doc.Path]
	}
}

// SortBy represents the sort field
type SortBy int

const (
	SortByDate SortBy = iota
	SortByName
	SortByWordCount
)

// SortDocuments sorts documents by the specified field
func SortDocuments(docs []*Document, sortBy SortBy, ascending bool) {
	switch sortBy {
	case SortByName:
		if ascending {
			sort.Slice(docs, func(i, j int) bool {
				return strings.ToLower(docs[i].Title) < strings.ToLower(docs[j].Title)
			})
		} else {
			sort.Slice(docs, func(i, j int) bool {
				return strings.ToLower(docs[i].Title) > strings.ToLower(docs[j].Title)
			})
		}
	case SortByWordCount:
		if ascending {
			sort.Slice(docs, func(i, j int) bool {
				return docs[i].WordCount < docs[j].WordCount
			})
		} else {
			sort.Slice(docs, func(i, j int) bool {
				return docs[i].WordCount > docs[j].WordCount
			})
		}
	default: // SortByDate
		if ascending {
			sort.Slice(docs, func(i, j int) bool {
				return docs[i].ModTime.Before(docs[j].ModTime)
			})
		} else {
			sort.Slice(docs, func(i, j int) bool {
				return docs[i].ModTime.After(docs[j].ModTime)
			})
		}
	}
}
