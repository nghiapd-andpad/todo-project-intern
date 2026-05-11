//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
)

func TestIntegration_TodoList_Create(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := "123"
	ctx := env.authCtx(userID)

	resp, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      fmt.Sprintf("users/%s", userID),
		DisplayName: "Task 1",
	})

	require.NoError(t, err)
	require.NotNil(t, resp.TodoList)

	assert.Equal(t, "Task 1", resp.TodoList.DisplayName)
	assert.NotEmpty(t, resp.TodoList.Name)
	assert.NotNil(t, resp.TodoList.CreatedAt)
}

func TestIntegration_TodoList_Create_ThenGet_MatchesExactly(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := "456"
	ctx := env.authCtx(userID)

	// Action: Create
	created, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      fmt.Sprintf("users/%s", userID),
		DisplayName: "Shopping List",
	})
	require.NoError(t, err)

	// Action: Get
	fetched, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{
		Name: created.TodoList.Name,
	})
	require.NoError(t, err)

	// Assert
	diff := cmp.Diff(
		created.TodoList,
		fetched.TodoList,
		protocmp.Transform(),
		protocmp.IgnoreFields(&todov1.TodoList{}, "created_at", "updated_at"),
	)
	assert.Empty(t, diff, "created and fetched todo list should match")
}

func TestIntegration_TodoList_Create_Unauthenticated(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)

	_, err := env.client.CreateTodoList(
		env.authCtx(""),
		&todov1.CreateTodoListRequest{
			Parent:      "users/1",
			DisplayName: "Work Tasks",
		},
	)

	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestIntegration_TodoList_Create_InvalidRequest(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	ctx := env.authCtx("1")

	_, err := env.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      "users/1",
		DisplayName: "",
	})

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}
