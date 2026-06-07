package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	Port     string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := Config{
		MongoURI: os.Getenv("MONGO_URI"),
		Port:     os.Getenv("PORT"),
	}
	if cfg.MongoURI == "" {
		return cfg, fmt.Errorf("MONGO_URI is not set")
	}
	if cfg.Port == "" {
		return cfg, fmt.Errorf("PORT is not set")
	}
	return cfg, nil
}