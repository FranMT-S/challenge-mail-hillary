package database

import (
	"database/sql"
	"fmt"
	"strings"

	"indexer/models"

	_ "github.com/lib/pq" // PostgreSQL driver compatible with CockroachDB
)

// Connection holds the database configuration
type Connection struct {
	Host     string
	Database string
	User     string
	Password string
	Port     string
	SSL      bool
	DB       *sql.DB
}

// NewConnection validates input and returns a Connection instance
func NewConnection(host, database, user, password, port string, ssl bool) (*Connection, error) {
	if host == "" {
		return nil, fmt.Errorf("host not specified")
	}

	if database == "" {
		return nil, fmt.Errorf("database not specified")
	}
	if user == "" {
		return nil, fmt.Errorf("user not specified")
	}

	if port == "" {
		return nil, fmt.Errorf("port not specified")
	}

	c := &Connection{
		Host:     host,
		Database: database,
		User:     user,
		Password: password,
		Port:     port,
		SSL:      ssl,
		DB:       nil,
	}

	db, err := c.Open()
	if err != nil {
		return nil, err
	}
	c.DB = db

	return c, nil
}

// NewConnectionWithPool validates input and returns a Connection instance
func NewConnectionWithPool(host, database, user, password, port string, ssl bool, maxOpenConns, maxIdleConns int) (*Connection, error) {
	if host == "" {
		return nil, fmt.Errorf("host not specified")
	}

	if database == "" {
		return nil, fmt.Errorf("database not specified")
	}
	if user == "" {
		return nil, fmt.Errorf("user not specified")
	}

	if port == "" {
		return nil, fmt.Errorf("port not specified")
	}

	c := &Connection{
		Host:     host,
		Database: database,
		User:     user,
		Password: password,
		Port:     port,
		SSL:      ssl,
		DB:       nil,
	}

	db, err := c.OpenWithPool(maxOpenConns, maxIdleConns)
	if err != nil {
		return nil, err
	}
	c.DB = db

	return c, nil
}

// Open opens the database connection without pool configuration
func (c *Connection) Open() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.BuildConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}
	c.DB = db
	return c.DB, nil
}

// OpenWithPool opens the database connection and configures connection pooling
func (c *Connection) OpenWithPool(maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := c.Open()
	if err != nil {
		return nil, err
	}

	c.DB = db
	c.DB.SetMaxOpenConns(maxOpenConns)
	c.DB.SetMaxIdleConns(maxIdleConns)
	return db, nil
}

// SendMails sends the emails to the database
// the mails is send in batches
func (c *Connection) SendMails(schemaName string, emails []models.Email) (int64, error) {
	if err := ValidateDBConnection(c.DB); err != nil {
		return 0, err
	}

	if len(emails) == 0 {
		return 0, nil
	}

	query, valueArgs, err := c.createInsertQuery(schemaName, emails)
	if err != nil {
		return 0, err
	}

	querySearchVector, valueArgsSearchVector, err := c.createSearchVectorQuery(schemaName, emails)
	if err != nil {
		return 0, err
	}

	tx, err := c.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	rows, err := tx.Exec(query, valueArgs...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch insert emails: %w", err)
	}

	_, err = tx.Exec(querySearchVector, valueArgsSearchVector...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch insert emails search vector: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	rowsInserted, err := rows.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsInserted, nil
}

// Ping checks if the database is reachable
func (c *Connection) Ping() error {
	return Ping(c.DB)
}

// Close safely closes the database connection
func (c *Connection) Close() error {
	if c.DB == nil {
		return nil
	}
	return c.DB.Close()
}

// BuildConnectionString builds the connection string for sql.Open
func (c *Connection) BuildConnectionString() string {
	if c.SSL {
		return fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=require",
			c.User, c.Password, c.Host, c.Port, c.Database,
		)
	}

	return fmt.Sprintf(
		"postgresql://%s@%s:%s/%s?sslmode=disable",
		c.User, c.Host, c.Port, c.Database,
	)
}

// CreateSchemaIfNotExist creates the schema if it does not exist
// validates tableName
func (c *Connection) CreateSchemaIfNotExist(schemaName string) error {

	if err := ValidateIsSafeString(schemaName); err != nil {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}

	_, err := c.DB.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schemaName))
	if err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}

	queries := []string{
		fmt.Sprintf(`
					CREATE TABLE IF NOT EXISTS "%s".emails (
							id INT PRIMARY KEY,
							date TIMESTAMP WITH TIME ZONE NOT NULL,
							subject TEXT DEFAULT '',
							"from" TEXT DEFAULT '',
							"to" TEXT DEFAULT '',
							content TEXT DEFAULT ''
					);
			`, schemaName),
		fmt.Sprintf(`
					CREATE TABLE IF NOT EXISTS "%s".emails_search (
							id INT PRIMARY KEY,
							search_vector TSVECTOR,
							FOREIGN KEY (id) REFERENCES "%s".emails(id)
					);
			`, schemaName, schemaName),
		fmt.Sprintf(`CREATE INVERTED INDEX IF NOT EXISTS idx_emails_search_vector
									ON "%s".emails_search (search_vector);`, schemaName),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_emails_date
									ON "%s".emails (date);`, schemaName),
	}

	for _, q := range queries {
		if _, err := c.DB.Exec(q); err != nil {
			return fmt.Errorf("error executing query: %w", err)
		}
	}

	return nil
}

func (c *Connection) IsSchemaCreated(schemaName string) (bool, error) {
	query := `SELECT EXISTS (
    SELECT 1
    FROM information_schema.schemata
    WHERE schema_name = $1
	) AS schema_exists;`

	var exists bool
	err := c.DB.QueryRow(query, schemaName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if schema exists: %w", err)
	}

	return exists, nil
}

// createInsertQuery creates the insert query and validate tableName is correct
func (c *Connection) createInsertQuery(schemaName string, emails []models.Email) (query string, valueArgs []any, err error) {
	valuesFlags := make([]string, 0, len(emails))
	valueArgs = make([]any, 0, len(emails)*6)

	if err := ValidateIsSafeString(schemaName); err != nil {
		return "", nil, err
	}

	// generate a bulkInsert query
	for i, e := range emails {
		// placeholders: ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12),...
		n := i*6 + 1
		id := n
		date := n + 1
		subject := n + 2
		from := n + 3
		to := n + 4
		content := n + 5

		mailValueFlag := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", id, date, subject, from, to, content)
		valuesFlags = append(valuesFlags, mailValueFlag)
		valueArgs = append(valueArgs, e.ID, e.Date, e.Subject, e.From, e.To, e.Content)
	}

	// build insert query with placeholders
	query = fmt.Sprintf(`
		INSERT INTO "%s".emails (id, date, subject, "from", "to", content)
		VALUES %s
		ON CONFLICT (id) DO NOTHING;
	`, schemaName, strings.Join(valuesFlags, ","))

	return query, valueArgs, nil
}

// createInsertQuery creates the insert query and validate tableName is correct
func (c *Connection) createSearchVectorQuery(schemaName string, emails []models.Email) (query string, valueArgs []any, err error) {
	valuesFlags := make([]string, 0, len(emails))
	valueArgs = make([]any, 0, len(emails)*5)

	if err := ValidateIsSafeString(schemaName); err != nil {
		return "", nil, err
	}

	// generate a bulkInsert query safe to insert search vector
	for i, e := range emails {
		n := i*5 + 1
		id := n
		subject := n + 1
		from := n + 2
		to := n + 3
		content := n + 4

		searchVectorValue := fmt.Sprintf("($%d, to_tsvector('english', $%d || ' ' || $%d || ' ' || $%d || ' ' || $%d))", id, subject, from, to, content)

		valuesFlags = append(valuesFlags, searchVectorValue)
		valueArgs = append(valueArgs, e.ID, e.Subject, e.From, e.To, e.Content)
	}

	// build insert query with placeholders
	query = fmt.Sprintf(`
		INSERT INTO "%s".emails_search (id, search_vector)
		VALUES %s
		ON CONFLICT (id) DO NOTHING;
	`, schemaName, strings.Join(valuesFlags, ","))

	return query, valueArgs, nil
}
