package main

import (
	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/recipe"
	shoppinglist "github.com/peppermint-recipes/peppermint-server/shopping-list"
	"github.com/peppermint-recipes/peppermint-server/weekplan"
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
	weekplanServer := weekplan.NewWeekplanServer()
	shoppingListServer := shoppinglist.NewShoppingListServer()

	router.GET("/recipes/:id", recipeServer.GetRecipeByIDHandler)
	router.GET("/recipes/", recipeServer.GetAllRecipesHandler)
	router.POST("/recipes/", recipeServer.CreateRecipeHandler)
	router.PUT("/recipes/", recipeServer.UpdateRecipeHandler)
	router.DELETE("/recipes/", recipeServer.DeleteRecipeHandler)

	router.GET("/weekplans/:id", weekplanServer.GetWeekplanByIDHandler)
	router.GET("/weekplans/", weekplanServer.GetAllWeekplansHandler)
	router.POST("/weekplans/", weekplanServer.CreateWeekplanHandler)
	router.PUT("/weekplans/", weekplanServer.UpdateWeekplanHandler)
	router.DELETE("/weekplans/", weekplanServer.DeleteWeekplanHandler)

	router.GET("/shopping-lists/:id", shoppingListServer.GetShoppingListsByIDHandler)
	router.GET("/shopping-lists/", shoppingListServer.GetAllWeekplansHandler)
	router.POST("/shopping-lists/", shoppingListServer.CreateWeekplanHandler)
	router.PUT("/shopping-lists/", shoppingListServer.UpdateWeekplanHandler)
	router.DELETE("/shopping-lists/", shoppingListServer.DeleteWeekplanHandler)

	router.Run(config.Web.Address + ":" + config.Web.Port)
}
