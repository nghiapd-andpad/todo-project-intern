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

func TestTodoListCommandsGateway_Create(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup    func(t *testing.T) *gorm.DB
		input    *entity.TodoList
		validate func(t *testing.T, got *entity.TodoList)
	}{
		"success: basic create": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
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
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
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
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := tt.setup(t)

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

	tests := map[string]struct {
		setup    func(t *testing.T) (*gorm.DB, *entity.TodoList)
		mutate   func(todoList *entity.TodoList)
		validate func(t *testing.T, got *entity.TodoList)
	}{
		"success: update name": {
			setup: func(t *testing.T) (*gorm.DB, *entity.TodoList) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodoList := testutil.CreateTodoList(
					t,
					db,
					"Old Name",
					entity.UserID(1),
				)

				return db, existingTodoList
			},
			mutate: func(todoList *entity.TodoList) {
				todoList.Name = "New Name"
			},
			validate: func(t *testing.T, got *entity.TodoList) {
				assert.Equal(t, "New Name", got.Name)
				assert.Equal(t, entity.UserID(1), got.OwnerID)
			},
		},

		"success: update owner": {
			setup: func(t *testing.T) (*gorm.DB, *entity.TodoList) {
				cfg := testutil.NewTestConfig(t)

				db := testutil.NewTestDB(t, cfg)

				existingTodoList := testutil.CreateTodoList(
					t,
					db,
					"Project A",
					entity.UserID(1),
				)

				return db, existingTodoList
			},
			mutate: func(todoList *entity.TodoList) {
				todoList.OwnerID = entity.UserID(2)
			},
			validate: func(t *testing.T, got *entity.TodoList) {
				assert.Equal(t, entity.UserID(2), got.OwnerID)
				assert.Equal(t, "Project A", got.Name)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, existingTodoList := tt.setup(t)

			repo := persistence.NewTodoListCommandsGateway(db)

			tt.mutate(existingTodoList)

			got, err := repo.Update(context.Background(), existingTodoList)

			require.NoError(t, err)
			require.NotNil(t, got)

			tt.validate(t, got)
		})
	}
}

func TestTodoListCommandsGateway_Delete(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup    func(t *testing.T) *gorm.DB
		testFunc func(
			t *testing.T,
			db *gorm.DB,
			repo *persistence.TodoListCommandsGateway,
		)
	}{
		"success: delete existing todo list": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			testFunc: func(
				t *testing.T,
				db *gorm.DB,
				repo *persistence.TodoListCommandsGateway,
			) {
				existingTodoList := testutil.CreateTodoList(
					t,
					db,
					"To Delete",
					entity.UserID(1),
				)

				err := repo.Delete(
					context.Background(),
					existingTodoList.ID,
				)

				require.NoError(t, err)
			},
		},

		"success: delete non-existent - no error": {
			setup: func(t *testing.T) *gorm.DB {
				cfg := testutil.NewTestConfig(t)

				return testutil.NewTestDB(t, cfg)
			},
			testFunc: func(
				t *testing.T,
				_ *gorm.DB,
				repo *persistence.TodoListCommandsGateway,
			) {
				err := repo.Delete(
					context.Background(),
					entity.TodoListID(9999),
				)

				assert.NoError(t, err)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := tt.setup(t)

			repo := persistence.NewTodoListCommandsGateway(db)

			tt.testFunc(t, db, repo)
		})
	}
}
