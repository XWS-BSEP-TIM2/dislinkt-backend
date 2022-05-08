package config

type Config struct {
	Port                  string
	PostDBHost            string
	PostDBPort            string
	AuthServiceHost       string
	AuthServicePort       string
	ProfileServiceHost    string
	ProfileServicePort    string
	ConnectionServiceHost string
	ConnectionServicePort string
}

func NewConfig() *Config {
	return &Config{
		//TODO: ENV varibale
		// mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
		Port:                  "8080",      //os.Getenv("POST_SERVICE_PORT"),
		PostDBHost:            "localhost", // os.Getenv("CATALOGUE_DB_HOST"),
		PostDBPort:            "27017",     // os.Getenv("CATALOGUE_DB_PORT"),
		AuthServiceHost:       "localhost",
		AuthServicePort:       "8081",
		ProfileServiceHost:    "localhost",
		ProfileServicePort:    "8082",
		ConnectionServiceHost: "localhost",
		ConnectionServicePort: "8001",
	}
}
