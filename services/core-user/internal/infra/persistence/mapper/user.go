package mapper

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
)

func ModelToEntity(m *model.User) *entity.User {
	if m == nil {
		return nil
	}
	return &entity.User{
		ID:             entity.UserID(m.ID),
		Username:       m.Username,
		Email:          m.Email,
		HashedPassword: m.HashedPassword,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func EntityToModel(e *entity.User) *model.User {
	if e == nil {
		return nil
	}
	return &model.User{
		ID:             int64(e.ID),
		Username:       e.Username,
		Email:          e.Email,
		HashedPassword: e.HashedPassword,
	}
}
