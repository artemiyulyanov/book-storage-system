package repositories

import (
	"book-storage-system/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

func (repo *BookRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	rows, err := repo.pool.Query(ctx, "SELECT id, title, description, author FROM books ORDER BY id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Book])

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (repo *BookRepository) GetBook(ctx context.Context, id int) (*models.Book, error) {
	var book models.Book

	err := repo.pool.QueryRow(ctx, "SELECT id, title, description, author FROM books WHERE id=$1", id).
		Scan(&book.ID, &book.Title, &book.Description, &book.Author)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		pool,
	}
}
