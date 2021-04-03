package shoppinglist

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/peppermint-recipes/peppermint-server/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errShoppingListIsNotValid = errors.New("shopping list is not valid")
)

type shoppingListServer struct {
	mongoClient *mongo.Client
}

func NewShoppingListServer() *shoppingListServer {
	mongoClient, _, _ := database.GetConnection()

	return &shoppingListServer{mongoClient: mongoClient}
}

func (sl *shoppingListServer) GetAllWeekplansHandler(context *gin.Context) {
	var shoppingLists []*shoppingList

	shoppingLists, err := getAllShoppingLists()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	// Return [] instead of null, if no elements found.
	if len(shoppingLists) == 0 {
		shoppingLists := make([]shoppingList, 0)
		context.JSON(http.StatusOK, shoppingLists)

		return
	}

	context.JSON(http.StatusOK, shoppingLists)
}

func (sl *shoppingListServer) GetShoppingListsByIDHandler(context *gin.Context) {
	shoppingListID := context.Param("id")

	var loadedWeekplan, err = getShoppingListByID(shoppingListID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, loadedWeekplan)
}

func (sl *shoppingListServer) CreateWeekplanHandler(context *gin.Context) {
	var shoppingList shoppingList

	fmt.Printf("%v", context)
	if err := context.ShouldBindJSON(&shoppingList); err != nil {
		context.String(http.StatusBadRequest, err.Error())

		return
	}

	if !shoppingList.isValid() {
		context.String(http.StatusBadRequest, errShoppingListIsNotValid.Error())

		return
	}

	shoppingList.LastUpdated = time.Now()

	createdShoppingList, err := createShoppingList(&shoppingList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, createdShoppingList)
}

func (sl *shoppingListServer) UpdateWeekplanHandler(context *gin.Context) {
	var shoppingList shoppingList

	if err := context.ShouldBindJSON(&shoppingList); err != nil {
		log.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": err})

		return
	}

	shoppingList.LastUpdated = time.Now()

	savedWeekplan, err := updateShoppingList(&shoppingList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, savedWeekplan)
}

func (sl *shoppingListServer) DeleteWeekplanHandler(context *gin.Context) {
	shoppingListID := context.Param("id")

	deletedShoppingList, err := deleteShoppingList(shoppingListID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, deletedShoppingList)
}
