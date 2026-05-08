package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/output"
)

type TodoLister struct {
	todoGateway gateway.TodoGateway
}

func NewTodoLister(todoGateway gateway.TodoGateway) *TodoLister {
	return &TodoLister{todoGateway: todoGateway}
}

func (u *TodoLister) ListTodoLists(ctx context.Context, parent string, opts input.ListTodoListsOptions) (*output.TodoListPage, error) {
	result, err := u.todoGateway.ListTodoLists(ctx, parent, mapper.ListTodoListsOptionsToGateway(opts))
	if err != nil {
		return nil, fmt.Errorf("TodoLister.ListTodoLists: %w", err)
	}

	return mapper.TodoListPageToUsecase(result), nil
}

func (u *TodoLister) ListTodos(ctx context.Context, parent string, opts input.ListTodosOptions) (*output.TodoPage, error) {
	result, err := u.todoGateway.ListTodos(ctx, parent, mapper.ListTodosOptionsToGateway(opts))
	if err != nil {
		return nil, fmt.Errorf("TodoLister.ListTodos: %w", err)
	}

	return mapper.TodoPageToUsecase(result), nil
}
