package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	todov1 "github.com/nghiaphunng18/todos/gen/todo/v1"
	"github.com/nghiaphunng18/todos/internal/domain/entity"
	"github.com/nghiaphunng18/todos/internal/handler/grpc/mapper"
	"github.com/nghiaphunng18/todos/internal/usecase/todos"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoHandler struct {
	todov1.UnimplementedTodosServiceServer
	TodoCreator    todos.TodoCreator
	TodoGetter     todos.TodoGetter
	TodoListReader todos.TodoLister
	TodoUpdater    todos.TodoUpdater
	TodoDeleter    todos.TodoDeleter
}

func NewTodoHandler(
	creator todos.TodoCreator,
	getter todos.TodoGetter,
	listReader todos.TodoLister,
	updater todos.TodoUpdater,
	deleter todos.TodoDeleter,
) *TodoHandler {
	return &TodoHandler{
		TodoCreator:    creator,
		TodoGetter:     getter,
		TodoListReader: listReader,
		TodoUpdater:    updater,
		TodoDeleter:    deleter,
	}
}

// Get a Todo by Resource Name: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
func (h *TodoHandler) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.Todo, error) {
	// Parse
	todoID, err := parseTodoID(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build Input
	in := &input.TodoGetter{
		ID: entity.TodoID(todoID),
	}

	// Execute
	out, err := h.TodoGetter.Get(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	// Map
	return mapper.TodoToPb(out.Todo), nil
}

// Get List of Todos by Parent Resource Name: users/{u_id}/todo-lists/{l_id}
func (h *TodoHandler) ListTodos(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	// Validate Parent
	if req.GetParent() == "" {
		return nil, status.Error(codes.InvalidArgument, "parent is required")
	}

	// Build Input
	in := &input.TodoLister{
		Parent: req.GetParent(),
	}

	// Execute
	out, err := h.TodoListReader.List(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list todos: %v", err)
	}

	// Map List
	pbTodos := make([]*todov1.Todo, len(out.Todos))
	for i, t := range out.Todos {
		pbTodos[i] = mapper.TodoToPb(t)
	}

	return &todov1.ListTodosResponse{
		Todos: pbTodos,
	}, nil
}

// Create a new Todo under Parent Resource Name: users/{u_id}/todo-lists/{l_id}
func (h *TodoHandler) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.Todo, error) {
	// Validate Parent
	if req.GetParent() == "" {
		return nil, status.Error(codes.InvalidArgument, "parent is required")
	}

	// Build Input
	in := &input.TodoCreator{
		Title:       req.GetTitle(),
		Description: &req.Description,
	}

	// Execute
	out, err := h.TodoCreator.Create(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create todo: %v", err)
	}

	// Map
	return mapper.TodoToPb(out.Todo), nil
}

// Update a Todo by Resource Name: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
func (h *TodoHandler) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.Todo, error) {
	if req.GetTodo() == nil {
		return nil, status.Error(codes.InvalidArgument, "todo resource is required")
	}

	// Parse ID từ field 'name' của object Todo gửi lên
	todoID, err := parseTodoID(req.GetTodo().GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build Input
	title := req.GetTodo().GetTitle()
	desc := req.GetTodo().GetDescription()

	in := &input.TodoUpdater{
		ID:          entity.TodoID(todoID),
		Title:       &title,
		Description: &desc,
	}

	// Execute
	out, err := h.TodoUpdater.Update(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update todo: %v", err)
	}

	// Map
	return mapper.TodoToPb(out.Todo), nil
}

// Delete a Todo by Resource Name: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
func (h *TodoHandler) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*emptypb.Empty, error) {
	// Parse
	todoID, err := parseTodoID(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build Input
	in := &input.TodoDeleter{
		ID: entity.TodoID(todoID),
	}

	// Execute
	_, err = h.TodoDeleter.Delete(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete todo: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Helper function to parse Todo ID from Resource Name
func parseTodoID(name string) (int64, error) {
	// Chuẩn: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
	parts := strings.Split(name, "/")
	if len(parts) != 6 || parts[4] != "todos" {
		return 0, fmt.Errorf("invalid resource name: %s. Expected format: users/{u_id}/todo-lists/{l_id}/todos/{t_id}", name)
	}

	id, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID in resource name: %v", err)
	}
	return id, nil
}
