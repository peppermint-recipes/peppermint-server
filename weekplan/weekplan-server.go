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
	mongoClient *mongo.Client
}

func NewWeekplanServer() *weekplanServer {
	mongoClient, _, _ := database.GetConnection()
	return &weekplanServer{mongoClient: mongoClient}
}

func (ws *weekplanServer) GetAllWeekplansHandler(context *gin.Context) {
	var loadedWeekplans, err = getAllWeekplans()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"weekplans": loadedWeekplans})
}

func (ws *weekplanServer) GetWeekplanByIDHandler(context *gin.Context) {
	weekplanID := context.Param("id")

	var loadedWeekplan, err = getWeekplanByID(weekplanID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}
	context.JSON(http.StatusOK, gin.H{"weekplan": loadedWeekplan})
}

func (ws *weekplanServer) CreateWeekplanHandler(context *gin.Context) {
	var weekplan weekPlan

	fmt.Printf("%v", context)
	if err := context.ShouldBindJSON(&weekplan); err != nil {
		context.String(http.StatusBadRequest, err.Error())

		return
	}

	if !weekplan.isValid() {
		context.String(http.StatusBadRequest, errWeekPlanIsNotValid.Error())

		return
	}

	weekplan.LastUpdated = time.Now()

	id, err := createWeekplan(&weekplan)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}
	context.JSON(http.StatusOK, gin.H{"id": id})
}

func (ws *weekplanServer) UpdateWeekplanHandler(context *gin.Context) {
	var weekplan weekPlan

	if err := context.ShouldBindJSON(&weekplan); err != nil {
		log.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	weekplan.LastUpdated = time.Now()

	savedWeekplan, err := updateWeekplan(&weekplan)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"weekplan": savedWeekplan})
}

func (ws *weekplanServer) DeleteWeekplanHandler(context *gin.Context) {
	weekplanID := context.Param("id")

	err := deleteWeekplan(weekplanID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, gin.H{"id": weekplanID})
}
