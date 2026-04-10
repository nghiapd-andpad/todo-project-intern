package persistence

import (
	"context"

	"github.com/nghiaphunng18/todos/internal/domain/entity"
	"github.com/nghiaphunng18/todos/internal/gateway"
	"github.com/nghiaphunng18/todos/internal/infra/persistence/mapper"
	"github.com/nghiaphunng18/todos/internal/infra/persistence/model"
	"gorm.io/gorm"
)

type todoCommandsGateway struct {
	db *gorm.DB
}

func NewTodoCommandsGateway(db *gorm.DB) gateway.TodoCommandsGateway {
	return &todoCommandsGateway{db: db}
}

func (g *todoCommandsGateway) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m := mapper.FromEntity(todo)

	if err := g.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}

	return mapper.ToEntity(m), nil
}

func (g *todoCommandsGateway) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	m := mapper.FromEntity(todo)

	if err := g.db.WithContext(ctx).Save(m).Error; err != nil {
		return nil, err
	}

	return mapper.ToEntity(m), nil
}

func (g *todoCommandsGateway) Delete(ctx context.Context, todoID entity.TodoID) error {
	if err := g.db.WithContext(ctx).Delete(&model.Todo{}, int64(todoID)).Error; err != nil {
		return err
	}
	return nil
}
