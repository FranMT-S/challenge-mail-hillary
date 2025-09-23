package scraper

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"indexer/models"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

// Scrapper represents the scraper
// parallelism: number of parallel requests
// delay: delay between requests in seconds
type Scrapper struct {
	parallelism int
	delay       int
}

// NewScrapper creates a new scrapper
// parallelism: number of parallel requests
// delay: delay between requests in seconds
func NewScrapper(parallelism int, delay int) *Scrapper {
	if parallelism < 1 {
		parallelism = 1
	}

	if delay < 1 {
		delay = 1
	}

	return &Scrapper{
		parallelism: parallelism,
		delay:       delay,
	}
}

// PaginationWikileaks represents the pagination parameter for the WikiLeaks API
type PaginationWikileaks int

const (
	PaginationWikileaks10  PaginationWikileaks = 10
	PaginationWikileaks25  PaginationWikileaks = 25
	PaginationWikileaks50  PaginationWikileaks = 50
	PaginationWikileaks100 PaginationWikileaks = 100
	PaginationWikileaks200 PaginationWikileaks = 200
)

func (p PaginationWikileaks) IsValid() bool {
	return p == PaginationWikileaks10 ||
		p == PaginationWikileaks25 ||
		p == PaginationWikileaks50 ||
		p == PaginationWikileaks100 ||
		p == PaginationWikileaks200
}

// ScrapeEmails scrapes emails from the WikiLeaks API
// fromPage: start page
// toPage: end page
// pagination: pagination parameter
// emailsQueue: channel to send the emails
// pagesQueueUpdater: channel to update the pages state
// intervalUpdate: interval to update the pages state
func (s *Scrapper) ScrapeEmails(fromPage int, toPage int, pagination PaginationWikileaks, emailsQueue chan<- models.EmailResult, pagesQueueUpdater chan<- models.PageResult, intervalUpdate int) error {

	if emailsQueue == nil {
		return errors.New("emailsQueue is nil")
	}

	err := s.validateParams(fromPage, toPage, int(pagination))
	if err != nil {
		return err
	}

	if !pagination.IsValid() {
		return fmt.Errorf("the pagination send is not supported by the api: %d", pagination)
	}

	semaphore := models.NewSemaphore(50)
	c := SetupWikileaksCollector(s.parallelism, s.delay, true)

	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting:", r.URL.String())
	})

	c.OnHTML(".table.search-result tbody", func(e *colly.HTMLElement) {
		var total int64 = 0
		var wg sync.WaitGroup

		// process rows
		e.ForEach("tr", func(i int, row *colly.HTMLElement) {
			semaphore.Acquire()
			wg.Add(1)
			go func() {
				defer semaphore.Release()
				defer wg.Done()
				current := atomic.AddInt64(&total, 1)

				// update page state every intervalUpdate rows
				if pagesQueueUpdater != nil && int(current)%intervalUpdate == 0 {
					pageStr := e.Request.Ctx.Get("page")
					page, err := strconv.Atoi(pageStr)

					if err == nil {
						s.updatePageState(pagesQueueUpdater, page, int(current), nil, models.PageResultStateProcessing)
					}
				}

				err := error(nil)
				email := models.Email{}

				// log every 10 rows just in debug
				pageStr := e.Request.Ctx.Get("page")

				row.ForEach("td", func(tdIndex int, columns *colly.HTMLElement) {
					// if there is an error, stop processing the row
					if err != nil {
						return
					}
					err = processRow(tdIndex, columns, &email)
					if err != nil {
						log.WithFields(log.Fields{"error": err, "id": email.ID}).Error("Error processing email")
					}

					// scrapt email content if it is the last column
					if tdIndex == 4 && email.ID > 0 {
						content, err := getMailContent(email.ID, c)
						if err != nil {
							log.WithFields(log.Fields{"error": err, "id": email.ID}).Error("Error getting email content")
						}
						email.Content = content
					}
				})

				// log every 10 rows just in debug
				if current%10 == 0 {
					log.Trace("Page ", pageStr, " processed ", current, " rows")
				}

				if err == nil {
					e.Request.Ctx.Put("error", "")
				} else {
					e.Request.Ctx.Put("error", err.Error())
				}

				emailsQueue <- models.EmailResult{Email: &email, Error: err}

			}()
		})

		wg.Wait()
		e.Request.Ctx.Put("total", strconv.Itoa(int(total)))
	})

	// update page state when the page is finished
	c.OnScraped(func(r *colly.Response) {
		var scrapeErrStr string = ""
		pageStr := r.Request.Ctx.Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Errorf("failed to convert page to int: page %s, error %v", pageStr, err)
			return
		}

		totalStr := r.Request.Ctx.Get("total")
		total, err := strconv.Atoi(totalStr)
		if err != nil {
			log.Errorf("failed to convert total to int: total %s, error %v", totalStr, err)
			return
		}

		if errStr := r.Request.Ctx.Get("error"); errStr != "" {
			scrapeErrStr = errStr
		}

		if pagesQueueUpdater != nil {
			pagesQueueUpdater <- models.PageResult{Page: page, Error: scrapeErrStr, Total: total, State: models.PageResultStateFinished}
		}
		log.Debug("Scraped page finish:", r.Request.URL.String())
	})

	log.Info("Setup Scraping pages from", fromPage, "to", toPage)
	for page := fromPage; page <= toPage; page++ {
		url := s.getPageUrlWithPagination(page, pagination)
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(page))

		err := c.Request("GET", url, nil, ctx, nil)
		if err != nil {
			log.Error("Error creating request:", err)
		}
	}
	log.Info("Setup Finish Scraping pages from ", fromPage, " to ", toPage)

	c.Wait()

	if pagesQueueUpdater != nil {
		close(pagesQueueUpdater)
	}

	if emailsQueue != nil {
		close(emailsQueue)
	}

	return nil
}

// GetLastPage gets the last page of the emails
// pagination: the pagination to use
// returns the last page and an error if any
func (s *Scrapper) GetLastPage(pagination PaginationWikileaks) (int, error) {
	if !pagination.IsValid() {
		return 0, fmt.Errorf("the pagination send is not supported by the api: %d", pagination)
	}

	c := SetupWikileaksCollector(1, 1, false)

	lastPage := 1

	c.OnHTML("ul.pagination.pagination", func(ul *colly.HTMLElement) {
		ul.ForEach("li", func(i int, li *colly.HTMLElement) {
			if i == ul.DOM.Find("li").Length()-2 { // penúltimo
				var err error
				lastPage, err = strconv.Atoi(li.Text)
				if err != nil {
					log.Error("failed to convert last page to int:", err)
					return
				}
			}
		})
	})

	err := c.Visit(s.getPageUrlWithPagination(1, pagination))
	if err != nil {
		return 0, err
	}

	return lastPage, nil
}

// getPageUrlWithPagination gets the url of the page with the pagination
func (s *Scrapper) getPageUrlWithPagination(page int, pagination PaginationWikileaks) string {
	return fmt.Sprintf("https://wikileaks.org/clinton-emails/?q=&mfrom=&mto=&title=&notitle=&date_from=&date_to=&nofrom=&noto=&sort=0&count=%d&page=%d#searchresult", pagination, page)
}

// processRow processes a row of the table
// tdIndex: index of the column
// column: the column element
// email: the email to fill
func processRow(tdIndex int, column *colly.HTMLElement, email *models.Email) error {
	switch tdIndex {
	case 0:
		id, err := strconv.Atoi(column.Text)
		if err != nil {
			log.Errorf("failed to convert '%s' to int: %v", column.Text, err)
			return err
		}

		email.ID = uint32(id)
	case 1:
		str := column.Text
		layout := "2006-01-02 15:04"

		t, err := time.Parse(layout, str)
		if err != nil {
			log.Errorf("failed to parse date '%s': %v", str, err)
			return err
		}

		email.Date = t.UTC()
	case 2:
		email.Subject = SanitizeHTML(column.Text)
	case 3:
		email.From = SanitizeHTML(column.Text)
	case 4:
		email.To = SanitizeHTML(column.Text)
	default:
		log.Warnf("Unknown column index: %d", tdIndex)
	}

	return nil
}

// Parse the mail content from the page
// id: the id of the email
// c: the collector to use
// returns the content of the email and an error if any
func getMailContent(id uint32, c *colly.Collector) (string, error) {
	collectorContent := c.Clone()
	collectorContent.Limit(&colly.LimitRule{
		DomainGlob: "wikileaks.org",
		Delay:      1 * time.Second,
	})

	collectorContent.Async = true

	var err error
	var html string
	collectorContent.OnHTML("div#content", func(e *colly.HTMLElement) {
		html, err = e.DOM.Html()
		if err != nil {
			log.Error("Error getting HTML:", err)
			return
		}
	})

	err = collectorContent.Visit("https://wikileaks.org/clinton-emails/emailid/" + fmt.Sprintf("%d", id))
	if err != nil {
		log.Error("Error visiting email content:", err)
		return "", err
	}

	collectorContent.Wait()

	return SanitizeHTML(html), nil
}

// validateParams validates the parameters for scraping
// fromPage: page to start scraping from
// toPage: page to end scraping at
// pagination: number of emails per page
// returns an error if any
func (s *Scrapper) validateParams(fromPage int, toPage int, pagination int) error {

	if fromPage < 1 {
		return fmt.Errorf("fromPage must be greater than 0")
	}

	if toPage < 1 {
		return fmt.Errorf("toPage must be greater than 0")
	}

	if pagination < 1 {
		return fmt.Errorf("paginationTotal must be greater than 0")
	}

	if fromPage > toPage {
		return fmt.Errorf("fromPage is greater than toPage")
	}

	return nil
}

func (s *Scrapper) updatePageState(pagesQueueUpdater chan<- models.PageResult, page int, total int, err error, state models.PageResultStateType) {
	// update page state every intervalUpdate rows

	strErr := ""
	if err != nil {
		strErr = err.Error()
	}

	if pagesQueueUpdater != nil {
		pagesQueueUpdater <- models.PageResult{Page: page, Error: strErr, Total: total, State: state}
	}
}
