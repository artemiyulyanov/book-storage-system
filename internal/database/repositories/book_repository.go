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

func (repo *BookRepository) CreateBook(ctx context.Context, book *models.Book) (int, error) {
	err := repo.pool.QueryRow(ctx, "INSERT INTO books (title, description, author) VALUES ($1, $2, $3) RETURNING id", book.Title, book.Description, book.Author).
		Scan(&book.ID)

	if err != nil {
		return book.ID, err
	}

	return book.ID, nil
}

func (repo *BookRepository) UpdateBook(ctx context.Context, id int, book *models.Book) (int64, error) {
	tag, err := repo.pool.Exec(ctx, "UPDATE books SET title=$1, description=$2, author=$3 WHERE id=$4", book.Title, book.Description, book.Author, id)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (repo *BookRepository) DeleteBook(ctx context.Context, id int) (int64, error) {
	tag, err := repo.pool.Exec(ctx, "DELETE FROM books WHERE id=$1", id)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		pool,
	}
}
