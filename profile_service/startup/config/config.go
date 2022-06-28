package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port                       string
	ProfileDBHost              string
	ProfileDBPort              string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string
}

func NewConfig() *Config {
	return &Config{
		Port: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ProfileDBHost: goDotEnvVariable("MONGO_DB_HOST"),
		ProfileDBPort: goDotEnvVariable("MONGO_DB_PORT"),

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
