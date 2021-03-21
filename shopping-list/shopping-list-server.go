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
	// store *taskstore.TaskStore
	mongoClient *mongo.Client
}

func NewShoppingListServer() *shoppingListServer {
	mongoClient, _, _ := database.GetConnection()
	return &shoppingListServer{mongoClient: mongoClient}
}

func (sl *shoppingListServer) GetAllWeekplansHandler(c *gin.Context) {
	var loadedShoppingLists, err = getAllShoppingLists()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shoppingLists": loadedShoppingLists})
}

func (sl *shoppingListServer) GetShoppingListsByIDHandler(c *gin.Context) {
	shoppingListID := c.Param("id")

	var loadedWeekplan, err = getShoppingListByID(shoppingListID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}
	c.JSON(http.StatusOK, gin.H{"shoppingList": loadedWeekplan})
}

func (sl *shoppingListServer) CreateWeekplanHandler(c *gin.Context) {
	var shoppingList shoppingList

	fmt.Printf("%v", c)
	if err := c.ShouldBindJSON(&shoppingList); err != nil {
		c.String(http.StatusBadRequest, err.Error())

		return
	}

	if !shoppingList.isValid() {
		c.String(http.StatusBadRequest, errShoppingListIsNotValid.Error())

		return
	}

	shoppingList.LastUpdated = time.Now()

	id, err := createShoppingList(&shoppingList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (sl *shoppingListServer) UpdateWeekplanHandler(c *gin.Context) {
	var shoppingList shoppingList

	if err := c.ShouldBindJSON(&shoppingList); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	shoppingList.LastUpdated = time.Now()

	savedWeekplan, err := updateShoppingList(&shoppingList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"shoppingList": savedWeekplan})
}

func (sl *shoppingListServer) DeleteWeekplanHandler(c *gin.Context) {
	shoppingListID := c.Param("id")

	err := deleteShoppingList(shoppingListID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	c.JSON(http.StatusOK, gin.H{"id": shoppingListID})
}
