package grpcclients

import (
	"context"

	pb "common/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(target string) (*UserClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UserClient{
		client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) GetUserByEmail(ctx context.Context, email string) (*pb.UserResponse, error) {
	return c.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: email})
}
