package persistence

import (
	"context"
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
	"gorm.io/gorm"
)

type userQueriesRepo struct {
	db *gorm.DB
}

func NewUserQueryGateway(db *gorm.DB) gateway.UserQueriesGateway {
	return &userQueriesRepo{db: db}
}

// Get user by username
func (r *userQueriesRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, entity.ErrInternal
	}
	return mapper.ModelToEntity(&m), nil
}

// Get user by email
func (r *userQueriesRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, entity.ErrInternal
	}
	return mapper.ModelToEntity(&m), nil
}

// Get user by ID
func (r *userQueriesRepo) GetByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, int64(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, entity.ErrInternal
	}
	return mapper.ModelToEntity(&m), nil
}
