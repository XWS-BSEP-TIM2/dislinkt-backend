package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	PostDBHost              string
	PostDBPort              string
	AuthServiceHost         string
	AuthServicePort         string
	ProfileServiceHost      string
	ProfileServicePort      string
	ConnectionServiceHost   string
	ConnectionServicePort   string
	NotificationServiceHost string
	NotificationServicePort string
}

func NewConfig() *Config {
	return &Config{
		// mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
		Port: goDotEnvVariable("POST_SERVICE_PORT"),

		PostDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		PostDBPort: goDotEnvVariable("MONGO_DB_PORT"),

		AuthServiceHost: goDotEnvVariable("AUTH_SERVICE_HOST"),
		AuthServicePort: goDotEnvVariable("AUTH_SERVICE_PORT"),

		ProfileServiceHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfileServicePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ConnectionServiceHost: goDotEnvVariable("CONNECTION_SERVICE_HOST"),
		ConnectionServicePort: goDotEnvVariable("CONNECTION_SERVICE_PORT"),

		NotificationServiceHost: goDotEnvVariable("NOTIFICATION_SERVICE_HOST"),
		NotificationServicePort: goDotEnvVariable("NOTIFICATION_SERVICE_PORT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
