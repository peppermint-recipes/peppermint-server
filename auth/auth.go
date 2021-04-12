package auth

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/peppermint-recipes/peppermint-server/user"
)

type User struct {
	UserName string
	UserID   string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

var IdentityKey = "id"

func RegisterAuthMiddleware(JWTSigningKey string) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "peppermint-server",
		Key:         []byte(JWTSigningKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKey: string(v.UserID),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserID: claims[IdentityKey].(string),
			}
		},
		Unauthorized: Unauthorized,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			user, err := user.IsUserAuthorized(userID, password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &User{
				UserName: user.Name,
				UserID:   user.ID.String(),
			}, nil
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value.
		// This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	return authMiddleware, err
}
