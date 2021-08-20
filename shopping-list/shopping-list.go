package shoppinglist

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShoppingListItem struct {
	Ingredient string `json:"ingredient"`
	Unit       string `json:"unit"`
	Amount     int    `json:"amount"`
}

// ShoppingList is a model of a shopping list
type ShoppingList struct {
	ID          primitive.ObjectID `json:"id"`
	UserID      string             `json:"userId"`
	Items       []ShoppingListItem `json:"items"`
	Deleted     bool               `json:"deleted"`
	LastUpdated time.Time          `json:"lastUpdated"`
}

func (sl *ShoppingList) isValid() bool {
	return sl.UserID != ""
}
