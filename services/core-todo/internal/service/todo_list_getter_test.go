package service_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/google/go-cmp/cmp"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"

// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
// )

// func TestTodoListGetter_Get(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		now         = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
// 		todoListID  = entity.TodoListID(1)
// 		requesterID = entity.UserID(10)

// 		todoList = &entity.TodoList{
// 			ID:        todoListID,
// 			Name:      "Work Tasks",
// 			OwnerID:   requesterID,
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		}

// 		validInput = &input.TodoListGetter{
// 			TodoListID:  todoListID,
// 			RequesterID: requesterID,
// 		}
// 	)

// 	type fields struct {
// 		mockQueries *mock.MockTodoListQueriesGateway
// 	}

// 	tests := map[string]struct {
// 		prepare  func(f *fields, ctx context.Context)
// 		input    *input.TodoListGetter
// 		expected *output.TodoListGetter
// 		wantErr  bool
// 		errCode  entity.ErrorCode
// 	}{
// 		"success: todo list found and requester is owner": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					Get(ctx, todoListID).
// 					Return(todoList, nil)
// 			},
// 			input:    validInput,
// 			expected: &output.TodoListGetter{TodoList: todoList},
// 		},

// 		"error: not found": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					Get(ctx, entity.TodoListID(999)).
// 					Return(nil, nil)
// 			},
// 			input: &input.TodoListGetter{
// 				TodoListID:  entity.TodoListID(999),
// 				RequesterID: requesterID,
// 			},
// 			wantErr: true,
// 			errCode: entity.ErrNotFound,
// 		},

// 		"error: query gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					Get(ctx, todoListID).
// 					Return(nil, fmt.Errorf("connection refused"))
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 		},

// 		"error: requester is not owner": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					Get(ctx, todoListID).
// 					Return(&entity.TodoList{
// 						ID:      todoListID,
// 						Name:    "Work Tasks",
// 						OwnerID: entity.UserID(999),
// 					}, nil)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 			errCode: entity.ErrAuthZ,
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

// 			sut := service.NewTodoListGetter(f.mockQueries)

// 			got, err := sut.Get(ctx, tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)

// 				if tt.errCode != "" {
// 					var appErr *entity.AppError
// 					assert.ErrorAs(t, err, &appErr)
// 					assert.Equal(t, tt.errCode, appErr.Code)
// 				}

// 				return
// 			}

// 			assert.NoError(t, err)

// 			if diff := cmp.Diff(tt.expected, got); diff != "" {
// 				t.Errorf("mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
