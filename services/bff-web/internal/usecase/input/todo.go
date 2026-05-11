// Package input defines the input structures for the use cases related to todo lists and todos, encapsulating the data required for creating and updating todo lists and todos.
package input

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"

type CreateTodoListInput struct {
	// Format: users/{user_id}
	Parent      string
	DisplayName string
}

type CreateTodoInput struct {
	// Format: users/{user_id}/todo-lists/{todo_list_id}
	Parent      string
	Title       string
	Description *string
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *string
}

type UpdateTodoInput struct {
	// Format: users/{user_id}/todo-lists/{todo_list_id}/todos/{todo_id}
	Name        string
	Title       *string
	Description *string
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *string
}

type UpdateTodoListInput struct {
	Name        string
	DisplayName *string
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
