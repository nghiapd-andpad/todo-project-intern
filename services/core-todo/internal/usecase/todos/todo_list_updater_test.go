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

func TestTodoListUpdater_Update(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		now        = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		todoListID = entity.TodoListID(1)
		ownerID    = entity.UserID(1)
		newName    = "New Name"
	)

	freshList := func() *entity.TodoList {
		return &entity.TodoList{
			ID:        todoListID,
			Name:      "Old Name",
			OwnerID:   ownerID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	type fields struct {
		mockCommands *mock.MockTodoListCommandsGateway
		mockQueries  *mock.MockTodoListQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoListUpdater
		expected *output.TodoListUpdater
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: update name": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(freshList(), nil),
					f.mockCommands.EXPECT().
						Update(gomock.Any(), &entity.TodoList{
							ID:        todoListID,
							Name:      "New Name",
							OwnerID:   ownerID,
							CreatedAt: now,
							UpdatedAt: now,
						}).
						Return(&entity.TodoList{
							ID:        todoListID,
							Name:      "New Name",
							OwnerID:   ownerID,
							CreatedAt: now,
							UpdatedAt: now,
						}, nil),
				)
			},
			input: &input.TodoListUpdater{ID: todoListID, Name: &newName},
			expected: &output.TodoListUpdater{TodoList: &entity.TodoList{
				ID:        todoListID,
				Name:      "New Name",
				OwnerID:   ownerID,
				CreatedAt: now,
				UpdatedAt: now,
			}},
			wantErr: false,
		},
		"success: nil name — name unchanged": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(freshList(), nil),
					f.mockCommands.EXPECT().
						Update(gomock.Any(), freshList()).
						Return(freshList(), nil),
				)
			},
			input:    &input.TodoListUpdater{ID: todoListID, Name: nil},
			expected: &output.TodoListUpdater{TodoList: freshList()},
			wantErr:  false,
		},
		"error: not found": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(nil, nil)
			},
			input:   &input.TodoListUpdater{ID: todoListID, Name: &newName},
			wantErr: true,
			errCode: entity.ErrNotFound,
		},
		"error: get db error": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(nil, fmt.Errorf("connection lost"))
			},
			input:   &input.TodoListUpdater{ID: todoListID, Name: &newName},
			wantErr: true,
		},
		"error: update db error": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(freshList(), nil),
					f.mockCommands.EXPECT().
						Update(gomock.Any(), gomock.Any()).
						Return(nil, fmt.Errorf("db error")),
				)
			},
			input:   &input.TodoListUpdater{ID: todoListID, Name: &newName},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{
				mockCommands: mock.NewMockTodoListCommandsGateway(ctrl),
				mockQueries:  mock.NewMockTodoListQueriesGateway(ctrl),
			}
			tt.prepare(f)

			sut := todos.NewTodoListUpdater(f.mockCommands, f.mockQueries)
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
