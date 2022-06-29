package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                       string
	UserDBHost                 string
	UserDBPort                 string
	ProfileServicePort         string
	ProfileServiceHost         string
	ApiGatewayHost             string
	ApiGatewayPort             string
	EmailHost                  string
	Email                      string
	PasswordEmail              string
	LoggingHost                string
	LoggingPort                string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string
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

		EmailHost:     goDotEnvVariable("EMAIL_HOST"),
		Email:         goDotEnvVariable("DISLINKT_EMAIL"),
		PasswordEmail: goDotEnvVariable("EMAIL_PASSWORD"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),

		NatsHost:                   goDotEnvVariable("NATS_HOST"),
		NatsPort:                   goDotEnvVariable("NATS_PORT"),
		NatsUser:                   goDotEnvVariable("NATS_USER"),
		NatsPass:                   goDotEnvVariable("NATS_PASS"),
		RegisterUserCommandSubject: goDotEnvVariable("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:   goDotEnvVariable("REGISTER_USER_REPLY_SUBJECT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
