package domain

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound = errors.New("Resource not found")
	ErrBadRequest = errors.New("Bad request error")
	ErrCacheKeyEmpty = errors.New("Resource not found")
	ErrCacheKeyExists = errors.New("Cache exists. Not writing to cache")
)
