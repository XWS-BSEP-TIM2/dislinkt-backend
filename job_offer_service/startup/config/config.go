package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port           string
	JobOfferDBHost string
	JobOfferDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:           goDotEnvVariable("JOB_OFFER_SERVICE_PORT"),
		JobOfferDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		JobOfferDBPort: goDotEnvVariable("MONGO_DB_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
