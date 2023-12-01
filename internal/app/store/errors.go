package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound         = errors.New("record not found")
	ErrEmailORPasswordInvalid = errors.New("invalid email or password")
	// ErrInvalidPassword        = errors.New("invalid password")
)
