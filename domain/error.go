package domain

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound = errors.New("Resource not found")
	ErrBadRequest = errors.New("Bad request error")
)
