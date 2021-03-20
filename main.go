package main

import (
	"log"

	"github.com/peppermint-recipes/peppermint-server/auth"
	"github.com/peppermint-recipes/peppermint-server/config"
	"github.com/peppermint-recipes/peppermint-server/database"
	"github.com/peppermint-recipes/peppermint-server/recipe"
	shoppinglist "github.com/peppermint-recipes/peppermint-server/shopping-list"
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

	router.GET("/livez", livezHandler)

	authMiddleware, err := auth.RegisterAuthMiddleware(JWTSigningKey)

	// authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
	// 	Realm:       "peppermint-server",
	// 	Key:         []byte(JWTSigningKey),
	// 	Timeout:     time.Hour,
	// 	MaxRefresh:  time.Hour,
	// 	IdentityKey: identityKey,
	// 	PayloadFunc: func(data interface{}) jwt.MapClaims {
	// 		if v, ok := data.(*User); ok {
	// 			return jwt.MapClaims{
	// 				identityKey: v.UserName,
	// 			}
	// 		}
	// 		return jwt.MapClaims{}
	// 	},
	// 	IdentityHandler: func(c *gin.Context) interface{} {
	// 		claims := jwt.ExtractClaims(c)
	// 		return &User{
	// 			UserName: claims[identityKey].(string),
	// 		}
	// 	},
	// 	// Authenticator: auth.Authenticator,
	// 	// Authorizator: auth.Authorizator,
	// 	Unauthorized: auth.Unauthorized,
	// 	Authenticator: func(c *gin.Context) (interface{}, error) {
	// 		var loginVals login
	// 		if err := c.ShouldBind(&loginVals); err != nil {
	// 			return "", jwt.ErrMissingLoginValues
	// 		}
	// 		userID := loginVals.Username
	// 		password := loginVals.Password

	// 		if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
	// 			return &User{
	// 				UserName:  userID,
	// 				LastName:  "Bo-Yi",
	// 				FirstName: "Wu",
	// 			}, nil
	// 		}

	// 		return nil, jwt.ErrFailedAuthentication
	// 	},
	// 	Authorizator: func(data interface{}, c *gin.Context) bool {
	// 		if v, ok := data.(*User); ok && v.UserName == "admin" {
	// 			return true
	// 		}

	// 		return false
	// 	},

	// 	// Unauthorized: func(c *gin.Context, code int, message string) {
	// 	// 	c.JSON(code, gin.H{
	// 	// 		"code":    code,
	// 	// 		"message": message,
	// 	// 	})
	// 	// },
	// 	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// 	// to extract token from the request.
	// 	// Optional. Default value "header:Authorization".
	// 	// Possible values:
	// 	// - "header:<name>"
	// 	// - "query:<name>"
	// 	// - "cookie:<name>"
	// 	// - "param:<name>"
	// 	TokenLookup: "header: Authorization, query: token, cookie: jwt",
	// 	// TokenLookup: "query:token",
	// 	// TokenLookup: "cookie:token",

	// 	// TokenHeadName is a string in the header. Default value is "Bearer"
	// 	TokenHeadName: "Bearer",

	// 	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	// 	TimeFunc: time.Now,
	// })

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
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
	}

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
