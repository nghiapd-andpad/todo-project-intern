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

// func TestIntegration_TodoList_Update_DisplayName(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(1)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP
// 	list := testutil.CreateTodoList(t, env.db, "Old Name", entity.UserID(userID))
// 	resourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION: Update
// 	updated, err := env.client.UpdateTodoList(ctx, &todov1.UpdateTodoListRequest{
// 		TodoList: &todov1.TodoList{
// 			Name:        resourceName,
// 			DisplayName: "New Name",
// 		},
// 		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
// 	})

// 	// ASSERT
// 	require.NoError(t, err)
// 	assert.Equal(t, "New Name", updated.TodoList.DisplayName)
// }

// func TestIntegration_TodoList_Update_Persisted(t *testing.T) {
// 	t.Parallel()

// 	env := newTestEnv(t)
// 	userID := int64(2)
// 	ctx := env.authCtx(fmt.Sprintf("%d", userID))

// 	// SETUP: Dùng Fixture
// 	list := testutil.CreateTodoList(t, env.db, "Old Name", entity.UserID(userID))
// 	resourceName := fmt.Sprintf("users/%d/todo-lists/%d", userID, list.ID)

// 	// ACTION: Update
// 	_, err := env.client.UpdateTodoList(ctx, &todov1.UpdateTodoListRequest{
// 		TodoList: &todov1.TodoList{
// 			Name:        resourceName,
// 			DisplayName: "New Name",
// 		},
// 		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
// 	})
// 	require.NoError(t, err)

// 	// ASSERT
// 	fetched, err := env.client.GetTodoList(ctx, &todov1.GetTodoListRequest{Name: resourceName})
// 	require.NoError(t, err)
// 	assert.Equal(t, "New Name", fetched.TodoList.DisplayName)
// }
