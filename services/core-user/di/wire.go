//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler/grpc/user"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/security"
	user_usecase "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user"
	"gorm.io/gorm"
)

func InitializeUserHandler(db *gorm.DB, cfg *config.Config) (*user.UserHandler, error) {
	wire.Build(
		persistence.WireSet,
		security.WireSet,
		user_usecase.WireSet,
		user.NewUserHandler,
	)
	return nil, nil
}
