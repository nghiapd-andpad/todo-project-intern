package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

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
		ctx        = context.Background()
		todoListID = entity.TodoListID(1)
		todoList   = []*entity.Todo{
			{ID: entity.TodoID(1), Title: "Todo 1"},
			{ID: entity.TodoID(2), Title: "Todo 2"},
		}
		pending = entity.TodoStatusPending
	)

	type fields struct {
		mockQueries *mock.MockTodoQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoLister
		expected *output.TodoLister
		wantErr  bool
	}{
		"success: list todos": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gatewayinput.ListTodosOptions{
						TodoListID: &todoListID,
						Limit:      20,
					}).
					Return(todoList, int64(2), nil)
			},
			input: &input.TodoLister{
				Opts: gatewayinput.ListTodosOptions{
					TodoListID: &todoListID,
					Limit:      20,
				},
			},
			expected: &output.TodoLister{Todos: todoList, Total: 2},
			wantErr:  false,
		},

		"success: empty list": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return([]*entity.Todo{}, int64(0), nil)
			},
			input:    &input.TodoLister{Opts: gatewayinput.ListTodosOptions{}},
			expected: &output.TodoLister{Todos: []*entity.Todo{}, Total: 0},
			wantErr:  false,
		},

		"success: list with status filter": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gatewayinput.ListTodosOptions{
						TodoListID: &todoListID,
						Status:     &pending,
						Limit:      10,
						Offset:     5,
					}).
					Return(todoList[:1], int64(1), nil)
			},
			input: &input.TodoLister{
				Opts: gatewayinput.ListTodosOptions{
					TodoListID: &todoListID,
					Status:     &pending,
					Limit:      10,
					Offset:     5,
				},
			},
			expected: &output.TodoLister{Todos: todoList[:1], Total: 1},
			wantErr:  false,
		},

		"error: db error": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), fmt.Errorf("db error"))
			},
			input:   &input.TodoLister{Opts: gatewayinput.ListTodosOptions{}},
			wantErr: true,
		},

		"error: validation failed (limit too high)": {
			prepare: func(f *fields) {},
			input: &input.TodoLister{
				Opts: gatewayinput.ListTodosOptions{
					Limit: 101,
				},
			},
			wantErr: true,
		},

		"success: list todos with default limit": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), &gatewayinput.ListTodosOptions{
						TodoListID: &todoListID,
						Limit:      20,
					}).
					Return(todoList, int64(2), nil)
			},
			input: &input.TodoLister{
				Opts: gatewayinput.ListTodosOptions{
					Limit: 0,
				},
			},
			expected: &output.TodoLister{Todos: todoList, Total: 2},
			wantErr:  false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockQueries: mock.NewMockTodoQueriesGateway(ctrl)}
			tt.prepare(f)

			sut := service.NewTodoLister(f.mockQueries)
			got, err := sut.List(ctx, tt.input)

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
