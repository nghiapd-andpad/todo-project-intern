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
		setup func(t *testing.T) (*gorm.DB, *entity.Todo)
		test  func(
			t *testing.T,
			repo *persistence.TodoQueriesGateway,
			existingTodo *entity.Todo,
		)
	}{
		"success: found todo": {
			setup: func(t *testing.T) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Unit Test Todo",
					entity.UserID(1),
				)

				return db, existingTodo
			},
			test: func(
				t *testing.T,
				repo *persistence.TodoQueriesGateway,
				existingTodo *entity.Todo,
			) {
				got, err := repo.Get(
					context.Background(),
					existingTodo.ID,
				)

				require.NoError(t, err)
				require.NotNil(t, got)

				assert.Equal(t, existingTodo.ID, got.ID)
				assert.Equal(t, "Unit Test Todo", got.Title)
			},
		},

		"not found: returns nil nil": {
			setup: func(t *testing.T) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				return db, nil
			},
			test: func(
				t *testing.T,
				repo *persistence.TodoQueriesGateway,
				_ *entity.Todo,
			) {
				got, err := repo.Get(
					context.Background(),
					entity.TodoID(9999),
				)

				assert.NoError(t, err)
				assert.Nil(t, got)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, existingTodo := tt.setup(t)

			repo := persistence.NewTodoQueriesGateway(db)

			tt.test(t, repo, existingTodo)
		})
	}
}

func TestTodoQueriesGateway_List(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup     func(t *testing.T) *gorm.DB
		opts      gatewayinput.ListTodosOptions
		wantLen   int
		wantTotal int64
	}{
		"success: list by todo_list_id": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				for i := 0; i < 3; i++ {
					testutil.CreateTodo(
						t,
						db,
						entity.TodoListID(1),
						fmt.Sprintf("Todo %d", i),
						entity.UserID(1),
					)
				}

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(2),
					"Other List",
					entity.UserID(1),
				)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: testutil.TodoListIDPtr(entity.TodoListID(1)),
				Limit:      10,
			},
			wantLen:   3,
			wantTotal: 3,
		},

		"success: filter by status pending": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				cmdRepo := persistence.NewTodoCommandsGateway(db)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Pending 1",
					entity.UserID(1),
				)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Pending 2",
					entity.UserID(1),
				)

				doneTodo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Done",
					entity.UserID(1),
				)

				doneTodo.Status = entity.TodoStatusDone

				_, err := cmdRepo.Update(
					context.Background(),
					doneTodo,
				)
				require.NoError(t, err)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: testutil.TodoListIDPtr(entity.TodoListID(1)),
				Status:     testutil.TodoStatusPtr(entity.TodoStatusPending),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: pagination": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				for i := 0; i < 5; i++ {
					testutil.CreateTodo(
						t,
						db,
						entity.TodoListID(1),
						fmt.Sprintf("Todo %d", i),
						entity.UserID(1),
					)
				}

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: testutil.TodoListIDPtr(entity.TodoListID(1)),
				Limit:      3,
				Offset:     3,
			},
			wantLen:   2,
			wantTotal: 5,
		},

		"success: title search": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"todo 1",
					entity.UserID(1),
				)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"todo 2",
					entity.UserID(1),
				)

				testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"task 3",
					entity.UserID(1),
				)

				return db
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID:  testutil.TodoListIDPtr(entity.TodoListID(1)),
				TitleSearch: testutil.StrPtr("todo"),
				Limit:       10,
			},
			wantLen:   2,
			wantTotal: 2,
		},

		"success: empty list": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			opts: gatewayinput.ListTodosOptions{
				TodoListID: testutil.TodoListIDPtr(entity.TodoListID(1)),
				Limit:      10,
			},
			wantLen:   0,
			wantTotal: 0,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := tt.setup(t)

			repo := persistence.NewTodoQueriesGateway(db)

			got, total, err := repo.List(
				context.Background(),
				&tt.opts,
			)

			require.NoError(t, err)

			assert.Len(t, got, tt.wantLen)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}
