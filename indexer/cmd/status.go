package cmd

import (
	"fmt"

	"indexer/models"
)

func (c *Cmd) Status() {

	if c.status.Len() == 0 {
		loadedStatus, err := c.LoadPageResults(statusDirectory, statusFilename)
		if err == nil {
			c.status.SetMap(loadedStatus)
		}
	}

	if c.status.Len() == 0 {
		fmt.Println("No pages to show status")
		return
	}

	fmt.Println("Checking status...")
	pendingPages := 0
	totalFinished := 0
	totalProcessing := 0
	c.status.Range(func(key int, value models.PageResult) {
		fmt.Printf("Page: %d, State: %s, Total: %d, Error: %s\n", key, value.State, value.Total, value.Error)
		switch value.State {
		case models.PageResultStatePending:
			pendingPages++
		case models.PageResultStateFinished:
			totalFinished++
		case models.PageResultStateProcessing:
			totalProcessing++
		}
	})
	fmt.Printf("Total pages: %d, Pending pages: %d, Processing pages: %d, Finished pages: %d\n", c.status.Len(), pendingPages, totalProcessing, totalFinished)
}
