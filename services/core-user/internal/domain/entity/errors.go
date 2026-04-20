package entity

import "errors"

var (
	// user errors
	ErrUsernameAlreadyExists = errors.New("username already taken")
	ErrEmailAlreadyExists    = errors.New("email already taken")
	ErrUserAlreadyExists     = errors.New("user already exists")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCredentials = errors.New("invalid username or password")

	// token erros
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token has expired")
	ErrUntrustedMethod = errors.New("unexpected signing method")
	ErrTokenSigning    = errors.New("error signing token")

	// general errors
	ErrInternal = errors.New("an internal error occurred")
)
