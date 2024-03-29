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
	shoppinglist "github.com/peppermint-recipes/peppermint-server/shopping-list"
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
		Name:         "my-recipe",
		ActiveTime:   1,
		TotalTime:    2,
		Servings:     2,
		Categories:   []string{"test"},
		Ingredients:  "banana",
		Instructions: "test all the things",
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

	testRecipe := recipe.Recipe{
		Name:         "my-recipe",
		ActiveTime:   1,
		TotalTime:    2,
		Servings:     2,
		Categories:   []string{"test"},
		Ingredients:  "banana",
		Instructions: "test all the things",
		UserID:       "1",
		Deleted:      false,
		Calories:     1337,
	}
	testRecipeAsJSON, err := json.Marshal(testRecipe)

	if err != nil {
		log.Fatal(err)
	}

	createResponse, err := http.Post(
		fmt.Sprintf("%s/recipes", testServer.URL),
		"application/json; charset=utf-8",
		bytes.NewBuffer(testRecipeAsJSON),
	)
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

	testRecipe.ID = p[0].ID
	testRecipe.LastUpdated = p[0].LastUpdated

	testDay1 := weekplan.Day{
		Breakfast: []recipe.Recipe{},
		Lunch:     []recipe.Recipe{testRecipe},
		Snack:     []recipe.Recipe{testRecipe},
		Dinner:    []recipe.Recipe{},
	}

	testDay2 := weekplan.Day{
		Breakfast: []recipe.Recipe{testRecipe},
		Lunch:     []recipe.Recipe{},
		Snack:     []recipe.Recipe{},
		Dinner:    []recipe.Recipe{testRecipe},
	}

	newWeekplan := weekplan.WeekPlan{
		UserID:       "1",
		CalendarWeek: 1,
		Year:         2021,
		Monday:       testDay1,
		Tuesday:      testDay2,
		Wednesday:    testDay1,
		Thursday:     testDay2,
		Friday:       testDay1,
		Saturday:     testDay2,
		Sunday:       testDay1,
	}

	testWeekPlanAsJSON, err := json.Marshal(newWeekplan)

	if err != nil {
		log.Fatal(err)
	}

	createResponseWeekPlan, err := http.Post(
		fmt.Sprintf("%s/weekplans", testServer.URL),
		"application/json; charset=utf-8",
		bytes.NewBuffer(testWeekPlanAsJSON),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, createResponseWeekPlan)

	getResponseWeekPlan, err := http.Get(fmt.Sprintf("%s/weekplans", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, getResponseWeekPlan)

	var w []weekplan.WeekPlan

	err = json.NewDecoder(getResponseWeekPlan.Body).Decode(&w)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	newWeekplan.ID = w[0].ID
	newWeekplan.LastUpdated = w[0].LastUpdated

	if !cmp.Equal(w[0], newWeekplan) {
		t.Fatalf("Expected weekplans %v to be equal %v", newWeekplan, w[0])
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

func TestShoppingListRouteCreate(t *testing.T) {
	beforeTestHook(t)
	defer afterTestHook(t)

	config := config.GetConfig()
	testServer := httptest.NewServer(setupServer(config.DB))
	defer testServer.Close()

	testItem1 := shoppinglist.ShoppingListItem{
		Ingredient: "Test1",
		Unit:       "my-unit",
		Amount:     100,
	}

	testItem2 := shoppinglist.ShoppingListItem{
		Ingredient: "Test2",
		Unit:       "my-unit-2",
		Amount:     200,
	}

	newShoppingList := shoppinglist.ShoppingList{
		UserID: "1",
		Items:  []shoppinglist.ShoppingListItem{testItem1, testItem2},
	}

	testWeekPlanAsJSON, err := json.Marshal(newShoppingList)
	if err != nil {
		log.Fatal(err)
	}

	createResponseShoppingList, err := http.Post(
		fmt.Sprintf("%s/shopping-lists", testServer.URL),
		"application/json; charset=utf-8",
		bytes.NewBuffer(testWeekPlanAsJSON),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, createResponseShoppingList)

	getResponseShoppingList, err := http.Get(fmt.Sprintf("%s/shopping-lists", testServer.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	validateHttp200Response(t, getResponseShoppingList)

	var w []shoppinglist.ShoppingList

	err = json.NewDecoder(getResponseShoppingList.Body).Decode(&w)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	newShoppingList.ID = w[0].ID
	newShoppingList.LastUpdated = w[0].LastUpdated

	if !cmp.Equal(w[0], newShoppingList) {
		t.Fatalf("Expected shopingLists %v to be equal %v", newShoppingList, w[0])
	}
}
