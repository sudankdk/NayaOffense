package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	PORT string
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	cfg := AppConfig{
		PORT: os.Getenv("PORT"),
	}
	switch {
	case cfg.PORT == "":
		return nil, errors.New("PORT missing in environment")
	}
	return &cfg, nil
}
