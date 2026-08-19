package repository

import (
	"common/models"
	"common/network/requests"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

func (repo *BookRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	rows, err := repo.pool.Query(ctx, "SELECT id, title, description, author_id FROM books ORDER BY id")

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

func (repo *BookRepository) GetBook(ctx context.Context, id int64) (*models.Book, error) {
	var book models.Book

	err := repo.pool.QueryRow(ctx, "SELECT id, title, description, author_id FROM books WHERE id=$1", id).
		Scan(&book.ID, &book.Title, &book.Description, &book.AuthorID)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (repo *BookRepository) CreateBook(ctx context.Context, req *requests.BookCreateRequest) (int64, error) {
	err := repo.pool.QueryRow(ctx, "INSERT INTO books (title, description, author_id) VALUES ($1, $2, $3) RETURNING id", req.Title, req.Description, req.AuthorID).
		Scan(&req.ID)

	if err != nil {
		return req.ID, err
	}

	return req.ID, nil
}

func (repo *BookRepository) UpdateBook(ctx context.Context, id int64, req *requests.BookUpdateRequest) (int64, error) {
	tag, err := repo.pool.Exec(ctx, "UPDATE books SET title=$1, description=$2 WHERE id=$3 AND author_id=$4", req.Title, req.Description, id, req.AuthorID)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (repo *BookRepository) DeleteBook(ctx context.Context, id int64, authorId int64) (int64, error) {
	tag, err := repo.pool.Exec(ctx, "DELETE FROM books WHERE id=$1 AND author_id=$2", id, authorId)

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
