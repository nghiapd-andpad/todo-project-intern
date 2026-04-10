package mapper

import (
	"github.com/nghiaphunng18/todos/internal/domain/entity"
	"github.com/nghiaphunng18/todos/internal/infra/persistence/model"
)

// convert from model to entity
func ToEntity(m *model.Todo) *entity.Todo {
	if m == nil {
		return nil
	}

	e := &entity.Todo{
		ID:          entity.TodoID(m.ID),
		Title:       m.Title,
		Description: m.Description,
		Status:      entity.TodoStatus(m.Status),
		Priority:    entity.Priority(m.Priority),
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	if m.TodoListID != nil {
		listID := entity.TodoListID(*m.TodoListID)
		e.TodoListID = &listID
	}

	if m.AssigneeID != nil {
		uID := entity.UserID(*m.AssigneeID)
		e.AssigneeID = &uID
	}

	return e
}

// convert from entity to model
func FromEntity(e *entity.Todo) *model.Todo {
	if e == nil {
		return nil
	}

	m := &model.Todo{
		ID:          int64(e.ID),
		Title:       e.Title,
		Description: e.Description,
		Status:      string(e.Status),
		Priority:    string(e.Priority),
		DueDate:     e.DueDate,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	if e.TodoListID != nil {
		id := int64(*e.TodoListID)
		m.TodoListID = &id
	}

	if e.AssigneeID != nil {
		id := int64(*e.AssigneeID)
		m.AssigneeID = &id
	}

	return m
}
