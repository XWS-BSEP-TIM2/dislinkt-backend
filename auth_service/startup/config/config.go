package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port               string
	UserDBHost         string
	UserDBPort         string
	ProfileServicePort string
	ProfileServiceHost string
}

func NewConfig() *Config {
	return &Config{
		//mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
		Port: goDotEnvVariable("AUTH_SERVICE_PORT"),

		UserDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		UserDBPort: goDotEnvVariable("MONGO_DB_PORT"),

		ProfileServiceHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfileServicePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),
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
