package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoLister struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoLister(todoQueriesGateway gateway.TodoQueriesGateway) *TodoLister {
	return &TodoLister{todoQueriesGateway: todoQueriesGateway}
}

func (s *TodoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	if in == nil {
		in = &input.TodoLister{}
	}

	s.applyDefaults(&in.Opts)
	if err := s.validate(&in.Opts); err != nil {
		return nil, err
	}

	todos, total, err := s.todoQueriesGateway.List(ctx, &in.Opts)
	if err != nil {
		return nil, fmt.Errorf("TodoLister.List: %w", err)
	}

	return &output.TodoLister{
		Todos: todos,
		Total: total,
	}, nil
}

func (s *TodoLister) applyDefaults(o *gatewayinput.ListTodosOptions) {

	if o.Limit == 0 {
		o.Limit = 20
	}
}

func (s *TodoLister) validate(o *gatewayinput.ListTodosOptions) error {
	if o.Limit < 0 {
		return entity.NewInvalidParameter("limit must be non-negative")
	}
	if o.Limit > 100 {
		return entity.NewInvalidParameter("limit must not exceed 100")
	}
	if o.Offset < 0 {
		return entity.NewInvalidParameter("offset must be non-negative")
	}
	if o.TitleSearch != nil && len(*o.TitleSearch) > 255 {
		return entity.NewInvalidParameter("title_search too long")
	}
	return nil
}
