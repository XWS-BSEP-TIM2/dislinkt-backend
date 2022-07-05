package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host          string
	Port          string
	Neo4jUri      string
	Neo4jHost     string
	Neo4jPort     string
	Neo4jUsername string
	Neo4jPassword string

	LoggingHost string
	LoggingPort string

	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string

	UpdateSkillsCommandSubject string
	UpdateSkillsReplySubject   string
}

func NewConfig() *Config {
	return &Config{
		Host: goDotEnvVariable("JOB_OFFER_SERVICE_HOST"),
		Port: goDotEnvVariable("JOB_OFFER_SERVICE_PORT"),

		Neo4jUri:      goDotEnvVariable("JOB_OFFER_NEO4J_URI"),
		Neo4jHost:     goDotEnvVariable("JOB_OFFER_NEO4J_HOST"),
		Neo4jPort:     goDotEnvVariable("JOB_OFFER_NEO4J_PORT"),
		Neo4jUsername: goDotEnvVariable("JOB_OFFER_NEO4J_USERNAME"),
		Neo4jPassword: goDotEnvVariable("JOB_OFFER_NEO4J_PASSWORD"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),

		NatsHost:                   goDotEnvVariable("NATS_HOST"),
		NatsPort:                   goDotEnvVariable("NATS_PORT"),
		NatsUser:                   goDotEnvVariable("NATS_USER"),
		NatsPass:                   goDotEnvVariable("NATS_PASS"),
		RegisterUserCommandSubject: goDotEnvVariable("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:   goDotEnvVariable("REGISTER_USER_REPLY_SUBJECT"),

		UpdateSkillsCommandSubject: goDotEnvVariable("UPDATE_SKILLS_COMMAND_SUBJECT"),
		UpdateSkillsReplySubject:   goDotEnvVariable("UPDATE_SKILLS_REPLY_SUBJECT"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
