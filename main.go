package main

import (
	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/recipe"
)

// func main() {
// 	router := gin.Default()
// 	server := NewTaskServer()

// 	router.POST("/task/", server.createTaskHandler)
// 	router.GET("/task/", server.getAllTasksHandler)
// 	router.DELETE("/task/", server.deleteAllTasksHandler)
// 	router.GET("/task/:id", server.getTaskHandler)
// 	router.DELETE("/task/:id", server.deleteTaskHandler)
// 	router.GET("/tag/:tag", server.tagHandler)
// 	router.GET("/due/:year/:month/:day", server.dueHandler)

// 	router.Run("localhost:" + os.Getenv("SERVERPORT"))
// }

func main() {

	config := config.GetConfig()
	router := gin.Default()
	recipeServer := recipe.NewRecipeServer()

	router.GET("/recipes/", recipeServer.GetAllRecipesHandler)
	router.GET("/recipes/:id", recipeServer.GetRecipeByIDHandler)
	router.POST("/recipes/", recipeServer.CreateRecipeHandler)
	router.PUT("/recipes/", recipeServer.UpdateRecipeHandler)
	router.DELETE("/recipes/", recipeServer.DeleteRecipeHandler)

	router.Run(config.Web.Address + ":" + config.Web.Port) // listen and serve on 0.0.0.0:8080
}
