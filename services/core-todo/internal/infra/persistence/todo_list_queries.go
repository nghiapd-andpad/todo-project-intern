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
	if opts.OwnerID != nil && opts.AssigneeID != nil {
		return g.listAll(ctx, opts)
	}
	return g.listSimple(ctx, opts)
}

func (g *TodoListQueriesGateway) listAll(
	ctx context.Context,
	opts *gatewayinput.ListTodoListsOptions,
) ([]*entity.TodoList, int64, error) {
	var total int64
	countArgs := []interface{}{int64(*opts.OwnerID), int64(*opts.AssigneeID)}
	countQuery := `
		SELECT COUNT(DISTINCT todo_lists.id)
		FROM todo_lists
		LEFT JOIN todos ON todos.todo_list_id = todo_lists.id AND todos.deleted_at IS NULL
		WHERE todo_lists.deleted_at IS NULL
		  AND (todo_lists.owner_id = ? OR todos.assignee_id = ?)`

	if opts.NameSearch != nil && *opts.NameSearch != "" {
		countQuery += " AND LOWER(todo_lists.name) LIKE LOWER(?)"
		countArgs = append(countArgs, "%"+*opts.NameSearch+"%")
	}

	if err := g.db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("db count todo lists (all): %w", err)
	}

	q := g.db.WithContext(ctx).
		Distinct("todo_lists.*").
		Model(&model.TodoList{}).
		Joins("LEFT JOIN todos ON todos.todo_list_id = todo_lists.id AND todos.deleted_at IS NULL").
		Where("todo_lists.owner_id = ? OR todos.assignee_id = ?",
			int64(*opts.OwnerID), int64(*opts.AssigneeID))

	if opts.NameSearch != nil && *opts.NameSearch != "" {
		q = q.Where("LOWER(todo_lists.name) LIKE LOWER(?)", "%"+*opts.NameSearch+"%")
	}

	q = q.Order("todo_lists.created_at DESC").
		Limit(opts.Limit).
		Offset(opts.Offset)

	var models []model.TodoList
	if err := q.Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("db list todo lists (all): %w", err)
	}

	return toEntities(models), total, nil
}

func (g *TodoListQueriesGateway) listSimple(
	ctx context.Context,
	opts *gatewayinput.ListTodoListsOptions,
) ([]*entity.TodoList, int64, error) {
	base := g.db.WithContext(ctx).Model(&model.TodoList{})

	if opts.OwnerID != nil {
		base = base.Where("todo_lists.owner_id = ?", int64(*opts.OwnerID))
	}
	if opts.AssigneeID != nil {
		base = base.Joins("INNER JOIN todos ON todos.todo_list_id = todo_lists.id AND todos.deleted_at IS NULL").
			Where("todos.assignee_id = ?", int64(*opts.AssigneeID))
	}
	if opts.NameSearch != nil && *opts.NameSearch != "" {
		base = base.Where("LOWER(todo_lists.name) LIKE LOWER(?)", "%"+*opts.NameSearch+"%")
	}

	var total int64
	countQuery := base.Session(&gorm.Session{})
	if opts.AssigneeID != nil {
		countQuery = countQuery.Distinct("todo_lists.id")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("db count todo lists: %w", err)
	}

	listQuery := base
	if opts.AssigneeID != nil {
		listQuery = listQuery.Distinct("todo_lists.*")
	}

	var models []model.TodoList
	if err := listQuery.Order("todo_lists.created_at DESC").
		Limit(opts.Limit).
		Offset(opts.Offset).
		Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("db list todo lists: %w", err)
	}

	return toEntities(models), total, nil
}

func (g *TodoListQueriesGateway) FindSoftDeletedTodoListIDs(
	ctx context.Context,
	cutoff time.Time,
	limit int,
) ([]entity.TodoListID, error) {
	conn := connFromContext(ctx, g.db)

	var ids []int64
	if err := conn.Unscoped().
		Model(&model.TodoList{}).
		Where("deleted_at IS NOT NULL").
		Where("deleted_at <= ?", cutoff).
		Order("deleted_at ASC").
		Limit(limit).
		Pluck("id", &ids).Error; err != nil {
		return nil, fmt.Errorf("db find soft deleted todo list ids: %w", err)
	}

	result := make([]entity.TodoListID, len(ids))
	for i, id := range ids {
		result[i] = entity.TodoListID(id)
	}

	return result, nil
}

func toEntities(models []model.TodoList) []*entity.TodoList {
	result := make([]*entity.TodoList, len(models))
	for i := range models {
		result[i] = mapper.TodoListToEntity(&models[i])
	}
	return result
}
