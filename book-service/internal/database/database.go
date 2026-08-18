package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("Could not resolve connectionString at database.NewDatabaseInstance()")
	}

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("Pooling database failed")
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Cannot establish a connection with database")
	}

	return pool, nil
}
