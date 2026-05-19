//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/cronjob"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	redisinfrastructure "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/redis"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

type App struct {
	GRPCServer *grpc.Server
	Scheduler  *cronjob.Scheduler
}

func InitializeApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(
		// OBSERVABILITY
		logutil.New,

		// INFRASTRUCTURE
		persistence.NewDatabase,

		persistence.NewTransactor,
		persistence.NewTodoCommandsGateway,
		persistence.NewTodoQueriesGateway,
		persistence.NewTodoListCommandsGateway,
		persistence.NewTodoListQueriesGateway,

		wire.Bind(new(gateway.Transactor), new(*persistence.Transactor)),
		wire.Bind(new(gateway.TodoCommandsGateway), new(*persistence.TodoCommandsGateway)),
		wire.Bind(new(gateway.TodoQueriesGateway), new(*persistence.TodoQueriesGateway)),
		wire.Bind(new(gateway.TodoListCommandsGateway), new(*persistence.TodoListCommandsGateway)),
		wire.Bind(new(gateway.TodoListQueriesGateway), new(*persistence.TodoListQueriesGateway)),

		// REDIS
		redisinfrastructure.NewClient,
		redisinfrastructure.NewDistributedLocker,
		wire.Bind(new(gateway.DistributedLocker), new(*redisinfrastructure.DistributedLocker)),

		// USE CASE
		service.NewTodoCreator,
		service.NewTodoGetter,
		service.NewTodoLister,
		service.NewTodoUpdater,
		service.NewTodoDeleter,

		service.NewTodoListCreator,
		service.NewTodoListGetter,
		service.NewTodoListLister,
		service.NewTodoListUpdater,
		service.NewTodoListDeleter,

		service.NewTodoOverdueMarker,

		wire.Bind(new(usecase.TodoCreator), new(*service.TodoCreator)),
		wire.Bind(new(usecase.TodoGetter), new(*service.TodoGetter)),
		wire.Bind(new(usecase.TodoLister), new(*service.TodoLister)),
		wire.Bind(new(usecase.TodoUpdater), new(*service.TodoUpdater)),
		wire.Bind(new(usecase.TodoDeleter), new(*service.TodoDeleter)),

		wire.Bind(new(usecase.TodoListCreator), new(*service.TodoListCreator)),
		wire.Bind(new(usecase.TodoListGetter), new(*service.TodoListGetter)),
		wire.Bind(new(usecase.TodoListLister), new(*service.TodoListLister)),
		wire.Bind(new(usecase.TodoListUpdater), new(*service.TodoListUpdater)),
		wire.Bind(new(usecase.TodoListDeleter), new(*service.TodoListDeleter)),

		wire.Bind(new(usecase.TodoOverdueMarker), new(*service.TodoOverdueMarker)),

		// CRONJOB
		cronjob.NewGoCronScheduler,
		cronjob.NewScheduler,

		// HANDLER
		handler.NewTodoHandler,

		// SERVER
		handler.NewGRPCServer,

		// App
		NewApp,
	)
	return nil, nil, nil
}

func NewApp(
	grpcServer *grpc.Server,
	scheduler *cronjob.Scheduler,
) *App {
	return &App{
		GRPCServer: grpcServer,
		Scheduler:  scheduler,
	}
}
