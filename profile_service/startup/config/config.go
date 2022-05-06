package config

type Config struct {
	Port          string
	ProfileDBHost string
	ProfileDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:          "8082",      //os.Getenv("POST_SERVICE_PORT"),
		ProfileDBHost: "localhost", // os.Getenv("CATALOGUE_DB_HOST"),
		ProfileDBPort: "27017",     // os.Getenv("CATALOGUE_DB_PORT"),
	}
}
