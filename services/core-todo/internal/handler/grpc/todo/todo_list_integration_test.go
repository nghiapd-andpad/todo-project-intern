//go:build integration

package todo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
)

func TestIntegration_TodoList_CreateAndGet(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	// Setup
	createResp, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      "users/1",
		DisplayName: "Work Tasks",
	})
	require.NoError(t, err)
	require.NotNil(t, createResp.TodoList)
	assert.Equal(t, "Work Tasks", createResp.TodoList.DisplayName)
	assert.NotEmpty(t, createResp.TodoList.Name)

	// Get and verify persisted correctly
	getResp, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{
		Name: createResp.TodoList.Name,
	})
	require.NoError(t, err)

	// Use protocmp for protobuf message comparison
	diff := cmp.Diff(createResp.TodoList, getResp.TodoList,
		protocmp.Transform(),
		protocmp.IgnoreFields(&todov1.TodoList{}, "created_at", "updated_at"),
	)
	assert.Empty(t, diff, "created and fetched todo list should match")
}

func TestIntegration_TodoList_UpdateDisplayName(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	// Setup: create
	created, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      "users/1",
		DisplayName: "Old Name",
	})
	require.NoError(t, err)

	// Update
	updated, err := env.client.UpdateTodoList(ctx, &todov1.UpdateTodoListRequest{
		TodoList: &todov1.TodoList{
			Name:        created.TodoList.Name,
			DisplayName: "New Name",
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
	})
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.TodoList.DisplayName)

	// Verify persisted
	fetched, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{
		Name: created.TodoList.Name,
	})
	require.NoError(t, err)
	assert.Equal(t, "New Name", fetched.TodoList.DisplayName)
}

func TestIntegration_TodoList_DeleteAndVerifyGone(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	// Create
	created, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      "users/1",
		DisplayName: "To Delete",
	})
	require.NoError(t, err)

	// Delete
	_, err = env.client.DeleteTodoList(ctx, &todov1.DeleteTodoListRequest{
		Name: created.TodoList.Name,
	})
	require.NoError(t, err)

	// Verify
	_, err = env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{
		Name: created.TodoList.Name,
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestIntegration_TodoList_List(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	// Create 3 todo lists
	names := []string{"Work", "Personal", "Shopping"}
	for _, name := range names {
		_, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
			Parent:      "users/1",
			DisplayName: name,
		})
		require.NoError(t, err)
	}

	resp, err := env.client.ListTodoLists(ctx, &todov1.ListTodoListsRequest{
		Parent:   "users/1",
		PageSize: 20,
	})
	require.NoError(t, err)
	assert.Len(t, resp.TodoLists, 3)
	assert.Equal(t, int64(3), resp.Total)
}

func TestIntegration_TodoList_List_WithNameSearch(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	// Create todo lists
	for _, name := range []string{"Work Tasks", "Work Projects", "Personal"} {
		_, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
			Parent:      "users/1",
			DisplayName: name,
		})
		require.NoError(t, err)
	}

	// Search for "Work"
	resp, err := env.client.ListTodoLists(ctx, &todov1.ListTodoListsRequest{
		Parent:     "users/1",
		PageSize:   20,
		NameSearch: "Work",
	})
	require.NoError(t, err)
	assert.Len(t, resp.TodoLists, 2)
	assert.Equal(t, int64(2), resp.Total)
}

func TestIntegration_TodoList_Get_NotFound(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	_, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{
		Name: "users/1/todo-lists/99999",
	})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestIntegration_TodoList_Unauthenticated(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)

	_, err := env.client.CreateTodoList(
		env.authCtx(""), // empty user ID
		&todov1.CreateTodoListRequest{
			Parent:      "users/1",
			DisplayName: "Work Tasks",
		},
	)

	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}
