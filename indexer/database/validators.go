package database

import (
	"database/sql"
	"fmt"
	"regexp"
)

// ValidateIsSafeString ensures the string is safe to use in SQL queries
func ValidateIsSafeString(name string) error {
	validName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid string: %s", name)
	}
	return nil
}

// ValidateDBConnection ensures the database connection is valid
func ValidateDBConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return nil
}
