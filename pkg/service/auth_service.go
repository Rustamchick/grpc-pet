package Service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authentification interface {
	Login(ctx context.Context, email, password string, appid int) (token string, err error)
	Register(ctx context.Context, email, password string) (user_id int64, err error)
	IsAdmin(ctx context.Context, userid int) (bool, error)
}

type AuthService struct {
	db sqlx.DB
}

func NewAuthService() *AuthService { // auth сервисного слоя
	return &AuthService{}
}

func (s *AuthService) Login(ctx context.Context, email, password string, appid int) (token string, err error) {
	return "Service return token", status.Errorf(codes.Unimplemented, "Service. Method Login is not implemented")
}

func (s *AuthService) Register(ctx context.Context, email, password string) (user_id int64, err error) {
	return -111, status.Errorf(codes.Unimplemented, "Service. Method Register is not implemented")
}

func (s *AuthService) IsAdmin(ctx context.Context, userid int) (bool, error) {
	return false, status.Errorf(codes.Unimplemented, "Service. Method IsAdmin is not implemented")
}
