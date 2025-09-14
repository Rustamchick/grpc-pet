package Service

import (
	"context"
	"errors"
	"grpc-pet/pkg/repository"
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

func NewService(log *logrus.Logger, repos repository.Authentification, AppProvider repository.AppProvider, TokenTTL time.Duration) *Service { // весь сервисный слой
	return &Service{
		Authentification: NewAuthService(log, repos, AppProvider, TokenTTL), // auth сервисного слоя

	}
}

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
)
