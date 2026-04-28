package todos_test

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
	usecase "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestTodoCreator_Create(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		now        = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todoListID = entity.TodoListID(2)
		creatorID  = entity.UserID(1)
		validDue   = "2026-05-01"

		createdEntity = &entity.Todo{
			ID:         entity.TodoID(10),
			TodoListID: todoListID,
			Title:      "Unit Test Create Todo",
			Status:     entity.TodoStatusPending,
			Priority:   entity.PriorityMedium,
			CreatorID:  creatorID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		validInput = &input.TodoCreator{
			TodoListID: todoListID,
			Title:      "Unit Test Create Todo",
			Priority:   entity.PriorityMedium,
			CreatorID:  creatorID,
		}
	)

	type fields struct {
		mockCommands *mock.MockTodoCommandsGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoCreator
		expected *output.TodoCreator
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		// Happy path — verify required fields
		"success: create with required fields": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&entity.Todo{})).
					DoAndReturn(func(_ context.Context, todo *entity.Todo) (*entity.Todo, error) {
						assert.Equal(t, entity.TodoStatusPending, todo.Status)
						assert.Nil(t, todo.DueDate)
						assert.Nil(t, todo.AssigneeID)
						return createdEntity, nil
					})
			},
			input:    validInput,
			expected: &output.TodoCreator{Todo: createdEntity},
		},

		// Happy path — parse DueDate
		"success: create with due_date": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&entity.Todo{})).
					DoAndReturn(func(_ context.Context, todo *entity.Todo) (*entity.Todo, error) {
						assert.NotNil(t, todo.DueDate)
						assert.Equal(t, 2026, todo.DueDate.Year())
						assert.Equal(t, time.May, todo.DueDate.Month())
						assert.Equal(t, 1, todo.DueDate.Day())
						entityWithDue := *createdEntity
						entityWithDue.DueDate = todo.DueDate
						return &entityWithDue, nil
					})
			},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "Unit Test Create Todo",
				Priority:   entity.PriorityMedium,
				CreatorID:  creatorID,
				DueDate:    &validDue,
			},
			expected: &output.TodoCreator{Todo: func() *entity.Todo {
				e := *createdEntity
				parsed, _ := time.Parse("2006-01-02", validDue)
				e.DueDate = &parsed
				return &e
			}()},
		},

		// Error path — DueDate wrong format
		"error: invalid due_date format": {
			prepare: func(f *fields) {
			},
			input: &input.TodoCreator{
				TodoListID: todoListID,
				Title:      "Unit Test Create Todo",
				CreatorID:  creatorID,
				DueDate:    func() *string { s := "01/05/2026"; return &s }(),
			},
			wantErr: true,
			errCode: entity.ErrInvalidParameter,
		},

		// Error path — gateway DB error
		"error: gateway db error": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("db connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockCommands: mock.NewMockTodoCommandsGateway(ctrl)}
			tt.prepare(f)

			sut := usecase.NewTodoCreator(f.mockCommands)
			got, err := sut.Create(ctx, tt.input)

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
