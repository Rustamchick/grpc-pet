package grpcapp

import (
	"fmt"
	"grpc-pet/pkg/handler"
	Service "grpc-pet/pkg/service"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	log        *logrus.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *logrus.Logger, port int, authService *Service.AuthService) *App {
	gRPCServer := grpc.NewServer()

	handler.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run()"

	log := a.log.WithField("op", op)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	// log.Info("Running gRPC server, addr:", lis.Addr().String())
	log.WithField("addr", lis.Addr().String()).Info("Running gRPC server")

	if err := a.gRPCServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop()"

	a.log.WithField("op", op).Info("Stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
