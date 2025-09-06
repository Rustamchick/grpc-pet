package handler

import (
	"context"
	Service "grpc-pet/pkg/service"

	grpcpetv1 "github.com/Rustamchick/protobuff/gen/go/pet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerApi struct {
	grpcpetv1.UnimplementedAuthServer
	auth Service.Authentification
}

func Register(gRPC *grpc.Server, auth Service.Authentification) {
	grpcpetv1.RegisterAuthServer(gRPC, &ServerApi{auth: auth})
}

func (s *ServerApi) Login(ctx context.Context, req *grpcpetv1.LoginRequest) (*grpcpetv1.LoginResponse, error) {
	if err := LoginIsValid(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error. auth.Login()") // dont give away internal errors to clients
	}

	return &grpcpetv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerApi) Register(ctx context.Context, req *grpcpetv1.RegisterRequest) (*grpcpetv1.RegisterResponse, error) {
	if err := RegisterIsValid(req); err != nil {
		return nil, err
	}

	user_id, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// to do handle variative errors (Double registration or invalid login password, etc)
		return nil, status.Error(codes.Internal, "Internal error. auth.Register()") // dont give away internal errors to clients
	}

	return &grpcpetv1.RegisterResponse{
		UserId: user_id,
	}, nil
}

func (s *ServerApi) IsAdmin(ctx context.Context, req *grpcpetv1.IsAdminRequest) (*grpcpetv1.IsAdminResponse, error) {
	if err := IsAdminIsValid(req); err != nil {
		return nil, err
	}

	IsAdmin, err := s.auth.IsAdmin(ctx, int(req.GetUserId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error. auth.IsAdmin()") // dont give away internal errors to clients
	}

	return &grpcpetv1.IsAdminResponse{
		IsAdmin: IsAdmin,
	}, nil
}

func LoginIsValid(req *grpcpetv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email is required.")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password is required.")
	}
	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "AppId is required.")
	}
	return nil
}

func RegisterIsValid(req *grpcpetv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email is required.")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password is required.")
	}
	return nil
}

func IsAdminIsValid(req *grpcpetv1.IsAdminRequest) error {
	if req.GetUserId() == 0 {
		return status.Error(codes.InvalidArgument, "Email is required.")
	}
	return nil
}
