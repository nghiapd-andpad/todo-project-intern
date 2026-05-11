//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestIntegration_Todo_Delete_Success(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// SETUP
	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
	todo := testutil.CreateTodo(t, env.db, list.ID, "To Delete", entity.UserID(userID))
	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

	// ACTION: DELETE
	_, err := env.client.DeleteTodo(ctx, &todov1.DeleteTodoRequest{Name: todoResourceName})

	// ASSERT
	require.NoError(t, err)
}

func TestIntegration_Todo_Delete_ThenGet_NotFound(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// SETUP
	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
	todo := testutil.CreateTodo(t, env.db, list.ID, "To Delete", entity.UserID(userID))
	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

	// ACTION: DELETE
	_, err := env.client.DeleteTodo(ctx, &todov1.DeleteTodoRequest{Name: todoResourceName})
	require.NoError(t, err)

	// ACTION: GET AGAIN
	_, err = env.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: todoResourceName})
	assert.Equal(t, codes.NotFound, status.Code(err))
}
