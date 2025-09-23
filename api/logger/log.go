package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// init initializes the logger
func init() {
	if err := os.MkdirAll("logs/api", os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.OpenFile("logs/api/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger = zerolog.New(file).With().Timestamp().Logger()
}

// Logger returns the logger singleton instance
func Logger() *zerolog.Logger {
	return &logger
}

// SetLogLevel sets the log level
func SetLogLevel(level string) {
	if level == "" {
		level = "info"
	}

	switch strings.ToLower(level) {
	case "debug":
		logger.Level(zerolog.DebugLevel)
	case "info":
		logger.Level(zerolog.InfoLevel)
	case "warn":
		logger.Level(zerolog.WarnLevel)
	case "error":
		logger.Level(zerolog.ErrorLevel)
	case "none":
		logger.Level(zerolog.NoLevel)
	default:
		logger.Level(zerolog.InfoLevel)
	}
}
