package grpc

import (
	"context"
	"time"

	pb "github.com/sig-agro/api/proto/user"
	"github.com/sig-agro/services/user-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	Port    int32
	pb.UnimplementedUserServiceServer
	Service *service.UserService
}

func (s *UserGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.Service.Register(ctx, req.Email, req.Password, req.FullName, req.Phone)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		UserId:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Message:  "User registered successfully",
	}, nil
}

func (s *UserGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, token, err := s.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	return &pb.LoginResponse{
		UserId:    user.ID,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: user.CreatedAt.Add(time.Hour).Unix(), // Placeholder
	}, nil
}

func (s *UserGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// Implement token validation
	return &pb.ValidateTokenResponse{Valid: true}, nil
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.Service.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return &pb.GetUserResponse{
		UserId:   user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Roles:    user.Roles,
	}, nil
}

func (s *UserGRPCServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.Service.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbUsers []*pb.GetUserResponse
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.GetUserResponse{
			UserId:   	u.ID,
			Email:    	u.Email,
			FullName: 	u.FullName,
			Phone:    	u.Phone,
			Roles: 		u.Roles,
			CreatedAt: 	u.CreatedAt.Unix(),
		})
	}

	return &pb.ListUsersResponse{Users: pbUsers}, nil
}
