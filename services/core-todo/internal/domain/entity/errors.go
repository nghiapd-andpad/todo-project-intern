package entity

import "fmt"

type ErrorCode string

const (
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrInvalidParameter ErrorCode = "INVALID_PARAMETER"
	ErrAuthZ            ErrorCode = "AUTHORIZATION"
	ErrAuthN            ErrorCode = "AUTHENTICATION"
	ErrInternal         ErrorCode = "INTERNAL"
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

func NewInternal(message string) *AppError {
	return &AppError{Code: ErrInternal, Message: message}
}
