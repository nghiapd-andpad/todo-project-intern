// Package mapper provides functions to map between different layers of the application, such as mapping input from usecases to gateway input.
package mapper

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	outputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/output"
	intputusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	outputusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

func ToGatewayListTodoListsOptions(opts *intputusecase.ListTodoListsOptions) inputgateway.ListTodoListsOptions {
	return inputgateway.ListTodoListsOptions{
		NameSearch: opts.NameSearch,
		Offset:     opts.Offset,
		Limit:      opts.Limit,
	}
}

func ToGatewayListTodosOptions(opts *intputusecase.ListTodosOptions) inputgateway.ListTodosOptions {
	return inputgateway.ListTodosOptions{
		Status:      opts.Status,
		Priority:    opts.Priority,
		TitleSearch: opts.TitleSearch,
		Offset:      opts.Offset,
		Limit:       opts.Limit,
	}
}

func ToTodoListPage(out *outputgateway.TodoListPage) *outputusecase.TodoListPage {
	if out == nil {
		return nil
	}

	res := make([]*outputusecase.TodoListOutput, 0, len(out.TodoLists))
	for _, item := range out.TodoLists {
		res = append(res, ToTodoListOutput(item))
	}

	return &outputusecase.TodoListPage{
		TodoLists: res,
		Total:     out.Total,
	}
}

func ToTodoPage(out *outputgateway.TodoPage) *outputusecase.TodoPage {
	if out == nil {
		return nil
	}

	res := make([]*outputusecase.TodoOutput, 0, len(out.Todos))
	for _, item := range out.Todos {
		res = append(res, ToTodoOutput(item))
	}

	return &outputusecase.TodoPage{
		Todos: res,
		Total: out.Total,
	}
}

func ToTodoListOutput(out *entity.TodoList) *outputusecase.TodoListOutput {
	if out == nil {
		return nil
	}

	return &outputusecase.TodoListOutput{
		ID:          out.ID,
		DisplayName: out.DisplayName,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	}
}

func ToTodoOutput(out *entity.Todo) *outputusecase.TodoOutput {
	if out == nil {
		return nil
	}

	return &outputusecase.TodoOutput{
		ID:          out.ID,
		TodoListID:  out.TodoListID,
		Title:       out.Title,
		Description: out.Description,
		Status:      out.Status,
		Priority:    out.Priority,
		DueDate:     out.DueDate,
		CreatorID:   out.CreatorID,
		AssigneeID:  out.AssigneeID,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	}
}
