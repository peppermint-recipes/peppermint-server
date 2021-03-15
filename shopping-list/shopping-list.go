package shoppinglist

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type shoppingListItem struct {
	Ingredient string `json:"ingredient,omitempty"`
	Unit       string `json:"unit,omitempty"`
	Amount     int    `json:"amount,omitempty"`
}

// ShoppingList is a model of a shopping list
type shoppingList struct {
	ID     primitive.ObjectID `json:"id,omitempty"`
	UserID string             `json:"user_id,omitempty"`
	Items  []shoppingListItem `json:"items,omitempty"`
}

func (sl *shoppingList) isValid() bool {
	if sl.UserID == "" {
		return false
	}

	return true
}