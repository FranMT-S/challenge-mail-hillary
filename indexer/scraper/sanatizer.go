package scraper

import "github.com/microcosm-cc/bluemonday"

// SanitizeHTML sanitizes HTML content
// input: the HTML content to sanitize
// returns the sanitized HTML content
func SanitizeHTML(input string) string {
	p := bluemonday.UGCPolicy() // UGC policy is safe for user generated content
	return p.Sanitize(input)
}
