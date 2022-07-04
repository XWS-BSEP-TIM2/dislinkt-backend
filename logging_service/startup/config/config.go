package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Host         string
	Port         string
	FilePath     string
	EventsDBHost string
	EventsDBPort string
}

func NewConfig() *Config {
	return &Config{
		Host:         goDotEnvVariable("LOGGING_SERVICE_HOST"),
		Port:         goDotEnvVariable("LOGGING_SERVICE_PORT"),
		FilePath:     goDotEnvVariable("LOGGING_FILE_PATH"),
		EventsDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		EventsDBPort: goDotEnvVariable("MONGO_DB_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
