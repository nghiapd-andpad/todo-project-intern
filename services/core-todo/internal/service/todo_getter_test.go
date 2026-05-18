package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

func TestTodoGetter_Get(t *testing.T) {
	t.Parallel()

	var (
		now         = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todoID      = entity.TodoID(1)
		todoListID  = entity.TodoListID(2)
		ownerID     = entity.UserID(10)
		assigneeID  = entity.UserID(20)
		requesterID = ownerID

		todo = &entity.Todo{
			ID:         todoID,
			TodoListID: todoListID,
			Title:      "Unit Test Get Todo",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityMedium,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		assignedTodo = &entity.Todo{
			ID:         todoID,
			TodoListID: todoListID,
			Title:      "Assigned Todo",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityMedium,
			AssigneeID: &assigneeID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		todoList = &entity.TodoList{
			ID:      todoListID,
			Name:    "Work Tasks",
			OwnerID: ownerID,
		}

		validInput = &input.TodoGetter{
			TodoID:      todoID,
			TodoListID:  todoListID,
			RequesterID: requesterID,
		}
	)

	type fields struct {
		mockTodoListQueries *mock.MockTodoListQueriesGateway
		mockTodoQueries     *mock.MockTodoQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields, ctx context.Context)
		input    *input.TodoGetter
		expected *output.TodoGetter
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: requester is todo list owner": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(todo, nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList, nil),
				)
			},
			input:    validInput,
			expected: &output.TodoGetter{Todo: todo},
		},

		"success: requester is todo assignee": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(assignedTodo, nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(&entity.TodoList{
							ID:      todoListID,
							Name:    "Work Tasks",
							OwnerID: entity.UserID(999),
						}, nil),
				)
			},
			input: &input.TodoGetter{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: assigneeID,
			},
			expected: &output.TodoGetter{Todo: assignedTodo},
		},

		"error: todo query gateway error": {
			prepare: func(f *fields, ctx context.Context) {
				f.mockTodoQueries.EXPECT().
					Get(ctx, todoID, todoListID).
					Return(nil, fmt.Errorf("connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},

		"error: todo not found": {
			prepare: func(f *fields, ctx context.Context) {
				f.mockTodoQueries.EXPECT().
					Get(ctx, todoID, todoListID).
					Return(nil, nil)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: todo list query gateway error": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(todo, nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(nil, fmt.Errorf("connection lost")),
				)
			},
			input:   validInput,
			wantErr: true,
		},

		"error: todo list not found": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(todo, nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(nil, nil),
				)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: requester is neither owner nor assignee": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(assignedTodo, nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(&entity.TodoList{
							ID:      todoListID,
							Name:    "Work Tasks",
							OwnerID: ownerID,
						}, nil),
				)
			},
			input: &input.TodoGetter{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: entity.UserID(999),
			},
			wantErr: true,
			errCode: entity.ErrAuthZ,
		},

		"error: requester is not owner and todo has no assignee": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(todo, nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList, nil),
				)
			},
			input: &input.TodoGetter{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: entity.UserID(999),
			},
			wantErr: true,
			errCode: entity.ErrAuthZ,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)

			f := &fields{
				mockTodoListQueries: mock.NewMockTodoListQueriesGateway(ctrl),
				mockTodoQueries:     mock.NewMockTodoQueriesGateway(ctrl),
			}

			tt.prepare(f, ctx)

			sut := service.NewTodoGetter(
				f.mockTodoListQueries,
				f.mockTodoQueries,
			)

			got, err := sut.Get(ctx, tt.input)

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

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
