package config

type Config struct {
	Port string

	PostHost string
	PostPort string

	AuthHost string
	AuthPort string

	ConnectionHost string
	ConnectionPort string
}

func NewConfig() *Config {
	return &Config{
		Port:           "9000",      //os.Getenv("GATEWAY_PORT"),
		PostHost:       "localhost", // os.Getenv("POST_SERVICE_HOST"),
		PostPort:       "8080",      //os.Getenv("POST_SERVICE_PORT"),
		AuthHost:       "localhost",
		AuthPort:       "8081",
		ConnectionHost: "localhost",
		ConnectionPort: "8001",
	}
}
