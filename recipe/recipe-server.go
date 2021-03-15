package recipe

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errRecipeIsNotValid = errors.New("recipe is not valid")
)

type recipeServer struct {
	// store *taskstore.TaskStore
	mongoClient *mongo.Client
}

func NewRecipeServer() *recipeServer {
	mongoClient, _, _ := database.GetConnection()
	return &recipeServer{mongoClient: mongoClient}
}

func (rs *recipeServer) GetAllRecipesHandler(c *gin.Context) {
	var loadedTasks, err = getAllRecipes()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipes": loadedTasks})
}

func (rs *recipeServer) GetRecipeByIDHandler(c *gin.Context) {
	recipeID := c.Param("id")

	var loadedRecipe, err = getRecipeByID(recipeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}
	c.JSON(http.StatusOK, gin.H{"Recipe": loadedRecipe})
}

func (rs *recipeServer) CreateRecipeHandler(c *gin.Context) {
	var recipe Recipe
	fmt.Printf("%v", c)
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.String(http.StatusBadRequest, err.Error())

		return
	}

	if !recipe.isValid() {
		c.String(http.StatusBadRequest, errRecipeIsNotValid.Error())

		return
	}

	id, err := createRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (rs *recipeServer) UpdateRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	savedRecipe, err := updateRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipe": savedRecipe})
}

// TODO: Fix
func (rs *recipeServer) DeleteRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	err := deleteRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
}
