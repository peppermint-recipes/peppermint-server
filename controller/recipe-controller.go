package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/peppermint-recipes/peppermint-server/models"
	"github.com/peppermint-recipes/peppermint-server/utils"
)

func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	recipe := &models.Recipe{}

	err := json.NewDecoder(r.Body).Decode(recipe)
	if err != nil {

		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	resp := recipe.Create()
	utils.Respond(w, resp)
}

func GetOneRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		fmt.Println(err)
	}
	id64 := uint(id)
	foundRecipe := models.GetRecipe(id64)

	resp := utils.Message(true, "success")
	resp["data"] = foundRecipe
	utils.Respond(w, resp)
}

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	foundRecipes := models.GetRecipes()

	resp := utils.Message(true, "success")
	resp["data"] = foundRecipes
	utils.Respond(w, resp)
}

func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	recipe := &models.Recipe{}

	recipeID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		fmt.Println(err)
	}

	id64 := uint(id)

	err = json.NewDecoder(r.Body).Decode(recipe)
	if err != nil {

		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	recipe.ID = id64

	resp := models.UpdateRecipe(recipe)
	utils.Respond(w, resp)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	recipe := &models.Recipe{}

	recipeID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		fmt.Println(err)
	}

	id64 := uint(id)

	err = json.NewDecoder(r.Body).Decode(recipe)
	if err != nil {

		fmt.Println(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	recipe.ID = id64

	resp := models.DeleteRecipe(recipe)
	utils.Respond(w, resp)

}
