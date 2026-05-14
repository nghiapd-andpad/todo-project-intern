//go:build integration

package integration

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"google.golang.org/protobuf/types/known/fieldmaskpb"

// 	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
// )

// func TestIntegration_Todo_Update_Title(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	todo := testutil.CreateTodo(t, env.db, list.ID, "Old Title", entity.UserID(userID))
// 	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

// 	// ACTION: only update Title
// 	updated, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
// 		Todo:       &todov1.Todo{Name: todoResourceName, Title: "New Title"},
// 		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"title"}},
// 	})

// 	// ASSERT
// 	require.NoError(t, err)
// 	assert.Equal(t, "New Title", updated.Todo.Title)
// 	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_PENDING, updated.Todo.Status)
// }

// func TestIntegration_Todo_Update_Title_Persisted(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	todo := testutil.CreateTodo(t, env.db, list.ID, "Old Title", entity.UserID(userID))
// 	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

// 	// ACTION: UPDATE
// 	_, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
// 		Todo:       &todov1.Todo{Name: todoResourceName, Title: "New Title"},
// 		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"title"}},
// 	})
// 	require.NoError(t, err)

// 	// ASSERT
// 	fetched, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: todoResourceName})
// 	require.NoError(t, err)
// 	assert.Equal(t, "New Title", fetched.Todo.Title)
// }

// func TestIntegration_Todo_Update_Status_ToInProgress(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	todo := testutil.CreateTodo(t, env.db, list.ID, "Task Original", entity.UserID(userID))
// 	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

// 	// ACTION: UPDATE status
// 	updated, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
// 		Todo:       &todov1.Todo{Name: todoResourceName, Status: todov1.TodoStatus_TODO_STATUS_IN_PROGRESS},
// 		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
// 	})

// 	// ASSERT
// 	require.NoError(t, err)
// 	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_IN_PROGRESS, updated.Todo.Status)
// 	assert.Equal(t, "Task Original", updated.Todo.Title)
// }

// func TestIntegration_Todo_Update_Status_ToDone(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
// 	todo := testutil.CreateTodo(t, env.db, list.ID, "Task", entity.UserID(userID))
// 	todoResourceName := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", userID, list.ID, todo.ID)

// 	// ACTION: Update to DONE
// 	updated, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
// 		Todo:       &todov1.Todo{Name: todoResourceName, Status: todov1.TodoStatus_TODO_STATUS_DONE},
// 		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
// 	})

// 	require.NoError(t, err)
// 	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_DONE, updated.Todo.Status)
// }
