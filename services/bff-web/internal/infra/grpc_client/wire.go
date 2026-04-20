package grpc_client

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewUserGRPCConn,
	NewUserServiceClient,
	NewAuthServiceClient,
)
