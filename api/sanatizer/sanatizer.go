package sanatizer

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"
)

// SanitizeForHTML escapes special characters to prevent XSS when rendering in HTML
func SanitizeForHTML(s string) string {
	// Escapes <, >, &, ' and "
	return html.EscapeString(s)
}

// SanitizeForSQL performs a basic cleanup
// This function:
// - removes SQL comments (-- ... and /* ... */)
// - removes null characters
// - removes semicolons
// - doubles single quotes
func SanitizeForSQL(s string, maxLen int) string {

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\x00", "")

	reLineComment := regexp.MustCompile(`(?i)--[^\n\r]*`)
	s = reLineComment.ReplaceAllString(s, "")

	reBlockComment := regexp.MustCompile(`(?is)/\*.*?\*/`)
	s = reBlockComment.ReplaceAllString(s, "")

	s = strings.ReplaceAll(s, ";", "")

	s = strings.ReplaceAll(s, "'", "''")

	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	if maxLen > 0 && utf8.RuneCountInString(s) > maxLen {
		runes := []rune(s)
		s = string(runes[:maxLen])
	}

	return s
}
