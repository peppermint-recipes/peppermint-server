package recipe

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errRecipeIsNotValid = errors.New("recipe is not valid")
)

type recipeServer struct {
	mongoClient *mongo.Client
}

func NewRecipeServer() *recipeServer {
	mongoClient, _, _ := database.GetConnection()

	return &recipeServer{mongoClient: mongoClient}
}

func (rs *recipeServer) GetAllRecipesHandler(context *gin.Context) {
	var recipes []*Recipe

	recipes, err := getAllRecipes()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	// Return [] instead of null, if no elements found.
	if len(recipes) == 0 {
		recipes := make([]Recipe, 0)
		context.JSON(http.StatusOK, recipes)

		return
	}

	context.JSON(http.StatusOK, recipes)
}

func (rs *recipeServer) GetRecipeByIDHandler(context *gin.Context) {
	recipeID := context.Param("id")

	var loadedRecipe, err = getRecipeByID(recipeID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, loadedRecipe)
}

func (rs *recipeServer) CreateRecipeHandler(context *gin.Context) {
	var recipe Recipe

	if err := context.ShouldBindJSON(&recipe); err != nil {
		context.String(http.StatusBadRequest, err.Error())

		return
	}

	if !recipe.isValid() {
		context.String(http.StatusBadRequest, errRecipeIsNotValid.Error())

		return
	}

	recipe.LastUpdated = time.Now()

	createdRecipe, err := createRecipe(&recipe)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, createdRecipe)
}

func (rs *recipeServer) UpdateRecipeHandler(context *gin.Context) {
	var recipe Recipe
	if err := context.ShouldBindJSON(&recipe); err != nil {
		log.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": err})

		return
	}

	recipe.LastUpdated = time.Now()

	savedRecipe, err := updateRecipe(&recipe)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, savedRecipe)
}

func (rs *recipeServer) DeleteRecipeHandler(context *gin.Context) {
	recipeID := context.Param("id")

	deletedRecipe, err := deleteRecipe(recipeID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, deletedRecipe)
}
