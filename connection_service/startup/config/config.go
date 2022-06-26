package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Host          string
	Neo4jUri      string
	Neo4jHost     string
	Neo4jPort     string
	Neo4jUsername string
	Neo4jPassword string

	LoggingHost string
	LoggingPort string

	NotificationServiceHost string
	NotificationServicePort string

	ProfileHost string
	ProfilePort string

	MessageHost string
	MessagePort string
}

func NewConfig() *Config {
	return &Config{
		Host: goDotEnvVariable("CONNECTION_SERVICE_HOST"),
		Port: goDotEnvVariable("CONNECTION_SERVICE_PORT"),

		Neo4jUri:      goDotEnvVariable("NEO4J_URI"),
		Neo4jHost:     goDotEnvVariable("NEO4J_HOST"),
		Neo4jPort:     goDotEnvVariable("NEO4J_PORT"),
		Neo4jUsername: goDotEnvVariable("NEO4J_USERNAME"),
		Neo4jPassword: goDotEnvVariable("NEO4J_PASSWORD"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),

		NotificationServiceHost: goDotEnvVariable("NOTIFICATION_SERVICE_HOST"),
		NotificationServicePort: goDotEnvVariable("NOTIFICATION_SERVICE_PORT"),

		ProfileHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfilePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		MessageHost: goDotEnvVariable("MESSAGE_SERVICE_HOST"),
		MessagePort: goDotEnvVariable("MESSAGE_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
