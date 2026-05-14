package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/pagination"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

func TestTodoLister_List(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		todoListID  = entity.TodoListID(1)
		requesterID = entity.UserID(10)
		ownerID     = entity.UserID(10)
		assigneeID  = entity.UserID(20)
		pending     = entity.TodoStatusPending
		high        = entity.PriorityHigh
		titleSearch = "Todo"

		todoList = &entity.TodoList{
			ID:      todoListID,
			Name:    "Work Tasks",
			OwnerID: ownerID,
		}

		todos = []*entity.Todo{
			{ID: entity.TodoID(1), TodoListID: todoListID, Title: "Todo 1"},
			{ID: entity.TodoID(2), TodoListID: todoListID, Title: "Todo 2"},
		}

		validInput = &input.TodoLister{
			TodoListID:  todoListID,
			RequesterID: requesterID,
			Limit:       20,
		}
	)

	type fields struct {
		mockTodoListQueries *mock.MockTodoListQueriesGateway
		mockTodoQueries     *mock.MockTodoQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoLister
		expected *output.TodoLister
		wantErr  bool
		errCode  entity.ErrorCode
	}{
		"success: owner lists all todos in todo list": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockTodoListQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(todoList, nil),

					f.mockTodoQueries.EXPECT().
						List(gomock.Any(), &gatewayinput.ListTodosOptions{
							TodoListID:   todoListID,
							AssigneeOnly: nil,
							Offset:       0,
							Limit:        20,
						}).
						Return(todos, int64(2), nil),
				)
			},
			input: validInput,
			expected: &output.TodoLister{
				Page: pagination.New(todos, int64(2), 0, 20),
			},
		},

		"success: non-owner lists only assigned todos": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockTodoListQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(&entity.TodoList{
							ID:      todoListID,
							Name:    "Work Tasks",
							OwnerID: ownerID,
						}, nil),

					f.mockTodoQueries.EXPECT().
						List(gomock.Any(), &gatewayinput.ListTodosOptions{
							TodoListID:   todoListID,
							AssigneeOnly: &assigneeID,
							Offset:       0,
							Limit:        20,
						}).
						Return(todos[:1], int64(1), nil),
				)
			},
			input: &input.TodoLister{
				TodoListID:  todoListID,
				RequesterID: assigneeID,
				Limit:       20,
			},
			expected: &output.TodoLister{
				Page: pagination.New(todos[:1], int64(1), 0, 20),
			},
		},

		"success: list with filters": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockTodoListQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(todoList, nil),

					f.mockTodoQueries.EXPECT().
						List(gomock.Any(), &gatewayinput.ListTodosOptions{
							TodoListID:   todoListID,
							Status:       &pending,
							Priority:     &high,
							TitleSearch:  &titleSearch,
							AssigneeOnly: nil,
							Offset:       5,
							Limit:        10,
						}).
						Return(todos[:1], int64(1), nil),
				)
			},
			input: &input.TodoLister{
				TodoListID:  todoListID,
				RequesterID: requesterID,
				Status:      &pending,
				Priority:    &high,
				TitleSearch: &titleSearch,
				Offset:      5,
				Limit:       10,
			},
			expected: &output.TodoLister{
				Page: pagination.New(todos[:1], int64(1), 5, 10),
			},
		},

		"success: empty result": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockTodoListQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(todoList, nil),

					f.mockTodoQueries.EXPECT().
						List(gomock.Any(), &gatewayinput.ListTodosOptions{
							TodoListID:   todoListID,
							AssigneeOnly: nil,
							Offset:       0,
							Limit:        20,
						}).
						Return([]*entity.Todo{}, int64(0), nil),
				)
			},
			input: validInput,
			expected: &output.TodoLister{
				Page: pagination.New([]*entity.Todo{}, int64(0), 0, 20),
			},
		},

		"error: todo list query gateway error": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(nil, fmt.Errorf("connection lost"))
			},
			input:   validInput,
			wantErr: true,
		},

		"error: todo list not found": {
			prepare: func(f *fields) {
				f.mockTodoListQueries.EXPECT().
					Get(gomock.Any(), todoListID).
					Return(nil, nil)
			},
			input:   validInput,
			wantErr: true,
			errCode: entity.ErrNotFound,
		},

		"error: todo query gateway error": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.mockTodoListQueries.EXPECT().
						Get(gomock.Any(), todoListID).
						Return(todoList, nil),

					f.mockTodoQueries.EXPECT().
						List(gomock.Any(), gomock.Any()).
						Return(nil, int64(0), fmt.Errorf("db error")),
				)
			},
			input:   validInput,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			f := &fields{
				mockTodoListQueries: mock.NewMockTodoListQueriesGateway(ctrl),
				mockTodoQueries:     mock.NewMockTodoQueriesGateway(ctrl),
			}

			tt.prepare(f)

			sut := service.NewTodoLister(
				f.mockTodoListQueries,
				f.mockTodoQueries,
			)

			got, err := sut.List(ctx, tt.input)

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
