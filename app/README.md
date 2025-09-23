# Hillary Clinton Emails App

This is the frontend application that provides the search functionality for the Hillary Clinton Emails Challenge, the mails are stored in a cockroachdb database, the data must be loaded from the indexer project.

## Enviroment variables
Create a .env file based on .env.template, fill the values with your own.

``` env
VITE_API_URL=http://localhost:8080
```


## Folder structure
```
├── assets: Assets files
├── components: Components files
├── constants: Constants files
├── models: Models files
├── plugins: Plugins to config the application
├── services: Services files to handle the business logic
├── store: pinia store files
```

## Tech stack
- Language: typescript
- Framework: vue 3
- State management: pinia
- CSS framework: tailwindcss 3
- UI components: vuetify

## Installation

```bash
yarn install
```

## Run
```bash
yarn dev
```

## build
```bash
yarn build
```

