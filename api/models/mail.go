package models

import (
	"api/config"
	"errors"
	"time"
)

// Email represents an email
type Email struct {
	ID      uint32    `json:"id"`
	Date    time.Time `json:"date"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Content string    `json:"content"`
}

// EmailRank represents an email with rank to inverted index search
type EmailRank struct {
	ID      uint32    `json:"id"`
	Date    time.Time `json:"date"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Content string    `json:"content"`
	Rank    float64   `json:"rank"`
}

func (e *Email) Validate() error {
	if e.ID == 0 {
		return errors.New("id is required")
	}
	if e.Date.IsZero() {
		return errors.New("date is required")
	}

	if e.From == "" {
		return errors.New("from is required")
	}

	if e.To == "" {
		return errors.New("to is required")
	}

	return nil
}

// TableName returns the table name for the model
// used in gorm to get data
func (Email) TableName() string {
	return config.GetConfig().SchemaName + "." + config.GetConfig().MailsTable
}
