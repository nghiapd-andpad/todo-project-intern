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

func TestTodoUpdater_Update(t *testing.T) {
	t.Parallel()

	var (
		todoID      = entity.TodoID(1)
		todoListID  = entity.TodoListID(2)
		ownerID     = entity.UserID(10)
		assigneeID  = entity.UserID(20)
		requesterID = ownerID

		newTitle       = "New title"
		newDescription = "New description"
		newStatus      = entity.TodoStatusInProgress
		newPriority    = entity.PriorityHigh
		newDueDate     = time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
		newAssigneeID  = entity.UserID(30)
	)

	oldTodo := func() *entity.Todo {
		return &entity.Todo{
			ID:         todoID,
			TodoListID: todoListID,
			Title:      "Old title",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityLow,
			AssigneeID: &assigneeID,
		}
	}

	oldTodoWithoutAssignee := func() *entity.Todo {
		return &entity.Todo{
			ID:         todoID,
			TodoListID: todoListID,
			Title:      "Old title",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityLow,
		}
	}

	todoList := func(ownerID entity.UserID) *entity.TodoList {
		return &entity.TodoList{
			ID:      todoListID,
			Name:    "Work Tasks",
			OwnerID: ownerID,
		}
	}

	ownerUpdatedTodo := func() *entity.Todo {
		return &entity.Todo{
			ID:          todoID,
			TodoListID:  todoListID,
			Title:       newTitle,
			Description: &newDescription,
			Status:      newStatus,
			Priority:    newPriority,
			DueDate:     &newDueDate,
			AssigneeID:  &newAssigneeID,
		}
	}

	assigneeUpdatedTodo := func() *entity.Todo {
		return &entity.Todo{
			ID:         todoID,
			TodoListID: todoListID,
			Title:      "Old title",
			Status:     newStatus,
			Priority:   entity.PriorityLow,
			AssigneeID: &assigneeID,
		}
	}

	validInput := &input.TodoUpdater{
		TodoID:      todoID,
		TodoListID:  todoListID,
		RequesterID: requesterID,
		Fields: input.UpdateTodoFields{
			Title:       &newTitle,
			Description: &newDescription,
			Status:      &newStatus,
			Priority:    &newPriority,
			DueDate:     &newDueDate,
			AssigneeID:  &newAssigneeID,
		},
	}

	type fields struct {
		mockTodoListQueries *mock.MockTodoListQueriesGateway
		mockTodoQueries     *mock.MockTodoQueriesGateway
		mockTodoCommands    *mock.MockTodoCommandsGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields, ctx context.Context)
		input    *input.TodoUpdater
		expected *output.TodoUpdater
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: owner updates all fields": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),

					f.mockTodoCommands.EXPECT().
						Update(ctx, ownerUpdatedTodo()).
						Return(ownerUpdatedTodo(), nil),
				)
			},
			input: validInput,
			expected: &output.TodoUpdater{
				Todo: ownerUpdatedTodo(),
			},
		},

		"success: owner updates no fields keeps todo unchanged": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),

					f.mockTodoCommands.EXPECT().
						Update(ctx, oldTodo()).
						Return(oldTodo(), nil),
				)
			},
			input: &input.TodoUpdater{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: ownerID,
				Fields:      input.UpdateTodoFields{},
			},
			expected: &output.TodoUpdater{
				Todo: oldTodo(),
			},
		},

		"success: assignee updates status only": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),

					f.mockTodoCommands.EXPECT().
						Update(ctx, assigneeUpdatedTodo()).
						Return(assigneeUpdatedTodo(), nil),
				)
			},
			input: &input.TodoUpdater{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: assigneeID,
				Fields: input.UpdateTodoFields{
					Title:    &newTitle,
					Status:   &newStatus,
					Priority: &newPriority,
				},
			},
			expected: &output.TodoUpdater{
				Todo: assigneeUpdatedTodo(),
			},
		},

		"success: assignee updates non-status fields only keeps todo unchanged": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),

					f.mockTodoCommands.EXPECT().
						Update(ctx, oldTodo()).
						Return(oldTodo(), nil),
				)
			},
			input: &input.TodoUpdater{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: assigneeID,
				Fields: input.UpdateTodoFields{
					Title:    &newTitle,
					Priority: &newPriority,
				},
			},
			expected: &output.TodoUpdater{
				Todo: oldTodo(),
			},
		},

		"error: todo query gateway error": {
			prepare: func(f *fields, ctx context.Context) {
				f.mockTodoQueries.EXPECT().
					Get(ctx, todoID, todoListID).
					Return(nil, fmt.Errorf("db error"))
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
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(nil, fmt.Errorf("db error")),
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
						Return(oldTodo(), nil),

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
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),
				)
			},
			input: &input.TodoUpdater{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: entity.UserID(999),
				Fields: input.UpdateTodoFields{
					Status: &newStatus,
				},
			},
			wantErr: true,
			errCode: entity.ErrAuthZ,
		},

		"error: requester is not owner and todo has no assignee": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(oldTodoWithoutAssignee(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),
				)
			},
			input: &input.TodoUpdater{
				TodoID:      todoID,
				TodoListID:  todoListID,
				RequesterID: assigneeID,
				Fields: input.UpdateTodoFields{
					Status: &newStatus,
				},
			},
			wantErr: true,
			errCode: entity.ErrAuthZ,
		},

		"error: update command gateway error": {
			prepare: func(f *fields, ctx context.Context) {
				gomock.InOrder(
					f.mockTodoQueries.EXPECT().
						Get(ctx, todoID, todoListID).
						Return(oldTodo(), nil),

					f.mockTodoListQueries.EXPECT().
						Get(ctx, todoListID).
						Return(todoList(ownerID), nil),

					f.mockTodoCommands.EXPECT().
						Update(ctx, ownerUpdatedTodo()).
						Return(nil, fmt.Errorf("db error")),
				)
			},
			input:   validInput,
			wantErr: true,
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
				mockTodoCommands:    mock.NewMockTodoCommandsGateway(ctrl),
			}

			tt.prepare(f, ctx)

			sut := service.NewTodoUpdater(
				f.mockTodoListQueries,
				f.mockTodoQueries,
				f.mockTodoCommands,
			)

			got, err := sut.Update(ctx, tt.input)

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
