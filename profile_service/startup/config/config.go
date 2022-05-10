package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port          string
	ProfileDBHost string
	ProfileDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ProfileDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		ProfileDBPort: goDotEnvVariable("MONGO_DB_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")

	/*
		if err != nil {
			log.Fatalf("Error loading .env file")
		} */
	return os.Getenv(key)
}
