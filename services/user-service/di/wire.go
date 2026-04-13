//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/handler/grpc/user"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/infra/security"
	user_usecase "github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/usecase/user"
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
