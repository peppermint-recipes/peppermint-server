package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peppermint-recipes/peppermint-server/config"
)

func TestLivezRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/livez", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

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
