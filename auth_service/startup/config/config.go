package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	UserDBHost         string
	UserDBPort         string
	ProfileServicePort string
	ProfileServiceHost string
	ApiGatewayHost     string
	ApiGatewayPort     string
	Email              string
	PasswordEmail      string
	LoggingHost        string
	LoggingPort        string
}

func NewConfig() *Config {
	return &Config{
		// mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
		Port: goDotEnvVariable("AUTH_SERVICE_PORT"),

		UserDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		UserDBPort: goDotEnvVariable("MONGO_DB_PORT"),

		ProfileServiceHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfileServicePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ApiGatewayHost: goDotEnvVariable("GATEWAY_HOST"),
		ApiGatewayPort: goDotEnvVariable("GATEWAY_PORT"),

		Email:         goDotEnvVariable("DISLINKT_EMAIL"),
		PasswordEmail: goDotEnvVariable("EMAIL_PASSWORD"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
