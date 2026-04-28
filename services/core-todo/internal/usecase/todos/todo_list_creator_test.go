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
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestTodoListCreator_Create(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		now     = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		ownerID = entity.UserID(1)

		validInput = &input.TodoListCreator{
			Name:    "Task 1",
			OwnerID: ownerID,
		}

		createdEntity = &entity.TodoList{
			ID:        entity.TodoListID(1),
			Name:      "Task 1",
			OwnerID:   ownerID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	)

	type fields struct {
		mockCommands *mock.MockTodoListCommandsGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoListCreator
		expected *output.TodoListCreator
		wantErr  bool
	}{
		"success: create todo list": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&entity.TodoList{})).
					DoAndReturn(func(_ context.Context, tl *entity.TodoList) (*entity.TodoList, error) {
						assert.Equal(t, "Task 1", tl.Name)
						assert.Equal(t, ownerID, tl.OwnerID)
						return createdEntity, nil
					})
			},
			input:    validInput,
			expected: &output.TodoListCreator{TodoList: createdEntity},
		},

		"success: create with empty name": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(&entity.TodoList{})).
					Return(&entity.TodoList{ID: 2, Name: "", OwnerID: ownerID}, nil)
			},
			input: &input.TodoListCreator{
				Name:    "",
				OwnerID: ownerID,
			},
			expected: &output.TodoListCreator{
				TodoList: &entity.TodoList{ID: 2, Name: "", OwnerID: ownerID},
			},
		},

		"error: db error": {
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
			f := &fields{
				mockCommands: mock.NewMockTodoListCommandsGateway(ctrl),
			}
			tt.prepare(f)

			sut := todos.NewTodoListCreator(f.mockCommands)
			got, err := sut.Create(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
