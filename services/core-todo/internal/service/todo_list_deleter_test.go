package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/mock"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func TestTodoListDeleter_Delete(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		todoListID = entity.TodoListID(1)
	)

	type fields struct {
		mockCommands *mock.MockTodoListCommandsGateway
	}

	tests := map[string]struct {
		prepare func(f *fields)
		input   *input.TodoListDeleter
		wantErr bool
	}{
		"success: delete todo list": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Delete(gomock.Any(), todoListID).
					Return(nil)
			},
			input:   &input.TodoListDeleter{ID: todoListID},
			wantErr: false,
		},
		"error: db error": {
			prepare: func(f *fields) {
				f.mockCommands.EXPECT().
					Delete(gomock.Any(), todoListID).
					Return(fmt.Errorf("db error"))
			},
			input:   &input.TodoListDeleter{ID: todoListID},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			f := &fields{mockCommands: mock.NewMockTodoListCommandsGateway(ctrl)}
			tt.prepare(f)

			sut := service.NewTodoListDeleter(f.mockCommands)
			_, err := sut.Delete(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
