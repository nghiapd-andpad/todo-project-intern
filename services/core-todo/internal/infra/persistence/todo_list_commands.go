package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type TodoListCommandsGateway struct {
	db *gorm.DB
}

func NewTodoListCommandsGateway(db *gorm.DB) *TodoListCommandsGateway {
	return &TodoListCommandsGateway{db: db}
}

var _ gateway.TodoListCommandsGateway = (*TodoListCommandsGateway)(nil)

func (g *TodoListCommandsGateway) Create(
	ctx context.Context,
	todoList *entity.TodoList,
) (*entity.TodoList, error) {
	conn := connFromContext(ctx, g.db)

	m := mapper.TodoListFromEntity(todoList)

	if err := conn.Create(m).Error; err != nil {
		return nil, fmt.Errorf("db create todo list: %w", err)
	}

	return mapper.TodoListToEntity(m), nil
}

func (g *TodoListCommandsGateway) Update(
	ctx context.Context,
	todoList *entity.TodoList,
) (*entity.TodoList, error) {
	conn := connFromContext(ctx, g.db)

	m := mapper.TodoListFromEntity(todoList)

	if err := conn.Save(m).Error; err != nil {
		return nil, fmt.Errorf("db update todo list: %w", err)
	}

	return mapper.TodoListToEntity(m), nil
}

func (g *TodoListCommandsGateway) Delete(
	ctx context.Context,
	todoListID entity.TodoListID,
) error {
	conn := connFromContext(ctx, g.db)

	if err := conn.Delete(&model.TodoList{}, int64(todoListID)).Error; err != nil {
		return fmt.Errorf("db delete todo list: %w", err)
	}

	return nil
}

func (g *TodoListCommandsGateway) HardDeleteTodoListsByIDs(
	ctx context.Context,
	ids []entity.TodoListID,
) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	conn := connFromContext(ctx, g.db)

	values := make([]int64, len(ids))
	for i, id := range ids {
		values[i] = int64(id)
	}

	result := conn.Unscoped().
		Where("id IN ?", values).
		Delete(&model.TodoList{})

	if result.Error != nil {
		return 0, fmt.Errorf("db hard delete todo lists by ids: %w", result.Error)
	}

	return result.RowsAffected, nil
}
