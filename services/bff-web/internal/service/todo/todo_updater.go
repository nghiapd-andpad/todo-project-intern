package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type TodoUpdater struct {
	todoGateway gateway.TodoGateway
}

func NewTodoUpdater(todoGateway gateway.TodoGateway) *TodoUpdater {
	return &TodoUpdater{todoGateway: todoGateway}
}

func (u *TodoUpdater) UpdateTodoList(ctx context.Context, input *input.UpdateTodoListInput) (*output.TodoListOutput, error) {
	result, err := u.todoGateway.UpdateTodoList(ctx, inputgateway.UpdateTodoListInput{
		Name:        input.Name,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("TodoUpdater.UpdateTodoList: %w", err)
	}

	return mapper.ToTodoListOutput(result), nil
}

func (u *TodoUpdater) UpdateTodo(ctx context.Context, input *input.UpdateTodoInput) (*output.TodoOutput, error) {
	result, err := u.todoGateway.UpdateTodo(ctx, inputgateway.UpdateTodoInput{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		AssigneeID:  input.AssigneeID,
	})
	if err != nil {
		return nil, fmt.Errorf("TodoUpdater.UpdateTodo: %w", err)
	}

	return mapper.ToTodoOutput(result), nil
}
