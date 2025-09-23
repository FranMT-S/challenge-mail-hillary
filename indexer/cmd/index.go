package cmd

import (
	"flag"
	"fmt"
	"strconv"

	"indexer/models"

	log "github.com/sirupsen/logrus"
)

// Indexer indexes emails from a range of pages
// The pages are specified with the --from and --to flags
// If --to is greater than the last page, it will use the last page
func (c *Cmd) Indexer(args []string) {
	var from, to int
	fs := flag.NewFlagSet("index", flag.ContinueOnError)
	fs.IntVar(&from, "from", 1, "page number to start indexing from")
	fs.IntVar(&to, "to", 1, "page number to end indexing at")

	// Parse the flags from the input
	err := fs.Parse(args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	if c.lastPage > 0 && to > c.lastPage {
		to = c.lastPage
	}

	for i := from; i <= to; i++ {
		key := strconv.Itoa(i)
		c.status.Set(key, models.PageResult{Page: i, State: models.PageResultStatePending, Total: 0, Error: ""})
	}
	err = c.SavePageResults(statusDirectory, statusFilename, c.status.GetCopy())
	if err != nil {
		log.Error("Error saving initial status:", err)
	}

	pageResultCh := make(chan models.PageResult)
	emailsCh := make(chan models.EmailResult)

	go func() {
		c.isScraping = true
		err = c.scrapper.ScrapeEmails(from, to, c.paginationSize, emailsCh, pageResultCh, c.intervalUpdate)
		if err != nil {
			fmt.Println("Error indexing:", err)
			log.Error("Error indexing:", err)
		}

		fmt.Println("Scraping finished")

		c.isScraping = false
	}()

	// Index emails in batches
	go func() {
		c.indexer.IndexEmail(emailsCh, c.batchSize)
	}()

	// Update status in real time
	go func() {
		for result := range pageResultCh {
			key := strconv.Itoa(result.Page)
			log.WithFields(log.Fields{"page": result.Page, "total": result.Total, "state": result.State}).Info("Update data page")
			c.status.Set(key, result)
			c.SavePageResults(statusDirectory, statusFilename, c.status.GetCopy())
		}
	}()

	fmt.Printf("Indexing from page %d to %d\n", from, to)
}
