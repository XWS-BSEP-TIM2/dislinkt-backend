package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port          string
	MessageDBHost string
	MessageDBPort string

	ProfileHost string
	ProfilePort string

	ConnectionHost string
	ConnectionPort string

	LoggingHost string
	LoggingPort string

	NotificationServiceHost string
	NotificationServicePort string
}

func NewConfig() *Config {
	return &Config{
		Port:          goDotEnvVariable("MESSAGE_SERVICE_PORT"),
		MessageDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		MessageDBPort: goDotEnvVariable("MONGO_DB_PORT"),

		ProfileHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfilePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ConnectionHost: goDotEnvVariable("CONNECTION_SERVICE_HOST"),
		ConnectionPort: goDotEnvVariable("CONNECTION_SERVICE_PORT"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),

		NotificationServiceHost: goDotEnvVariable("NOTIFICATION_SERVICE_HOST"),
		NotificationServicePort: goDotEnvVariable("NOTIFICATION_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
