# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Start development (Docker — runs API + Postgres with hot reload via Air)
make dev

# Build the production image (the Dockerfile's runtime stage)
make build

# Run all tests
make test

# Run a single test package
go test ./internal/fizzbuzz/...

# Generate SQLC type-safe DB code from SQL queries
make sql

# Regenerate the OpenAPI spec in docs/ from the handler annotations
# (installs the swag generator on first use; no PATH setup needed)
make swagger

# Database migrations
make migrate-up
make migrate-down
make migrate-create name="migration_name"
make migrate-goto version=14
make migrate-force version=5   # fix dirty migration state

# Database dump / restore
make sql-dump
make sql-reset file=dump.sql
```

Requires only `cp .env.example .env` before the first `make dev`; compose creates its
own network and volumes.

### Docker layout

One `Dockerfile` with three stages. `docker-compose.yml` builds `--target dev`
(toolchain + pinned Air, working tree bind-mounted at `/app`); `make build` builds the
default `runtime` stage (Alpine, non-root, static binary + migrations). Both share the
`builder` stage, so the shipped image is exercised by an ordinary make target rather
than only by a deploy. Air is pinned in the Dockerfile — `@latest` now requires Go 1.26
and would break the 1.25 build. The `api` service waits on the `db` healthcheck via
`condition: service_healthy`, because Postgres accepts TCP before it accepts queries and
the API exits on a failed connection.

## Architecture

Go 1.25 REST API using:
- **Gin** — HTTP server and middleware
- **SQLC** — type-safe SQL → Go code generation
- **golang-migrate** — database migrations
- **PostgreSQL 17** — database (runs in Docker)

### Request flow

```
HTTP → Gin → Middleware → api.Handler → db.ExecTx → SQLC queries
```

### Key directories

| Path | Purpose |
|------|---------|
| `cmd/api/` | Entrypoint — wiring, HTTP server, graceful shutdown |
| `docs/` | Auto-generated OpenAPI spec — do not edit, run `make swagger` |
| `internal/api/` | REST handler implementations (fizzbuzz) |
| `internal/db/query/*.sql` | SQLC query definitions |
| `internal/db/*.sql.go` | Auto-generated SQLC code — do not edit |
| `internal/db/schema/schema.sql` | Current DB schema (updated only after successful migrations) |
| `internal/db/migrations/` | golang-migrate up/down SQL files |
| `internal/fizzbuzz/` | FizzBuzz domain logic — pure, no HTTP or DB |
| `internal/config/config.go` | All env var loading |
| `internal/middleware/` | Gin context injection |
| `internal/utils/` | Utility functions — all must have unit tests |
| `.dockerignore` | Keeps `.env` and `.git` out of the build context — never exclude `docs/` or `internal/db/**.sql` |

### REST endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/fizzbuzz` | — | FizzBuzz |
| GET | `/api/v1/fizzbuzz/stats` | — | FizzBuzz stats |
| GET | `/swagger/index.html` | — | Browsable API documentation |

### Adding a new endpoint (full workflow)

1. Create a migration if schema changes: `make migrate-create name="..."`, edit the file, run `make migrate-up`
2. Update `internal/db/schema/schema.sql` only after the migration succeeds
3. Add/edit SQL queries in `internal/db/query/<feature>.sql`
4. Run `make sql` to regenerate SQLC code
5. Add handler methods to `internal/api/<feature>.go` on the `Handler` struct
6. Register routes in `internal/api/handler.go` `RegisterRoutes()`
7. Annotate the handler (`@Summary`, `@Param`, `@Success`, `@Failure`) and run `make swagger`

### Conventions

- URL paths: kebab-case, plural nouns (`/api/v1/users`)
- SQL: snake_case table and column names, plural table names
- Always use `db.ExecTx` for multi-step database operations to ensure transactional safety
- Soft deletes: filter `WHERE deleted_at IS NULL`
- All utility functions in `internal/utils/` must have unit tests
- SQLC models in `internal/db/` are returned directly as JSON where possible

### Environment variables

All config is loaded from `.env` (read by `docker-compose.yml` and `Makefile`). Key vars: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`, `POSTGRES_HOST`, `API_PORT`, `API_SECRET`, `ENV` (dev/prod), `API_TZ`, `API_CORS`.
