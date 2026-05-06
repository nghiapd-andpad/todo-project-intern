package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/testutil"
)

func TestTodoCommandsGateway_Create(t *testing.T) {
	t.Parallel()

	// load config
	cfg, err := config.New()
	require.NoError(t, err)

	tests := map[string]struct {
		input    *entity.Todo
		validate func(t *testing.T, got *entity.Todo)
	}{
		"success: required fields only": {
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
			input: &entity.Todo{
				TodoListID:  entity.TodoListID(1),
				Title:       "Unit Test Todo with Optional Fields",
				Description: strPtr("Unit test to verify create with optional fields"),
				Status:      entity.TodoStatusPending,
				Priority:    entity.PriorityHigh,
				CreatorID:   entity.UserID(1),
				AssigneeID:  userIDPtr(entity.UserID(2)),
			},
			validate: func(t *testing.T, got *entity.Todo) {
				assert.NotZero(t, got.ID)
				assert.NotNil(t, got.Description)
				assert.Equal(t, "Unit test to verify create with optional fields", *got.Description)
				assert.NotNil(t, got.AssigneeID)
				assert.Equal(t, entity.UserID(2), *got.AssigneeID)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutil.NewTestDB(t, cfg)
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

	// load config
	cfg, err := config.New()
	require.NoError(t, err)

	tests := map[string]struct {
		setup    func(t *testing.T, db *gorm.DB) *entity.Todo
		mutate   func(todo *entity.Todo)
		validate func(t *testing.T, got *entity.Todo)
	}{
		"success: update title": {
			setup: func(t *testing.T, db *gorm.DB) *entity.Todo {
				return testutil.CreateTodo(t, db, entity.TodoListID(1), "Old Title", entity.UserID(1))
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
			setup: func(t *testing.T, db *gorm.DB) *entity.Todo {
				return testutil.CreateTodo(t, db, entity.TodoListID(1), "Task 1", entity.UserID(1))
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutil.NewTestDB(t, cfg)
			repo := persistence.NewTodoCommandsGateway(db)

			existing := tt.setup(t, db)
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

	// load config
	cfg, err := config.New()
	require.NoError(t, err)

	t.Run("success: soft delete", func(t *testing.T) {
		t.Parallel()

		db := testutil.NewTestDB(t, cfg)
		cmdRepo := persistence.NewTodoCommandsGateway(db)
		queryRepo := persistence.NewTodoQueriesGateway(db)

		// Setup create a todo to be deleted
		todo := testutil.CreateTodo(t, db, entity.TodoListID(1), "To Be Deleted", entity.UserID(1))

		// Verify exists before delete
		before, err := queryRepo.Get(context.Background(), todo.ID)
		require.NoError(t, err)
		require.NotNil(t, before)

		// Delete
		err = cmdRepo.Delete(context.Background(), todo.ID)
		require.NoError(t, err)

		// Verify soft deleted
		after, err := queryRepo.Get(context.Background(), todo.ID)
		require.NoError(t, err)
		assert.Nil(t, after)
	})

	t.Run("success: delete not existent, no error", func(t *testing.T) {
		t.Parallel()

		db := testutil.NewTestDB(t, cfg)
		repo := persistence.NewTodoCommandsGateway(db)

		err := repo.Delete(context.Background(), entity.TodoID(999))
		assert.NoError(t, err)
	})
}

// helper functions
func strPtr(s string) *string                   { return &s }
func userIDPtr(id entity.UserID) *entity.UserID { return &id }
