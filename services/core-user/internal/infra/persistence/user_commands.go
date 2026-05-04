package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/mapper"
)

type userCommandsRepo struct {
	db *gorm.DB
}

func NewUserCommandsGateway(db *gorm.DB) gateway.UserCommandsGateway {
	return &userCommandsRepo{db: db}
}

func (r *userCommandsRepo) Create(ctx context.Context, e *entity.User) (*entity.User, error) {
	m := mapper.EntityToModel(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, fmt.Errorf("db create user: %w", err)
	}
	return mapper.ModelToEntity(m), nil
}
