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

// func TestIntegration_TodoList_Get_Found(t *testing.T) {
// 	t.Parallel()
// 	env := newTestEnv(t)
// 	userID := int64(100)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "My List", entity.UserID(userID))

// 	// ACTION
// 	resourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)
// 	resp, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{Name: resourceName})

// 	// ASSERT
// 	require.NoError(t, err)
// 	assert.Equal(t, "My List", resp.TodoList.DisplayName)
// }

// func TestIntegration_TodoList_Get_NotFound(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	ctx := env.authCtx("1")

// 	_, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{
// 		Name: "users/1/todo-lists/99999",
// 	})

// 	assert.Equal(t, codes.NotFound, status.Code(err))
// }
