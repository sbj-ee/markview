package gui

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// SpellChecker provides spell checking functionality
type SpellChecker struct {
	enabled     bool
	cache       map[string]bool // word -> isCorrect
	cacheMutex  sync.RWMutex
	hasAspell   bool
	customWords map[string]bool // user-added words
}

// NewSpellChecker creates a new spell checker
func NewSpellChecker() *SpellChecker {
	sc := &SpellChecker{
		enabled:     true,
		cache:       make(map[string]bool),
		customWords: make(map[string]bool),
	}

	// Check if aspell is available
	_, err := exec.LookPath("aspell")
	sc.hasAspell = err == nil

	return sc
}

// IsEnabled returns whether spell checking is enabled
func (sc *SpellChecker) IsEnabled() bool {
	return sc.enabled && sc.hasAspell
}

// SetEnabled enables or disables spell checking
func (sc *SpellChecker) SetEnabled(enabled bool) {
	sc.enabled = enabled
}

// HasAspell returns whether aspell is available
func (sc *SpellChecker) HasAspell() bool {
	return sc.hasAspell
}

// AddCustomWord adds a word to the custom dictionary
func (sc *SpellChecker) AddCustomWord(word string) {
	sc.cacheMutex.Lock()
	defer sc.cacheMutex.Unlock()
	sc.customWords[strings.ToLower(word)] = true
	sc.cache[strings.ToLower(word)] = true
}

// CheckWord checks if a single word is spelled correctly
func (sc *SpellChecker) CheckWord(word string) bool {
	if !sc.enabled || !sc.hasAspell {
		return true
	}

	// Normalize word
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" || len(word) < 2 {
		return true
	}

	// Skip words that look like code, URLs, or paths
	if isCodeLikeWord(word) {
		return true
	}

	// Check cache first
	sc.cacheMutex.RLock()
	if correct, found := sc.cache[word]; found {
		sc.cacheMutex.RUnlock()
		return correct
	}
	sc.cacheMutex.RUnlock()

	// Check custom words
	sc.cacheMutex.RLock()
	if sc.customWords[word] {
		sc.cacheMutex.RUnlock()
		return true
	}
	sc.cacheMutex.RUnlock()

	// Use aspell to check
	correct := sc.checkWithAspell(word)

	// Cache result
	sc.cacheMutex.Lock()
	sc.cache[word] = correct
	sc.cacheMutex.Unlock()

	return correct
}

// checkWithAspell uses aspell to check a word
func (sc *SpellChecker) checkWithAspell(word string) bool {
	cmd := exec.Command("aspell", "-a")
	cmd.Stdin = strings.NewReader(word)
	output, err := cmd.Output()
	if err != nil {
		return true // Assume correct on error
	}

	// aspell returns "*" for correct words, "&" or "#" for incorrect
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "*") {
			return true
		}
		if strings.HasPrefix(line, "&") || strings.HasPrefix(line, "#") {
			return false
		}
	}
	return true
}

// GetMisspelledWords returns a list of misspelled words with their positions
func (sc *SpellChecker) GetMisspelledWords(text string) []MisspelledWord {
	if !sc.enabled || !sc.hasAspell {
		return nil
	}

	var misspelled []MisspelledWord
	wordRegex := regexp.MustCompile(`\b[a-zA-Z']+\b`)

	// Find all words and their positions
	matches := wordRegex.FindAllStringIndex(text, -1)
	for _, match := range matches {
		word := text[match[0]:match[1]]
		if !sc.CheckWord(word) {
			misspelled = append(misspelled, MisspelledWord{
				Word:  word,
				Start: match[0],
				End:   match[1],
			})
		}
	}

	return misspelled
}

// GetSuggestions returns spelling suggestions for a misspelled word
func (sc *SpellChecker) GetSuggestions(word string) []string {
	if !sc.hasAspell {
		return nil
	}

	cmd := exec.Command("aspell", "-a")
	cmd.Stdin = strings.NewReader(word)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	// Parse aspell output for suggestions
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "&") {
			// Format: & word count offset: suggestion1, suggestion2, ...
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				suggestions := strings.Split(strings.TrimSpace(parts[1]), ", ")
				if len(suggestions) > 5 {
					suggestions = suggestions[:5] // Limit to 5 suggestions
				}
				return suggestions
			}
		}
	}
	return nil
}

// MisspelledWord represents a misspelled word with its position
type MisspelledWord struct {
	Word  string
	Start int
	End   int
}

// isCodeLikeWord returns true if the word looks like code
func isCodeLikeWord(word string) bool {
	// Skip if contains numbers
	for _, c := range word {
		if c >= '0' && c <= '9' {
			return true
		}
	}

	// Skip common code patterns
	codePatterns := []string{
		"http", "https", "www", "com", "org", "io",
		"func", "var", "const", "struct", "interface",
		"int", "string", "bool", "float", "nil", "null",
		"true", "false", "err", "ctx", "fmt", "json",
	}
	wordLower := strings.ToLower(word)
	for _, pattern := range codePatterns {
		if wordLower == pattern {
			return true
		}
	}

	// Skip camelCase or snake_case
	if strings.Contains(word, "_") {
		return true
	}

	// Skip if mixed case (camelCase)
	hasLower := false
	hasUpper := false
	for i, c := range word {
		if i > 0 && c >= 'A' && c <= 'Z' {
			hasUpper = true
		}
		if c >= 'a' && c <= 'z' {
			hasLower = true
		}
	}
	if hasLower && hasUpper {
		return true
	}

	return false
}
