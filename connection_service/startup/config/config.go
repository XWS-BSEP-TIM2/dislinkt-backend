package config

type Config struct {
	Port          string
	Host          string
	Neo4jUri      string
	Neo4jUsername string
	Neo4jPassword string
}

func NewConfig() *Config {
	return &Config{
		Port:          "8001",
		Host:          "localhost",
		Neo4jUri:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "connection"}
}
