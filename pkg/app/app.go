package app

import (
	grpcapp "grpc-pet/pkg/app/grpc"
	Service "grpc-pet/pkg/service"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *logrus.Logger, grpcPort int, storagePath string, tokenTTL time.Duration, auth *Service.AuthService) *App {
	// TODO: init storage

	// TODO: init auth service

	grpcApp := grpcapp.New(log, grpcPort, auth)

	return &App{
		GRPCApp: grpcApp,
	}
}
