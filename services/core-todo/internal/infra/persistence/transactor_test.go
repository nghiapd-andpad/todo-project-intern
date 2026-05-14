package persistence_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/testutil"
)

func TestTransactor_Transaction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup func(t *testing.T) (
			*persistence.Transactor,
			*persistence.TodoListCommandsGateway,
			*persistence.TodoCommandsGateway,
			*persistence.TodoListQueriesGateway,
			*persistence.TodoQueriesGateway,
		)
		test func(
			t *testing.T,
			transactor *persistence.Transactor,
			todoListCmdRepo *persistence.TodoListCommandsGateway,
			todoCmdRepo *persistence.TodoCommandsGateway,
			todoListQueryRepo *persistence.TodoListQueriesGateway,
			todoQueryRepo *persistence.TodoQueriesGateway,
		)
	}{
		"success: commit transaction": {
			setup: func(t *testing.T) (
				*persistence.Transactor,
				*persistence.TodoListCommandsGateway,
				*persistence.TodoCommandsGateway,
				*persistence.TodoListQueriesGateway,
				*persistence.TodoQueriesGateway,
			) {
				cfg := testutil.NewTestConfig(t)
				db := testutil.NewTestDB(t, cfg)

				return persistence.NewTransactor(db),
					persistence.NewTodoListCommandsGateway(db),
					persistence.NewTodoCommandsGateway(db),
					persistence.NewTodoListQueriesGateway(db),
					persistence.NewTodoQueriesGateway(db)
			},
			test: func(
				t *testing.T,
				transactor *persistence.Transactor,
				todoListCmdRepo *persistence.TodoListCommandsGateway,
				todoCmdRepo *persistence.TodoCommandsGateway,
				todoListQueryRepo *persistence.TodoListQueriesGateway,
				todoQueryRepo *persistence.TodoQueriesGateway,
			) {
				var createdTodoListID entity.TodoListID

				err := transactor.Transaction(context.Background(), func(txCtx context.Context) error {
					todoList, err := todoListCmdRepo.Create(txCtx, &entity.TodoList{
						Name:    "Committed Todo List",
						OwnerID: entity.UserID(1),
					})
					if err != nil {
						return err
					}

					createdTodoListID = todoList.ID

					_, err = todoCmdRepo.Create(txCtx, &entity.Todo{
						TodoListID: todoList.ID,
						Title:      "Committed Todo",
						Status:     entity.TodoStatusPending,
						Priority:   entity.PriorityMedium,
					})
					return err
				})

				require.NoError(t, err)

				todoList, err := todoListQueryRepo.Get(context.Background(), createdTodoListID)
				require.NoError(t, err)
				require.NotNil(t, todoList)
				assert.Equal(t, "Committed Todo List", todoList.Name)

				todos, total, err := todoQueryRepo.List(context.Background(), &gatewayinput.ListTodosOptions{
					TodoListID: createdTodoListID,
					Limit:      10,
				})
				require.NoError(t, err)
				assert.Equal(t, int64(1), total)
				require.Len(t, todos, 1)
				assert.Equal(t, "Committed Todo", todos[0].Title)
			},
		},

		"failure: rollback transaction": {
			setup: func(t *testing.T) (
				*persistence.Transactor,
				*persistence.TodoListCommandsGateway,
				*persistence.TodoCommandsGateway,
				*persistence.TodoListQueriesGateway,
				*persistence.TodoQueriesGateway,
			) {
				cfg := testutil.NewTestConfig(t)
				db := testutil.NewTestDB(t, cfg)

				return persistence.NewTransactor(db),
					persistence.NewTodoListCommandsGateway(db),
					persistence.NewTodoCommandsGateway(db),
					persistence.NewTodoListQueriesGateway(db),
					persistence.NewTodoQueriesGateway(db)
			},
			test: func(
				t *testing.T,
				transactor *persistence.Transactor,
				todoListCmdRepo *persistence.TodoListCommandsGateway,
				todoCmdRepo *persistence.TodoCommandsGateway,
				todoListQueryRepo *persistence.TodoListQueriesGateway,
				_ *persistence.TodoQueriesGateway,
			) {
				errRollback := errors.New("force rollback")

				err := transactor.Transaction(context.Background(), func(txCtx context.Context) error {
					todoList, err := todoListCmdRepo.Create(txCtx, &entity.TodoList{
						Name:    "Rolled Back Todo List",
						OwnerID: entity.UserID(1),
					})
					if err != nil {
						return err
					}

					_, err = todoCmdRepo.Create(txCtx, &entity.Todo{
						TodoListID: todoList.ID,
						Title:      "Rolled Back Todo",
						Status:     entity.TodoStatusPending,
						Priority:   entity.PriorityMedium,
					})
					if err != nil {
						return err
					}

					return errRollback
				})

				require.Error(t, err)
				assert.ErrorIs(t, err, errRollback)

				todoLists, total, err := todoListQueryRepo.List(context.Background(), &gatewayinput.ListTodoListsOptions{
					OwnerID: testutil.UserIDPtr(entity.UserID(1)),
					Limit:   10,
				})
				require.NoError(t, err)
				assert.Equal(t, int64(0), total)
				assert.Empty(t, todoLists)
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			transactor,
				todoListCmdRepo,
				todoCmdRepo,
				todoListQueryRepo,
				todoQueryRepo := tt.setup(t)

			tt.test(
				t,
				transactor,
				todoListCmdRepo,
				todoCmdRepo,
				todoListQueryRepo,
				todoQueryRepo,
			)
		})
	}
}
