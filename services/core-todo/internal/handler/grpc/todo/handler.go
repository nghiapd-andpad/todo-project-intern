package todo

import (
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
)

type TodoHandler struct {
	todov1.UnimplementedTodosServiceServer
	TodoCreator    todos.TodoCreator
	TodoGetter     todos.TodoGetter
	TodoListReader todos.TodoLister
	TodoUpdater    todos.TodoUpdater
	TodoDeleter    todos.TodoDeleter
}

func NewTodoHandler(
	creator todos.TodoCreator,
	getter todos.TodoGetter,
	listReader todos.TodoLister,
	updater todos.TodoUpdater,
	deleter todos.TodoDeleter,
) *TodoHandler {
	return &TodoHandler{
		TodoCreator:    creator,
		TodoGetter:     getter,
		TodoListReader: listReader,
		TodoUpdater:    updater,
		TodoDeleter:    deleter,
	}
}
