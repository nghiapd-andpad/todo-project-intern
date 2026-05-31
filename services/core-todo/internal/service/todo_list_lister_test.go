package service_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"

// 	"github.com/nghiapd-andpad/todo-project-intern/pkg/pagination"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
// 	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
// )

// func TestTodoListLister_List(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		requesterID = entity.UserID(1)
// 		nameSearch  = "Work"

// 		todoLists = []*entity.TodoList{
// 			{ID: entity.TodoListID(1), Name: "Work Tasks", OwnerID: requesterID},
// 			{ID: entity.TodoListID(2), Name: "Personal", OwnerID: requesterID},
// 		}
// 	)

// 	type fields struct {
// 		mockQueries *mock.MockTodoListQueriesGateway
// 	}

// 	tests := map[string]struct {
// 		prepare  func(f *fields, ctx context.Context)
// 		input    *input.TodoListLister
// 		expected *output.TodoListLister
// 		wantErr  bool
// 	}{
// 		"success: default filter lists owned or assigned todo lists": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					List(ctx, &gatewayinput.ListTodoListsOptions{
// 						OwnerID:    &requesterID,
// 						AssigneeID: &requesterID,
// 						Offset:     0,
// 						Limit:      20,
// 					}).
// 					Return(todoLists, int64(2), nil)
// 			},
// 			input: &input.TodoListLister{
// 				RequesterID: requesterID,
// 				Offset:      0,
// 				Limit:       20,
// 			},
// 			expected: &output.TodoListLister{
// 				Page: pagination.New(todoLists, int64(2), 0, 20),
// 			},
// 		},

// 		"success: owned filter": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					List(ctx, &gatewayinput.ListTodoListsOptions{
// 						OwnerID: &requesterID,
// 						Offset:  0,
// 						Limit:   20,
// 					}).
// 					Return(todoLists, int64(2), nil)
// 			},
// 			input: &input.TodoListLister{
// 				RequesterID: requesterID,
// 				Filter:      input.TodoListFilterOwned,
// 				Offset:      0,
// 				Limit:       20,
// 			},
// 			expected: &output.TodoListLister{
// 				Page: pagination.New(todoLists, int64(2), 0, 20),
// 			},
// 		},

// 		"success: assigned filter": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					List(ctx, &gatewayinput.ListTodoListsOptions{
// 						AssigneeID: &requesterID,
// 						Offset:     0,
// 						Limit:      20,
// 					}).
// 					Return(todoLists, int64(2), nil)
// 			},
// 			input: &input.TodoListLister{
// 				RequesterID: requesterID,
// 				Filter:      input.TodoListFilterAssigned,
// 				Offset:      0,
// 				Limit:       20,
// 			},
// 			expected: &output.TodoListLister{
// 				Page: pagination.New(todoLists, int64(2), 0, 20),
// 			},
// 		},

// 		"success: with name search": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					List(ctx, &gatewayinput.ListTodoListsOptions{
// 						OwnerID:    &requesterID,
// 						AssigneeID: &requesterID,
// 						NameSearch: &nameSearch,
// 						Offset:     10,
// 						Limit:      10,
// 					}).
// 					Return(todoLists[:1], int64(1), nil)
// 			},
// 			input: &input.TodoListLister{
// 				RequesterID: requesterID,
// 				NameSearch:  &nameSearch,
// 				Offset:      10,
// 				Limit:       10,
// 			},
// 			expected: &output.TodoListLister{
// 				Page: pagination.New(todoLists[:1], int64(1), 10, 10),
// 			},
// 		},

// 		"success: empty result": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					List(ctx, &gatewayinput.ListTodoListsOptions{
// 						OwnerID:    &requesterID,
// 						AssigneeID: &requesterID,
// 						Offset:     0,
// 						Limit:      20,
// 					}).
// 					Return([]*entity.TodoList{}, int64(0), nil)
// 			},
// 			input: &input.TodoListLister{
// 				RequesterID: requesterID,
// 				Offset:      0,
// 				Limit:       20,
// 			},
// 			expected: &output.TodoListLister{
// 				Page: pagination.New([]*entity.TodoList{}, int64(0), 0, 20),
// 			},
// 		},

// 		"error: query gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					List(ctx, gomock.Any()).
// 					Return(nil, int64(0), fmt.Errorf("db error"))
// 			},
// 			input: &input.TodoListLister{
// 				RequesterID: requesterID,
// 				Offset:      0,
// 				Limit:       20,
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			t.Parallel()

// 			ctx := context.Background()

// 			ctrl := gomock.NewController(t)

// 			f := &fields{
// 				mockQueries: mock.NewMockTodoListQueriesGateway(ctrl),
// 			}

// 			tt.prepare(f, ctx)

// 			sut := service.NewTodoListLister(f.mockQueries)

// 			got, err := sut.List(ctx, tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}

// 			assert.NoError(t, err)

// 			if diff := cmp.Diff(tt.expected, got); diff != "" {
// 				t.Errorf("mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
