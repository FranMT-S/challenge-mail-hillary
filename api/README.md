## Hillary Clinton Emails API

This is the API that provides the search functionality for the Hillary Clinton Emails Challenge, the mails are stored in a cockroachdb database, the data must be loaded from the indexer project.


## Tech stack
- Language: golang 
- API framework: chi
- Database: cockroachdb
- ORM: gorm
- Logger: zerolog
- Environment variables: godotenv

## Folder structure

```
├── config: class files to config the application
├── controllers: Controller HTTP files
├── database: connection to the database
├── http: examples how use the API Endpoints
├── logger: Logger files  
├── logs: directory where the logs are stored
├── middleware: Middleware to handle the requests
├── models: Data models
├── routes: Api endpoints
├── sanatizer: Utility to sanitize the data
├── server: Main server file
├── services: Utility to handle the business logic
```

## Enviroment variables
Create a .env file based on .env.template, fill the values with your own.

``` env
DB_HOST=localhost # Database host
DB_NAME=defaultdb # Database name
DB_USER=root # Database user
DB_PASSWORD= # Database password
DB_PORT=26257 # Database port
DB_SSL=false # Database SSL mode
API_PORT=8080 # API port
CLIENT_HOST="http://localhost:5173" # Client host
LOG_LEVEL=DEBUG # Log level, Options: trace, debug, info, warn, error, dpanic, panic, fatal
LOG_DB=false # Log database, used to debug queries
```

## Installation

```bash
go mod download
go mod init
go mod tidy
```

## Run
Configure the environment variables and run the application.

```bash
go run main.go
```


## Endpoints

### POST /mails/search
Search for emails in the database.

``` http
POST /mails/search
{
  "query": "hillary", // Search query
  "type": "AND", // AND, OR
  "page": 1, // Page number
  "limit": 20, // Number of results per page
  "dateSearch": {
    "date": "2019-07-01T06:00:00.000Z", // Date to search
    "operator": "<=" //<=, >=, =, <, >
  },
  "orderBy": "desc" // asc, desc
}
```