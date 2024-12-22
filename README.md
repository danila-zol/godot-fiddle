Repo of the Godot Fiddle project

### Backend deployment

> [!IMPORTANT]
> Requires `Go >= 1.22.0`

Install dependencies and generate Swagger docs:

```bash
$ go mod download
$ swag init
$ swag init -g ./server/handlers.go     # docs are mainly in handlers
```

Create a `.env` file with configuration appending the migrator config:

```
MIGRATE_DB=true
VERSION_COLUMN="schema_version"
EXPECTED_VERSION=1
```

> [!CAUTION]
> Be careful with migrations in production environment!

Build the project:

`$ go build`

or run it:

`$ go run .`
