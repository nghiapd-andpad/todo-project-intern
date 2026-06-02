package persistence

import (
	"context"
	"fmt"
	"time"

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

	updates := map[string]any{
		"title":       todo.Title,
		"description": todo.Description,
		"status":      string(todo.Status),
		"priority":    string(todo.Priority),
		"due_date":    todo.DueDate,
		"assignee_id": todo.AssigneeID,
	}

	result := conn.
		Model(&model.Todo{}).
		Where("id = ?", int64(todo.ID)).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("db update todo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, entity.NewConflict(
			"todo was updated by another request",
		)
	}

	return todo, nil
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

func (g *TodoCommandsGateway) MarkOverdueByIDs(ctx context.Context, ids []entity.TodoID, markedAt time.Time) (int64, error) {
	conn := connFromContext(ctx, g.db)

	if len(ids) == 0 {
		return 0, nil
	}

	rawIDs := make([]int64, len(ids))
	for i, id := range ids {
		rawIDs[i] = int64(id)
	}

	result := conn.
		Model(&model.Todo{}).
		Where("deleted_at IS NULL").
		Where("id IN ?", rawIDs).
		Where("status IN ?", []string{
			string(entity.TodoStatusPending),
			string(entity.TodoStatusInProgress),
		}).
		Updates(map[string]any{
			"status":     string(entity.TodoStatusOverdue),
			"updated_at": markedAt,
		})

	if result.Error != nil {
		return 0, fmt.Errorf("db mark overdue todos by ids: %w", result.Error)
	}

	return result.RowsAffected, nil
}

func (g *TodoCommandsGateway) HardDeleteTodosByIDs(ctx context.Context, ids []entity.TodoID) (int64, error) {
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
		Delete(&model.Todo{})

	if result.Error != nil {
		return 0, fmt.Errorf("db hard delete todos by ids: %w", result.Error)
	}

	return result.RowsAffected, nil
}

func (g *TodoCommandsGateway) HardDeleteTodosByTodoListIDs(ctx context.Context, todoListIDs []entity.TodoListID) (int64, error) {
	if len(todoListIDs) == 0 {
		return 0, nil
	}

	conn := connFromContext(ctx, g.db)

	values := make([]int64, len(todoListIDs))
	for i, id := range todoListIDs {
		values[i] = int64(id)
	}

	result := conn.Unscoped().
		Where("todo_list_id IN ?", values).
		Delete(&model.Todo{})

	if result.Error != nil {
		return 0, fmt.Errorf("db hard delete todos by todo list ids: %w", result.Error)
	}

	return result.RowsAffected, nil
}
