package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestTodoQueriesGateway_Get(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup func(t *testing.T, ctx context.Context) (*gorm.DB, *entity.Todo)
		test  func(
			t *testing.T,
			ctx context.Context,
			repo *persistence.TodoQueriesGateway,
			existingTodo *entity.Todo,
		)
	}{
		"success: found todo by id and todo list id": {
			setup: func(t *testing.T, ctx context.Context) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Unit Test Todo",
				)

				return db, existingTodo
			},
			test: func(
				t *testing.T,
				ctx context.Context,
				repo *persistence.TodoQueriesGateway,
				existingTodo *entity.Todo,
			) {
				got, err := repo.Get(
					ctx,
					existingTodo.ID,
					existingTodo.TodoListID,
				)

				require.NoError(t, err)
				require.NotNil(t, got)

				assert.Equal(t, existingTodo.ID, got.ID)
				assert.Equal(t, existingTodo.TodoListID, got.TodoListID)
				assert.Equal(t, "Unit Test Todo", got.Title)
			},
		},

		"not found: returns nil nil": {
			setup: func(t *testing.T, ctx context.Context) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				return db, nil
			},
			test: func(
				t *testing.T,
				ctx context.Context,
				repo *persistence.TodoQueriesGateway,
				_ *entity.Todo,
			) {
				got, err := repo.Get(
					ctx,
					entity.TodoID(9999),
					entity.TodoListID(1),
				)

				assert.NoError(t, err)
				assert.Nil(t, got)
			},
		},

		"not found: same todo id but different todo list id": {
			setup: func(t *testing.T, ctx context.Context) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Unit Test Todo",
				)

				return db, existingTodo
			},
			test: func(
				t *testing.T,
				ctx context.Context,
				repo *persistence.TodoQueriesGateway,
				existingTodo *entity.Todo,
			) {
				got, err := repo.Get(
					ctx,
					existingTodo.ID,
					entity.TodoListID(999),
				)

				assert.NoError(t, err)
				assert.Nil(t, got)
			},
		},

		"not found: soft-deleted todo returns nil nil": {
			setup: func(t *testing.T, ctx context.Context) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Deleted Todo",
				)

				cmdRepo := persistence.NewTodoCommandsGateway(db)

				err := cmdRepo.Delete(ctx, existingTodo.ID)
				require.NoError(t, err)

				return db, existingTodo
			},
			test: func(
				t *testing.T,
				ctx context.Context,
				repo *persistence.TodoQueriesGateway,
				existingTodo *entity.Todo,
			) {
				got, err := repo.Get(
					ctx,
					existingTodo.ID,
					existingTodo.TodoListID,
				)

				assert.NoError(t, err)
				assert.Nil(t, got)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			db, existingTodo := tt.setup(t, ctx)

			repo := persistence.NewTodoQueriesGateway(db)

			tt.test(t, ctx, repo, existingTodo)
		})
	}
}

func TestTodoQueriesGateway_List(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup     func(t *testing.T, ctx context.Context) *gorm.DB
		opts      gatewayinput.ListTodosOptions
		wantLen   int
		wantTotal int64
	}{
		"success: list by todo_list_id": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				for i := 0; i < 3; i++ {
					testutil.CreateTodo(
						t,
						db,
						entity.TodoListID(1),
						fmt.Sprintf("Todo %d", i),
					)
				}

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(2),
					"Other List",
				)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: entity.TodoListID(1),
				Limit:      10,
			},
			wantLen:   3,
			wantTotal: 3,
		},

		"success: filter by assignee": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				testutil.CreateTodoWithAssignee(
					t,
					db,
					entity.TodoListID(1),
					"Assigned to 1",
					entity.UserID(1),
				)
				testutil.CreateTodoWithAssignee(
					t,
					db,
					entity.TodoListID(1),
					"Assigned to 2",
					entity.UserID(2),
				)
				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Unassigned",
				)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID:   entity.TodoListID(1),
				AssigneeOnly: testutil.UserIDPtr(entity.UserID(1)),
				Limit:        10,
			},
			wantLen:   1,
			wantTotal: 1,
		},

		"success: filter by status pending": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				cmdRepo := persistence.NewTodoCommandsGateway(db)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Pending 1",
				)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Pending 2",
				)

				doneTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Done",
				)

				doneTodo.Status = entity.TodoStatusDone

				_, err := cmdRepo.Update(
					ctx,
					doneTodo,
				)
				require.NoError(t, err)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: entity.TodoListID(1),
				Status:     testutil.TodoStatusPtr(entity.TodoStatusPending),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: filter by priority": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				cmdRepo := persistence.NewTodoCommandsGateway(db)

				urgentTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Urgent Todo",
				)
				urgentTodo.Priority = entity.PriorityUrgent

				_, err := cmdRepo.Update(
					ctx,
					urgentTodo,
				)
				require.NoError(t, err)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Medium Todo",
				)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: entity.TodoListID(1),
				Priority:   testutil.PriorityPtr(entity.PriorityUrgent),
				Limit:      10,
			},
			wantLen:   1,
			wantTotal: 1,
		},

		"success: pagination": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				for i := 0; i < 5; i++ {
					testutil.CreateTodo(
						t,
						db,
						entity.TodoListID(1),
						fmt.Sprintf("Todo %d", i),
					)
				}

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: entity.TodoListID(1),
				Limit:      3,
				Offset:     3,
			},
			wantLen:   2,
			wantTotal: 5,
		},

		"success: title search case-insensitive": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"TODO 1",
				)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"todo 2",
				)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"task 3",
				)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID:  entity.TodoListID(1),
				TitleSearch: testutil.StrPtr("todo"),
				Limit:       10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: empty list": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: entity.TodoListID(1),
				Limit:      10,
			},
			wantLen:   0,
			wantTotal: 0,
		},

		"success: soft-deleted todos are excluded": {
			setup: func(t *testing.T, ctx context.Context) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				cmdRepo := persistence.NewTodoCommandsGateway(db)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Visible Todo",
				)

				deletedTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Deleted Todo",
				)

				err := cmdRepo.Delete(ctx, deletedTodo.ID)
				require.NoError(t, err)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: entity.TodoListID(1),
				Limit:      10,
			},
			wantLen:   1,
			wantTotal: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			db := tt.setup(t, ctx)

			repo := persistence.NewTodoQueriesGateway(db)

			got, total, err := repo.List(
				ctx,
				&tt.opts,
			)

			require.NoError(t, err)

			assert.Len(t, got, tt.wantLen)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}
