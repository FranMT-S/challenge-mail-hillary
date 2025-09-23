## Indexer 

this is a indexer to index the emails in the app, the mails is stored in a cockroachdb database, this project scrapes the emails from the [Emails Hillary Clinton](https://wikileaks.org/clinton-emails/).

## Tech stack
- Language: golang
- scraper: colly
- Database: cockroachdb

### Folder Structure
```
├── cmd: CLI application
├── config: Configuration files
├── data: directory when the scraper status is saved
├── database: Database connection, contains the script.sql file to create the database
├── logger: Logger files config
├── logs: directory where the logs are stored
├── models: Data models
├── scraper: scraper functions
```

## Enviroment variables
Create a .env file based on .env.template, fill the values with your own.

``` env
ENVIRONMENT=dev # Options: dev, prod
DB_HOST=localhost # CockroachDB host
DB_NAME=defaultdb # CockroachDB database name
DB_USER=root # CockroachDB user
DB_PASSWORD= # CockroachDB password
DB_PORT=26257 # CockroachDB port
DB_SSL=false # CockroachDB SSL
SCRAPPER_PARALLELISM=10 # Number of parallel scrapers
SCRAPPER_DELAY=1 # Delay between request of scrapers
LOG_LEVEL=trace # Log level Options: trace, debug, info, warn, error, dpanic, panic, fatal
```

## Installation

```go
go mod download
go mod init
go mod tidy

```

## Run
Configure the environment variables and run the application.

```go
go run main.go
```

## CLI Commands
```
index --from=N --to=M   Start indexing from page N to M (default 1)
status                  Show current status of the indexer, show the status of the last page indexed
exit                    Exit the CLI
help                    Show  command help message
```