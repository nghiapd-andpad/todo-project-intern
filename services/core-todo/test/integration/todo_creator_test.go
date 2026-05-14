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

// func TestIntegration_Todo_Create_RequiredFieldsOnly(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION
// 	resp, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
// 		Parent: listResourceName,
// 		Title:  "Sub Work 1",
// 	})

// 	// ASSERT
// 	require.NoError(t, err)
// 	require.NotNil(t, resp.Todo)
// 	assert.Equal(t, "Sub Work 1", resp.Todo.Title)
// 	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_PENDING, resp.Todo.Status)
// 	assert.NotEmpty(t, resp.Todo.Name)
// }

// func TestIntegration_Todo_Create_ThenGet_Matches(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION 1: CREATE
// 	resp, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
// 		Parent: listResourceName,
// 		Title:  "Write tests",
// 	})
// 	require.NoError(t, err)
// 	todoName := resp.Todo.Name

// 	// ACTION 2: GET
// 	fetched, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: todoName})

// 	// ASSERT
// 	require.NoError(t, err)
// 	assert.Equal(t, "Write tests", fetched.Todo.Title)
// 	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_PENDING, fetched.Todo.Status)
// }

// func TestIntegration_Todo_Create_Unauthenticated(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION: GET WITH NOT TOKEN
// 	_, err := env.client.CreateTodo(
// 		env.authCtx(""),
// 		&todov1.CreateTodoRequest{
// 			Parent: listResourceName,
// 			Title:  "Unauthorized Task",
// 		},
// 	)

// 	// ASSERT
// 	assert.Equal(t, codes.Unauthenticated, status.Code(err))
// }
