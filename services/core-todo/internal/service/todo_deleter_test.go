package service_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"

// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
// 	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
// )

// func TestTodoDeleter_Delete(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		todoID      = entity.TodoID(1)
// 		todoListID  = entity.TodoListID(2)
// 		requesterID = entity.UserID(10)
// 		ownerID     = entity.UserID(10)

// 		todo = &entity.Todo{
// 			ID:         todoID,
// 			TodoListID: todoListID,
// 			Title:      "Unit Test Delete Todo",
// 		}

// 		todoList = &entity.TodoList{
// 			ID:      todoListID,
// 			Name:    "Work Tasks",
// 			OwnerID: ownerID,
// 		}

// 		validInput = &input.TodoDeleter{
// 			TodoID:      todoID,
// 			TodoListID:  todoListID,
// 			RequesterID: requesterID,
// 		}
// 	)

// 	type fields struct {
// 		mockTodoListQueries *mock.MockTodoListQueriesGateway
// 		mockTodoQueries     *mock.MockTodoQueriesGateway
// 		mockTodoCommands    *mock.MockTodoCommandsGateway
// 	}

// 	tests := map[string]struct {
// 		prepare func(f *fields, ctx context.Context)
// 		input   *input.TodoDeleter
// 		wantErr bool
// 		errCode entity.ErrorCode
// 	}{
// 		"success: todo exists, todo list exists, requester is owner": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockTodoQueries.EXPECT().
// 						Get(ctx, todoID, todoListID).
// 						Return(todo, nil),

// 					f.mockTodoListQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(todoList, nil),

// 					f.mockTodoCommands.EXPECT().
// 						Delete(ctx, todoID).
// 						Return(nil),
// 				)
// 			},
// 			input: validInput,
// 		},

// 		"error: todo query gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockTodoQueries.EXPECT().
// 					Get(ctx, todoID, todoListID).
// 					Return(nil, fmt.Errorf("connection lost"))
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 		},

// 		"error: todo not found": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				f.mockTodoQueries.EXPECT().
// 					Get(ctx, todoID, todoListID).
// 					Return(nil, nil)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 			errCode: entity.ErrNotFound,
// 		},

// 		"error: todo list query gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockTodoQueries.EXPECT().
// 						Get(ctx, todoID, todoListID).
// 						Return(todo, nil),

// 					f.mockTodoListQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(nil, fmt.Errorf("connection lost")),
// 				)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 		},

// 		"error: todo list not found": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockTodoQueries.EXPECT().
// 						Get(ctx, todoID, todoListID).
// 						Return(todo, nil),

// 					f.mockTodoListQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(nil, nil),
// 				)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 			errCode: entity.ErrNotFound,
// 		},

// 		"error: requester is not todo list owner": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockTodoQueries.EXPECT().
// 						Get(ctx, todoID, todoListID).
// 						Return(todo, nil),

// 					f.mockTodoListQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(&entity.TodoList{
// 							ID:      todoListID,
// 							Name:    "Work Tasks",
// 							OwnerID: entity.UserID(999),
// 						}, nil),
// 				)
// 			},
// 			input:   validInput,
// 			wantErr: true,
// 			errCode: entity.ErrAuthZ,
// 		},

// 		"error: delete command gateway error": {
// 			prepare: func(f *fields, ctx context.Context) {
// 				gomock.InOrder(
// 					f.mockTodoQueries.EXPECT().
// 						Get(ctx, todoID, todoListID).
// 						Return(todo, nil),

// 					f.mockTodoListQueries.EXPECT().
// 						Get(ctx, todoListID).
// 						Return(todoList, nil),

// 					f.mockTodoCommands.EXPECT().
// 						Delete(ctx, todoID).
// 						Return(fmt.Errorf("db error")),
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
// 				mockTodoListQueries: mock.NewMockTodoListQueriesGateway(ctrl),
// 				mockTodoQueries:     mock.NewMockTodoQueriesGateway(ctrl),
// 				mockTodoCommands:    mock.NewMockTodoCommandsGateway(ctrl),
// 			}

// 			tt.prepare(f, ctx)

// 			sut := service.NewTodoDeleter(
// 				f.mockTodoListQueries,
// 				f.mockTodoQueries,
// 				f.mockTodoCommands,
// 			)

// 			got, err := sut.Delete(ctx, tt.input)

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
// 			assert.NotNil(t, got)
// 		})
// 	}
// }
