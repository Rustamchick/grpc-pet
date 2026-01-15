package main

import (
	"context"
	"grpc-pet/pkg/config"
	"net"
	"strconv"
	"testing"

	grpcpetv1 "github.com/Rustamchick/protobuff/gen/go/pet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestClient struct {
	*testing.T                      // Потребуется для вызова методов *testing.T внутри Suite
	Cfg        *config.Config       // Конфигурация приложения
	AuthClient grpcpetv1.AuthClient // Клиент для взаимодействия с gRPC-сервером
}

const (
	grpcHost = "localhost"
)

func NewTestClient(t *testing.T) (context.Context, *TestClient) {
	t.Helper()
	t.Parallel()

	cfg := config.InitConfigByPath(configPath())
	// fmt.Printf("Config: %+v", cfg)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	hostport := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))

	cc, err := grpc.NewClient(hostport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &TestClient{
		T:          t,
		Cfg:        &cfg,
		AuthClient: grpcpetv1.NewAuthClient(cc),
	}
}

func configPath() string {
	// const key = "CONFIG_PATH"

	// if v := os.Getenv(key); v != "" {
	// 	return v
	// }

	return "../../configs/local_tests.yaml"
}
