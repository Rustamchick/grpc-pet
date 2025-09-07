package Service

import (
	"context"
	"grpc-pet/pkg/repository"
)

type AuthService struct {
	repos repository.Authentification
}

func NewAuthService(repos repository.Authentification) *AuthService {
	return &AuthService{
		repos: repos,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string, appid int) (string, error) {
	token, err := s.repos.Login(ctx, email, password, appid)

	return token, err
	// "Service return token", status.Errorf(codes.Unimplemented, "Service. Method Login is not implemented")
}

func (s *AuthService) Register(ctx context.Context, email, password string) (int64, error) {
	userid, err := s.repos.Register(ctx, email, password)

	return userid, err
	// -111, status.Errorf(codes.Unimplemented, "Service. Method Register is not implemented")
}

func (s *AuthService) IsAdmin(ctx context.Context, userid int) (bool, error) {
	isAdmin, err := s.repos.IsAdmin(ctx, userid)

	return isAdmin, err
	// false, status.Errorf(codes.Unimplemented, "Service. Method IsAdmin is not implemented")
}
