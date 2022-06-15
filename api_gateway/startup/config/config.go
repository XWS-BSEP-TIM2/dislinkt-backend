package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	PostHost string
	PostPort string

	AuthHost string
	AuthPort string

	ProfileHost string
	ProfilePort string

	ConnectionHost string
	ConnectionPort string

	JobOfferHost string
	JobOfferPort string

	MessageHost string
	MessagePort string

	LoggingHost string
	LoggingPort string

	CertificatePath           string
	CertificatePrivateKeyPath string
}

func NewConfig() *Config {
	return &Config{
		Port: goDotEnvVariable("GATEWAY_PORT"),

		PostHost: goDotEnvVariable("POST_SERVICE_HOST"),
		PostPort: goDotEnvVariable("POST_SERVICE_PORT"),

		AuthHost: goDotEnvVariable("AUTH_SERVICE_HOST"),
		AuthPort: goDotEnvVariable("AUTH_SERVICE_PORT"),

		ProfileHost: goDotEnvVariable("PROFILE_SERVICE_HOST"),
		ProfilePort: goDotEnvVariable("PROFILE_SERVICE_PORT"),

		ConnectionHost: goDotEnvVariable("CONNECTION_SERVICE_HOST"),
		ConnectionPort: goDotEnvVariable("CONNECTION_SERVICE_PORT"),

		JobOfferHost: goDotEnvVariable("JOB_OFFER_SERVICE_HOST"),
		JobOfferPort: goDotEnvVariable("JOB_OFFER_SERVICE_PORT"),

		MessageHost: goDotEnvVariable("MESSAGE_SERVICE_HOST"),
		MessagePort: goDotEnvVariable("MESSAGE_SERVICE_PORT"),

		LoggingHost: goDotEnvVariable("LOGGING_SERVICE_HOST"),
		LoggingPort: goDotEnvVariable("LOGGING_SERVICE_PORT"),

		CertificatePath:           goDotEnvVariable("CERTIFICATE_PATH"),
		CertificatePrivateKeyPath: goDotEnvVariable("CERTIFICATE_PRIVATE_KEY_PATH"),
	}
}

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")
	return os.Getenv(key)
}
