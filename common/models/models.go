package models

type Book struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title" validate:"required,max=100"`
	Description string `json:"description" db:"description" validate:"required,max=200"`
	AuthorID    int64  `json:"author_id" db:"author_id" validate:"required,min=1"`
}

type User struct {
	ID           int64  `json:"id" db:"id"`
	FirstName    string `json:"first_name" db:"first_name" validate:"required,max=100"`
	LastName     string `json:"last_name" db:"last_name" validate:"required,max=100"`
	Email        string `json:"email" db:"email" validate:"required,email,max=200"`
	PasswordHash string `json:"-" db:"password_hash"`
}
