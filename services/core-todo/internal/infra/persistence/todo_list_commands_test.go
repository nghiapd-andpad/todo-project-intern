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

func TestTodoListCommandsGateway_Create(t *testing.T) {
	t.Parallel()

	cfg, err := config.New()
	require.NoError(t, err)

	tests := map[string]struct {
		input    *entity.TodoList
		validate func(t *testing.T, got *entity.TodoList)
	}{
		"success: basic create": {
			input: &entity.TodoList{
				Name:    "Work Tasks",
				OwnerID: entity.UserID(1),
			},
			validate: func(t *testing.T, got *entity.TodoList) {
				assert.NotZero(t, got.ID)
				assert.Equal(t, "Work Tasks", got.Name)
				assert.Equal(t, entity.UserID(1), got.OwnerID)
				assert.NotZero(t, got.CreatedAt)
				assert.NotZero(t, got.UpdatedAt)
			},
		},

		"success: different owner": {
			input: &entity.TodoList{
				Name:    "Personal",
				OwnerID: entity.UserID(2),
			},
			validate: func(t *testing.T, got *entity.TodoList) {
				assert.Equal(t, "Personal", got.Name)
				assert.Equal(t, entity.UserID(2), got.OwnerID)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutil.NewTestDB(t, cfg)
			repo := persistence.NewTodoListCommandsGateway(db)

			got, err := repo.Create(context.Background(), tt.input)

			require.NoError(t, err)
			require.NotNil(t, got)

			tt.validate(t, got)
		})
	}
}

func TestTodoListCommandsGateway_Update(t *testing.T) {
	t.Parallel()

	cfg, err := config.New()
	require.NoError(t, err)

	tests := map[string]struct {
		setup    func(t *testing.T, db *gorm.DB) *entity.TodoList
		mutate   func(tl *entity.TodoList)
		validate func(t *testing.T, got *entity.TodoList)
	}{
		"success: update name": {
			setup: func(t *testing.T, db *gorm.DB) *entity.TodoList {
				return testutil.CreateTodoList(t, db, "Old Name", entity.UserID(1))
			},
			mutate: func(tl *entity.TodoList) {
				tl.Name = "New Name"
			},
			validate: func(t *testing.T, got *entity.TodoList) {
				assert.Equal(t, "New Name", got.Name)
				assert.Equal(t, entity.UserID(1), got.OwnerID)
			},
		},

		"success: update owner": {
			setup: func(t *testing.T, db *gorm.DB) *entity.TodoList {
				return testutil.CreateTodoList(t, db, "Project A", entity.UserID(1))
			},
			mutate: func(tl *entity.TodoList) {
				tl.OwnerID = entity.UserID(2)
			},
			validate: func(t *testing.T, got *entity.TodoList) {
				assert.Equal(t, entity.UserID(2), got.OwnerID)
				assert.Equal(t, "Project A", got.Name)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := testutil.NewTestDB(t, cfg)
			repo := persistence.NewTodoListCommandsGateway(db)

			existing := tt.setup(t, db)
			tt.mutate(existing)

			got, err := repo.Update(context.Background(), existing)

			require.NoError(t, err)
			require.NotNil(t, got)

			tt.validate(t, got)
		})
	}
}

func TestTodoListCommandsGateway_Delete(t *testing.T) {
	t.Parallel()

	cfg, err := config.New()
	require.NoError(t, err)

	t.Run("success: delete existing todo list", func(t *testing.T) {
		t.Parallel()

		db := testutil.NewTestDB(t, cfg)
		repo := persistence.NewTodoListCommandsGateway(db)

		tl := testutil.CreateTodoList(t, db, "To Delete", entity.UserID(1))

		err := repo.Delete(context.Background(), tl.ID)
		require.NoError(t, err)
	})

	t.Run("success: delete non-existent - no error", func(t *testing.T) {
		t.Parallel()

		db := testutil.NewTestDB(t, cfg)
		repo := persistence.NewTodoListCommandsGateway(db)

		err := repo.Delete(context.Background(), entity.TodoListID(9999))
		assert.NoError(t, err)
	})
}
