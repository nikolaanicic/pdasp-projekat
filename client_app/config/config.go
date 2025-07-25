package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string `long:"host" env:"SDK_APP_HOST"`
	Port        string `long:"port" env:"SDK_APP_PORT"`
	TokenSecret string `long:"secret" env:"JWT_TOKEN_SECRET"`
}

func Load() (*Config, error) {

	var config Config

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load the config:", err)
	}

	if err := setEnv("DISCOVERY_AS_LOCALHOST", "true"); err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

	config.Host = os.Getenv("SDK_APP_HOST")
	config.Port = os.Getenv("SDK_APP_PORT")
	config.TokenSecret = os.Getenv("JWT_TOKEN_SECRET")

	return &config, nil
}

func setEnv(env string, value string) error {
	if err := os.Setenv(env, value); err != nil {
		return err
	}
	return nil
}
