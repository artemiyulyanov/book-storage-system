package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=200"`
	Password string `json:"password" validate:"required,max=100"`
}

type RegisterRequest struct {
	FirstName      string `json:"first_name" validate:"required,max=100"`
	LastName       string `json:"last_name" validate:"required,max=100"`
	Email          string `json:"email" validate:"required,max=200"`
	Password       string `json:"password" validate:"required,max=100"`
	PasswordRepeat string `json:"password_repeat" validate:"required,max=100"`
}

type BookCreateRequest struct {
	ID          int64  `json:"-" db:"id"`
	Title       string `json:"title" db:"title" validate:"required,max=100"`
	Description string `json:"description" db:"description" validate:"required,max=200"`
	AuthorID    int64  `json:"-" db:"author_id"`
}

type BookUpdateRequest struct {
	ID          int64  `json:"-" db:"id"`
	Title       string `json:"title" db:"title" validate:"required,max=100"`
	Description string `json:"description" db:"description" validate:"required,max=200"`
	AuthorID    int64  `json:"-" db:"author_id"`
}
