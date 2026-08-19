package repository

import (
	"common/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func (repo *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	rows, err := repo.pool.Query(ctx, "SELECT * FROM users ORDER BY id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	var user models.User

	err := repo.pool.QueryRow(ctx, "SELECT id, first_name, last_name, email, password_hash FROM users WHERE id=$1", id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool,
	}
}
