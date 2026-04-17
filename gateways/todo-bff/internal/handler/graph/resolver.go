package graph

import "github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/usecase/auth"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	AuthUseCase auth.AuthUseCase
}

func NewResolver(authUC auth.AuthUseCase) *Resolver {
	return &Resolver{
		AuthUseCase: authUC,
	}
}
