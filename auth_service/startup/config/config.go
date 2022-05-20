package config

type Config struct {
	Port               string
	UserDBHost         string
	UserDBPort         string
	ProfileServicePort string
	ProfileServiceHost string
	ApiGatwayHost      string
	ApiGatwayPort      string
	Email              string
	PasswordEmail      string
}

func NewConfig() *Config {
	return &Config{
		//TODO: ENV varibale
		// mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
		Port:               "8081",      //os.Getenv("POST_SERVICE_PORT"),
		UserDBHost:         "localhost", // os.Getenv("CATALOGUE_DB_HOST"),
		UserDBPort:         "27017",     // os.Getenv("CATALOGUE_DB_PORT"),
		ProfileServiceHost: "localhost",
		ProfileServicePort: "8082",
		ApiGatwayHost:      "localhost",
		ApiGatwayPort:      "9000",
		Email:              "dislinkt@outlook.com",
		PasswordEmail:      "disXWSpass2022",
	}
}
