package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type TodoLister struct {
	todoGateway gateway.TodoGateway
}

func NewTodoLister(todoGateway gateway.TodoGateway) *TodoLister {
	return &TodoLister{todoGateway: todoGateway}
}

func (u *TodoLister) ListTodoLists(ctx context.Context, parent string, opts *input.ListTodoListsOptions) (*output.TodoListPage, error) {
	result, err := u.todoGateway.ListTodoLists(ctx, parent, mapper.ToGatewayListTodoListsOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("TodoLister.ListTodoLists: %w", err)
	}

	return mapper.ToTodoListPage(result), nil
}

func (u *TodoLister) ListTodos(ctx context.Context, parent string, opts *input.ListTodosOptions) (*output.TodoPage, error) {
	result, err := u.todoGateway.ListTodos(ctx, parent, mapper.ToGatewayListTodosOptions(opts))
	if err != nil {
		return nil, fmt.Errorf("TodoLister.ListTodos: %w", err)
	}

	return mapper.ToTodoPage(result), nil
}
