package input

import "github.com/nghiaphunng18/todos/internal/domain/entity"

type TodoGetter struct {
	ID entity.TodoID
}

type TodoLister struct {
	Parent string // users/{u_id}/todo-lists/{l_id}
}

type TodoCreator struct {
	Title       string
	Description *string
}

type TodoUpdater struct {
	ID          entity.TodoID
	Title       *string
	Description *string
	Status      *entity.TodoStatus
	Priority    *entity.Priority
}

type TodoDeleter struct {
	ID entity.TodoID
}
