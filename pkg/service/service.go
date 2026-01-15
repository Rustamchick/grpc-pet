package Service

import (
	"context"
	"grpc-pet/pkg/repository"
	AuthService "grpc-pet/pkg/service/auth"
	"time"

	"github.com/sirupsen/logrus"
)

type Authentification interface {
	Login(ctx context.Context, email, password string, appid int) (token string, err error)
	Register(ctx context.Context, email, password string) (user_id int64, err error)
	IsAdmin(ctx context.Context, userid int) (bool, error)
}

type Service struct {
	Authentification
}

func NewService(log *logrus.Logger, repos repository.Authentification, AppProvider repository.AppProvider, TokenTTL time.Duration) *Service {
	return &Service{
		Authentification: AuthService.NewAuthService(log, repos, AppProvider, TokenTTL),
	}
}
