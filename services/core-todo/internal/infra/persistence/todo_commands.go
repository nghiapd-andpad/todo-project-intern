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

type TodoCommandsGateway struct {
	db *gorm.DB
}

var _ gateway.TodoCommandsGateway = (*TodoCommandsGateway)(nil)

func NewTodoCommandsGateway(db *gorm.DB) *TodoCommandsGateway {
	return &TodoCommandsGateway{db: db}
}

func (g *TodoCommandsGateway) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	conn := connFromContext(ctx, g.db)

	m := mapper.TodoFromEntity(todo)

	if err := conn.Create(m).Error; err != nil {
		return nil, fmt.Errorf("db create todo: %w", err)
	}

	return mapper.TodoToEntity(m), nil
}

func (g *TodoCommandsGateway) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	conn := connFromContext(ctx, g.db)

	m := mapper.TodoFromEntity(todo)

	if err := conn.Save(m).Error; err != nil {
		return nil, fmt.Errorf("db update todo: %w", err)
	}

	return mapper.TodoToEntity(m), nil
}

func (g *TodoCommandsGateway) Delete(ctx context.Context, todoID entity.TodoID) error {
	conn := connFromContext(ctx, g.db)

	// soft delete
	result := conn.Delete(&model.Todo{}, int64(todoID))

	if result.Error != nil {
		return fmt.Errorf("db delete todo: %w", result.Error)
	}

	return nil
}

func (g *TodoCommandsGateway) DeleteByTodoListID(ctx context.Context, todoListID entity.TodoListID) error {
	conn := connFromContext(ctx, g.db)

	if err := conn.Where("todo_list_id = ?", int64(todoListID)).Delete(&model.Todo{}).Error; err != nil {
		return fmt.Errorf("db delete todos by list id: %w", err)
	}

	return nil
}
