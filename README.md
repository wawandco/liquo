# Liquo

Liquo is an experimental Liquibase plugin for oxpecker. It's aimed to make the development of Buffalo apps with Liquibase migrations a bit more easy but avoiding the installation of Java tools for the development workflow.
## Installation

Liquo is provided as a library you can pull it with

```
go get github.com/wawandco/liquo
```

An then in your cmd/ox/main.go use it by appending it to the CLi plugins.

```
...
cl.Plugins = append(cl.Plugins, liquo.Plugins()...)
...
```

## Limited Functionallity

Liquo still experimental, it does not provide the same amount of statements, formats or databases that liquibase supports. 

Liquo ONLY supports:

- PostgresSQL Database
- Liquibase XML format, only the following statements:
    - sql
    - rollback

While is possible to add the rest of statements this is where the tool is at the moment.
## License

Liquo is released under the [MIT License](LICENSE).