package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is a model of a user
type User struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `form:"name" json:"name" binding:"required"`
	Password    string             `form:"password" json:"password" binding:"required"`
	Deleted     bool               `json:"deleted"`
	LastUpdated time.Time          `json:"lastUpdated"`
}

func (user *User) isValid() bool {
	if user.Name == "" {
		return false
	}

	return true
}
