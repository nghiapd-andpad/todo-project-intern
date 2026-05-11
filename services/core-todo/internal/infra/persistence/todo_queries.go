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

type TodoQueriesGateway struct {
	db *gorm.DB
}

func NewTodoQueriesGateway(db *gorm.DB) *TodoQueriesGateway {
	return &TodoQueriesGateway{db: db}
}

var _ gateway.TodoQueriesGateway = (*TodoQueriesGateway)(nil)

func (g *TodoQueriesGateway) Get(ctx context.Context, todoID entity.TodoID) (*entity.Todo, error) {
	var m model.Todo

	err := g.db.WithContext(ctx).First(&m, int64(todoID)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db get todo: %w", err)
	}

	return mapper.TodoToEntity(&m), nil
}

func (g *TodoQueriesGateway) List(ctx context.Context, opts *gatewayinput.ListTodosOptions) ([]*entity.Todo, int64, error) {
	// Build base query
	q := g.db.WithContext(ctx).Model(&model.Todo{})

	// Apply optional filters
	if opts.TodoListID != nil {
		q = q.Where("todo_list_id = ?", int64(*opts.TodoListID))
	}
	if opts.Status != nil {
		q = q.Where("status = ?", string(*opts.Status))
	}
	if opts.Priority != nil {
		q = q.Where("priority = ?", string(*opts.Priority))
	}
	if opts.CreatorID != nil {
		q = q.Where("creator_id = ?", int64(*opts.CreatorID))
	}
	if opts.AssigneeID != nil {
		q = q.Where("assignee_id = ?", int64(*opts.AssigneeID))
	}
	if opts.TitleSearch != nil && *opts.TitleSearch != "" {
		q = q.Where("LOWER(title) LIKE LOWER(?)", "%"+*opts.TitleSearch+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("db count todos: %w", err)
	}

	// Apply pagination
	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		q = q.Offset(opts.Offset)
	}

	var models []model.Todo
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("db list todos: %w", err)
	}

	entities := make([]*entity.Todo, len(models))
	for i := range models {
		entities[i] = mapper.TodoToEntity(&models[i])
	}

	return entities, total, nil
}
