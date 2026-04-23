package grpc_client

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewTodoGateway,
	NewAuthGateway,
	NewUserGateway,
)
