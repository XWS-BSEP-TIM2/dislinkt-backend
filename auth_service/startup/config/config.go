package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	UserDBHost         string
	UserDBPort         string
	ProfileServicePort string
	ProfileServiceHost string
	ApiGatwayHost      string
	ApiGatwayPort      string
	Email              string
	PasswordEmail      string
}

func NewConfig() *Config {
	return &Config{
		//TODO: ENV varibale
		// mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
		Port: goDotEnvVariable("AUTH_SERVICE_PORT"),

		UserDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		UserDBPort: goDotEnvVariable("MONGO_DB_PORT"),

		ProfileServiceHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfileServicePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ApiGatwayHost:      goDotEnvVariable("GATEWAY_HOST"),
		ApiGatwayPort:      goDotEnvVariable("GATEWAY_PORT"),

		Email:              "dislinkt@outlook.com",
		PasswordEmail:      "disXWSpass2022",

	}
}

func goDotEnvVariable(key string) string {
	var err = godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	} 

	return os.Getenv(key)
}
