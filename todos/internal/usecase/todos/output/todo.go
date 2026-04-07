package output

import "github.com/nghiaphunng18/todos/internal/domain/entity"

type TodoGetter struct {
	Todo *entity.Todo
}

type TodoListReader struct {
	Todos []*entity.Todo
}

type TodoCreator struct {
	Todo *entity.Todo
}

type TodoUpdater struct {
	Todo *entity.Todo
}

type TodoDeleter struct {
	Success bool
}
