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
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
