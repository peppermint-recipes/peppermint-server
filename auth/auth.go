package auth

import (
	jwt "github.com/appleboy/gin-jwt"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserName  string
	FirstName string
	LastName  string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &User{
			UserName:  userID,
			LastName:  "Bo-Yi",
			FirstName: "Wu",
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok && v.UserName == "admin" {
		return true
	}

	return false
}

// func PayloadFunc(data interface{}) jwt.MapClaims {
// 	if v, ok := data.(*User); ok {
// 		return jwt.MapClaims{
// 			identityKey: v.UserName,
// 		}
// 	}
// 	return jwt.MapClaims{}
// }

// func IdentityHandler(c *gin.Context) interface{} {
// 	claims := jwt.ExtractClaims(c)
// 	return &User{
// 		UserName: claims[identityKey].(string),
// 	}
// }

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
