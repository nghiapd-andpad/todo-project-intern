//go:build integration

// Package integration contains integration tests for the core-todo service.
// These tests run the full flow: gRPC client -> bufconn -> auth interceptor -> handler -> service -> persistence -> real DB.
package integration

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

const bufSize = 1024 * 1024

// testEnv holds a complete integration test environment.
type testEnv struct {
	client todov1.TodosServiceClient
	db     *gorm.DB
	// authCtx injects user ID via gRPC outgoing metadata.
	authCtx func(userID string) context.Context
}

// newTestEnv wires up a complete test environment.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	cfg := testutil.NewTestConfig(t)

	db := testutil.NewTestDB(t, cfg)

	// INFRASTRUCTURE
	todoCmds := persistence.NewTodoCommandsGateway(db)
	todoQueries := persistence.NewTodoQueriesGateway(db)
	todoListCmds := persistence.NewTodoListCommandsGateway(db)
	todoListQrs := persistence.NewTodoListQueriesGateway(db)

	// USECASE
	todoCreator := service.NewTodoCreator(todoCmds, cfg)
	todoGetter := service.NewTodoGetter(todoQueries)
	todoLister := service.NewTodoLister(todoQueries)
	todoUpdater := service.NewTodoUpdater(todoCmds, todoQueries)
	todoDeleter := service.NewTodoDeleter(todoCmds, todoQueries)

	todoListCreator := service.NewTodoListCreator(todoListCmds)
	todoListGetter := service.NewTodoListGetter(todoListQrs)
	todoListLister := service.NewTodoListLister(todoListQrs)
	todoListUpdater := service.NewTodoListUpdater(todoListCmds, todoListQrs)
	todoListDeleter := service.NewTodoListDeleter(todoListCmds)

	// HANDLER
	h := handler.NewTodoHandler(
		todoCreator,
		todoGetter,
		todoLister,
		todoUpdater,
		todoDeleter,
		todoListCreator,
		todoListGetter,
		todoListLister,
		todoListUpdater,
		todoListDeleter,
	)

	// Setup In-memory gRPC server
	srv, err := handler.NewGRPCServer(cfg, h)
	require.NoError(t, err)

	lis := bufconn.Listen(bufSize)

	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			fmt.Printf("gRPC server serve error: %v\n", err)
		}
	}()
	t.Cleanup(func() {
		srv.Stop()
		lis.Close()
	})

	// Setup client
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return &testEnv{
		client: todov1.NewTodosServiceClient(conn),
		db:     db,
		authCtx: func(userID string) context.Context {
			md := metadata.Pairs("x-user-id", userID)
			return metadata.NewOutgoingContext(context.Background(), md)
		},
	}
}
