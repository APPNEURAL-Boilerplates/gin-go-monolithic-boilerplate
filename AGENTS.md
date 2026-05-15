# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project intent

This is a Gin-based Go modular monolith. Keep it as one deployable service, organized by feature modules.

## Architecture rules

- HTTP handlers live in `internal/modules/<feature>/handler.go`.
- Business logic lives in `internal/modules/<feature>/service.go`.
- Persistence contracts and implementations live in `internal/modules/<feature>/repository.go`.
- Shared API envelopes live in `internal/common/response`.
- Shared application errors live in `internal/common/apperror`.
- Cross-cutting HTTP concerns live in `internal/http/middleware`.

## Before finishing changes

Run:

```bash
gofmt -w .
go mod tidy
go vet ./...
go test ./...
```

## Do not

- Add a database dependency unless the task asks for it.
- Add a large framework on top of Gin unless there is a clear need.
- Put business logic inside Gin handlers.
- Return raw internal errors to clients.
- Commit `.env` files or secrets.
