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

	if !areRecipeEqual(t, values, createdRecipe[0]) {
		t.Fatalf("Expected %+v, got %+v", values, createdRecipe[0])
	}

	cleanupRecipes(testServer)
}

func TestDeleteRecipeRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	cleanupRecipes(testServer)

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

	body, err := ioutil.ReadAll(createResponse.Body)
	defer createResponse.Body.Close()

	var createdRecipe recipe.Recipe
	err = json.Unmarshal(body, &createdRecipe)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf(
		"%s/recipes/%s", testServer.URL, createdRecipe.ID.Hex()), nil)
	deleteResponse, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if deleteResponse.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", createResponse.StatusCode)
	}

	response, err := http.Get(fmt.Sprintf("%s/recipes/", testServer.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", response.StatusCode)
	}

	body, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var allRecipes []recipe.Recipe
	err = json.Unmarshal(body, &allRecipes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, recipe := range allRecipes {
		if recipe.Deleted != true {
			t.Fatalf("Expected no recipes, got at least one %+v", recipe)
		}
	}

	cleanupRecipes(testServer)
}

func TestUpdateRecipeRoute(t *testing.T) {
	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	initialRecipe := recipe.Recipe{
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
	initialRecipeAsJSON, err := json.Marshal(initialRecipe)

	if err != nil {
		log.Fatal(err)
	}

	createResponse, err := http.Post(fmt.Sprintf(
		"%s/recipes/", testServer.URL),
		"application/json",
		bytes.NewBuffer(initialRecipeAsJSON),
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createResponse.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", createResponse.StatusCode)
	}

	body, err := ioutil.ReadAll(createResponse.Body)
	defer createResponse.Body.Close()

	var createdRecipe recipe.Recipe
	err = json.Unmarshal(body, &createdRecipe)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	createdRecipe.Name = "Test-Test"

	updatedRecipeAsJSON, err := json.Marshal(createdRecipe)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf(
		"%s/recipes/", testServer.URL), bytes.NewBuffer(updatedRecipeAsJSON))
	putResponse, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if putResponse.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", createResponse.StatusCode)
	}

	response, err := http.Get(fmt.Sprintf("%s/recipes/%s", testServer.URL, createdRecipe.ID.Hex()))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", response.StatusCode)
	}

	body, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var updatedRecipe recipe.Recipe
	err = json.Unmarshal(body, &updatedRecipe)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedRecipe.Name != "Test-Test" {
		t.Fatalf("Expected name Test-Test, got %s", updatedRecipe.Name)
	}

	cleanupRecipes(testServer)
}

func areRecipeEqual(t *testing.T, recipe1 recipe.Recipe, recipe2 recipe.Recipe) bool {
	if recipe1.Name != recipe2.Name {
		t.Fatalf("Expected %s, got %s", recipe1.Name, recipe2.Name)
		return false
	}
	if recipe1.ActiveTime != recipe2.ActiveTime {
		t.Fatalf("Expected %d, got %d", recipe1.ActiveTime, recipe2.ActiveTime)
		return false
	}
	if recipe1.TotalTime != recipe2.TotalTime {
		t.Fatalf("Expected %d, got %d", recipe1.TotalTime, recipe2.TotalTime)
		return false
	}
	if recipe1.Servings != recipe2.Servings {
		t.Fatalf("Expected %d, got %d", recipe1.TotalTime, recipe2.TotalTime)
		return false
	}
	if recipe1.Ingredients != recipe2.Ingredients {
		t.Fatalf("Expected %s, got %s", recipe1.Ingredients, recipe2.Ingredients)
		return false
	}
	if recipe1.Instructions != recipe2.Instructions {
		t.Fatalf("Expected %s, got %s", recipe1.Instructions, recipe2.Instructions)
		return false
	}
	if recipe1.UserID != recipe2.UserID {
		t.Fatalf("Expected %s, got %s", recipe1.UserID, recipe2.UserID)
		return false
	}
	if recipe1.Calories != recipe2.Calories {
		t.Fatalf("Expected %d, got %d", recipe1.Calories, recipe2.Calories)
		return false
	}

	for index, category := range recipe1.Categories {
		if category != recipe1.Categories[index] {
			t.Fatalf("Expected %s, got %s", category, recipe2.Categories[index])
			return false
		}
	}

	return true
}

func cleanupRecipes(testServer *httptest.Server) {
	getRecipesResponse, err := http.Get(fmt.Sprintf("%s/recipes/", testServer.URL))

	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	if getRecipesResponse.StatusCode != 200 {
		log.Fatalf("Expected status code 200, got %v", getRecipesResponse.StatusCode)
	}

	body, err := ioutil.ReadAll(getRecipesResponse.Body)
	defer getRecipesResponse.Body.Close()

	var allRecipes []recipe.Recipe
	err = json.Unmarshal(body, &allRecipes)
	if err != nil {
		log.Fatalf("Expected no error, got %v", err)
	}

	for _, recipe := range allRecipes {
		// if recipe.Deleted != false {
		req, err := http.NewRequest("DELETE", fmt.Sprintf(
			"%s/recipes/%s", testServer.URL, recipe.ID.Hex()), nil)
		deleteResponse, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Fatalf("Expected no error, got %v", err)
		}

		if deleteResponse.StatusCode != 200 {
			log.Fatalf("Expected status code 200, got %v", req)
		}
		// }

	}
}
