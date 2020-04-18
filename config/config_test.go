package config

import (
	"os"
	"reflect"
	"testing"
)

var testingEnvVaribles = map[string]string{
	"DATABASE_DIALECT":  "dialect-mysql",
	"DATABASE_HOST":     "10.13.37.1",
	"DATABASE_PORT":     "1337",
	"DATABASE_USER":     "database-root",
	"DATABASE_PASSWORD": "database-password",
	"DATABASE_NAME":     "database-name",
}

func setTestEnvVariables() {
	for key, value := range testingEnvVaribles {
		os.Setenv(key, value)
	}
}

func unsetTestEnvVariables() {
	for key, _ := range testingEnvVaribles {
		os.Unsetenv(key)
	}
}

func TestGetFromEnvAsString(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestDefault", args{"wasd", "default"}, "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFromEnvAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFromEnvAsInt(t *testing.T) {
	type args struct {
		key          string
		defaultValue int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"TestDefault", args{"wasd", 1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFromEnvAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigWithDefaultValues(t *testing.T) {
	unsetTestEnvVariables()
	config := GetConfig()

	if config.DB.Dialect != "mysql" {
		t.Errorf("GetConfig() = %s; want mysql", config.DB.Dialect)
	}
	if config.DB.Host != "127.0.0.1" {
		t.Errorf("GetConfig() = %s; want 127.0.0.1", config.DB.Host)
	}
	if config.DB.Port != 3306 {
		t.Errorf("GetConfig() = %d; want 3306", config.DB.Port)
	}
	if config.DB.Username != "root" {
		t.Errorf("GetConfig() = %s; want root", config.DB.Username)
	}
	if config.DB.Password != "Guest0000!" {
		t.Errorf("GetConfig() = %s; want Guest0000!", config.DB.Password)
	}
	if config.DB.Name != "peppermint" {
		t.Errorf("GetConfig() = %s; want peppermint", config.DB.Name)
	}
}

func TestGetConfigWithValuesFromEnvironment(t *testing.T) {
	setTestEnvVariables()
	config := GetConfig()

	if config.DB.Dialect != "dialect-mysql" {
		t.Errorf("GetConfig() = %s; want dialect-mysql", config.DB.Dialect)
	}
	if config.DB.Host != "10.13.37.1" {
		t.Errorf("GetConfig() = %s; want 10.13.37.1", config.DB.Host)
	}
	if config.DB.Port != 1337 {
		t.Errorf("GetConfig() = %d; want 1337", config.DB.Port)
	}
	if config.DB.Username != "database-root" {
		t.Errorf("GetConfig() = %s; want database-root", config.DB.Username)
	}
	if config.DB.Password != "database-password" {
		t.Errorf("GetConfig() = %s; want database-password", config.DB.Password)
	}
	if config.DB.Name != "database-name" {
		t.Errorf("GetConfig() = %s; want database-name", config.DB.Name)
	}

	unsetTestEnvVariables()
}
