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

func (us *userServer) GetUserByIDHandler(context *gin.Context) {
	userID := context.Param("id")

	var loadedUser, err = getUserByID(userID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err})

		return
	}

	loadedUser.Password = "nope"

	context.JSON(http.StatusOK, gin.H{"User": loadedUser})
}

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

func (us *userServer) UpdateUserHandler(context *gin.Context) {
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	user.LastUpdated = time.Now()

	savedUser, err := updateUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	savedUser.Password = "nope"

	context.JSON(http.StatusOK, gin.H{"user": savedUser})
}

func (us *userServer) DeleteUserHandler(context *gin.Context) {
	userID := context.Param("id")

	deletedUser, err := deleteUser(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	deletedUser.Password = "nope"

	context.JSON(http.StatusOK, gin.H{"user": deletedUser})
}
