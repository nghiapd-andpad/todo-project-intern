package entity

import "fmt"

type ErrorCode string

const (
	// Common error codes
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrInvalidParameter ErrorCode = "INVALID_PARAMETER"
	ErrAuthZ            ErrorCode = "AUTHORIZATION"
	ErrAuthN            ErrorCode = "AUTHENTICATION"
	ErrAlreadyHandled   ErrorCode = "AlreadyHandled"
	ErrInternal         ErrorCode = "INTERNAL"

	// User-specific error codes
	ErrInvalidCredentials    ErrorCode = "INVALID_CREDENTIALS"
	ErrUsernameAlreadyExists ErrorCode = "USERNAME_ALREADY_EXISTS"
	ErrEmailAlreadyExists    ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrInvalidToken          ErrorCode = "INVALID_TOKEN"
	ErrExpiredToken          ErrorCode = "EXPIRED_TOKEN"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details map[string]string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) WithDetail(key, value string) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// General constructors
func NewNotFound(message string) *AppError {
	return &AppError{Code: ErrNotFound, Message: message}
}

func NewInvalidParameter(message string) *AppError {
	return &AppError{Code: ErrInvalidParameter, Message: message}
}

func NewAuthZ(message string) *AppError {
	return &AppError{Code: ErrAuthZ, Message: message}
}

func NewAuthN(message string) *AppError {
	return &AppError{Code: ErrAuthN, Message: message}
}

func NewAlreadyHandled(message string) *AppError {
	return &AppError{Code: ErrAlreadyHandled, Message: message}
}

func NewInternal(message string) *AppError {
	return &AppError{Code: ErrInternal, Message: message}
}

func NewInvalidCredentials() *AppError {
	return &AppError{Code: ErrInvalidCredentials, Message: "invalid username or password"}
}

func NewUsernameAlreadyExists() *AppError {
	return &AppError{Code: ErrUsernameAlreadyExists, Message: "username already taken"}
}

func NewEmailAlreadyExists() *AppError {
	return &AppError{Code: ErrEmailAlreadyExists, Message: "email already taken"}
}

func NewInvalidToken() *AppError {
	return &AppError{Code: ErrInvalidToken, Message: "invalid token"}
}

func NewExpiredToken() *AppError {
	return &AppError{Code: ErrExpiredToken, Message: "token has expired"}
}
