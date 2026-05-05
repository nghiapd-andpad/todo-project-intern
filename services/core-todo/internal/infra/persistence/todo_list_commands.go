package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type TodoListCommandsGateway struct {
	db *gorm.DB
}

func NewTodoListCommandsGateway(db *gorm.DB) *TodoListCommandsGateway {
	return &TodoListCommandsGateway{db: db}
}

func (g *TodoListCommandsGateway) Create(
	ctx context.Context,
	todoList *entity.TodoList,
) (*entity.TodoList, error) {
	m := mapper.TodoListFromEntity(todoList)

	if err := g.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, fmt.Errorf("db create todo list: %w", err)
	}

	return mapper.TodoListToEntity(m), nil
}

func (g *TodoListCommandsGateway) Update(
	ctx context.Context,
	todoList *entity.TodoList,
) (*entity.TodoList, error) {
	m := mapper.TodoListFromEntity(todoList)

	if err := g.db.WithContext(ctx).Save(m).Error; err != nil {
		return nil, fmt.Errorf("db update todo list: %w", err)
	}

	return mapper.TodoListToEntity(m), nil
}

func (g *TodoListCommandsGateway) Delete(
	ctx context.Context,
	todoListID entity.TodoListID,
) error {
	result := g.db.WithContext(ctx).Delete(&model.TodoList{}, int64(todoListID))

	if result.Error != nil {
		return fmt.Errorf("db delete todo list: %w", result.Error)
	}

	return nil
}
