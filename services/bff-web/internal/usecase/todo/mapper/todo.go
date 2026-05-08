// Package mapper provides functions to map between different layers of the application, such as mapping input from usecases to gateway input.
package mapper

import (
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	outputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/output"
	intputusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/input"
	outputusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/output"
)

func ListTodoListsOptionsToGateway(opts intputusecase.ListTodoListsOptions) inputgateway.ListTodoListsOptions {
	return inputgateway.ListTodoListsOptions{
		NameSearch: opts.NameSearch,
		Offset:     opts.Offset,
		Limit:      opts.Limit,
	}
}

func ListTodosOptionsToGateway(opts intputusecase.ListTodosOptions) inputgateway.ListTodosOptions {
	return inputgateway.ListTodosOptions{
		Status:      opts.Status,
		Priority:    opts.Priority,
		TitleSearch: opts.TitleSearch,
		Offset:      opts.Offset,
		Limit:       opts.Limit,
	}
}

func TodoListPageToUsecase(page *outputgateway.TodoListPage) *outputusecase.TodoListPage {
	if page == nil {
		return nil
	}

	return &outputusecase.TodoListPage{
		TodoLists: page.TodoLists,
		Total:     page.Total,
	}
}

func TodoPageToUsecase(page *outputgateway.TodoPage) *outputusecase.TodoPage {
	if page == nil {
		return nil
	}

	return &outputusecase.TodoPage{
		Todos: page.Todos, // Giữ nguyên slice entity
		Total: page.Total,
	}
}
