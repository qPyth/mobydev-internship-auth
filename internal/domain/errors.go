package domain

import "errors"

var (
	ErrEmailExists        = errors.New("user with email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidBDay        = errors.New("invalid birthdate")
	ErrInvalidPhone       = errors.New("invalid phone number")
)
