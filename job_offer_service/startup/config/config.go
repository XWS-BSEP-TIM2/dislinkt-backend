package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Host          string
	Port          string
	Neo4jUri      string
	Neo4jHost     string
	Neo4jPort     string
	Neo4jUsername string
	Neo4jPassword string

	LoggingHost string
	LoggingPort string
}

func NewConfig() *Config {
	return &Config{
		Host: goDotEnvVariable("JOB_OFFER_SERVICE_HOST"),
		Port: goDotEnvVariable("JOB_OFFER_SERVICE_PORT"),

		Neo4jUri:      goDotEnvVariable("JOB_OFFER_NEO4J_URI"),
		Neo4jHost:     goDotEnvVariable("JOB_OFFER_NEO4J_HOST"),
		Neo4jPort:     goDotEnvVariable("JOB_OFFER_NEO4J_PORT"),
		Neo4jUsername: goDotEnvVariable("JOB_OFFER_NEO4J_USERNAME"),
		Neo4jPassword: goDotEnvVariable("JOB_OFFER_NEO4J_PASSWORD"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
