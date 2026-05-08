// Package input defines the input structures for the gateway operations.
package input

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"

type CreateTodoListInput struct {
	// Parent is the resource name of the parent, in the format "users/{user_id}".
	Parent      string
	DisplayName string
}

type UpdateTodoListInput struct {
	// Name is the resource name of the todo list, in the format "users/{user_id}/todo-lists/{list_id}".
	Name        string
	DisplayName *string
}

type CreateTodoInput struct {
	// Parent is the resource name of the parent, in the format "users/{user_id}/todo-lists/{list_id}".
	Parent      string
	Title       string
	Description *string
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *string
}

type UpdateTodoInput struct {
	// Name is the resource name of the todo, in the format "users/{user_id}/todo-lists/{list_id}/todos/{todo_id}".
	Name        string
	Title       *string
	Description *string
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *string
}

type ListTodoListsOptions struct {
	NameSearch *string
	Offset     int
	Limit      int
}

type ListTodosOptions struct {
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	TitleSearch *string
	Offset      int
	Limit       int
}
