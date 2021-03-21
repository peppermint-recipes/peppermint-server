package weekplan

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
	errWeekPlanIsNotValid = errors.New("weekplan is not valid")
)

type weekplanServer struct {
	// store *taskstore.TaskStore
	mongoClient *mongo.Client
}

func NewWeekplanServer() *weekplanServer {
	mongoClient, _, _ := database.GetConnection()
	return &weekplanServer{mongoClient: mongoClient}
}

func (rs *weekplanServer) GetAllWeekplansHandler(c *gin.Context) {
	var loadedWeekplans, err = getAllWeekplans()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"weekplans": loadedWeekplans})
}

func (rs *weekplanServer) GetWeekplanByIDHandler(c *gin.Context) {
	weekplanID := c.Param("id")

	var loadedWeekplan, err = getWeekplanByID(weekplanID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}
	c.JSON(http.StatusOK, gin.H{"weekplan": loadedWeekplan})
}

func (rs *weekplanServer) CreateWeekplanHandler(c *gin.Context) {
	var weekplan weekPlan

	fmt.Printf("%v", c)
	if err := c.ShouldBindJSON(&weekplan); err != nil {
		c.String(http.StatusBadRequest, err.Error())

		return
	}

	if !weekplan.isValid() {
		c.String(http.StatusBadRequest, errWeekPlanIsNotValid.Error())

		return
	}

	weekplan.LastUpdated = time.Now()

	id, err := createWeekplan(&weekplan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (rs *weekplanServer) UpdateWeekplanHandler(c *gin.Context) {
	var weekplan weekPlan

	if err := c.ShouldBindJSON(&weekplan); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	weekplan.LastUpdated = time.Now()

	savedWeekplan, err := updateWeekplan(&weekplan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"weekplan": savedWeekplan})
}

func (rs *weekplanServer) DeleteWeekplanHandler(c *gin.Context) {
	weekplanID := c.Param("id")

	err := deleteWeekplan(weekplanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	c.JSON(http.StatusOK, gin.H{"id": weekplanID})
}
