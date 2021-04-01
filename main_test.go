package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/recipe"
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

func TestGetAllRecipesRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	values := recipe.Recipe{
		Name:         "Test",
		ActiveTime:   1,
		TotalTime:    2,
		Servings:     3,
		Categories:   []string{"nice"},
		Ingredients:  "wasd",
		Instructions: "wasd",
		UserID:       "1337",
		Calories:     1,
	}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	createResponse, err := http.Post(fmt.Sprintf(
		"%s/recipes/", testServer.URL),
		"application/json",
		bytes.NewBuffer(json_data),
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createResponse.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", createResponse.StatusCode)
	}

	response, err := http.Get(fmt.Sprintf("%s/recipes/", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var createdRecipe []recipe.Recipe
	err = json.Unmarshal(body, &createdRecipe)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdRecipe[0].Name != "Test" {
		t.Fatalf("Expected Test, got %s", createdRecipe[0].Name)
	}

	if createdRecipe[0].ActiveTime != 1 {
		t.Fatalf("Expected Test, got %d", createdRecipe[0].ActiveTime)
	}

	if createdRecipe[0].TotalTime != 2 {
		t.Fatalf("Expected Test, got %d", createdRecipe[0].TotalTime)
	}

	if createdRecipe[0].Ingredients != "wasd" {
		t.Fatalf("Expected Test, got %s", createdRecipe[0].Ingredients)
	}

	if createdRecipe[0].Ingredients != "wasd" {
		t.Fatalf("Expected Test, got %s", createdRecipe[0].Ingredients)
	}

	if createdRecipe[0].Instructions != "wasd" {
		t.Fatalf("Expected Test, got %s", createdRecipe[0].Instructions)
	}

	if createdRecipe[0].UserID != "1337" {
		t.Fatalf("Expected Test, got %s", createdRecipe[0].UserID)
	}

	fmt.Printf("%+v", response.Body)
}
