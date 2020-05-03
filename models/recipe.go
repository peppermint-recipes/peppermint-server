package models

import (
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

func (recipe *Recipe) Validate() (string, bool) {
	if recipe.Name == "" {
		return "Recipe name should be on the payload", false
	}

	if recipe.Yield == "" {
		return "Recipe yield should be on the payload", false
	}

	if recipe.ActiveTime == "" {
		return "Recipe activeTime should be on the payload", false
	}

	if recipe.TotalTime == "" {
		return "Recipe totalTime should be on the payload", false
	}

	if recipe.Ingredients == "" {
		return "Recipe ingredients should be on the payload", false
	}

	if recipe.Instructions == "" {
		return "Recipe instructions should be on the payload", false
	}

	//All the required parameters are present
	return "success", true
}

func (recipe *Recipe) Create() map[string]interface{} {
	db.Create(recipe)

	response := utils.PrepareReturn()
	response["recipe"] = recipe
	return response
}

func GetRecipe(id uint) *Recipe {
	recipe := &Recipe{}
	db.First(&recipe, id)

	return recipe
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
