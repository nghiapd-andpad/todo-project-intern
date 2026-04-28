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

func TestTodoGetter_Get(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		now  = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todo = &entity.Todo{
			ID:        entity.TodoID(1),
			Title:     "Unit Test Get Todo",
			Status:    entity.TodoStatusPending,
			Priority:  entity.PriorityMedium,
			CreatorID: entity.UserID(1),
			CreatedAt: now,
			UpdatedAt: now,
		}
	)

	type fields struct {
		mockQueries *mock.MockTodoQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoGetter
		expected *output.TodoGetter
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: todo found": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), entity.TodoID(1)).
					Return(todo, nil)
			},
			input:    &input.TodoGetter{ID: entity.TodoID(1)},
			expected: &output.TodoGetter{Todo: todo},
		},

		"error: todo not found — infra returns nil,nil": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), entity.TodoID(999)).
					Return(nil, nil)
			},

			input:   &input.TodoGetter{ID: entity.TodoID(999)},
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: db error": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("connection lost"))
			},
			input:   &input.TodoGetter{ID: entity.TodoID(1)},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockQueries: mock.NewMockTodoQueriesGateway(ctrl)}
			tt.prepare(f)

			sut := todos.NewTodoGetter(f.mockQueries)
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
