package main

import (
	"log"

	"github.com/peppermint-recipes/peppermint-server/auth"
	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/database"
	"github.com/peppermint-recipes/peppermint-server/recipe"
	shoppinglist "github.com/peppermint-recipes/peppermint-server/shopping-list"
	"github.com/peppermint-recipes/peppermint-server/user"
	"github.com/peppermint-recipes/peppermint-server/weekplan"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func livezHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "healthy",
	})
}

func setupServer(dbConfig *config.DBConfig, JWTSigningKey string) *gin.Engine {
	database.RegisterConnection(dbConfig.Username, dbConfig.Password, dbConfig.Endpoint)

	router := gin.Default()
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
	auth.GET("/hello", helloHandler)

	router.POST("/register", userServer.CreateUserHandler)

	recipes := router.Group("/recipes")
	recipes.Use(authMiddleware.MiddlewareFunc())
	recipes.GET("/:id", recipeServer.GetRecipeByIDHandler)
	recipes.GET("/", recipeServer.GetAllRecipesHandler)
	recipes.POST("/", recipeServer.CreateRecipeHandler)
	recipes.PUT("/", recipeServer.UpdateRecipeHandler)
	recipes.DELETE("/:id", recipeServer.DeleteRecipeHandler)

	weekplans := router.Group("/weekplans")
	// weekplans.Use(authMiddleware.MiddlewareFunc())
	weekplans.GET("/:id", weekplanServer.GetWeekplanByIDHandler)
	weekplans.GET("/", weekplanServer.GetAllWeekplansHandler)
	weekplans.POST("/", weekplanServer.CreateWeekplanHandler)
	weekplans.PUT("/", weekplanServer.UpdateWeekplanHandler)
	weekplans.DELETE("/:id", weekplanServer.DeleteWeekplanHandler)

	router.GET("/shopping-lists/:id", shoppingListServer.GetShoppingListsByIDHandler)
	router.GET("/shopping-lists/", shoppingListServer.GetAllWeekplansHandler)
	router.POST("/shopping-lists/", shoppingListServer.CreateWeekplanHandler)
	router.PUT("/shopping-lists/", shoppingListServer.UpdateWeekplanHandler)
	router.DELETE("/shopping-lists/:id", shoppingListServer.DeleteWeekplanHandler)

	return router
}

func main() {
	config := config.GetConfig()
	setupServer(
		config.DB,
		config.Web.JWTSigningKey,
	).Run(config.Web.Address + ":" + config.Web.Port)
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}
