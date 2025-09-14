package database

import (
	"database/sql"

	"indexer/models"
)

type IDatabase interface {
	SendMails(tableName string, emails []models.Email) (int64, error)
	IsTableCreated(tableName string) (bool, error)
	CreateTableIfNotExists(tableName string) error
	DropTable(tableName string) error
	Open() (*sql.DB, error)
	OpenWithPool(maxOpenConns, maxIdleConns int) (*sql.DB, error)
	Ping() error
	Close() error
}
