package main

import (
	"api/config"
	"api/database"
	"api/logger"
	"api/server"
)

func main() {

	config := config.GetConfig()
	if config == nil {
		panic("config not initialized")
	}

	logger.SetLogLevel(config.LogLevel)

	// Setup routes
	db, err := database.InitDB(database.DatabaseConfig{
		Host:       config.Host,
		Port:       config.Port,
		User:       config.User,
		Password:   config.Password,
		DBName:     config.DBName,
		SSLMode:    config.SSLMode,
		SchemaName: config.SchemaName,
	})

	if err != nil {
		panic(err)
	}

	defer database.CloseDB()

	server := server.NewServer(db, config.ApiPort)
	server.Start()

}
