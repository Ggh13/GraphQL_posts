package config

import (
	"qraphQL_posts/pkg/postgres"

	"context"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PostgresCFG postgres.Config `env:"POSTGRES" env-default:"POSTGRES" yaml:"POSTGRES"`
	RestHost    string          `env:"REST_HOST" env-default:"REST_HOST" yaml:"REST_HOST"`
	TypeDB      string          `env:"TYPE_DB" env-default:"TYPE_DB" yaml:"TYPE_DB"`
	PORT        string          `env:"PORT" env-default:"PORT" yaml:"PORT"`
}

func NewConfig(ctx context.Context) (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config/NewConfig: %w", err)
	}

	return &cfg, nil
}
