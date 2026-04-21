package todos

import "github.com/google/wire"

var WireSet = wire.NewSet(
	// Todo usecases
	NewTodoCreator,
	NewTodoGetter,
	NewTodoLister,
	NewTodoUpdater,
	NewTodoDeleter,
	// TodoList usecases
	NewTodoListCreator,
	NewTodoListGetter,
	NewTodoListLister,
	NewTodoListUpdater,
	NewTodoListDeleter,
)
