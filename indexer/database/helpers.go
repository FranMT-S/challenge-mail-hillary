package database

import (
	"database/sql"
	"fmt"
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
