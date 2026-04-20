package entity

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already taken")
	ErrEmailAlreadyExists    = errors.New("email already taken")
	ErrUserAlreadyExists     = errors.New("user already exists")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCredentials = errors.New("invalid username or password")

	ErrInternal = errors.New("an internal error occurred")
)
