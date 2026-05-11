package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
)

type UserQueriesGateway struct {
	db *gorm.DB
}

func NewUserQueryGateway(db *gorm.DB) *UserQueriesGateway {
	return &UserQueriesGateway{db: db}
}

func (r *UserQueriesGateway) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db get user by username: %w", err)
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *UserQueriesGateway) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db get user by email: %w", err)
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *UserQueriesGateway) GetByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).First(&m, int64(id)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db get user by id: %w", err)
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *UserQueriesGateway) GetByIDs(ctx context.Context, ids []entity.UserID) ([]*entity.User, error) {
	var models []model.User

	rawIDs := make([]int64, len(ids))
	for i, id := range ids {
		rawIDs[i] = int64(id)
	}

	err := r.db.WithContext(ctx).Where("id IN ?", rawIDs).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("db get users by ids: %w", err)
	}

	entities := make([]*entity.User, len(models))
	for i, m := range models {
		entities[i] = mapper.ModelToEntity(&m)
	}

	return entities, nil
}
