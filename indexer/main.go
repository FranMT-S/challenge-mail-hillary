package main

import (
	"indexer/cmd"
	"indexer/config"
	"indexer/database"
	"indexer/logger"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.InitConfigEnviroment()
	file, err := logger.InitLogger("logs", "app.log", config.GetConfig().LogLevel)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	defer file.Close()

	db, err := database.InitDb(*config.GetConfig())
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer db.Close()

	c := cmd.NewCmd(db, config.GetConfig().Scrapper.Parallelism, config.GetConfig().Scrapper.Delay)
	c.Execute()

}
