//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestIntegration_Todo_List_All(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// SETUP
	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
	titles := []string{"Task A", "Task B", "Task C"}
	for _, title := range titles {
		testutil.CreateTodo(t, env.db, list.ID, title, entity.UserID(userID))
	}

	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

	// ACTION
	resp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:   listResourceName,
		PageSize: 20,
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.Todos, 3)
	assert.Equal(t, int64(3), resp.Total)
}

func TestIntegration_Todo_List_FilterByStatus(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

	// SETUP
	testutil.CreateTodo(t, env.db, list.ID, "Pending 1", entity.UserID(userID))
	testutil.CreateTodo(t, env.db, list.ID, "Pending 2", entity.UserID(userID))

	doneTodoEntity := testutil.CreateTodo(t, env.db, list.ID, "Done Task", entity.UserID(userID))
	doneTodoName := fmt.Sprintf("%s/todos/%d", listResourceName, doneTodoEntity.ID)

	_, err := env.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
		Todo: &todov1.Todo{
			Name:   doneTodoName,
			Status: todov1.TodoStatus_TODO_STATUS_DONE,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
	})
	require.NoError(t, err)

	// ACTION: only GET task PENDING
	resp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:       listResourceName,
		PageSize:     20,
		StatusFilter: todov1.TodoStatus_TODO_STATUS_PENDING,
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.Todos, 2)
	assert.Equal(t, int64(2), resp.Total)
	for _, td := range resp.Todos {
		assert.Equal(t, todov1.TodoStatus_TODO_STATUS_PENDING, td.Status)
	}
}

func TestIntegration_Todo_List_FilterByTitleSearch(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	list := testutil.CreateTodoList(t, env.db, "Work", entity.UserID(userID))
	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

	// SETUP
	testutil.CreateTodo(t, env.db, list.ID, "Buy groceries", entity.UserID(userID))
	testutil.CreateTodo(t, env.db, list.ID, "Buy medicine", entity.UserID(userID))
	testutil.CreateTodo(t, env.db, list.ID, "Write tests", entity.UserID(userID))

	// ACTION: Search keyword "Buy"
	resp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:      listResourceName,
		PageSize:    20,
		TitleSearch: "Buy",
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.Todos, 2)
	assert.Equal(t, int64(2), resp.Total)
}

func TestIntegration_Todo_List_Empty(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(1103)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	list := testutil.CreateTodoList(t, env.db, "Empty List", entity.UserID(userID))
	listResourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

	// ACTION
	resp, err := env.client.ListTodos(ctx, &todov1.ListTodosRequest{
		Parent:   listResourceName,
		PageSize: 20,
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.Todos, 0)
	assert.Equal(t, int64(0), resp.Total)
}
