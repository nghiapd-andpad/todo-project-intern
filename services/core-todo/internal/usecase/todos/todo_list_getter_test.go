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

func TestTodoListGetter_Get(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		now        = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todoListID = entity.TodoListID(1)
		todoList   = &entity.TodoList{
			ID:        todoListID,
			Name:      "Work Tasks",
			OwnerID:   entity.UserID(1),
			CreatedAt: now,
			UpdatedAt: now,
		}
	)

	type fields struct {
		mockQueries *mock.MockTodoListQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoListGetter
		expected *output.TodoListGetter
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: todo list found": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(todoList, nil)
			},
			input:    &input.TodoListGetter{ID: todoListID},
			expected: &output.TodoListGetter{TodoList: todoList},
		},

		"error: not found — infra returns nil,nil": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), entity.TodoListID(999)).
					Return(nil, nil)
			},
			input:   &input.TodoListGetter{ID: entity.TodoListID(999)},
			wantErr: true,
			errCode: entity.ErrNotFound,
		},
		"error: db error": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("connection refused"))
			},
			input:   &input.TodoListGetter{ID: todoListID},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockQueries: mock.NewMockTodoListQueriesGateway(ctrl)}
			tt.prepare(f)

			sut := todos.NewTodoListGetter(f.mockQueries)
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
