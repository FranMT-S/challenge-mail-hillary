package database

import (
	"indexer/models"

	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	db IDatabase
}

func NewIndexer(db IDatabase) *Indexer {
	return &Indexer{db: db}
}

// IndexEmail indexes emails from a mailsCh channel
// mailsCh: channel of ScraperResult
// tableName: name of the table to index emails into
// batchSize: number of emails to index at once
func (i *Indexer) IndexEmail(mailsCh <-chan models.EmailResult, tableName string, batchSize int) error {
	batch := make([]models.Email, 0, batchSize)
	totalInserted := 0
	totalError := 0
	for result := range mailsCh {
		if result.Error != nil {
			totalError++
			continue
		}

		if result.Email.ID == 0 {
			continue
		}

		batch = append(batch, *result.Email)

		if len(batch) == batchSize {
			if inserted, err := i.db.SendMails(tableName, batch); err != nil {
				log.Error(err)
			} else {
				totalInserted += int(inserted)
				totalError += len(batch) - int(inserted)
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if inserted, err := i.db.SendMails(tableName, batch); err != nil {
			return err
		} else {
			totalInserted += int(inserted)
			totalError += len(batch) - int(inserted)
		}
	}

	log.Info("Batch inserted: ", totalInserted, " Batch total errors: ", totalError)
	return nil
}
