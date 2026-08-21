package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(connectionString, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connectionString,
	)

	if err != nil {
		return fmt.Errorf("Failed to init migrate: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Migrations: no changes to apply")
			return nil
		}

		return fmt.Errorf("Failed to apply migrations: %w", err)
	}

	log.Println("Migrations: applied successfully!")
	return nil
}
