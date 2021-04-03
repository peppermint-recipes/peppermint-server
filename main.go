package main

import (
	"log"

	"net/http"

	"github.com/peppermint-recipes/peppermint-server/auth"

	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/database"
	"github.com/peppermint-recipes/peppermint-server/recipe"
	shoppinglist "github.com/peppermint-recipes/peppermint-server/shopping-list"
	"github.com/peppermint-recipes/peppermint-server/user"
	"github.com/peppermint-recipes/peppermint-server/weekplan"

	jwt "github.com/appleboy/gin-jwt/v2"
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

func setupServer(dbConfig *config.DBConfig, JWTSigningKey string) *gin.Engine {
	database.RegisterConnection(dbConfig.Username, dbConfig.Password, dbConfig.Endpoint)

	router := gin.Default()
	router.Use(CORSMiddleware)
	recipeServer := recipe.NewRecipeServer()
	weekplanServer := weekplan.NewWeekplanServer()
	shoppingListServer := shoppinglist.NewShoppingListServer()
	userServer := user.NewUserServer()

	router.GET("/livez", livezHandler)

	authMiddleware, err := auth.RegisterAuthMiddleware(JWTSigningKey)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	router.POST("/login", authMiddleware.LoginHandler)
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := router.Group("/auth")
	// Refresh time can be longer than token timeout, therefore it must be registered BEFORE the auth middleware.
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())

	router.POST("/register", userServer.CreateUserHandler)

	recipes := router.Group("/recipes")
	recipes.Use(authMiddleware.MiddlewareFunc())
	recipes.GET("/:id", recipeServer.GetRecipeByIDHandler)
	recipes.GET("/", recipeServer.GetAllRecipesHandler)
	recipes.POST("/", recipeServer.CreateRecipeHandler)
	recipes.PUT("/", recipeServer.UpdateRecipeHandler)
	recipes.DELETE("/:id", recipeServer.DeleteRecipeHandler)

	weekplans := router.Group("/weekplans")
	weekplans.Use(authMiddleware.MiddlewareFunc())
	weekplans.GET("/:id", weekplanServer.GetWeekplanByIDHandler)
	weekplans.GET("/", weekplanServer.GetAllWeekplansHandler)
	weekplans.POST("/", weekplanServer.CreateWeekplanHandler)
	weekplans.PUT("/", weekplanServer.UpdateWeekplanHandler)
	weekplans.DELETE("/:id", weekplanServer.DeleteWeekplanHandler)

	shoppingLists := router.Group("/shopping-lists")
	shoppingLists.Use(authMiddleware.MiddlewareFunc())
	shoppingLists.GET("/:id", shoppingListServer.GetShoppingListsByIDHandler)
	shoppingLists.GET("/", shoppingListServer.GetAllWeekplansHandler)
	shoppingLists.POST("/", shoppingListServer.CreateWeekplanHandler)
	shoppingLists.PUT("/", shoppingListServer.UpdateWeekplanHandler)
	shoppingLists.DELETE("/:id", shoppingListServer.DeleteWeekplanHandler)

	return router
}

func main() {
	config := config.GetConfig()
	setupServer(
		config.DB,
		config.Web.JWTSigningKey,
	).Run(config.Web.Address + ":" + config.Web.Port)
}
