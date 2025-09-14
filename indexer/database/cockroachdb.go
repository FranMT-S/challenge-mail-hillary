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
func (c *Connection) SendMails(tableName string, emails []models.Email) (int64, error) {
	if err := ValidateDBConnection(c.DB); err != nil {
		return 0, err
	}

	if len(emails) == 0 {
		return 0, nil
	}

	query, valueArgs, err := c.createInsertQuery(tableName, emails)
	if err != nil {
		return 0, err
	}

	res, err := c.DB.Exec(query, valueArgs...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch insert emails: %w", err)
	}

	// get affected rows
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return rows, nil
}

// IsTableCreated checks if the table exists in the database
func (c *Connection) IsTableCreated(tableName string) (bool, error) {
	if err := ValidateDBConnection(c.DB); err != nil {
		return false, err
	}

	if err := ValidateTableName(tableName); err != nil {
		return false, err
	}

	var exists bool
	row := c.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", tableName)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if table exists: %w", err)
	}
	return exists, nil
}

// DropTable drops the table if it exists
func (c *Connection) DropTable(tableName string) error {
	if err := ValidateDBConnection(c.DB); err != nil {
		return err
	}

	query, err := c.createDropTableQuery(tableName)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec(query)
	return err
}

// CreateTableIfNotExists creates the table if it does not exist
func (c *Connection) CreateTableIfNotExists(tableName string) error {
	if err := ValidateDBConnection(c.DB); err != nil {
		return err
	}

	query, err := c.createTableIfNotExistQuery(tableName)
	if err != nil {
		return err
	}

	_, err = c.DB.Exec(query)
	return err

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

// createTableIfNotExistQuery creates the create table query
// validates tableName
func (c *Connection) createTableIfNotExistQuery(tableName string) (string, error) {
	if err := ValidateTableName(tableName); err != nil {
		return "", err
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INT PRIMARY KEY,
		date TIMESTAMP WITH TIME ZONE NOT NULL,
		subject TEXT,
		"from" TEXT NOT NULL,
		"to" TEXT NOT NULL,
		content TEXT
	);`, tableName)

	return query, nil
}

// createDropTableQuery creates the drop table query and validate tableName is correct
func (c *Connection) createDropTableQuery(tableName string) (string, error) {
	if err := ValidateTableName(tableName); err != nil {
		return "", err
	}
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName), nil
}

// createInsertQuery creates the insert query and validate tableName is correct
func (c *Connection) createInsertQuery(tableName string, emails []models.Email) (query string, valueArgs []any, err error) {
	if err := ValidateTableName(tableName); err != nil {
		return "", nil, err
	}

	valueStrings := make([]string, 0, len(emails))
	valueArgs = make([]any, 0, len(emails)*6)
	for i, e := range emails {
		// placeholders: ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12), ...
		n := i*6 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", n, n+1, n+2, n+3, n+4, n+5))
		valueArgs = append(valueArgs, e.ID, e.Date, e.Subject, e.From, e.To, e.Content)
	}

	// build insert query with placeholders
	query = fmt.Sprintf(`
		INSERT INTO %s (id, date, subject, "from", "to", content)
		VALUES %s
		ON CONFLICT (id) DO NOTHING;
	`, tableName, strings.Join(valueStrings, ","))

	return query, valueArgs, nil
}
