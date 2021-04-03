package userSettings

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

type userSettingsServer struct {
	mongoClient *mongo.Client
}

func NewUserSettingsServer() *userSettingsServer {
	mongoClient, _, _ := database.GetConnection()
	return &userSettingsServer{mongoClient: mongoClient}
}

func (us *userSettingsServer) GetUserSettingsForUserHandler(context *gin.Context) {
	userID := context.Param("id")

	var loadedUserSettings, err = getUserSettingsOfUser(userID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	context.JSON(http.StatusOK, gin.H{"User": loadedUserSettings})
}

func (us *userSettingsServer) CreateUserSettingsHandler(context *gin.Context) {
	var userSettings UserSettings

	if err := context.ShouldBindJSON(&userSettings); err != nil {
		context.String(http.StatusBadRequest, err.Error())

		return
	}

	if !userSettings.isValid() {
		context.String(http.StatusBadRequest, errUserIsNotValid.Error())

		return
	}

	userSettings.LastUpdated = time.Now()

	createdUser, err := createUserSettings(&userSettings)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"recipe": createdUser})
}

func (us *userSettingsServer) UpdateUserSettingsHandler(context *gin.Context) {
	var userSettings UserSettings
	if err := context.ShouldBindJSON(&userSettings); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	userSettings.LastUpdated = time.Now()

	savedUser, err := updateUserSettings(&userSettings)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": savedUser})
}

func (us *userSettingsServer) DeleteUserSettingsHandler(context *gin.Context) {
	userID := context.Param("id")

	deletedUser, err := deleteUserSettingsForUser(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": deletedUser})
}
