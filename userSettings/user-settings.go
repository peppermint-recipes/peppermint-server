package userSettings

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserSettings is a model of a user settings
type UserSettings struct {
	ID          primitive.ObjectID `json:"id"`
	UserID      string             `json:"userId"`
	Deleted     bool               `json:"deleted"`
	LastUpdated time.Time          `json:"lastUpdated"`
}

func (userSettings *UserSettings) isValid() bool {
	if userSettings.UserID == "" {
		return false
	}

	return true
}
