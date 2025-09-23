package database

import (
	"api/config"
	"fmt"
	"log"
	"os"
	"time"

	apiLogger "api/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// DatabaseConfig stores the configuration for the database
type DatabaseConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	SSLMode    bool
	SchemaName string
}

// isValid check if the Host, port, user and db name are not empty
func (dc *DatabaseConfig) isValid() bool {
	return dc.Host != "" && dc.Port != 0 && dc.User != "" && dc.DBName != ""
}

// InitDB initializes the database connection
func InitDB(cfg DatabaseConfig) (*gorm.DB, error) {

	if !cfg.isValid() {
		return nil, fmt.Errorf("invalid database configuration")
	}

	sslMode := "disable"
	if cfg.SSLMode {
		sslMode = "require"
	}

	var newLogger logger.Interface = nil
	fmt.Println("LogDB: ", config.GetConfig().LogDB)
	if config.GetConfig().LogDB {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             2 * time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info,     // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,            // Don't include params in the SQL log
				Colorful:                  true,            // Disable color
			},
		)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, sslMode, cfg.SchemaName)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: cfg.SchemaName + ".",
		},
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	apiLogger.Logger().Info().Msg("Successfully connected to the database")
	return DB, nil
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		apiLogger.Logger().Error().Msg("Failed to close database connection: " + err.Error())
		return
	}

	sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
