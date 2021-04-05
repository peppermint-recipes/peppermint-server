package shoppinglist

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type shoppingListItem struct {
	Ingredient string `json:"ingredient"`
	Unit       string `json:"unit"`
	Amount     int    `json:"amount"`
}

// ShoppingList is a model of a shopping list
type shoppingList struct {
	ID          primitive.ObjectID `json:"id"`
	UserID      string             `json:"userId"`
	Items       []shoppingListItem `json:"items"`
	Deleted     bool               `json:"deleted"`
	LastUpdated time.Time          `json:"lastUpdated"`
}

func (sl *shoppingList) isValid() bool {
	return sl.UserID != ""
}
