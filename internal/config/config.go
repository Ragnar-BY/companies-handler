package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is config struct
type Config struct {
	PostgreSQLAddr     string `env:"POSTGRESQL_ADDRESS"`
	PostgreSQLUser     string `env:"POSTGRESQL_USER"`
	PostgreSQLPassword string `env:"POSTGRESQL_PASSWORD"`
	PostgreSQLDatabase string `env:"POSTGRESQL_DATABASE"`

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
