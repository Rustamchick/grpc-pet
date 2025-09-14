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

func New(log *logrus.Logger, port int, service *Service.Service) *App {
	gRPCServer := grpc.NewServer()

	handler.Register(gRPCServer, service)

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
	const loc = "grpcapp.Run()"

	log := a.log.WithField(loc, loc)

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
	const loc = "grpcapp.Stop()"

	a.log.WithField("loc", loc).Info("Stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
