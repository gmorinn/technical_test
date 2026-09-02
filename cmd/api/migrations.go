package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gmorinn/technical_test/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func performMigrations(config *config.API) error {
	var migrationSourceURL = "file://internal/db/migrations"
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("could not connect to postgres: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("❌ ⛔️ could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationSourceURL,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("❌ migration failed: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("❌ migration up failed: %w", err)
	}

	log.Println("🧠 Database migrated successfully.")
	return nil
}
