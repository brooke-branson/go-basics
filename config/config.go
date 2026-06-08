package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	Port     string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	cfg := Config{
		MongoURI: os.Getenv("MONGO_URI"),
		Port:     os.Getenv("PORT"),
	}
	if cfg.MongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	if cfg.Port == "" {
		log.Fatal("PORT is not set")
	}
	return cfg
}
