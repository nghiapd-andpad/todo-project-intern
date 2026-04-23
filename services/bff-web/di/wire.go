//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/graph"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpc_client"
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
		grpc_client.WireSet,
		jwt.WireSet,
		authusecase.WireSet,
		todousecase.WireSet,
		userusecase.WireSet,
		graph.NewResolver,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
