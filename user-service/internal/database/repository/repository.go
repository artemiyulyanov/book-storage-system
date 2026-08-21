package repository

import (
	"common/auth"
	"common/models"
	"common/network/errors"
	"context"

	pb "common/proto/user"

	apperrors "common/network/errors"
	stderrors "errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5/pgconn"
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

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := repo.pool.QueryRow(ctx, "SELECT id, first_name, last_name, email, password_hash FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (int64, error) {
	if !auth.PasswordsMatch(req.Password, req.PasswordRepeat) {
		return 0, errors.ErrPasswordsMismatch
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	var id int64

	err = repo.pool.QueryRow(ctx, "INSERT INTO users (first_name, last_name, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id", req.FirstName, req.LastName, req.Email, passwordHash).
		Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if stderrors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrEmailTaken
		}
		return 0, err
	}

	return id, nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool,
	}
}
