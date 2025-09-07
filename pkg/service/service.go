package Service

import (
	"context"
	"grpc-pet/pkg/repository"
)

type Authentification interface {
	Login(ctx context.Context, email, password string, appid int) (token string, err error)
	Register(ctx context.Context, email, password string) (user_id int64, err error)
	IsAdmin(ctx context.Context, userid int) (bool, error)
}

type Service struct {
	Authentification
}

func NewService(repos repository.Authentification) *Service { // весь сервисный слой
	return &Service{
		Authentification: NewAuthService(repos), // auth сервисного слоя
		// могут добавиться еще функции в сервисный слой
	}
}
