package todos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

func TestTodoListLister_List(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		ownerID    = entity.UserID(1)
		nameSearch = "Work"

		todoLists = []*entity.TodoList{
			{ID: entity.TodoListID(1), Name: "Work Tasks", OwnerID: ownerID},
			{ID: entity.TodoListID(2), Name: "Personal", OwnerID: ownerID},
		}
	)

	type fields struct {
		mockQueries *mock.MockTodoListQueriesGateway
	}

	tests := map[string]struct {
		prepare  func(f *fields)
		input    *input.TodoListLister
		expected *output.TodoListLister
		wantErr  bool
	}{
		"success: list all": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gateway.ListTodoListsOptions{
						OwnerID: &ownerID,
						Limit:   20,
					}).
					Return(todoLists, int64(2), nil)
			},
			input: &input.TodoListLister{
				Opts: gateway.ListTodoListsOptions{
					OwnerID: &ownerID,
					Limit:   20,
				},
			},
			expected: &output.TodoListLister{TodoLists: todoLists, Total: 2},
		},
		"success: empty": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return([]*entity.TodoList{}, int64(0), nil)
			},
			input:    &input.TodoListLister{Opts: gateway.ListTodoListsOptions{}},
			expected: &output.TodoListLister{TodoLists: []*entity.TodoList{}, Total: 0},
		},
		"success: with name search": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gateway.ListTodoListsOptions{
						OwnerID:    &ownerID,
						NameSearch: &nameSearch,
						Limit:      10,
					}).
					Return(todoLists[:1], int64(1), nil)
			},
			input: &input.TodoListLister{
				Opts: gateway.ListTodoListsOptions{
					OwnerID:    &ownerID,
					NameSearch: &nameSearch,
					Limit:      10,
				},
			},
			expected: &output.TodoListLister{TodoLists: todoLists[:1], Total: 1},
		},
		"error: db error": {
			prepare: func(f *fields) {
				f.mockQueries.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), fmt.Errorf("db error"))
			},
			input:   &input.TodoListLister{Opts: gateway.ListTodoListsOptions{}},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockQueries: mock.NewMockTodoListQueriesGateway(ctrl)}
			tt.prepare(f)

			sut := todos.NewTodoListLister(f.mockQueries)
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
