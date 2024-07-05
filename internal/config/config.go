package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Address    string `env:"ADDRESS" env-required:"true"`
	DBHost     string `env:"DB_HOST" env-required:"true"`
	DBPort     int    `env:"DB_PORT" env-required:"true" env-default:"5432"`
	DBName     string `env:"DB_NAME" env-required:"true"`
	DBUser     string `env:"DB_USER" env-required:"true"`
	DBPassword string `env:"DB_PWD" env-required:"true"`
	InfoAPIURL string `env:"API" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
