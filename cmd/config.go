package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"ADDR" required:"true" default:":9090"`
	Postgres
}

type Postgres struct {
	DSN string `envconfig:"POSTGRES_DSN" default:"postgres://postgres:123456@localhost:5432/karma8?connect_timeout=5&sslmode=disable" required:"true"`
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
