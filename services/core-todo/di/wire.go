//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/cronjob"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/redis"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/worker"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/worker/job"
)

type ServerApp struct {
	GRPCServer *grpc.Server
	Logger     *zap.Logger
}

type WorkerApp struct {
	Worker *worker.Worker
	Logger *zap.Logger
}

func InitializeServer(cfg *config.Config) (*ServerApp, func(), error) {
	wire.Build(
		// OBSERVABILITY
		logutil.New,

		// INFRASTRUCTURE
		persistence.NewDatabase,

		persistence.NewTransactor,
		persistence.NewIdempotencyGateway,
		persistence.NewTodoCommandsGateway,
		persistence.NewTodoQueriesGateway,
		persistence.NewTodoListCommandsGateway,
		persistence.NewTodoListQueriesGateway,

		persistence.NewOutboxEventCommandsGateway,

		wire.Bind(new(gateway.Transactor), new(*persistence.Transactor)),
		wire.Bind(new(gateway.IdempotencyGateway), new(*persistence.IdempotencyGateway)),
		wire.Bind(new(gateway.TodoCommandsGateway), new(*persistence.TodoCommandsGateway)),
		wire.Bind(new(gateway.TodoQueriesGateway), new(*persistence.TodoQueriesGateway)),
		wire.Bind(new(gateway.TodoListCommandsGateway), new(*persistence.TodoListCommandsGateway)),
		wire.Bind(new(gateway.TodoListQueriesGateway), new(*persistence.TodoListQueriesGateway)),

		wire.Bind(new(gateway.OutboxEventCommandsGateway), new(*persistence.OutboxEventCommandsGateway)),

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

		// HANDLER
		handler.NewTodoHandler,

		// SERVER
		handler.NewGRPCServer,

		// App
		NewServerApp,
	)
	return nil, nil, nil
}

func InitializeWorker(cfg *config.Config) (*WorkerApp, func(), error) {
	wire.Build(
		logutil.New,

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

		redis.NewClient,
		redis.NewDistributedLocker,
		wire.Bind(new(gateway.DistributedLocker), new(*redis.DistributedLocker)),

		cronjob.NewScheduler,
		wire.Bind(new(gateway.Scheduler), new(*cronjob.Scheduler)),

		service.NewTodoOverdueMarker,
		wire.Bind(new(usecase.TodoOverdueMarker), new(*service.TodoOverdueMarker)),

		service.NewTodoSoftDeletedCleaner,
		wire.Bind(new(usecase.TodoSoftDeletedCleaner), new(*service.TodoSoftDeletedCleaner)),

		job.NewTodoOverdueMarkerJob,
		job.NewTodoSoftDeletedCleanupJob,

		worker.NewWorker,

		NewWorkerApp,
	)
	return nil, nil, nil
}

func NewWorkerApp(worker *worker.Worker, zapLogger *zap.Logger) *WorkerApp {
	return &WorkerApp{
		Worker: worker,
		Logger: zapLogger,
	}
}

func NewServerApp(grpcServer *grpc.Server, zapLogger *zap.Logger) *ServerApp {
	return &ServerApp{
		GRPCServer: grpcServer,
		Logger:     zapLogger,
	}
}
