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
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

type WebServerConfig struct {
	AddressAndPort string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  getFromEnvAsString("DATABASE_DIALECT", "mysql"),
			Host:     getFromEnvAsString("DATABASE_HOST", "127.0.0.1"),
			Port:     getFromEnvAsInt("DATABASE_PORT", 3306),
			Username: getFromEnvAsString("DATABASE_USER", "root"),
			Password: getFromEnvAsString("DATABASE_PASSWORD", "Guest0000!"),
			Name:     getFromEnvAsString("DATABASE_NAME", "peppermint"),
		},
		Web: &WebServerConfig{
			AddressAndPort: getFromEnvAsString("WEBSERVER_ADDRESS_PORT", "127.0.0.1:1337"),
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
