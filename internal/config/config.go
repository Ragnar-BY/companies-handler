package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is config struct
type Config struct {
	PostgresAddress  string `env:"POSTGRES_ADDRESS"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`

	ServerAddress string `env:"SERVER_ADDRESS"`

	JWTKey string `env:"JWT_KEY"`
}

// LoadConfig loads config from .env file
func LoadConfig(path string) (config Config, err error) {
	err = cleanenv.ReadConfig(path, &config)
	if errors.Is(err, os.ErrNotExist) {
		err = cleanenv.ReadEnv(&config)
	}
	return
}
