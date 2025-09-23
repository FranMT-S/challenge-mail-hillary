package config

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// Config represents the configuration for the indexer
// Env: Environment
// DBConfig: Database configuration
// Scrapper: Scraper configuration
// LogLevel: Log level
type Config struct {
	Env      string
	DBConfig struct {
		Host     string
		Name     string
		User     string
		Password string
		Port     string
		SSL      bool
	}
	Scrapper struct {
		Parallelism int
		Delay       int
	}
	LogLevel string
}

var config *Config

func InitConfigEnviroment() {
	config = &Config{Env: os.Getenv("ENVIRONMENT")}
	config.DBConfig.Host = os.Getenv("DB_HOST")
	config.DBConfig.Name = os.Getenv("DB_NAME")
	config.DBConfig.User = os.Getenv("DB_USER")
	config.DBConfig.Password = os.Getenv("DB_PASSWORD")
	config.DBConfig.Port = os.Getenv("DB_PORT")
	config.DBConfig.SSL = os.Getenv("DB_SSL") == "true"
	config.LogLevel = os.Getenv("LOG_LEVEL")

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	if strings.ToLower(config.Env) == "prod" {
		config.LogLevel = "warn"
	}

	var err error
	config.Scrapper.Parallelism, err = strconv.Atoi(os.Getenv("SCRAPPER_PARALLELISM"))
	if err != nil {
		config.Scrapper.Parallelism = 10
	}
	config.Scrapper.Delay, err = strconv.Atoi(os.Getenv("SCRAPPER_DELAY"))
	if err != nil {
		config.Scrapper.Delay = 1
	}
}

func GetConfig() *Config {
	if config == nil {
		panic("config not initialized")
	}

	if config.Env == "" {
		config.Env = "dev"
	}

	if config.DBConfig.Host == "" {
		panic("db_host not specified")
	}

	if config.DBConfig.Name == "" {
		panic("db_name not specified")
	}

	if config.DBConfig.User == "" {
		panic("db_user not specified")
	}

	if config.DBConfig.Port == "" {
		panic("db_port not specified")
	}

	if config.DBConfig.SSL && config.DBConfig.Password == "" {
		panic("db_password not specified when db_ssl is true")
	}

	return config
}
