package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB  *DBConfig
	Web *WebServerConfig
}

type DBConfig struct {
	Endpoint string
	Username string
	Password string
}

type WebServerConfig struct {
	Address string
	Port    string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Endpoint: getFromEnvAsString("DATABASE_ENDPOINT", "127.0.0.1:27017"),
			Username: getFromEnvAsString("DATABASE_USER", "root"),
			Password: getFromEnvAsString("DATABASE_PASSWORD", "example"),
		},
		Web: &WebServerConfig{
			Address: getFromEnvAsString("WEBSERVER_ADDRESS", "127.0.0.1"),
			Port:    getFromEnvAsString("WEBSERVER_PORT", "1337"),
		},
	}
}

func getFromEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getFromEnvAsInt(name string, defaultValue int) int {
	valueStr := getFromEnvAsString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
