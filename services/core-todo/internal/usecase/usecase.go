// Package usecase defines interfaces for todo application use cases.
package usecase

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoCreator interface {
	Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error)
}

type TodoGetter interface {
	Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}

type TodoLister interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}

type TodoUpdater interface {
	Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error)
}

type TodoDeleter interface {
	Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error)
}

type TodoListCreator interface {
	Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error)
}

type TodoListGetter interface {
	Get(ctx context.Context, in *input.TodoListGetter) (*output.TodoListGetter, error)
}

type TodoListLister interface {
	List(ctx context.Context, in *input.TodoListLister) (*output.TodoListLister, error)
}

type TodoListUpdater interface {
	Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error)
}

type TodoListDeleter interface {
	Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error)
}

type TodoOverdueMarker interface {
	MarkOverdue(ctx context.Context, in *input.TodoOverdueMarker) (*output.TodoOverdueMarker, error)
}
