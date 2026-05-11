package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestTodoCommandsGateway_Create(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup    func(t *testing.T) *gorm.DB
		input    *entity.Todo
		validate func(t *testing.T, got *entity.Todo)
	}{
		"success: required fields only": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			input: &entity.Todo{
				TodoListID: entity.TodoListID(1),
				Title:      "Unit Test Todo",
				Status:     entity.TodoStatusPending,
				Priority:   entity.PriorityMedium,
				CreatorID:  entity.UserID(1),
			},
			validate: func(t *testing.T, got *entity.Todo) {
				assert.NotZero(t, got.ID)
				assert.Equal(t, "Unit Test Todo", got.Title)
				assert.Equal(t, entity.TodoStatusPending, got.Status)
				assert.Equal(t, entity.PriorityMedium, got.Priority)

				assert.Nil(t, got.Description)
				assert.Nil(t, got.DueDate)
				assert.Nil(t, got.AssigneeID)

				assert.NotZero(t, got.CreatedAt)
				assert.NotZero(t, got.UpdatedAt)
			},
		},

		"success: with optional fields": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			input: &entity.Todo{
				TodoListID:  entity.TodoListID(1),
				Title:       "Unit Test Todo with Optional Fields",
				Description: testutil.StrPtr("Unit test to verify create with optional fields"),
				Status:      entity.TodoStatusPending,
				Priority:    entity.PriorityHigh,
				CreatorID:   entity.UserID(1),
				AssigneeID:  testutil.UserIDPtr(entity.UserID(2)),
			},
			validate: func(t *testing.T, got *entity.Todo) {
				assert.NotZero(t, got.ID)

				require.NotNil(t, got.Description)
				assert.Equal(t, "Unit test to verify create with optional fields", *got.Description)

				require.NotNil(t, got.AssigneeID)
				assert.Equal(t, entity.UserID(2), *got.AssigneeID)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := tt.setup(t)
			repo := persistence.NewTodoCommandsGateway(db)

			got, err := repo.Create(context.Background(), tt.input)

			require.NoError(t, err)
			require.NotNil(t, got)

			tt.validate(t, got)
		})
	}
}

func TestTodoCommandsGateway_Update(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup    func(t *testing.T) (*gorm.DB, *entity.Todo)
		mutate   func(todo *entity.Todo)
		validate func(t *testing.T, got *entity.Todo)
	}{
		"success: update title": {
			setup: func(t *testing.T) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				todo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Old Title",
					entity.UserID(1),
				)

				return db, todo
			},
			mutate: func(todo *entity.Todo) {
				todo.Title = "New Title"
			},
			validate: func(t *testing.T, got *entity.Todo) {
				assert.Equal(t, "New Title", got.Title)
				assert.Equal(t, entity.TodoStatusPending, got.Status)
			},
		},

		"success: update status": {
			setup: func(t *testing.T) (*gorm.DB, *entity.Todo) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				todo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"Task 1",
					entity.UserID(1),
				)

				return db, todo
			},
			mutate: func(todo *entity.Todo) {
				todo.Status = entity.TodoStatusInProgress
			},
			validate: func(t *testing.T, got *entity.Todo) {
				assert.Equal(t, entity.TodoStatusInProgress, got.Status)
				assert.Equal(t, "Task 1", got.Title)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, existing := tt.setup(t)

			repo := persistence.NewTodoCommandsGateway(db)

			tt.mutate(existing)

			got, err := repo.Update(context.Background(), existing)

			require.NoError(t, err)
			require.NotNil(t, got)

			tt.validate(t, got)
		})
	}
}

func TestTodoCommandsGateway_Delete(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup    func(t *testing.T) *gorm.DB
		testFunc func(
			t *testing.T,
			db *gorm.DB,
			cmdRepo *persistence.TodoCommandsGateway,
			queryRepo *persistence.TodoQueriesGateway,
		)
	}{
		"success: soft delete": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			testFunc: func(
				t *testing.T,
				db *gorm.DB,
				cmdRepo *persistence.TodoCommandsGateway,
				queryRepo *persistence.TodoQueriesGateway,
			) {
				todo := testutil.CreateTodo(
					t,
					db,
					entity.TodoListID(1),
					"To Be Deleted",
					entity.UserID(1),
				)

				before, err := queryRepo.Get(context.Background(), todo.ID)
				require.NoError(t, err)
				require.NotNil(t, before)

				err = cmdRepo.Delete(context.Background(), todo.ID)
				require.NoError(t, err)

				after, err := queryRepo.Get(context.Background(), todo.ID)
				require.NoError(t, err)

				assert.Nil(t, after)
			},
		},

		"success: delete not existent, no error": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			testFunc: func(
				t *testing.T,
				_ *gorm.DB,
				cmdRepo *persistence.TodoCommandsGateway,
				_ *persistence.TodoQueriesGateway,
			) {
				err := cmdRepo.Delete(context.Background(), entity.TodoID(999))

				assert.NoError(t, err)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := tt.setup(t)

			cmdRepo := persistence.NewTodoCommandsGateway(db)
			queryRepo := persistence.NewTodoQueriesGateway(db)

			tt.testFunc(t, db, cmdRepo, queryRepo)
		})
	}
}
