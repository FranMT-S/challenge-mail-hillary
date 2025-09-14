package scraper

import (
	"log"
	"net/http/cookiejar"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// SetupWikileaksCollector sets up a collector for WikiLeaks
// parallelism: number of parallel requests
// delayInSeconds: delay between requests
// returns the collector
func SetupWikileaksCollector(parallelism int, delayInSeconds int, async bool) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("wikileaks.org"),
		colly.AllowURLRevisit(),
		colly.Async(async),
	)

	jar, err := cookiejar.New(nil)
	if err == nil {
		c.SetCookieJar(jar)
	}

	c.Limit(&colly.LimitRule{
		DomainGlob:  "wikileaks.org",
		Parallelism: parallelism,
		Delay:       time.Duration(delayInSeconds) * time.Second,
	})

	// avoid blocking in a request
	c.SetRequestTimeout(30 * time.Second)

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("error in %s: %v", r.Request.URL, err)
	})

	extensions.RandomUserAgent(c) // changes the UA randomly

	return c
}
