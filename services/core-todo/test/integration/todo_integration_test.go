//go:build integration

package todo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
)

// createListHelper creates a todo list and returns its name.
func createListHelper(t *testing.T, env *testEnv, userID, displayName string) string {
	t.Helper()
	resp, err := env.client.CreateTodoList(env.authCtx(userID), &todov1.CreateTodoListRequest{
		Parent:      "users/" + userID,
		DisplayName: displayName,
	})
	require.NoError(t, err)
	return resp.TodoList.Name
}

func TestIntegration_Todo_CreateAndGet(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	// Create todo list
	listName := createListHelper(t, env, "1", "Work Tasks")

	// Create todo
	createResp, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
		Parent: listName,
		Title:  "Buy groceries",
	})
	require.NoError(t, err)
	require.NotNil(t, createResp.Todo)
	assert.Equal(t, "Buy groceries", createResp.Todo.Title)
	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_PENDING, createResp.Todo.Status) // default
	assert.NotEmpty(t, createResp.Todo.Name)

	// Get and verify persisted
	getResp, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{
		Name: createResp.Todo.Name,
	})
	require.NoError(t, err)
	assert.Equal(t, "Buy groceries", getResp.Todo.Title)
}

func TestIntegration_Todo_UpdateTitle(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	listName := createListHelper(t, env, "1", "Work")

	created, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
		Parent: listName,
		Title:  "Old Title",
	})
	require.NoError(t, err)

	// Update title
	updated, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
		Todo: &todov1.Todo{
			Name:  created.Todo.Name,
			Title: "New Title",
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"title"}},
	})
	require.NoError(t, err)
	assert.Equal(t, "New Title", updated.Todo.Title)
	// Status unchanged
	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_PENDING, updated.Todo.Status)

	// Verify persisted
	fetched, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: created.Todo.Name})
	require.NoError(t, err)
	assert.Equal(t, "New Title", fetched.Todo.Title)
}

func TestIntegration_Todo_UpdateStatus(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	listName := createListHelper(t, env, "1", "Work")

	created, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
		Parent: listName,
		Title:  "Task",
	})
	require.NoError(t, err)

	// Update status to DONE
	updated, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
		Todo: &todov1.Todo{
			Name:   created.Todo.Name,
			Status: todov1.TodoStatus_TODO_STATUS_DONE,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
	})
	require.NoError(t, err)
	assert.Equal(t, todov1.TodoStatus_TODO_STATUS_DONE, updated.Todo.Status)
	assert.Equal(t, "Task", updated.Todo.Title) // unchanged
}

func TestIntegration_Todo_DeleteAndVerifyGone(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	listName := createListHelper(t, env, "1", "Work")

	// Create
	created, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
		Parent: listName,
		Title:  "To Delete",
	})
	require.NoError(t, err)

	// Delete
	_, err = env.client.DeleteTodo(ctx, &todov1.DeleteTodoRequest{
		Name: created.Todo.Name,
	})
	require.NoError(t, err)

	// Verify gone
	_, err = env.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: created.Todo.Name})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestIntegration_Todo_ListWithFilters(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	listName := createListHelper(t, env, "1", "Work")

	// Create: 2 PENDING + 1 DONE
	for _, title := range []string{"Pending 1", "Pending 2"} {
		_, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
			Parent: listName,
			Title:  title,
		})
		require.NoError(t, err)
	}

	done, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
		Parent: listName,
		Title:  "Done Task",
	})
	require.NoError(t, err)

	// Update one to DONE
	_, err = env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
		Todo: &todov1.Todo{
			Name:   done.Todo.Name,
			Status: todov1.TodoStatus_TODO_STATUS_DONE,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
	})
	require.NoError(t, err)

	allResp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:   listName,
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.Len(t, allResp.Todos, 3)
	assert.Equal(t, int64(3), allResp.Total)

	pendingResp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:       listName,
		PageSize:     20,
		StatusFilter: todov1.TodoStatus_TODO_STATUS_PENDING,
	})
	require.NoError(t, err)
	assert.Len(t, pendingResp.Todos, 2)
	assert.Equal(t, int64(2), pendingResp.Total)
}

func TestIntegration_Todo_ListWithTitleSearch(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	listName := createListHelper(t, env, "1", "Work")

	for _, title := range []string{"Buy groceries", "Buy medicine", "Write tests"} {
		_, err := env.client.CreateTodo(ctx, &todov1.CreateTodoRequest{
			Parent: listName,
			Title:  title,
		})
		require.NoError(t, err)
	}

	resp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:      listName,
		PageSize:    20,
		TitleSearch: "Buy",
	})
	require.NoError(t, err)
	assert.Len(t, resp.Todos, 2)
	assert.Equal(t, int64(2), resp.Total)
}

func TestIntegration_Todo_Get_NotFound(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	_, err := env.client.GetTodo(ctx, &todov1.GetTodoRequest{
		Name: "users/1/todo-lists/1/todos/99999",
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
}
