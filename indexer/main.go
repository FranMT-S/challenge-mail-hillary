package main

import (
	"indexer/cmd"
	"indexer/config"
	"indexer/database"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.InitConfigEnviroment()
	file, err := initLogger(config.GetConfig().Env == "dev")
	if err != nil {
		log.Error(err)
		panic(err)
	}

	defer file.Close()

	db, err := initDb(*config.GetConfig())
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer db.Close()

	c := cmd.NewCmd(db, config.GetConfig().Scrapper.Parallelism, config.GetConfig().Scrapper.Delay)
	c.Execute()

}

// initLogger initializes the logger
// dev: if true, sets the log level to debug, otherwise info
// returns the file and an error if any
func initLogger(dev bool) (*os.File, error) {
	// create logs directory if it doesn't exist
	os.Mkdir("logs", 0755)
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetLevel(log.InfoLevel)
	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// is development mode
	if dev {
		log.SetLevel(log.DebugLevel)
	}

	return file, nil
}

func initDb(cfg config.Config) (*database.Connection, error) {
	db, err := database.NewConnection(cfg.DBConfig.Host, cfg.DBConfig.Name, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Port, cfg.DBConfig.SSL)
	if err != nil {
		return nil, err
	}

	_, err = db.OpenWithPool(10, 5)
	if err != nil {
		return nil, err
	}

	err = db.CreateTableIfNotExists(database.DBTableName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
