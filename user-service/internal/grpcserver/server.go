package grpcserver

import (
	"context"
	"database/sql"
	"errors"

	pb "common/proto/user"

	"user-service/internal/database/repository"

	repositoryErrors "common/network/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var blankCreateUserResponse = pb.UserCreatedResponse{
	Id: 0,
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserServer(repo *repository.UserRepository) *UserServer {
	return &UserServer{repo: repo}
}

func (s *UserServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		Id:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserCreatedResponse, error) {
	id, err := s.repo.CreateUser(ctx, req)

	if err != nil {
		if errors.Is(err, repositoryErrors.ErrPasswordsMismatch) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserCreatedResponse{Id: id}, nil
}
