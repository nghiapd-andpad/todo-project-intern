//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	userHandler "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler/grpc/user"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/security"
	userUsecase "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user"
	"google.golang.org/grpc"
)

func InitializeApp() (*grpc.Server, func(), error) {
	wire.Build(
		config.New,
		persistence.NewDatabase,
		persistence.WireSet,
		security.WireSet,
		userUsecase.WireSet,        // Dùng alias mới
		userHandler.NewUserHandler, // Dùng alias mới
		userHandler.NewGRPCServer,  // Dùng alias mới
	)
	return nil, nil, nil
}
