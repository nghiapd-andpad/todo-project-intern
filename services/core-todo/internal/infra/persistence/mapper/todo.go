package mapper

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

func TodoToEntity(m *model.Todo) *entity.Todo {
	if m == nil {
		return nil
	}

	e := &entity.Todo{
		ID:          entity.TodoID(m.ID),
		TodoListID:  entity.TodoListID(m.TodoListID),
		Title:       m.Title,
		Description: m.Description,
		Status:      entity.TodoStatus(m.Status),
		Priority:    entity.Priority(m.Priority),
		CreatorID:   entity.UserID(m.CreatorID),
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	if m.AssigneeID != nil {
		uID := entity.UserID(*m.AssigneeID)
		e.AssigneeID = &uID
	}

	return e
}

func TodoFromEntity(e *entity.Todo) *model.Todo {
	if e == nil {
		return nil
	}

	m := &model.Todo{
		ID:          int64(e.ID),
		TodoListID:  int64(e.TodoListID),
		Title:       e.Title,
		Description: e.Description,
		Status:      string(e.Status),
		Priority:    string(e.Priority),
		CreatorID:   int64(e.CreatorID),
		DueDate:     e.DueDate,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	if e.AssigneeID != nil {
		id := int64(*e.AssigneeID)
		m.AssigneeID = &id
	}

	return m
}
