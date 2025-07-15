package models

import "github.com/pkg/errors"

var (
	ErrInvalidInput = errors.New("Invalid input")
	ErrNotFound     = errors.New("Not found")
	ErrUnauthorized = errors.New("Unauthorized")
)
