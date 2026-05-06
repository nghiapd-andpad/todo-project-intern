package persistence_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/testutil"
)

func TestTodoQueriesGateway_Get(t *testing.T) {
	t.Parallel()

	t.Run("success: found todo", func(t *testing.T) {
		t.Parallel()

		// load config
		cfg, err := config.New()
		require.NoError(t, err)

		db := testutil.NewTestDB(t, cfg)
		queryRepo := persistence.NewTodoQueriesGateway(db)

		created := testutil.CreateTodo(t, db, entity.TodoListID(1), "Unit Test Todo", entity.UserID(1))

		got, err := queryRepo.Get(context.Background(), created.ID)

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, "Unit Test Todo", got.Title)
	})

	t.Run("not found: returns nil nil", func(t *testing.T) {
		t.Parallel()

		// load config
		cfg, err := config.New()
		require.NoError(t, err)

		db := testutil.NewTestDB(t, cfg)
		queryRepo := persistence.NewTodoQueriesGateway(db)

		got, err := queryRepo.Get(context.Background(), entity.TodoID(9999))

		// not error
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestTodoQueriesGateway_List(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup     func(t *testing.T, db *gorm.DB)
		opts      gateway.ListTodosOptions
		wantLen   int
		wantTotal int64
	}{
		"success: list by todo_list_id": {
			setup: func(t *testing.T, db *gorm.DB) {
				for i := 0; i < 3; i++ {
					testutil.CreateTodo(t, db,
						entity.TodoListID(1),
						fmt.Sprintf("Todo %d", i),
						entity.UserID(1),
					)
				}
				testutil.CreateTodo(t, db, entity.TodoListID(2), "Other List", entity.UserID(1))
			},
			opts: gateway.ListTodosOptions{
				TodoListID: todoListIDPtr(entity.TodoListID(1)),
				Limit:      10,
			},
			wantLen:   3,
			wantTotal: 3,
		},
		"success: filter by status PENDING": {
			setup: func(t *testing.T, db *gorm.DB) {
				cmdRepo := persistence.NewTodoCommandsGateway(db)

				testutil.CreateTodo(t, db, entity.TodoListID(1), "Pending 1", entity.UserID(1))
				testutil.CreateTodo(t, db, entity.TodoListID(1), "Pending 2", entity.UserID(1))
				done := testutil.CreateTodo(t, db, entity.TodoListID(1), "Done", entity.UserID(1))

				// update status to DONE
				done.Status = entity.TodoStatusDone

				cmdRepo.Update(context.Background(), done)
			},
			opts: gateway.ListTodosOptions{
				TodoListID: todoListIDPtr(entity.TodoListID(1)),
				Status:     todoStatusPtr(entity.TodoStatusPending),
				Limit:      10,
			},
			wantLen:   2,
			wantTotal: 2,
		},
		"success: pagination": {
			setup: func(t *testing.T, db *gorm.DB) {
				for i := 0; i < 5; i++ {
					testutil.CreateTodo(t, db,
						entity.TodoListID(1),
						fmt.Sprintf("Todo %d", i),
						entity.UserID(1),
					)
				}
			},
			opts: gateway.ListTodosOptions{
				TodoListID: todoListIDPtr(entity.TodoListID(1)),
				Limit:      3,
				Offset:     3,
			},
			wantLen:   2,
			wantTotal: 5,
		},
		"success: title search": {
			setup: func(t *testing.T, db *gorm.DB) {
				testutil.CreateTodo(t, db, entity.TodoListID(1), "todo 1", entity.UserID(1))
				testutil.CreateTodo(t, db, entity.TodoListID(1), "todo 2", entity.UserID(1))
				testutil.CreateTodo(t, db, entity.TodoListID(1), "task 3", entity.UserID(1))
			},
			opts: gateway.ListTodosOptions{
				TodoListID:  todoListIDPtr(entity.TodoListID(1)),
				TitleSearch: strPtr("todo"),
				Limit:       10,
			},
			wantLen:   2,
			wantTotal: 2,
		},
		"success: empty list": {
			setup: func(t *testing.T, db *gorm.DB) {},
			opts: gateway.ListTodosOptions{
				TodoListID: todoListIDPtr(entity.TodoListID(1)),
				Limit:      10,
			},
			wantLen:   0,
			wantTotal: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// load config
			cfg, err := config.New()
			require.NoError(t, err)

			db := testutil.NewTestDB(t, cfg)
			tt.setup(t, db)
			queryRepo := persistence.NewTodoQueriesGateway(db)

			got, total, err := queryRepo.List(context.Background(), tt.opts)

			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

// helper functions
func todoListIDPtr(id entity.TodoListID) *entity.TodoListID { return &id }
func todoStatusPtr(s entity.TodoStatus) *entity.TodoStatus  { return &s }
