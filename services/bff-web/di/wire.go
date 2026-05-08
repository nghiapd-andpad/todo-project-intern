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
	authusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth"
	todousecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo"
	userusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/user"
)

type App struct {
	Resolver   *graph.Resolver
	JwtManager *jwt.JwtManager
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
		authusecase.NewRegisterer,
		wire.Bind(new(graph.AuthRegistererUsecase), new(*authusecase.Registerer)),

		authusecase.NewLoginer,
		wire.Bind(new(graph.AuthLoginerUsecase), new(*authusecase.Loginer)),

		todousecase.NewTodoCreator,
		wire.Bind(new(graph.TodoCreatorUsecase), new(*todousecase.TodoCreator)),

		todousecase.NewTodoGetter,
		wire.Bind(new(graph.TodoGetterUsecase), new(*todousecase.TodoGetter)),

		todousecase.NewTodoLister,
		wire.Bind(new(graph.TodoListerUsecase), new(*todousecase.TodoLister)),

		todousecase.NewTodoUpdater,
		wire.Bind(new(graph.TodoUpdaterUsecase), new(*todousecase.TodoUpdater)),

		todousecase.NewTodoDeleter,
		wire.Bind(new(graph.TodoDeleterUsecase), new(*todousecase.TodoDeleter)),

		userusecase.NewUserGetter,
		wire.Bind(new(graph.UserGetterUsecase), new(*userusecase.UserGetter)),

		// Handler
		graph.NewResolver,

		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
