// Package grpcclient provides gRPC client implementations for the core service, allowing the BFF to communicate with the core service over gRPC.
package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpcclient/mapper"
)

type TodoGateway struct {
	client todov1.TodosServiceClient
}

var _ gateway.TodoGateway = (*TodoGateway)(nil)

func NewTodoGateway(cfg *config.Config) (*TodoGateway, func(), error) {
	conn, err := grpc.Dial(
		cfg.TodoServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial todo service: %w", err)
	}
	client := todov1.NewTodosServiceClient(conn)
	return &TodoGateway{client: client}, func() { conn.Close() }, nil
}

func (g *TodoGateway) GetTodoList(ctx context.Context, name string) (*entity.TodoList, error) {
	resp, err := g.client.GetTodoList(ctx, &todov1.GetTodoListRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return mapper.TodoListFromPb(resp.TodoList), nil
}

func (g *TodoGateway) ListTodoLists(ctx context.Context, parent string, opts input.ListTodoListsOptions) (*output.TodoListPage, error) {
	req := &todov1.ListTodoListsRequest{
		Parent:   parent,
		PageSize: int32(opts.Limit),
		Offset:   int32(opts.Offset),
	}
	if opts.NameSearch != nil {
		req.NameSearch = *opts.NameSearch
	}

	resp, err := g.client.ListTodoLists(ctx, req)
	if err != nil {
		return nil, err
	}

	lists := make([]*entity.TodoList, len(resp.TodoLists))
	for i, tl := range resp.TodoLists {
		lists[i] = mapper.TodoListFromPb(tl)
	}

	return &output.TodoListPage{TodoLists: lists, Total: resp.Total}, nil
}

func (g *TodoGateway) CreateTodoList(ctx context.Context, input input.CreateTodoListInput) (*entity.TodoList, error) {
	resp, err := g.client.CreateTodoList(ctx, &todov1.CreateTodoListRequest{
		Parent:      input.Parent,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, err
	}

	return mapper.TodoListFromPb(resp.TodoList), nil
}

func (g *TodoGateway) UpdateTodoList(ctx context.Context, input input.UpdateTodoListInput) (*entity.TodoList, error) {
	req := &todov1.UpdateTodoListRequest{
		TodoList:   &todov1.TodoList{Name: input.Name},
		UpdateMask: &fieldmaskpb.FieldMask{},
	}
	if input.DisplayName != nil {
		req.TodoList.DisplayName = *input.DisplayName
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "display_name")
	}

	resp, err := g.client.UpdateTodoList(ctx, req)
	if err != nil {
		return nil, err
	}

	return mapper.TodoListFromPb(resp.TodoList), nil
}

func (g *TodoGateway) DeleteTodoList(ctx context.Context, name string) error {
	_, err := g.client.DeleteTodoList(ctx, &todov1.DeleteTodoListRequest{Name: name})
	return err
}

func (g *TodoGateway) GetTodo(ctx context.Context, name string) (*entity.Todo, error) {
	resp, err := g.client.GetTodo(ctx, &todov1.GetTodoRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return mapper.TodoFromPb(resp.Todo), nil
}

func (g *TodoGateway) ListTodos(ctx context.Context, parent string, opts input.ListTodosOptions) (*output.TodoPage, error) {
	req := &todov1.ListTodosRequest{
		Parent:   parent,
		PageSize: int32(opts.Limit),
		Offset:   int32(opts.Offset),
	}
	if opts.Status != nil {
		req.StatusFilter = mapper.TodoStatusToPb(*opts.Status)
	}
	if opts.Priority != nil {
		req.PriorityFilter = mapper.PriorityToPb(*opts.Priority)
	}
	if opts.TitleSearch != nil {
		req.TitleSearch = *opts.TitleSearch
	}
	resp, err := g.client.ListTodos(ctx, req)
	if err != nil {
		return nil, err
	}

	todos := make([]*entity.Todo, len(resp.Todos))
	for i, t := range resp.Todos {
		todos[i] = mapper.TodoFromPb(t)
	}

	return &output.TodoPage{Todos: todos, Total: resp.Total}, nil
}

func (g *TodoGateway) CreateTodo(ctx context.Context, input input.CreateTodoInput) (*entity.Todo, error) {
	req := &todov1.CreateTodoRequest{
		Parent: input.Parent,
		Title:  input.Title,
	}
	if input.Description != nil {
		req.Description = *input.Description
	}
	if input.Priority != nil {
		req.Priority = mapper.PriorityToPb(*input.Priority)
	}
	if input.DueDate != nil {
		req.DueDate = *input.DueDate
	}
	if input.AssigneeID != nil {
		var aID int64
		fmt.Sscanf(*input.AssigneeID, "%d", &aID)
		req.AssigneeId = aID
	}

	resp, err := g.client.CreateTodo(ctx, req)
	if err != nil {
		return nil, err
	}

	return mapper.TodoFromPb(resp.Todo), nil
}

func (g *TodoGateway) UpdateTodo(ctx context.Context, input input.UpdateTodoInput) (*entity.Todo, error) {
	req := &todov1.UpdateTodoRequest{
		Todo:       &todov1.Todo{Name: input.Name},
		UpdateMask: &fieldmaskpb.FieldMask{},
	}
	if input.Title != nil {
		req.Todo.Title = *input.Title
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "title")
	}
	if input.Description != nil {
		req.Todo.Description = *input.Description
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "description")
	}
	if input.Status != nil {
		req.Todo.Status = mapper.TodoStatusToPb(*input.Status)
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "status")
	}
	if input.Priority != nil {
		req.Todo.Priority = mapper.PriorityToPb(*input.Priority)
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "priority")
	}
	if input.DueDate != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "due_date")
	}
	if input.AssigneeID != nil {
		var aID int64
		fmt.Sscanf(*input.AssigneeID, "%d", &aID)
		req.Todo.AssigneeId = aID
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "assignee_id")
	}

	resp, err := g.client.UpdateTodo(ctx, req)
	if err != nil {
		return nil, err
	}

	return mapper.TodoFromPb(resp.Todo), nil
}

func (g *TodoGateway) DeleteTodo(ctx context.Context, name string) error {
	_, err := g.client.DeleteTodo(ctx, &todov1.DeleteTodoRequest{Name: name})
	return err
}
