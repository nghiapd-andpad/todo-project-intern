package todo

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewTodoGetter,
	NewTodoLister,
	NewTodoCreator,
	NewTodoUpdater,
	NewTodoDeleter,
)
