//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/cronjob"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/email"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/rabbitmq"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/security"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/utils/logger"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/worker"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/worker/job"
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
		logutil.New,

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

		NewServerApp,
	)
	return nil, nil, nil
}

func InitializeWorker(cfg *config.Config) (*WorkerApp, func(), error) {
	wire.Build(
		logutil.New,

		// DB + persistence
		persistence.NewDatabase,
		persistence.NewTransactor,
		persistence.NewOutboxEventCommandsGateway,
		persistence.NewOutboxEventQueriesGateway,
		persistence.NewNotificationCommandsGateway,
		persistence.NewUserQueryGateway,
		persistence.NewProcessedEventGateway,

		wire.Bind(new(gateway.Transactor), new(*persistence.Transactor)),
		wire.Bind(new(gateway.OutboxEventCommandsGateway), new(*persistence.OutboxEventCommandsGateway)),
		wire.Bind(new(gateway.OutboxEventQueriesGateway), new(*persistence.OutboxEventQueriesGateway)),
		wire.Bind(new(gateway.NotificationCommandsGateway), new(*persistence.NotificationCommandsGateway)),
		wire.Bind(new(gateway.UserQueriesGateway), new(*persistence.UserQueriesGateway)),
		wire.Bind(new(gateway.ProcessedEventGateway), new(*persistence.ProcessedEventGateway)),


		// Email Sender
		email.NewSMTPEmailSender,
		wire.Bind(new(gateway.EmailSender), new(*email.SMTPEmailSender)),

		// RabbitMQ
		rabbitmq.NewConnection,
		rabbitmq.NewPublisher,
		wire.Bind(new(gateway.EventPublisher), new(*rabbitmq.Publisher)),
		rabbitmq.NewNotificationConsumer,
		rabbitmq.NewEmailConsumer,

		// Scheduler
		cronjob.NewScheduler,
		wire.Bind(new(gateway.Scheduler), new(*cronjob.Scheduler)),

		// Cron jobs
		job.NewOutboxPublisherJob,

		worker.NewWorker,
		NewWorkerApp,
	)
	return nil, nil, nil
}

func NewWorkerApp(w *worker.Worker, zapLogger *zap.Logger) *WorkerApp {
	return &WorkerApp{Worker: w, Logger: zapLogger}
}

func NewServerApp(grpcServer *grpc.Server, zapLogger *zap.Logger) *ServerApp {
	return &ServerApp{GRPCServer: grpcServer, Logger: zapLogger}
}
