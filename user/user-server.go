package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errUserIsNotValid = errors.New("user is not valid")
)

type userServer struct {
	mongoClient *mongo.Client
}

func NewUserServer() *userServer {
	mongoClient, _, _ := database.GetConnection()
	return &userServer{mongoClient: mongoClient}
}

// func (rs *recipeServer) GetRecipeByIDHandler(context *gin.Context) {
// 	recipeID := context.Param("id")

// 	var loadedRecipe, err = getRecipeByID(recipeID)
// 	if err != nil {
// 		context.JSON(http.StatusNotFound, gin.H{"message": err})

// 		return
// 	}
// 	context.JSON(http.StatusOK, gin.H{"Recipe": loadedRecipe})
// }

func (us *userServer) CreateUserHandler(context *gin.Context) {
	var user User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.String(http.StatusBadRequest, err.Error())

		return
	}

	if !user.isValid() {
		context.String(http.StatusBadRequest, errUserIsNotValid.Error())

		return
	}

	user.LastUpdated = time.Now()

	createdUser, err := createUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"recipe": createdUser})
}

// func (rs *recipeServer) UpdateRecipeHandler(context *gin.Context) {
// 	var recipe Recipe
// 	if err := context.ShouldBindJSON(&recipe); err != nil {
// 		log.Print(err)
// 		context.JSON(http.StatusBadRequest, gin.H{"message": err})
// 		return
// 	}

// 	recipe.LastUpdated = time.Now()

// 	savedRecipe, err := updateRecipe(&recipe)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
// 		return
// 	}
// 	context.JSON(http.StatusOK, gin.H{"recipe": savedRecipe})
// }

// func (rs *recipeServer) DeleteRecipeHandler(context *gin.Context) {
// 	recipeID := context.Param("id")

// 	deletedRecipe, err := deleteRecipe(recipeID)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
// 		return
// 	}
// 	context.JSON(http.StatusOK, gin.H{"recipe": deletedRecipe})
// }
