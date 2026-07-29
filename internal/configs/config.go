package configs

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	Database  DatabaseConfig
	Migration MigrationConfig
	Redis     RedisConfig
	GRPC      GrpcConfig
	Server    ServerConfig
}

func LoadConfig(log *logrus.Logger) (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Warn("Failed to load .env file, reading environment variables directly.")
	}
	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	log.Info("Configuration loaded successfully")
	return cfg, nil
}
