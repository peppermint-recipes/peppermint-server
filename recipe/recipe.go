package recipe

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe is a model of a recipe
type Recipe struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	ActiveTime   int                `json:"activeTime"`
	TotalTime    int                `json:"totalTime"`
	Servings     int                `json:"servings"`
	Categories   []string           `json:"categories"`
	Ingredients  []string           `json:"ingredients"`
	Instructions string             `json:"instructions"`
	UserID       string             `json:"userId"`
	Deleted      bool               `json:"deleted"`
	LastUpdated  time.Time          `json:"lastUpdated"`
	Calories     int                `json:"calories"`
}

func (recipe *Recipe) isValid() bool {
	if recipe.UserID == "" {
		return false
	}

	return true
}
