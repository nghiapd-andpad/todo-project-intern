package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type TodoGetter struct {
	todoGateway gateway.TodoGateway
}

func NewTodoGetter(todoGateway gateway.TodoGateway) *TodoGetter {
	return &TodoGetter{todoGateway: todoGateway}
}

func (u *TodoGetter) GetTodoList(ctx context.Context, name string) (*output.TodoListOutput, error) {
	result, err := u.todoGateway.GetTodoList(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("TodoGetter.GetTodoList: %w", err)
	}

	return mapper.ToTodoListOutput(result), nil
}

func (u *TodoGetter) GetTodo(ctx context.Context, name string) (*output.TodoOutput, error) {
	result, err := u.todoGateway.GetTodo(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("TodoGetter.GetTodo: %w", err)
	}

	return mapper.ToTodoOutput(result), nil
}
