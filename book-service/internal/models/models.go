package models

type Book struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title" validate:"required,max=100"`
	Description string `json:"description" db:"description" validate:"required,max=200"`
	AuthorID    int64  `json:"author_id" db:"author_id" validate:"required,min=1"`
}
