package persistence

import (
	"context"
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/infra/persistence/model"
	"gorm.io/gorm"
)

type todoQueriesGateway struct {
	db *gorm.DB
}

func NewTodoQueriesGateway(db *gorm.DB) gateway.TodoQueriesGateway {
	return &todoQueriesGateway{db: db}
}

func (todoQueriesGateway *todoQueriesGateway) Get(ctx context.Context, todoID entity.TodoID) (*entity.Todo, error) {
	var todo model.Todo
	if err := todoQueriesGateway.db.WithContext(ctx).First(&todo, todoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToEntity(&todo), nil
}

func (todoQueriesGateway *todoQueriesGateway) List(ctx context.Context) ([]*entity.Todo, error) {
	var models []model.Todo
	if err := todoQueriesGateway.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*entity.Todo, len(models))
	for i := range models {
		entities[i] = mapper.ToEntity(&models[i])
	}
	return entities, nil
}
