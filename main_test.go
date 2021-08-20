package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/database"
	"github.com/peppermint-recipes/peppermint-server/recipe"

	"github.com/peppermint-recipes/peppermint-server/weekplan"

	"github.com/google/go-cmp/cmp"
)

func beforeTestHook(t *testing.T) {
	clearDatabase(t)
}

func afterTestHook(t *testing.T) {
	clearDatabase(t)
}

func clearDatabase(t *testing.T) {
	client, ctx, _ := database.GetConnection()

	err := client.Database(database.DatabaseName).Drop(ctx)
	if err != nil {
		t.Fatalf("Could not drop database %v", err)
	}
}

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

func TestRecipeRouteGet(t *testing.T) {
	beforeTestHook(t)
	defer afterTestHook(t)

	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/recipes", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}

func TestRecipeRouteCreate(t *testing.T) {
	beforeTestHook(t)
	defer afterTestHook(t)

	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	newRecipe := recipe.Recipe{
		Name:         "",
		ActiveTime:   1,
		TotalTime:    2,
		Servings:     2,
		Categories:   []string{"test"},
		Ingredients:  "Krams",
		Instructions: "Test",
		UserID:       "1",
		Deleted:      false,
		Calories:     1337,
	}
	json_data, err := json.Marshal(newRecipe)

	if err != nil {
		log.Fatal(err)
	}

	createResponse, err := http.Post(fmt.Sprintf("%s/recipes", testServer.URL), "application/json; charset=utf-8", bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, createResponse)

	getResponse, err := http.Get(fmt.Sprintf("%s/recipes", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, getResponse)

	var p []recipe.Recipe

	err = json.NewDecoder(getResponse.Body).Decode(&p)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	newRecipe.ID = p[0].ID
	newRecipe.LastUpdated = p[0].LastUpdated

	if !cmp.Equal(p[0], newRecipe) {
		t.Fatalf("Expected recipe %v to be equal %v", newRecipe, p[0])
	}
}

func TestWeekplanRoute(t *testing.T) {
	beforeTestHook(t)
	defer afterTestHook(t)

	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/weekplans", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}

func TestWeekplanRouteCreate(t *testing.T) {
	beforeTestHook(t)
	defer afterTestHook(t)

	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	newWeekplan := weekplan.Weekplan{
		Name:         "",
		ActiveTime:   1,
		TotalTime:    2,
		Servings:     2,
		Categories:   []string{"test"},
		Ingredients:  "Krams",
		Instructions: "Test",
		UserID:       "1",
		Deleted:      false,
		Calories:     1337,
	}
	json_data, err := json.Marshal(newWeekplan)

	if err != nil {
		log.Fatal(err)
	}

	createResponse, err := http.Post(fmt.Sprintf("%s/recipes", testServer.URL), "application/json; charset=utf-8", bytes.NewBuffer(json_data))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, createResponse)

	getResponse, err := http.Get(fmt.Sprintf("%s/recipes", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, getResponse)

	var p []recipe.Recipe

	err = json.NewDecoder(getResponse.Body).Decode(&p)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	newWeekplan.ID = p[0].ID
	newWeekplan.LastUpdated = p[0].LastUpdated

	if !cmp.Equal(p[0], newWeekplan) {
		t.Fatalf("Expected recipe %v to be equal %v", newWeekplan, p[0])
	}
}

func TestShoppingListRoute(t *testing.T) {
	beforeTestHook(t)
	defer afterTestHook(t)

	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	response, err := http.Get(fmt.Sprintf("%s/shopping-lists", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	validateHttp200Response(t, response)
}
