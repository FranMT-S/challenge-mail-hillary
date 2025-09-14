package database

import (
	"database/sql"
	"fmt"
	"regexp"
)

// ValidateTableName ensures the table name is safe to use in SQL queries
func ValidateTableName(name string) error {
	validName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid table name: %s", name)
	}
	return nil
}

func ValidateDBConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return nil
}
