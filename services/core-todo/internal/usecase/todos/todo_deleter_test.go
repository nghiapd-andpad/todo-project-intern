package todos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
)

func TestTodoDeleter_Delete(t *testing.T) {
	t.Parallel()

	var (
		ctx    = context.Background()
		todoID = entity.TodoID(1)
		todo   = &entity.Todo{ID: todoID, Title: "Unit Test Delete Todo"}
	)

	type fields struct {
		mockCommands *mock.MockTodoCommandsGateway
		mockQueries  *mock.MockTodoQueriesGateway
	}

	tests := map[string]struct {
		prepare func(f *fields)
		input   *input.TodoDeleter
		wantErr bool
		errCode entity.ErrorCode
	}{
		"success: todo exists and deleted": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoID).
						Return(todo, nil),
					f.mockCommands.EXPECT().
						Delete(gomock.Any(), todoID).
						Return(nil),
				)
			},
			input:   &input.TodoDeleter{ID: todoID},
			wantErr: false,
		},

		"error: todo not found — delete not called": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoID).
					Return(nil, nil)
			},
			input:   &input.TodoDeleter{ID: todoID},
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: todo exists and delete db error": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoID).
						Return(todo, nil),
					f.mockCommands.EXPECT().
						Delete(gomock.Any(), todoID).
						Return(fmt.Errorf("db error")),
				)
			},
			input:   &input.TodoDeleter{ID: todoID},
			wantErr: true,
		},

		"error: get db error - delete not called": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoID).
					Return(nil, fmt.Errorf("connection lost"))
			},
			input:   &input.TodoDeleter{ID: todoID},
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

			sut := todos.NewTodoDeleter(f.mockCommands, f.mockQueries)
			_, err := sut.Delete(ctx, tt.input)

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
		})
	}
}
