package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/peppermint-recipes/peppermint-server/utils"
)

type Recipe struct {
	gorm.Model
	Name         string `gorm:"" json:"name"`
	Yield        string `gorm:"" json:"yield"`
	ActiveTime   string `gorm:"" json:"activeTime"`
	TotalTime    string `gorm:"" json:"totalTime"`
	Ingredients  string `gorm:"" json:"ingredients"`
	Instructions string `gorm:"" json:"instructions"`
}

func (recipe *Recipe) Validate() error {
	if recipe.Name == "" {
		return errors.New("Recipe name should be on the payload")
	}
	if recipe.Yield == "" {
		return errors.New("Recipe yield should be on the payload")
	}
	if recipe.ActiveTime == "" {
		return errors.New("Recipe activeTime should be on the payload")
	}
	if recipe.TotalTime == "" {
		return errors.New("Recipe totalTime should be on the payload")
	}
	if recipe.Ingredients == "" {
		return errors.New("Recipe ingredients should be on the payload")
	}
	if recipe.Instructions == "" {
		return errors.New("Recipe instructions should be on the payload")
	}

	return nil
}

func (recipe *Recipe) Create() map[string]interface{} {
	db.Create(recipe)

	response := utils.PrepareReturn()
	response["recipe"] = recipe
	return response
}

func GetRecipe(id uint) (error, *Recipe) {
	recipe := &Recipe{}
	db.First(&recipe, id)

	if recipe.ID == 0 {
		return errors.New("no recipe found"), nil
	}

	return nil, recipe
}

func GetRecipes() []Recipe {
	recipes := []Recipe{}
	db.Find(&recipes)

	return recipes
}

func UpdateRecipe(recipe *Recipe) map[string]interface{} {
	db.Save(&recipe)
	response := utils.PrepareReturn()
	response["recipe"] = recipe

	return response
}

func DeleteRecipe(recipe *Recipe) map[string]interface{} {
	db.Delete(&recipe)
	response := utils.PrepareReturn()

	return response
}
