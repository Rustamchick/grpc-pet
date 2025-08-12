package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string        `yaml:"env" env-default:"local"`
	Storage_path string        `yaml:"storage_path" env-required:"true"`
	TokenTTL     time.Duration `yaml:"token_ttl" env-required:"12h"`
	GRPC         GrpcConfig    `yaml:"grpc"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustInitConfig() Config {
	path := getConfigPath()
	if path == "" {
		panic("Empty config path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("There is no config file in " + path)
	}

	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("Error reading config " + err.Error())
	}
	return *cfg
}

func getConfigPath() string {
	// в будущем может сделать дополнительную возможность прописывать путь через флаги
	res := os.Getenv("CONFIG_PATH")
	return res
}
