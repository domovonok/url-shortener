package model

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrCodeNotFound = errors.New("code not found")
)
