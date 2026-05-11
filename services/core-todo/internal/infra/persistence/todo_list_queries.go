package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type TodoListQueriesGateway struct {
	db *gorm.DB
}

func NewTodoListQueriesGateway(db *gorm.DB) *TodoListQueriesGateway {
	return &TodoListQueriesGateway{db: db}
}

var _ gateway.TodoListQueriesGateway = (*TodoListQueriesGateway)(nil)

func (g *TodoListQueriesGateway) Get(
	ctx context.Context,
	todoListID entity.TodoListID,
) (*entity.TodoList, error) {
	var m model.TodoList

	err := g.db.WithContext(ctx).First(&m, int64(todoListID)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db get todo list: %w", err)
	}

	return mapper.TodoListToEntity(&m), nil
}

func (g *TodoListQueriesGateway) List(
	ctx context.Context,
	opts *gatewayinput.ListTodoListsOptions,
) ([]*entity.TodoList, int64, error) {
	q := g.db.WithContext(ctx).Model(&model.TodoList{})

	if opts.OwnerID != nil {
		q = q.Where("owner_id = ?", int64(*opts.OwnerID))
	}
	if opts.NameSearch != nil && *opts.NameSearch != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?)", "%"+*opts.NameSearch+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, entity.NewInternal("failed to count todo lists")
	}

	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		q = q.Offset(opts.Offset)
	}

	var models []model.TodoList
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, entity.NewInternal("failed to list todo lists")
	}

	entities := make([]*entity.TodoList, len(models))
	for i := range models {
		entities[i] = mapper.TodoListToEntity(&models[i])
	}

	return entities, total, nil
}
