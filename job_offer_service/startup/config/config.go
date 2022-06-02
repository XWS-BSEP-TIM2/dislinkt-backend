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
		Port:          goDotEnvVariable("JOB_OFFER_SERVICE_PORT"),
		ProfileDBHost: goDotEnvVariable("JOB_OFFER_SERVICE_HOST"),
		ProfileDBPort: goDotEnvVariable("MONGO_DB_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
