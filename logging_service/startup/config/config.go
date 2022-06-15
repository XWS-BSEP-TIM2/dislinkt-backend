package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Host string
	Port string
}

func NewConfig() *Config {
	return &Config{
		Host: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		Port: goDotEnvVariable("LOGGING_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
