package config

import "os"

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

	CertificatePath           string
	CertificatePrivateKeyPath string
}

func NewConfig() *Config {
	return &Config{
		Port:                      getEnv("GATEWAY_PORT", "9000"),
		PostHost:                  getEnv("POST_HOST", "localhost"),
		PostPort:                  getEnv("POST_PORT", "8080"),
		AuthHost:                  getEnv("AUTH_HOST", "localhost"),
		AuthPort:                  getEnv("AUTH_PORT", "8081"),
		ProfileHost:               getEnv("PROFILE_HOST", "localhost"),
		ProfilePort:               getEnv("PROFILE_PORT", "8082"),
		ConnectionHost:            getEnv("CONNECTION_HOST", "localhost"),
		ConnectionPort:            getEnv("CONNECTION_PORT", "8082"),
		CertificatePath:           getEnv("CERTIFICATE_PATH", "certificates/dislinkt_gateway.crt"),
		CertificatePrivateKeyPath: getEnv("CERTIFICATE_PRIVATE_KEY_PATH", "certificates/dislinkt_gateway.key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
