package scraper

import "github.com/microcosm-cc/bluemonday"

func SanitizeHTML(input string) string {
	p := bluemonday.UGCPolicy() // UGC policy is safe for user generated content
	return p.Sanitize(input)
}
