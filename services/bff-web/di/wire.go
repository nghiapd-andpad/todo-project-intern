//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/graph"
	grpcclient "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpcclient"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/jwt"
	authservice "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/service/auth"
	todoservice "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/service/todo"
	userservice "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/service/user"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase"
)

type App struct {
	Resolver   *graph.Resolver
	JwtManager *jwt.JWTManager
}

func InitializeApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(
		// Infra
		grpcclient.NewAuthGateway,
		wire.Bind(new(gateway.AuthGateway), new(*grpcclient.AuthGateway)),

		grpcclient.NewTodoGateway,
		wire.Bind(new(gateway.TodoGateway), new(*grpcclient.TodoGateway)),

		grpcclient.NewUserGateway,
		wire.Bind(new(gateway.UserGateway), new(*grpcclient.UserGateway)),

		jwt.NewJwtManager,

		// ── Usecase
		authservice.NewRegisterer,
		wire.Bind(new(usecase.AuthRegisterer), new(*authservice.Registerer)),

		authservice.NewLoginer,
		wire.Bind(new(usecase.AuthLoginer), new(*authservice.Loginer)),

		todoservice.NewTodoCreator,
		wire.Bind(new(usecase.TodoCreator), new(*todoservice.TodoCreator)),

		todoservice.NewTodoGetter,
		wire.Bind(new(usecase.TodoGetter), new(*todoservice.TodoGetter)),

		todoservice.NewTodoLister,
		wire.Bind(new(usecase.TodoLister), new(*todoservice.TodoLister)),

		todoservice.NewTodoUpdater,
		wire.Bind(new(usecase.TodoUpdater), new(*todoservice.TodoUpdater)),

		todoservice.NewTodoDeleter,
		wire.Bind(new(usecase.TodoDeleter), new(*todoservice.TodoDeleter)),

		userservice.NewUserGetter,
		wire.Bind(new(usecase.UserGetter), new(*userservice.UserGetter)),

		// Handler
		graph.NewResolver,

		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
