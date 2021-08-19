package config

import (
	"os"
	"reflect"
	"testing"
)

var testingEnvVaribles = map[string]string{
	"DATABASE_ENDPOINT": "127.0.0.1:1337",
	"DATABASE_USER":     "my-user",
	"DATABASE_PASSWORD": "example-pw",
	"WEBSERVER_ADDRESS": "10.13.37.1",
	"WEBSERVER_PORT":    "1338",
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

	if config.DB.Endpoint != "127.0.0.1:27017" {
		t.Errorf("config.DB.Endpoint = %s; want 127.0.0.1:27017", config.DB.Endpoint)
	}
	if config.DB.Username != "root" {
		t.Errorf("config.DB.Username = %s; want root", config.DB.Username)
	}
	if config.DB.Password != "example" {
		t.Errorf("config.DB.Password = %s; want example", config.DB.Password)
	}
	if config.Web.Address != "127.0.0.1" {
		t.Errorf("config.Web.Address = %s; want 127.0.0.1", config.Web.Address)
	}
	if config.Web.Port != "1337" {
		t.Errorf("config.Web.Port = %s; want 1337", config.Web.Port)
	}
}

func TestGetConfigWithValuesFromEnvironment(t *testing.T) {
	setTestEnvVariables()
	config := GetConfig()

	if config.DB.Endpoint != "127.0.0.1:1337" {
		t.Errorf("config.DB.Endpoint = %s; want 127.0.0.1:1337", config.DB.Endpoint)
	}
	if config.DB.Username != "my-user" {
		t.Errorf("config.DB.Username = %s; want my-user", config.DB.Username)
	}
	if config.DB.Password != "example-pw" {
		t.Errorf("config.DB.Password = %s; want example-pw", config.DB.Password)
	}
	if config.Web.Address != "10.13.37.1" {
		t.Errorf("config.Web.Address = %s; want 10.13.37.1", config.Web.Address)
	}
	if config.Web.Port != "1338" {
		t.Errorf("config.Web.Port = %s; want 1338", config.Web.Port)
	}
	unsetTestEnvVariables()
}
