# Liquo

Liquo is an experimental Liquibase plugin for oxpecker. It's aimed to make the development of Buffalo apps with Liquibase migrations a bit more easy but avoiding the installation of Java tools for the development workflow.
## Installation

Liquo is provided as a library you can pull with

```
go get github.com/wawandco/liquo
```

An then in your `cmd/ox/main.go` use it by appending it to the CLI plugins.

```
...
cli.Remove("pop/generate-migration") // make sure pop plugin is not being used
cli.Use(liquo.Plugins()...)
...
```

## Limited Functionality

Liquo still experimental, it does not provide the same amount of statements, formats or databases that liquibase supports. 

Liquo ONLY supports:

- PostgresSQL Database
- Liquibase XML format, only the following statements:
    - sql
    - rollback

While is possible to add the rest of statements this is where the tool is at the moment.
## Usage
Generate migration file in `./migrations` default directory:
- `ox generate migration create-users-table`

Generate migration file in `./my/directory` directory:
- `ox generate migration ./my/directory/create-users-table`

Run migrations:
- `ox db migrate`

Run one single migration:
- `ox db migrate up`

Rollback one single migration:
- `ox db migrate down`

Usage notes:
1. Generating a migration file auto-adds the import path in the `changelog.xml` file.
2. If no `--conn` flag is provided, liquo assumes `development` as its standard DB connection.

## License

Liquo is released under the [MIT License](LICENSE).