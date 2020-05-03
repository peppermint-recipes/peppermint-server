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

func CreateRecipe(writer http.ResponseWriter, request *http.Request) {
	recipe := &models.Recipe{}

	err := json.NewDecoder(request.Body).Decode(recipe)
	if err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage("Error while decoding request body"))

		return
	}

	if err = recipe.Validate(); err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage(err.Error()))

		return
	}

	resp := recipe.Create()
	utils.Respond(writer, resp)
}

func GetOneRecipe(writer http.ResponseWriter, request *http.Request) {
	recipeID := mux.Vars(request)["id"]

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		fmt.Println(err)
	}
	id64 := uint(id)

	err, foundRecipe := models.GetRecipe(id64)
	if err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage(err.Error()))

		return
	}

	resp := utils.PrepareReturn()
	resp["data"] = foundRecipe
	utils.Respond(writer, resp)
}

func GetRecipes(writer http.ResponseWriter, request *http.Request) {
	foundRecipes := models.GetRecipes()

	resp := utils.PrepareReturn()
	resp["data"] = foundRecipes
	utils.Respond(writer, resp)
}

func UpdateRecipe(writer http.ResponseWriter, request *http.Request) {
	recipe := &models.Recipe{}

	// recipeID := mux.Vars(request)["id"]

	// id, err := strconv.Atoi(recipeID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// id64 := uint(id)

	err := json.NewDecoder(request.Body).Decode(recipe)
	if err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage("Error while decoding request body"))

		return
	}

	if err = recipe.Validate(); err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage(err.Error()))

		return
	}

	resp := models.UpdateRecipe(recipe)
	utils.Respond(writer, resp)
}

func DeleteRecipe(writer http.ResponseWriter, request *http.Request) {
	recipe := &models.Recipe{}

	recipeID := mux.Vars(request)["id"]

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		fmt.Println(err)
	}

	id64 := uint(id)

	err = json.NewDecoder(request.Body).Decode(recipe)
	if err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage("Error while decoding request body"))

		return
	}

	if err = recipe.Validate(); err != nil {
		utils.RespondWithError(writer, 500, utils.ErrorMessage(err.Error()))

		return
	}

	recipe.ID = id64

	resp := models.DeleteRecipe(recipe)
	utils.Respond(writer, resp)
}
