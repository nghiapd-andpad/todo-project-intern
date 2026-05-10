//go:build integration

package todo_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	todoHandler "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/testutil"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
)

const bufSize = 1024 * 1024

// testEnv holds all dependencies for integration tests.
type testEnv struct {
	client todov1.TodosServiceClient
	// authCtx returns context with user ID injected via gRPC metadata
	authCtx func(userID string) context.Context
}

// newTestEnv creates a full integration test environment: real DB -> real usecases -> real handler -> bufconn gRPC server -> client.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	// Load config per test
	cfg, err := config.New()
	require.NoError(t, err)

	// Clone DB — each test gets isolated database
	db := testutil.NewTestDB(t, cfg)

	// Infra
	todoCmds := persistence.NewTodoCommandsGateway(db)
	todoQueries := persistence.NewTodoQueriesGateway(db)
	todoListCmds := persistence.NewTodoListCommandsGateway(db)
	todoListQrs := persistence.NewTodoListQueriesGateway(db)

	// Usecases
	handler := todoHandler.NewTodoHandler(
		todos.NewTodoCreator(todoCmds),
		todos.NewTodoGetter(todoQueries),
		todos.NewTodoLister(todoQueries),
		todos.NewTodoUpdater(todoCmds, todoQueries),
		todos.NewTodoDeleter(todoCmds, todoQueries),
		todos.NewTodoListCreator(todoListCmds),
		todos.NewTodoListGetter(todoListQrs),
		todos.NewTodoListLister(todoListQrs),
		todos.NewTodoListUpdater(todoListCmds, todoListQrs),
		todos.NewTodoListDeleter(todoListCmds),
	)

	// bufconn gRPC server — in-memory
	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)
	todov1.RegisterTodosServiceServer(srv, handler)

	go func() {
		if err := srv.Serve(lis); err != nil && err.Error() != "closed" {
			t.Logf("bufconn server error: %v", err)
		}
	}()

	t.Cleanup(func() {
		srv.Stop()
		lis.Close()
	})

	// gRPC client
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
		authCtx: func(userID string) context.Context {
			md := metadata.Pairs("x-user-id", userID)
			return metadata.NewOutgoingContext(context.Background(), md)
		},
	}
}
