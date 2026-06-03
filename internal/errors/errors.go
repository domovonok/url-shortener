package errors

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidUrl   = errors.New("invalid url")
	ErrInvalidCode  = errors.New("invalid code")
	ErrUrlNotFound  = errors.New("url not found")
)
