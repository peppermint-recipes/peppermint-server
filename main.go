package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/database"
	"github.com/peppermint-recipes/peppermint-server/recipe"
	shoppinglist "github.com/peppermint-recipes/peppermint-server/shopping-list"
	"github.com/peppermint-recipes/peppermint-server/weekplan"
)

func livezHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "healthy",
	})
}

// Based on https://asanchez.dev/blog/cors-golang-options/
func CORSMiddleware(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "*")
	context.Header("Access-Control-Allow-Headers", "*")
	context.Header("Content-Type", "application/json; charset=utf-8")

	if context.Request.Method != "OPTIONS" {
		context.Next()
	} else {
		context.AbortWithStatus(http.StatusNoContent)
	}
}

func setupServer(dbConfig *config.DBConfig) *gin.Engine {
	database.RegisterConnection(dbConfig.Username, dbConfig.Password, dbConfig.Endpoint)

	router := gin.Default()
	router.Use(CORSMiddleware)
	recipeServer := recipe.NewRecipeServer()
	weekplanServer := weekplan.NewWeekplanServer()
	shoppingListServer := shoppinglist.NewShoppingListServer()

	router.GET("/livez", livezHandler)

	router.GET("/recipes/:id", recipeServer.GetRecipeByIDHandler)
	router.GET("/recipes/", recipeServer.GetAllRecipesHandler)
	router.POST("/recipes/", recipeServer.CreateRecipeHandler)
	router.PUT("/recipes/", recipeServer.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipeServer.DeleteRecipeHandler)

	router.GET("/weekplans/:id", weekplanServer.GetWeekplanByIDHandler)
	router.GET("/weekplans/", weekplanServer.GetAllWeekplansHandler)
	router.POST("/weekplans/", weekplanServer.CreateWeekplanHandler)
	router.PUT("/weekplans/", weekplanServer.UpdateWeekplanHandler)
	router.DELETE("/weekplans/:id", weekplanServer.DeleteWeekplanHandler)

	router.GET("/shopping-lists/:id", shoppingListServer.GetShoppingListsByIDHandler)
	router.GET("/shopping-lists/", shoppingListServer.GetAllWeekplansHandler)
	router.POST("/shopping-lists/", shoppingListServer.CreateWeekplanHandler)
	router.PUT("/shopping-lists/", shoppingListServer.UpdateWeekplanHandler)
	router.DELETE("/shopping-lists/:id", shoppingListServer.DeleteWeekplanHandler)

	return router
}

func main() {
	config := config.GetConfig()
	setupServer(config.DB).Run(config.Web.Address + ":" + config.Web.Port)
}
