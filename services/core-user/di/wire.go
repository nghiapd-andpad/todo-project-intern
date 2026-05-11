//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/security"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase"
)

func InitializeApp(cfg *config.Config) (*grpc.Server, func(), error) {
	wire.Build(
		// INFRASTRUCTURE
		persistence.NewDatabase,

		persistence.NewUserCommandsGateway,
		persistence.NewUserQueryGateway,
		wire.Bind(new(gateway.UserQueriesGateway), new(*persistence.UserQueriesGateway)),
		wire.Bind(new(gateway.UserCommandsGateway), new(*persistence.UserCommandsGateway)),

		security.NewJWTManager,
		wire.Bind(new(gateway.TokenManager), new(*security.JWTManager)),

		// USE CASE
		service.NewUserAuthenticator,
		service.NewUserCreator,
		service.NewUserGetter,
		wire.Bind(new(usecase.UserAuthenticator), new(*service.UserAuthenticator)),
		wire.Bind(new(usecase.UserCreator), new(*service.UserCreator)),
		wire.Bind(new(usecase.UserGetter), new(*service.UserGetter)),

		// HANDLER
		handler.NewUserHandler,

		// SERVER
		handler.NewGRPCServer,
	)
	return nil, nil, nil
}
