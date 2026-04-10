package todos

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewTodoCreator,
	NewTodoGetter,
	NewTodoListReader,
	NewTodoUpdater,
	NewTodoDeleter,
)
