package database

import (
	"database/sql"
	"fmt"
	"indexer/config"
)

// General Ping checks if the database is reachable
func Ping(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

// General Close safely closes the database connection
func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}

// InitDb initializes the database connection
func InitDb(cfg config.Config) (*Connection, error) {
	db, err := NewConnection(cfg.DBConfig.Host, cfg.DBConfig.Name, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Port, cfg.DBConfig.SSL)
	if err != nil {
		return nil, err
	}

	_, err = db.OpenWithPool(10, 5)
	if err != nil {
		return nil, err
	}

	err = db.CreateSchemaIfNotExist(DBSchemaName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
