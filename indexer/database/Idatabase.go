package database

import (
	"database/sql"

	"indexer/models"
)

// IDatabase represents the database interface
// SendMails: Sends emails to the database
// CreateSchemaIfNotExist: Creates the schema if it doesn't exist
// IsSchemaCreated: Checks if the schema exists
// Open: Opens the database connection
// OpenWithPool: Opens the database connection with a pool
// Ping: Pings the database connection
// Close: Closes the database connection
type IDatabase interface {
	SendMails(schemaName string, emails []models.Email) (int64, error)
	CreateSchemaIfNotExist(schemaName string) error
	IsSchemaCreated(schemaName string) (bool, error)
	Open() (*sql.DB, error)
	OpenWithPool(maxOpenConns, maxIdleConns int) (*sql.DB, error)
	Ping() error
	Close() error
}
