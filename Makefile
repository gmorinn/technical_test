include .env
export

DIR=$(notdir $(shell pwd))
export DIR

GO_BIN_DIR := $(shell go env GOBIN)
ifeq ($(GO_BIN_DIR),)
GO_BIN_DIR := $(shell go env GOPATH)/bin
endif

SWAG_VERSION := v1.16.6
SWAG := $(shell command -v swag 2>/dev/null || echo $(GO_BIN_DIR)/swag)

migrate-up:
	@echo "\n\t⬆️\n"
	migrate -path internal/db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migrate-down:
	@echo "\n\t⬇️\n"
	migrate -path internal/db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down

migrate-create:
	@echo "\n\t📝\n"
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

test:
	@echo "\n\t❌\n"
	@go test -cover ./...

sql:
	@echo "\n\t🧠\n"
	@sqlc generate

swagger:
	@echo "\n\t📗\n"
	@test -x "$(SWAG)" || (echo "installing swag $(SWAG_VERSION)…" && go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION))
	@$(SWAG) init -g cmd/api/main.go --parseInternal

dev:
	@echo "\n\t💣\n"
	docker compose -p $(DIR) up --build --force-recreate --remove-orphans

build:
	@echo "\n\t🐳\n"
	docker build -t $(DIR):latest .

sql-dump:
	@echo "\n\t🧠 Creating SQL Dump\n"
	@PGPASSWORD=$(POSTGRES_PASSWORD) pg_dump -h localhost -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -F p -v -f ./dump.sql

sql-reset:
	@echo "\n\t🧠 Resetting Database and Restoring from SQL Dump\n"
	$(eval DUMP_FILE := $(if $(file),$(file),./dump.sql))
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -h localhost -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -h localhost -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f $(DUMP_FILE)

migrate-goto:
	@echo "\n\t🧠\n"
	@if [ -z "$(version)" ]; then \
		echo "Error: Please specify a version number. Usage: make migrate-goto version=<number>"; \
		exit 1; \
	fi
	migrate -path internal/db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" goto $(version)

migrate-force:
	@echo "\n\t🧠\n"
	migrate -path internal/db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" force $(version)

.PHONY: dev build test migrate-up migrate-down migrate-create sql swagger sql-dump sql-reset migrate-goto migrate-force
