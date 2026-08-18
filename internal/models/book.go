package models

type Book struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" validate:"required,max=100"`
	Description string `json:"description" db:"description" validate:"required,max=200"`
	Author      string `json:"author" db:"author" validate:"required,max=100"`
}
