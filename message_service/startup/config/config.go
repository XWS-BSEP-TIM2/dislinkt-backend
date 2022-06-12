package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port          string
	MessageDBHost string
	MessageDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:          goDotEnvVariable("MESSAGE_SERVICE_PORT"),
		MessageDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		MessageDBPort: goDotEnvVariable("MONGO_DB_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
