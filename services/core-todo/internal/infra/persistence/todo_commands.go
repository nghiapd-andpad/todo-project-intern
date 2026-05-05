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
	m := mapper.TodoFromEntity(todo)

	if err := g.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, fmt.Errorf("db create todo: %w", err)
	}

	return mapper.TodoToEntity(m), nil
}

func (g *TodoCommandsGateway) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m := mapper.TodoFromEntity(todo)

	if err := g.db.WithContext(ctx).Save(m).Error; err != nil {
		return nil, fmt.Errorf("db update todo: %w", err)
	}

	return mapper.TodoToEntity(m), nil
}

func (g *TodoCommandsGateway) Delete(ctx context.Context, todoID entity.TodoID) error {
	// soft delete
	result := g.db.WithContext(ctx).Delete(&model.Todo{}, int64(todoID))

	if result.Error != nil {
		return fmt.Errorf("db delete todo: %w", result.Error)
	}

	return nil
}
