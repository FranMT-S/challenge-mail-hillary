package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config stores the configuration for the API
type Config struct {
	Host            string // Database host
	Port            int    // Database port
	User            string // Database user
	Password        string // Database password
	DBName          string // Database name
	SSLMode         bool   // Database SSL mode
	ClientHost      string // Client host
	SchemaName      string // Schema name
	MailsTable      string // Mails table name
	MailSearchTable string // Mail search table name
	LogLevel        string // Log level
	LogDB           bool   // Log database
	ApiPort         int    // API port
}

var config *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file
	_ = godotenv.Load()

	// Parse port from environment
	port, err := strconv.Atoi(getEnv("DB_PORT", "26257"))
	if err != nil {
		port = 26257 // Default port if conversion fails
	}
	sslstr := strings.ToLower(getEnv("DB_SSL", "false"))

	ssl := sslstr == "true"
	config = &Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            port,
		User:            getEnv("DB_USER", "root"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "defaultdb"),
		SSLMode:         ssl,
		ClientHost:      getEnv("CLIENT_HOST", "localhost:5173"),
		SchemaName:      "emails_hillary",
		MailsTable:      "emails",
		MailSearchTable: "emails_search",
		LogLevel:        strings.ToLower(getEnv("LOG_LEVEL", "info")),
		LogDB:           strings.ToLower(getEnv("LOG_DB", "false")) == "true",
	}

	// config apiPort
	apiPortStr := getEnv("API_PORT", "8080")
	apiPort, err := strconv.Atoi(apiPortStr)
	if err != nil {
		apiPort = 8080 // Default port if conversion fails
	}

	config.ApiPort = apiPort

	if config.MailsTable == "" {
		panic("MAILS_TABLE not specified")
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetConfig() *Config {

	if config == nil {
		_, err := LoadConfig()
		if err != nil {
			panic(err)
		}
	}

	return config
}
