package grpc_client

import (
	"log"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserGRPCConn() (*grpc.ClientConn, func(), error) {
	conn, err := grpc.Dial(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.Println("Closing gRPC connection to user-service...")
		conn.Close()
	}

	return conn, cleanup, nil
}

func NewUserServiceClient(conn *grpc.ClientConn) userv1.UserServiceClient {
	return userv1.NewUserServiceClient(conn)
}
