//go:build integration

package integration

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"

// 	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
// )

// func TestIntegration_TodoList_Delete_Success(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "To Delete", entity.UserID(userID))
// 	resourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION: DELETE
// 	_, err := env.client.DeleteTodoList(ctx, &todov1.DeleteTodoListRequest{Name: resourceName})

// 	// ASSERT
// 	require.NoError(t, err)
// }

// func TestIntegration_TodoList_Delete_ThenGet_NotFound(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "To Delete", entity.UserID(userID))
// 	resourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION: DELETE
// 	_, err := env.client.DeleteTodoList(ctx, &todov1.DeleteTodoListRequest{Name: resourceName})
// 	require.NoError(t, err)

// 	// ACTION: GET AGAIN
// 	_, err = env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{Name: resourceName})
// 	assert.Equal(t, codes.NotFound, status.Code(err))
// }
