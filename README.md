## Hillary Clinton Emails Challenge

This is a challenge to index the emails of Hillary Clinton from the WikiLeaks this project use colly to index the mails to the cockroachdb database and then can search the mails using a rest api.

## Requirements

* [CockroachDB](https://www.cockroachlabs.com/docs/v25.3/install-cockroachdb-windows.html)
* [Go](https://golang.org/dl/)
* [Vue 3](https://vuejs.org/guide/introduction.html)
* [Node.js](https://nodejs.org/en/download/)
* [NPM](https://www.npmjs.com/)
* [Yarn](https://yarnpkg.com/)

## Database

The database is a [cockroachdb](https://www.cockroachlabs.com/docs/v25.3/install-cockroachdb-windows.html) can view the documentation to learn how install also can use docker, the schema is in the [script.sql](./indexer/database/script.sql) file.

## Tech stack
- API: golang
- Indexer: golang
- Database: cockroachdb
- Client: vue 3 + pinia + vite

## API
* Language: golang
* Framework: chi
* Database: cockroachdb
* ORM: gorm
* Logger: zerolog
* Environment variables: godotenv

# Client
* Language: vue 3
* Framework: vite
* State management: pinia
* Router: vue-router
* Environment variables: dotenv

# Indexer
* Language: golang
* Framework: colly
* Database: cockroachdb
* Environment variables: godotenv

## Folder structure
```
├── api: API files
├── indexer: Indexer files
├── client: Client files
```

## Documentation
* [API](./api/README.md)
* [Indexer](./indexer/README.md)
* [Client](./client/README.md)