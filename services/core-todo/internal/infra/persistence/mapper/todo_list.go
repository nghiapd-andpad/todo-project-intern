// Package mapper provides functions to convert between domain entities and persistence models.
package mapper

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

func TodoListToEntity(m *model.TodoList) *entity.TodoList {
	if m == nil {
		return nil
	}

	return &entity.TodoList{
		ID:        entity.TodoListID(m.ID),
		Name:      m.Name,
		OwnerID:   entity.UserID(m.OwnerID),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Version:   m.Version,
	}
}

func TodoListstoEntities(models []model.TodoList) []*entity.TodoList {
	result := make([]*entity.TodoList, len(models))
	for i := range models {
		result[i] = TodoListToEntity(&models[i])
	}
	return result
}

func TodoListFromEntity(e *entity.TodoList) *model.TodoList {
	if e == nil {
		return nil
	}

	return &model.TodoList{
		ID:        int64(e.ID),
		Name:      e.Name,
		OwnerID:   int64(e.OwnerID),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		Version:   e.Version,
	}
}
