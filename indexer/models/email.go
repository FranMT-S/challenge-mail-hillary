package models

import "time"

// Email represents an email
type Email struct {
	ID      uint32    `json:"id"`
	Date    time.Time `json:"date"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Content string    `json:"content"`
}
