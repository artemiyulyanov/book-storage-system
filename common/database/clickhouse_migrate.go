package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func RunClickHouseMigrations(databaseURL, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to init clickhouse migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("clickhouse migrations: no changes to apply")
			return nil
		}
		return fmt.Errorf("failed to apply clickhouse migrations: %w", err)
	}

	log.Println("clickhouse migrations: applied successfully")
	return nil
}
