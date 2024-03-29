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
	var weekplans []*WeekPlan
	weekplans, err := getAllWeekplans()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	// Return [] instead of null, if no elements found.
	if len(weekplans) == 0 {
		weekplans := make([]WeekPlan, 0)
		context.JSON(http.StatusOK, weekplans)

		return
	}

	context.JSON(http.StatusOK, weekplans)
}

func (ws *weekplanServer) GetWeekplanByIDHandler(context *gin.Context) {
	weekplanID := context.Param("id")

	var loadedWeekplan, err = getWeekplanByID(weekplanID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, loadedWeekplan)
}

func (ws *weekplanServer) CreateWeekplanHandler(context *gin.Context) {
	var weekplan WeekPlan

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

	createdWeekplan, err := createWeekplan(&weekplan)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, createdWeekplan)
}

func (ws *weekplanServer) UpdateWeekplanHandler(context *gin.Context) {
	var weekplan WeekPlan

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

	context.JSON(http.StatusOK, savedWeekplan)
}

func (ws *weekplanServer) DeleteWeekplanHandler(context *gin.Context) {
	weekplanID := context.Param("id")

	deletedWeekPlan, err := deleteWeekplan(weekplanID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, deletedWeekPlan)
}
