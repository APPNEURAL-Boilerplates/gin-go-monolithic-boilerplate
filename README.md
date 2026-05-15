# Gin Monolithic Boilerplate

A production-minded Gin + Go modular monolith boilerplate.

It is a single deployable backend service organized by feature modules:

```txt
handler -> service -> repository
```

## Features

- Go 1.25+
- Gin 1.12+
- Modular monolith package layout
- Root endpoint
- Health endpoint
- Users example module
- Handler -> service -> repository pattern
- In-memory repository for instant local startup
- DTO binding and validation
- Consistent JSON response envelope
- Centralized application error type
- Request ID middleware
- Structured logging with `log/slog`
- Security headers middleware
- CORS middleware
- Graceful shutdown
- Tests with `net/http/httptest`
- Dockerfile
- docker-compose.yml
- Makefile
- AGENTS.md for Codex or other coding agents

## Requirements

- Go 1.25 or newer. Go 1.26.3 is recommended.
- Docker, optional.

## Getting started

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api
```

Open:

```txt
http://localhost:8080
http://localhost:8080/health
http://localhost:8080/api/v1/users
```

## Commands

```bash
make dev        # run locally
make test       # run tests
make vet        # run go vet
make check      # format, tidy, vet, and test
make build      # build binary into ./bin/api
make docker-up  # run with Docker Compose
```

Equivalent raw Go commands:

```bash
go run ./cmd/api
go test ./...
go vet ./...
go build -o bin/api ./cmd/api
```

## API routes

```txt
GET    /                         Root metadata
GET    /health                   Health check
GET    /api/v1/users             List users
POST   /api/v1/users             Create user
GET    /api/v1/users/:id         Get user by ID
PATCH  /api/v1/users/:id         Update user
DELETE /api/v1/users/:id         Delete user
```

## Example request

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "content-type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
```

## Response format

Success:

```json
{
  "ok": true,
  "data": {},
  "requestId": "..."
}
```

Error:

```json
{
  "ok": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "The requested endpoint was not found."
  },
  "requestId": "..."
}
```

## Project structure

```txt
gin-monolithic-boilerplate/
├─ cmd/
│  └─ api/
│     └─ main.go
├─ internal/
│  ├─ app/
│  │  └─ router.go
│  ├─ common/
│  │  ├─ apperror/
│  │  ├─ id/
│  │  ├─ requestid/
│  │  └─ response/
│  ├─ config/
│  ├─ http/
│  │  └─ middleware/
│  └─ modules/
│     ├─ health/
│     ├─ root/
│     └─ users/
├─ test/
├─ http/
├─ Dockerfile
├─ docker-compose.yml
├─ Makefile
├─ AGENTS.md
└─ go.mod
```

## Adding a new module

For a new feature, create:

```txt
internal/modules/<feature>/
├─ dto.go
├─ handler.go
├─ model.go
├─ repository.go
└─ service.go
```

Then register its routes in `internal/app/router.go`.

## Replacing the in-memory repository

The users module currently stores data in memory so the app runs instantly.

For production, replace `InMemoryRepository` with a database-backed implementation while keeping the `Repository` interface stable.

Good options:

- PostgreSQL with `database/sql`
- PostgreSQL with `pgx`
- GORM
- sqlc
- Ent
- MongoDB

## Environment variables

See `.env.example`.

```txt
APP_NAME=gin-monolithic-boilerplate
APP_ENV=local
PORT=8080
GIN_MODE=debug
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=*
SHUTDOWN_TIMEOUT_SECONDS=10
```

Do not commit real `.env` files or secrets.

## Docker

```bash
docker compose up --build
```

The API runs on:

```txt
http://localhost:8080
```
