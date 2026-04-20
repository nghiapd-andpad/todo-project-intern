package persistence

import (
	"context"
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	appErrors "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/mapper"
	"gorm.io/gorm"
)

type userCommandsRepo struct {
	db *gorm.DB
}

func NewUserCommandsGateway(db *gorm.DB) gateway.UserCommandsGateway {
	return &userCommandsRepo{db: db}
}

// Create new user
func (r *userCommandsRepo) Create(ctx context.Context, e *entity.User) error {
	m := mapper.EntityToModel(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return appErrors.ParseDuplicateField(err)
		}
		return entity.ErrInternal
	}
	e.ID = entity.UserID(m.ID)
	return nil
}
