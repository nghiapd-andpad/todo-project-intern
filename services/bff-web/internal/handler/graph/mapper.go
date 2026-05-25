package graph

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

func ToGraphQLTodoList(in *output.TodoListOutput) *TodoList {
	if in == nil {
		return nil
	}

	return &TodoList{
		ID:          in.ID,
		DisplayName: in.DisplayName,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func ToGraphQLTodo(in *output.TodoOutput) *Todo {
	if in == nil {
		return nil
	}

	return &Todo{
		ID:          in.ID,
		TodoListID:  in.TodoListID,
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.TodoStatus(in.Status),
		Priority:    entity.Priority(in.Priority),
		DueDate:     in.DueDate,
		AssigneeID:  in.AssigneeID,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func ToGraphQLTodoListPage(in *output.TodoListPage) *TodoListPage {
	if in == nil {
		return nil
	}

	res := make([]*TodoList, 0, len(in.TodoLists))
	for _, item := range in.TodoLists {
		res = append(res, ToGraphQLTodoList(item))
	}

	return &TodoListPage{
		TodoLists: res,
		Total:     int(in.Total),
	}
}

func ToGraphQLTodoPage(in *output.TodoPage) *TodoPage {
	if in == nil {
		return nil
	}

	res := make([]*Todo, 0, len(in.Todos))
	for _, item := range in.Todos {
		res = append(res, ToGraphQLTodo(item))
	}

	return &TodoPage{
		Todos: res,
		Total: int(in.Total),
	}
}

func ToGraphQLUser(in *output.UserOutput) *User {
	if in == nil {
		return nil
	}

	return &User{
		ID:       in.ID,
		Username: in.Username,
		Email:    in.Email,
	}
}
