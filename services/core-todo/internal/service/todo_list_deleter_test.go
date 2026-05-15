package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func TestTodoListDeleter_Delete(t *testing.T) {
	t.Parallel()

	type txKey struct{}

	var (
		ctx         = context.Background()
		txCtx       = context.WithValue(ctx, txKey{}, "tx-test")
		todoListID  = entity.TodoListID(1)
		requesterID = entity.UserID(10)
		ownerID     = entity.UserID(10)

		todoList = &entity.TodoList{
			ID:      todoListID,
			Name:    "Work Tasks",
			OwnerID: ownerID,
		}

		validInput = &input.TodoListDeleter{
			TodoListID:  todoListID,
			RequesterID: requesterID,
		}
	)

	type fields struct {
		mockTransactor       *mock.MockTransactor
		mockTodoListQueries  *mock.MockTodoListQueriesGateway
		mockTodoListCommands *mock.MockTodoListCommandsGateway
		mockTodoCommands     *mock.MockTodoCommandsGateway
	}

	tests := map[string]struct {
		prepare func(f *fields)
		input   *input.TodoListDeleter
		wantErr bool
		errCode entity.ErrorCode
	}{
		"success: delete todos and todo list in transaction": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(todoList, nil)

				f.mockTransactor.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})

				gomock.InOrder(
					f.mockTodoCommands.EXPECT().
						DeleteByTodoListID(txCtx, todoListID).
						Return(nil),

					f.mockTodoListCommands.EXPECT().
						Delete(txCtx, todoListID).
						Return(nil),
				)
			},
			input: validInput,
		},

		"error: todo list query gateway error": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(nil, fmt.Errorf("connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},

		"error: todo list not found": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(nil, nil)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: requester is not todo list owner": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(&entity.TodoList{
						ID:      todoListID,
						Name:    "Work Tasks",
						OwnerID: entity.UserID(999),
					}, nil)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrAuthZ,
		},

		"error: delete todos in transaction failed": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(todoList, nil)

				f.mockTransactor.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})

				f.mockTodoCommands.EXPECT().
					DeleteByTodoListID(txCtx, todoListID).
					Return(fmt.Errorf("db error"))
			},
			input:   validInput,
			wantErr: true,
		},

		"error: delete todo list in transaction failed": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(todoList, nil)

				f.mockTransactor.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})

				gomock.InOrder(
					f.mockTodoCommands.EXPECT().
						DeleteByTodoListID(txCtx, todoListID).
						Return(nil),

					f.mockTodoListCommands.EXPECT().
						Delete(txCtx, todoListID).
						Return(fmt.Errorf("db error")),
				)
			},
			input:   validInput,
			wantErr: true,
		},

		"error: transactor returns error": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(ctx, todoListID).
					Return(todoList, nil)

				f.mockTransactor.EXPECT().
					Transaction(ctx, gomock.Any()).
					Return(fmt.Errorf("transaction failed"))
			},
			input:   validInput,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			f := &fields{
				mockTransactor:       mock.NewMockTransactor(ctrl),
				mockTodoListQueries:  mock.NewMockTodoListQueriesGateway(ctrl),
				mockTodoListCommands: mock.NewMockTodoListCommandsGateway(ctrl),
				mockTodoCommands:     mock.NewMockTodoCommandsGateway(ctrl),
			}

			tt.prepare(f)

			sut := service.NewTodoListDeleter(
				f.mockTransactor,
				f.mockTodoListQueries,
				f.mockTodoListCommands,
				f.mockTodoCommands,
			)

			got, err := sut.Delete(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)

				if tt.errCode != "" {
					var appErr *entity.AppError
					assert.ErrorAs(t, err, &appErr)
					assert.Equal(t, tt.errCode, appErr.Code)
				}

				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}
