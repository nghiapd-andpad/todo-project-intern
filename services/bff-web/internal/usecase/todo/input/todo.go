package input

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"

type CreateTodoListInput struct {
	Parent      string
	DisplayName string
}

type CreateTodoInput struct {
	Title       string
	Description *string
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *string
}

type UpdateTodoInput struct {
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
