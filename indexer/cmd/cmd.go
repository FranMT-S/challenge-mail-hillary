package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"indexer/database"
	"indexer/models"
	"indexer/scraper"
)

type Cmd struct {
	lastPage       int
	scrapper       *scraper.Scrapper
	isScraping     bool
	status         *models.SafeMap
	indexer        *database.Indexer
	mu             *sync.Mutex
	intervalUpdate int
	paginationSize scraper.PaginationWikileaks
	batchSize      int
}

const statusDirectory = "data"
const statusFilename = "data200pag.json"

func NewCmd(db *database.Connection, parallelism, delayRequest int) *Cmd {
	status := models.NewSafeMap()
	var mu sync.Mutex

	return &Cmd{
		scrapper:       scraper.NewScrapper(parallelism, delayRequest),
		isScraping:     false,
		status:         status,
		indexer:        database.NewIndexer(db),
		mu:             &mu,
		intervalUpdate: 50,
		paginationSize: scraper.PaginationWikileaks200,
		batchSize:      100,
	}
}

// SetBatchSize sets the batch size for the indexer
// The minimum value is 100
func (c *Cmd) SetBatchSize(size int) {
	if size < 100 {
		size = 100
	}
	c.batchSize = size
}

// SetPaginationSize sets the pagination size for the scraper
// The minimum value is 10 and the maximum is 200
func (c *Cmd) SetPaginationSize(size scraper.PaginationWikileaks) {
	if !size.IsValid() {
		size = scraper.PaginationWikileaks200
	}

	c.paginationSize = size
}

// SetIntervalUpdate sets the interval in rows to update the status
// The minimum value is 50
func (c *Cmd) SetIntervalUpdate(interval int) {
	if interval < 50 {
		interval = 50
	}

	c.intervalUpdate = interval
}

func (c *Cmd) Execute() {

	statusLoaded, err := c.LoadPageResults(statusDirectory, statusFilename)
	if err == nil {
		c.status.SetMap(statusLoaded)
	}

	c.lastPage = -1
	scanner := bufio.NewScanner(os.Stdin)

	c.lastPage, err = c.scrapper.GetLastPage(c.paginationSize)
	if err != nil {
		c.lastPage = -1
	}

	c.printHelp()

	for {
		fmt.Print("> write a command: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)

		switch args[0] {
		case "index":
			if c.isScraping {
				fmt.Println("Already indexing")
				continue
			}

			// create map to save the status
			go func() {
				c.Indexer(args)
			}()
		case "status":
			// load the status
			// update the status every time we receive a page result

			c.Status()
		case "help":
			c.printHelp()

		case "exit":
			fmt.Println("Bye!")
			return

		default:
			fmt.Println("Unknown command:", args[0])
			c.printHelp()
		}
	}
}

func (c *Cmd) printHelp() {
	indexMessage := "index --from=N --to=M   Start indexing from page N to M (default 1)"
	if c.lastPage > 0 {
		indexMessage += " (last page: " + strconv.Itoa(c.lastPage) + ")"
	}
	fmt.Println("Available commands:")
	fmt.Println(indexMessage)
	fmt.Println("  status                  Show current status")
	fmt.Println("  exit                    Exit the CLI")
	fmt.Println("  help                    Show this help message")
}

func (c *Cmd) SavePageResults(directory, filename string, data map[string]models.PageResult) error {
	// create directory if it doesn't exist
	if err := ensureDir(directory); err != nil {
		return fmt.Errorf("error creating directory %s: %w", directory, err)
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	c.mu.Lock()
	err = os.WriteFile(directory+"/"+filename, bytes, 0644)
	c.mu.Unlock()

	return err
}

func (c *Cmd) LoadPageResults(directory, filename string) (map[string]models.PageResult, error) {

	c.mu.Lock()
	bytes, err := os.ReadFile(directory + "/" + filename)
	c.mu.Unlock()
	if err != nil {
		return nil, err
	}

	result := make(map[string]models.PageResult)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Ensure the directory exists, creating it if necessary
func ensureDir(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	return err
}
