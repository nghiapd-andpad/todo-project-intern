package persistence

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) gateway.UserCommandsGateway {
	return &userRepository{db: db}
}

func NewUserQueryRepository(db *gorm.DB) gateway.UserQueriesGateway {
	return &userRepository{db: db}
}

// Create new user
func (r *userRepository) Create(ctx context.Context, e *entity.User) error {
	m := mapper.EntityToModel(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	e.ID = entity.UserID(m.ID)
	return nil
}

// Get user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

// Get user by ID
func (r *userRepository) GetByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	var m model.User
	if err := r.db.WithContext(ctx).First(&m, int64(id)).Error; err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}
