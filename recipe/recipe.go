package recipe

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recipe is a model of a recipe
type Recipe struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	ActiveTime   int                `json:"active_time,omitempty"`
	TotalTime    int                `json:"total_time,omitempty"`
	Servings     int                `json:"servings,omitempty"`
	Categories   []string           `json:"categories,omitempty"`
	Ingredients  []string           `json:"ingredients,omitempty"`
	Instructions string             `json:"instructions,omitempty"`
	UserID       string             `json:"user_id"`
	Deleted      bool               `json:"deleted,omitempty"`
	LastUpdated  time.Time          `json:"last_updated,omitempty"`
	Calories     int                `json:"calories,omitempty"`
}

func (recipe *Recipe) isValid() bool {
	if recipe.UserID == "" {
		return false
	}

	return true
}
