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

func TestIntegration_Todo_Get_Found(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// SETUP
	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
	todo := testutil.CreateTodo(t, env.db, list.ID, "Sub Work 1", entity.UserID(userID))

	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

	// ACTION
	resp, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: todoResourceName})

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "Sub Work 1", resp.Todo.Title)
	assert.Equal(t, todoResourceName, resp.Todo.Name)
}

func TestIntegration_Todo_Get_NotFound(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	_, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{
		Name: fmt.Sprintf("users/%d/todo-lists/1/todos/99999", userID),
	})

	// ASSERT
	assert.Equal(t, codes.NotFound, status.Code(err))
}
