package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peppermint-recipes/peppermint-server/config"
)

func validateHttp200Response(t *testing.T, response *http.Response) {
	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", response.StatusCode)
	}

	val, ok := response.Header["Content-Type"]

	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestLivezRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/livez", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}

func TestRecipeRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/recipes", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}

func TestWeekplanRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/weekplans", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}

func TestShoppingListRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/shopping-lists", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}
