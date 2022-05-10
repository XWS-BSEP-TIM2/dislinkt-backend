package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port string

	PostHost string
	PostPort string

	AuthHost string
	AuthPort string

	ProfileHost    string
	ProfilePort    string
	ConnectionHost string
	ConnectionPort string
}

func NewConfig() *Config {
	return &Config{
		Port: goDotEnvVariable("GATEWAY_PORT"),

		PostHost: goDotEnvVariable("POST_SERVICE_HOST"),
		PostPort: goDotEnvVariable("POST_SERVICE_PORT"),

		AuthHost: goDotEnvVariable("AUTH_SERVICE_HOST"),
		AuthPort: goDotEnvVariable("AUTH_SERVICE_PORT"),

		ProfileHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfilePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ConnectionHost: goDotEnvVariable("CONNECTION_SERVICE_HOST"),
		ConnectionPort: goDotEnvVariable("CONNECTION_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")

	/*
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		} */
	return os.Getenv(key)
}
