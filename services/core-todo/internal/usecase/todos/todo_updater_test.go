package todos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestTodoUpdater_Update(t *testing.T) {
	t.Parallel()

	var (
		ctx    = context.Background()
		todoID = entity.TodoID(1)

		newTitle  = "New title"
		newStatus = entity.TodoStatusInProgress
		badDue    = "01-05-2026"
	)

	oldTodo := func() *entity.Todo {
		return &entity.Todo{
			ID:       todoID,
			Title:    "Old title",
			Status:   entity.TodoStatusPending,
			Priority: entity.PriorityLow,
		}
	}

	type fields struct {
		mockCommands *mock.MockTodoCommandsGateway
		mockQueries  *mock.MockTodoQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoUpdater
		expected *output.TodoUpdater
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: update title only": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoID).
						Return(oldTodo(), nil),
					f.mockCommands.EXPECT().
						Update(gomock.Any(), &entity.Todo{
							ID:       todoID,
							Title:    "New title",
							Status:   entity.TodoStatusPending,
							Priority: entity.PriorityLow,
						}).
						Return(&entity.Todo{
							ID:       todoID,
							Title:    "New title",
							Status:   entity.TodoStatusPending,
							Priority: entity.PriorityLow,
						}, nil),
				)
			},
			input: &input.TodoUpdater{
				ID:    todoID,
				Title: &newTitle,
			},
			expected: &output.TodoUpdater{Todo: &entity.Todo{
				ID:       todoID,
				Title:    "New title",
				Status:   entity.TodoStatusPending,
				Priority: entity.PriorityLow,
			}},
		},

		"success: update status only": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoID).
						Return(oldTodo(), nil),
					f.mockCommands.EXPECT().
						Update(gomock.Any(), &entity.Todo{
							ID:       todoID,
							Title:    "Old title",
							Status:   entity.TodoStatusInProgress,
							Priority: entity.PriorityLow,
						}).
						Return(&entity.Todo{
							ID:       todoID,
							Title:    "Old title",
							Status:   entity.TodoStatusInProgress,
							Priority: entity.PriorityLow,
						}, nil),
				)
			},
			input: &input.TodoUpdater{
				ID:     todoID,
				Status: &newStatus,
			},
			expected: &output.TodoUpdater{Todo: &entity.Todo{
				ID:       todoID,
				Title:    "Old title",
				Status:   entity.TodoStatusInProgress,
				Priority: entity.PriorityLow,
			}},
		},

		"error: invalid due_date format": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoID).
					Return(oldTodo(), nil)
			},
			input: &input.TodoUpdater{
				ID:      todoID,
				DueDate: &badDue,
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		"error: todo not found": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoID).
					Return(nil, nil)
			},
			input:   &input.TodoUpdater{ID: todoID},
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: todo query error": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoID).
					Return(nil, fmt.Errorf("db error"))
			},
			input:   &input.TodoUpdater{ID: todoID},
			wantErr: true,
		},

		"error: db error": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoID).
						Return(oldTodo(), nil),
					f.mockCommands.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(nil, fmt.Errorf("db error")),
				)
			},
			input:   &input.TodoUpdater{ID: todoID, Title: &newTitle},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{
				mockCommands: mock.NewMockTodoCommandsGateway(ctrl),
				mockQueries:  mock.NewMockTodoQueriesGateway(ctrl),
			}
			tt.prepare(f)

			sut := todos.NewTodoUpdater(f.mockCommands, f.mockQueries)
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
