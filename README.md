Repo of the Godot Fiddle project

### Backend deployment

> [!IMPORTANT]
> Requires `Go >= 1.24.2`

Install dependencies and generate Swagger docs:

```bash
$ go mod download
$ swag init --dir ./cmd/api/,./internal/delivery/http/v1,./internal/domain/models --parseInternal
```

Create a `.env` file with configuration as stated in the `.env.example`

Build the project:

```bash
$ cd backend
$ go build -v -o ./bin/ ./...
```

or run it:

`$ go run ./...`

or containerize it — Dockerfile is in the `/backend/`
