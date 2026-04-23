package graph

import (
	authusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth"
	todousecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo"
	userusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/user"
)

type Resolver struct {
	// Auth
	authUsecase  authusecase.Registerer
	loginUsecase authusecase.Loginer

	// Todo
	todoGetter  todousecase.TodoGetter
	todoLister  todousecase.TodoLister
	todoCreator todousecase.TodoCreator
	todoUpdater todousecase.TodoUpdater
	todoDeleter todousecase.TodoDeleter

	// User
	userGetter userusecase.UserGetter
}

func NewResolver(
	authUsecase authusecase.Registerer,
	loginUsecase authusecase.Loginer,
	todoGetter todousecase.TodoGetter,
	todoLister todousecase.TodoLister,
	todoCreator todousecase.TodoCreator,
	todoUpdater todousecase.TodoUpdater,
	todoDeleter todousecase.TodoDeleter,
	userGetter userusecase.UserGetter,
) *Resolver {
	return &Resolver{
		authUsecase:  authUsecase,
		loginUsecase: loginUsecase,
		todoGetter:   todoGetter,
		todoLister:   todoLister,
		todoCreator:  todoCreator,
		todoUpdater:  todoUpdater,
		todoDeleter:  todoDeleter,
		userGetter:   userGetter,
	}
}

func (r *Resolver) GetUserGetter() userusecase.UserGetter {
	return r.userGetter
}
