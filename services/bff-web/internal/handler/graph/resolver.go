package graph

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase"

type Resolver struct {
	todoCreator usecase.TodoCreator
	todoGetter  usecase.TodoGetter
	todoLister  usecase.TodoLister
	todoUpdater usecase.TodoUpdater
	todoDeleter usecase.TodoDeleter

	authRegisterer usecase.AuthRegisterer
	authLoginer    usecase.AuthLoginer

	userGetter usecase.UserGetter
}

func NewResolver(
	todoCreator usecase.TodoCreator,
	todoGetter usecase.TodoGetter,
	todoLister usecase.TodoLister,
	todoUpdater usecase.TodoUpdater,
	todoDeleter usecase.TodoDeleter,
	authRegisterer usecase.AuthRegisterer,
	authLoginer usecase.AuthLoginer,
	userGetter usecase.UserGetter,
) *Resolver {
	return &Resolver{
		todoCreator:    todoCreator,
		todoGetter:     todoGetter,
		todoLister:     todoLister,
		todoUpdater:    todoUpdater,
		todoDeleter:    todoDeleter,
		authRegisterer: authRegisterer,
		authLoginer:    authLoginer,
		userGetter:     userGetter,
	}
}

func (r *Resolver) UserGetter() usecase.UserGetter {
	return r.userGetter
}
