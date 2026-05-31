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

// func TestTodoListUpdater_Update(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		now         = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
// 		todoListID  = entity.TodoListID(1)
// 		requesterID = entity.UserID(10)
// 		newName     = "New Name"
// 	)

// 	freshList := func() *entity.TodoList {
// 		return &entity.TodoList{
// 			ID:        todoListID,
// 			Name:      "Old Name",
// 			OwnerID:   requesterID,
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		}
// 	}

// 	updatedList := func() *entity.TodoList {
// 		return &entity.TodoList{
// 			ID:        todoListID,
// 			Name:      "New Name",
// 			OwnerID:   requesterID,
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		}
// 	}

// 	validInput := &input.TodoListUpdater{
// 		TodoListID:  todoListID,
// 		Name:        &newName,
// 		RequesterID: requesterID,
// 	}

// 	type fields struct {
// 		mockCommands *mock.MockTodoListCommandsGateway
// 		mockQueries  *mock.MockTodoListQueriesGateway
// 	}

// 	tests := map[string]struct {
// 		prepare  func(f *fields, ctx context.Context)
// 		input    *input.TodoListUpdater
// 		expected *output.TodoListUpdater
// 		wantErr  bool
// 		errCode  entity.ErrorCode
// 	}{
// 		"success: update name": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(freshList(), nil),

// 					f.mockCommands.EXPECT().
// 						Update(ctx, updatedList()).
// 						Return(updatedList(), nil),
// 				)
// 			},
// 			input: validInput,
// 			expected: &output.TodoListUpdater{
// 				TodoList: updatedList(),
// 			},
// 		},

// 		"success: nil name keeps original name": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(freshList(), nil),

// 					f.mockCommands.EXPECT().
// 						Update(ctx, freshList()).
// 						Return(freshList(), nil),
// 				)
// 			},
// 			input: &input.TodoListUpdater{
// 				TodoListID:  todoListID,
// 				Name:        nil,
// 				RequesterID: requesterID,
// 			},
// 			expected: &output.TodoListUpdater{
// 				TodoList: freshList(),
// 			},
// 		},

// 		"error: todo list not found": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					Get(ctx, todoListID).
// 					Return(nil, nil)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 			errCode: entity.ErrNotFound,
// 		},

// 		"error: query gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockQueries.EXPECT().
// 					Get(ctx, todoListID).
// 					Return(nil, fmt.Errorf("connection lost"))
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
// 						Name:    "Old Name",
// 						OwnerID: entity.UserID(999),
// 					}, nil)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 			errCode: entity.ErrAuthZ,
// 		},

// 		"error: update command gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(freshList(), nil),

// 					f.mockCommands.EXPECT().
// 						Update(ctx, updatedList()).
// 						Return(nil, fmt.Errorf("db error")),
// 				)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 		},
// 	}

// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			t.Parallel()

// 			ctx := context.Background()

// 			ctrl := gomock.NewController(t)

// 			f := &fields{
// 				mockCommands: mock.NewMockTodoListCommandsGateway(ctrl),
// 				mockQueries:  mock.NewMockTodoListQueriesGateway(ctrl),
// 			}

// 			tt.prepare(f, ctx)

// 			sut := service.NewTodoListUpdater(
// 				f.mockCommands,
// 				f.mockQueries,
// 			)

// 			got, err := sut.Update(ctx, tt.input)

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
