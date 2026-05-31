package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (g *TodoQueriesGateway) Get(ctx context.Context, todoID entity.TodoID, todoListID entity.TodoListID) (*entity.Todo, error) {
	conn := connFromContext(ctx, g.db)

	var m model.Todo

	err := conn.Where("id = ? AND todo_list_id = ?", int64(todoID), int64(todoListID)).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db get todo: %w", err)
	}

	return mapper.TodoToEntity(&m), nil
}

func (g *TodoQueriesGateway) List(ctx context.Context, opts *gatewayinput.ListTodosOptions) ([]*entity.Todo, int64, error) {
	conn := connFromContext(ctx, g.db)

	q := conn.Model(&model.Todo{}).
		Where("todo_list_id = ?", int64(opts.TodoListID))

	if opts.AssigneeOnly != nil {
		q = q.Where("assignee_id = ?", int64(*opts.AssigneeOnly))
	}
	if opts.Status != nil {
		q = q.Where("status = ?", string(*opts.Status))
	}
	if opts.Priority != nil {
		q = q.Where("priority = ?", string(*opts.Priority))
	}
	if opts.TitleSearch != nil && *opts.TitleSearch != "" {
		q = q.Where("LOWER(title) LIKE LOWER(?)", "%"+*opts.TitleSearch+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("db count todos: %w", err)
	}

	if opts.CursorCreatedAt != nil && opts.CursorID != nil {
		q = q.Where(
			"(created_at < ? OR (created_at = ? AND id < ?))",
			*opts.CursorCreatedAt,
			*opts.CursorCreatedAt,
			int64(*opts.CursorID),
		)
	}

	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}

	if opts.Offset > 0 && opts.CursorCreatedAt == nil && opts.CursorID == nil {
		q = q.Offset(opts.Offset)
	}

	var models []model.Todo
	if err := q.
		Order("created_at DESC").
		Order("id DESC").
		Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("db list todos: %w", err)
	}

	entities := make([]*entity.Todo, len(models))
	for i := range models {
		entities[i] = mapper.TodoToEntity(&models[i])
	}

	return entities, total, nil
}

func (g *TodoQueriesGateway) FindOverdueTodoIDs(ctx context.Context, asOf time.Time, limit int) ([]entity.TodoID, error) {
	conn := connFromContext(ctx, g.db)

	if limit <= 0 {
		return []entity.TodoID{}, nil
	}

	var ids []int64
	if err := conn.Model(&model.Todo{}).
		Where("deleted_at IS NULL").
		Where("due_date IS NOT NULL").
		Where("due_date <= ?", asOf).
		Where("status IN ?", []string{
			string(entity.TodoStatusPending),
			string(entity.TodoStatusInProgress),
		}).
		Order("due_date ASC").
		Limit(limit).
		Pluck("id", &ids).Error; err != nil {
		return nil, fmt.Errorf("db find overdue todo ids: %w", err)
	}

	result := make([]entity.TodoID, len(ids))
	for i, id := range ids {
		result[i] = entity.TodoID(id)
	}

	return result, nil
}

func (g *TodoQueriesGateway) FindSoftDeletedTodoIDs(ctx context.Context, cutoff time.Time, limit int) ([]entity.TodoID, error) {
	conn := connFromContext(ctx, g.db)

	var ids []int64
	if err := conn.Unscoped().
		Model(&model.Todo{}).
		Where("deleted_at IS NOT NULL").
		Where("deleted_at <= ?", cutoff).
		Order("deleted_at ASC").
		Limit(limit).
		Pluck("id", &ids).Error; err != nil {
		return nil, fmt.Errorf("db find soft deleted todo ids: %w", err)
	}

	result := make([]entity.TodoID, len(ids))
	for i, id := range ids {
		result[i] = entity.TodoID(id)
	}

	return result, nil
}
