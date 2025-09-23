package logger

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// initLogger initializes the logger
// dev: if true, sets the log level to debug, otherwise info
// returns the file and an error if any
func InitLogger(directory, filename string, level string) (*os.File, error) {
	// create logs directory if it doesn't exist
	os.Mkdir(directory, 0755)
	file, err := os.OpenFile(directory+"/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	setLogLevel(level)

	return file, nil
}

func setLogLevel(level string) {
	if level == "" {
		level = "info"
	}

	level = strings.ToLower(level)
	level = strings.TrimSpace(level)
	switch level {
	case "trace":
		fmt.Println("Setting log level to trace")
		log.SetLevel(log.TraceLevel)
	case "debug":
		fmt.Println("Setting log level to debug")
		log.SetLevel(log.DebugLevel)
	case "info":
		fmt.Println("Setting log level to info")
		log.SetLevel(log.InfoLevel)
	case "warn":
		fmt.Println("Setting log level to warn")
		log.SetLevel(log.WarnLevel)
	case "error":
		fmt.Println("Setting log level to error")
		log.SetLevel(log.ErrorLevel)
	case "panic":
		fmt.Println("Setting log level to panic")
		log.SetLevel(log.PanicLevel)
	case "fatal":
		fmt.Println("Setting log level to fatal")
		log.SetLevel(log.FatalLevel)
	default:
		fmt.Println("Setting log level to info")
		log.SetLevel(log.InfoLevel)
	}
}
