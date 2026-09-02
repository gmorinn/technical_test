# GM REST API

A REST API built with Go, Gin, SQLC, and PostgreSQL.

## Setup

One step before the first run — create your `.env`:

```bash
cp .env.example .env
```

`.env` is gitignored, and the `Makefile` opens with `include .env` — without this file
*every* make target fails, not just `make dev`. The API also validates its environment
at startup and will refuse to boot listing any variable that is missing or malformed.

**Stack:**

- **Gin** — HTTP server and routing
- **SQLC** — type-safe SQL → Go code generation
- **golang-migrate** — database migrations
- **PostgreSQL 17** — database (runs in Docker)
- **swaggo** — OpenAPI spec generated from handler annotations, served at `/swagger/index.html`
- **golang-jwt** — JWT authentication

## Requirements

Go 1.25 or higher (see `go.mod`)

Install **go migrate**:

On MacOS:

```bash
brew install golang-migrate
```

## How to launch the API

```bash
make dev
```

## REST Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/fizzbuzz` | — | FizzBuzz |
| GET | `/api/v1/fizzbuzz/stats` | — | FizzBuzz stats |
| GET | `/swagger/index.html` | — | Browsable API documentation |

## How to add a new endpoint (study case: Get user by ID)

1. **Create a migration** if the schema changes:
   ```bash
   make migrate-create name="add_companies_table"
   ```
   Edit the generated file in `internal/db/migrations/`, then apply it:
   ```bash
   make migrate-up
   ```

2. **Update `internal/db/schema/schema.sql`** only after the migration succeeds.

3. **Add a SQL query** in `internal/db/query/<feature>.sql`:
   ```sql
   -- name: GetUserByID :one
   SELECT * FROM users
   WHERE id = $1
   AND deleted_at IS NULL
   LIMIT 1;
   ```

4. **Regenerate SQLC code**:
   ```bash
   make sql
   ```

5. **Add a handler** in `internal/api/<feature>.go`:
   ```go
   // Each feature declares the slice of the database it needs, next to the
   // handlers that use it, so handlers can be tested against a fake.
   type UserStore interface {
       GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
   }

   // GetUser godoc
   //
   //	@Summary	Fetch a user
   //	@Tags		users
   //	@Produce	json
   //	@Param		id	path		string	true	"User ID"
   //	@Success	200	{object}	db.User
   //	@Failure	400	{object}	map[string]string
   //	@Failure	404	{object}	map[string]string
   //	@Router		/users/{id} [get]
   func (h *Handler) GetUser(c *gin.Context) {
       id, err := uuid.Parse(c.Param("id"))
       if err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
           return
       }

       user, err := h.DB.GetUserByID(c.Request.Context(), id)
       if err != nil {
           if errors.Is(err, sql.ErrNoRows) {
               c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
               return
           }
           c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
           return
       }

       c.JSON(http.StatusOK, user)
   }
   ```

   Then run `make swagger` to regenerate `docs/` from the annotations.

6. **Register the route** in `internal/api/handler.go` `RegisterRoutes()`:
   ```go
   users.GET("/:id", h.GetUser)
   ```

## Architecture

```
.
├── Dockerfile                # dev stage (compose) + runtime stage (make build)
├── .dockerignore
├── docker-compose.yml
├── Makefile
├── Readme.md
├── air.toml
├── cmd/
│   └── api/
│       ├── main.go           # Wiring, HTTP server, graceful shutdown
│       └── migrations.go     # Runs golang-migrate on boot
├── internal/                 # Import-fenced by the compiler
│   ├── api/
│   │   ├── fizzbuzz.go       # FizzBuzz handlers
│   │   └── handler.go        # Handler struct, RegisterRoutes
│   ├── config/config.go      # Env var loading (fails fast)
│   ├── db/
│   │   ├── migrations/       # golang-migrate up/down SQL files
│   │   ├── query/            # SQLC query definitions
│   │   │   └── fizzbuzz.sql
│   │   ├── schema/
│   │   │   └── schema.sql    # Current DB schema
│   │   ├── *.sql.go          # Auto-generated SQLC code — do not edit
│   │   ├── db.go
│   │   ├── models.go
│   │   └── repository.go
│   ├── fizzbuzz/             # Domain logic — pure, no HTTP or DB
│   │   ├── fizzbuzz.go
│   │   └── fizzbuzz_test.go
│   ├── middleware/           # Gin context injection
│   └── utils/                # Utility functions (all must have unit tests)
└── sqlc.yaml
```

### Request flow

```
HTTP → Gin → Middleware → api.Handler → db.ExecTx → SQLC queries
```

## Commands

```bash
# Start development (Docker — API + Postgres with hot reload via Air)
make dev

# Build the production image — the runtime stage of the same Dockerfile
make build

# Run all tests
make test

# Run a single test package
go test ./internal/fizzbuzz/...

# Generate SQLC type-safe DB code from SQL queries
make sql

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

## DB Dev/Prod To Local

```bash
pg_dump -v --inserts -d postgresql://{USER}:{PASSWORD}@{HOST}:{PORT}/{DB} -F t > ./backup.dump && \
pg_restore --verbose --clean --no-acl --no-owner -d postgresql://postgres:postgres_docker@0.0.0.0:54320 backup.dump
```

## Clear dirty database

```bash
migrate -path internal/db/migrations -database "postgresql://postgres_user:postgres_password@localhost:54320/postgres_db?sslmode=disable" -verbose force {DIRTY_VERSION}
make migrate-down
```

## Go to migration version

```bash
migrate -path internal/db/migrations -database "postgresql://postgres_user:postgres_password@localhost:54320/postgres_db?sslmode=disable" goto 14
```

## Deployment (if cannot connect to DB with correct credentials)

```bash
sudo iptables -A INPUT -p tcp --dport 5432 -j ACCEPT
sudo apt-get install iptables-persistent
```

## Conventions

- URL paths: kebab-case, plural nouns (`/api/v1/users`)
- SQL: snake_case table and column names, plural table names
- Always use `db.ExecTx` for multi-step database operations
- Soft deletes: filter `WHERE deleted_at IS NULL`
- All utility functions in `internal/utils/` must have unit tests
- SQLC models in `internal/db/` are returned directly as JSON where possible

## Environment variables

All config is loaded from `.env` (read by `docker-compose.yml` and `Makefile`).

| Variable | Description |
|----------|-------------|
| `POSTGRES_USER` | DB username |
| `POSTGRES_PASSWORD` | DB password |
| `POSTGRES_DB` | DB name |
| `POSTGRES_PORT` | DB port |
| `POSTGRES_HOST` | DB host |
| `API_PORT` | API listen port |
| `API_SECRET` | JWT signing secret |
| `ENV` | `dev` or `prod` |
| `API_TZ` | Timezone |
| `API_CORS` | Allowed CORS origins |
