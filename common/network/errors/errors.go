package errors

import "errors"

var (
	ErrEmailTaken        = errors.New("email already registered")
	ErrPasswordsMismatch = errors.New("passwords do not match")
)
