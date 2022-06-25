package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port               string
	NotificationDBHost string
	NotificationDBPort string

	ProfileHost string
	ProfilePort string

	ConnectionHost string
	ConnectionPort string

	LoggingHost string
	LoggingPort string
}

func NewConfig() *Config {
	return &Config{
		Port:               goDotEnvVariable("NOTIFICATION_SERVICE_PORT"),
		NotificationDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		NotificationDBPort: goDotEnvVariable("MONGO_DB_PORT"),

		ProfileHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfilePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ConnectionHost: goDotEnvVariable("CONNECTION_SERVICE_HOST"),
		ConnectionPort: goDotEnvVariable("CONNECTION_SERVICE_PORT"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
