//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestIntegration_TodoList_List_All(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(600)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// SETUP
	names := []string{"Work", "Personal", "Shopping"}
	for _, name := range names {
		testutil.CreateTodoList(t, env.db, name, entity.UserID(userID))
	}

	// ACTION
	resp, err := env.client.ListTodoLists(ctx, &todov1.ListTodoListsRequest{
		Parent:   fmt.Sprintf("users/%d", userID),
		PageSize: 20,
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.TodoLists, 3)
	assert.Equal(t, int64(3), resp.Total)
}

func TestIntegration_TodoList_List_WithNameSearch(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(601)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// SETUP
	testutil.CreateTodoList(t, env.db, "Work Tasks", entity.UserID(userID))
	testutil.CreateTodoList(t, env.db, "Work Projects", entity.UserID(userID))
	testutil.CreateTodoList(t, env.db, "Personal", entity.UserID(userID))

	// ACTION: Search with keyword "Work"
	resp, err := env.client.ListTodoLists(ctx, &todov1.ListTodoListsRequest{
		Parent:     fmt.Sprintf("users/%d", userID),
		PageSize:   20,
		NameSearch: "Work",
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.TodoLists, 2)
	assert.Equal(t, int64(2), resp.Total)
}

func TestIntegration_TodoList_List_IsolatedByOwner(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	user1 := int64(602)
	user2 := int64(603)

	// SETUP
	testutil.CreateTodoList(t, env.db, "User1 List A", entity.UserID(user1))
	testutil.CreateTodoList(t, env.db, "User1 List B", entity.UserID(user1))
	testutil.CreateTodoList(t, env.db, "User2 List", entity.UserID(user2))

	// ACTION: User 1 GET List
	resp, err := env.client.ListTodoLists(env.authCtx(fmt.Sprintf("%d", user1)), &todov1.ListTodoListsRequest{
		Parent:   fmt.Sprintf("users/%d", user1),
		PageSize: 20,
	})

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, resp.TodoLists, 2)
	assert.Equal(t, int64(2), resp.Total)
	for _, list := range resp.TodoLists {
		assert.Contains(t, list.DisplayName, "User1")
	}
}

func TestIntegration_TodoList_List_Empty(t *testing.T) {
	t.Parallel()

	env := newTestEnv(t)
	userID := int64(604)
	ctx := env.authCtx(fmt.Sprintf("%d", userID))

	// ACTION
	resp, err := env.client.ListTodoLists(ctx, &todov1.ListTodoListsRequest{
		Parent:   fmt.Sprintf("users/%d", userID),
		PageSize: 20,
	})

	// ASSERT: return empty, not nil
	require.NoError(t, err)
	assert.Len(t, resp.TodoLists, 0)
	assert.Equal(t, int64(0), resp.Total)
}
