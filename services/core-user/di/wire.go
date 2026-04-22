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

func InitializeApp(cfg *config.Config) (*grpc.Server, func(), error) {
	wire.Build(
		// Infrastructure
		persistence.NewDatabase,
		persistence.WireSet,
		security.WireSet,

		// Usecases
		userUsecase.WireSet,

		// Handler and gRPC server
		userHandler.NewUserHandler,
		userHandler.NewGRPCServer,
	)
	return nil, nil, nil
}
