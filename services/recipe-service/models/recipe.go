package models

import (
	"github.com/jinzhu/gorm"
	"github.com/theErikss0n/peppermint-server/services/recipe-service/utils"
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

func (recipe *Recipe) Validate() (map[string]interface{}, bool) {

	if recipe.Name == "" {
		return utils.Message(false, "Recipe name should be on the payload"), false
	}

	if recipe.Yield == "" {
		return utils.Message(false, "Recipe yield should be on the payload"), false
	}

	if recipe.ActiveTime == "" {
		return utils.Message(false, "Recipe activeTime should be on the payload"), false
	}

	if recipe.TotalTime == "" {
		return utils.Message(false, "Recipe totalTime should be on the payload"), false
	}

	if recipe.Ingredients == "" {
		return utils.Message(false, "Recipe ingredients should be on the payload"), false
	}

	if recipe.Instructions == "" {
		return utils.Message(false, "Recipe instructions should be on the payload"), false
	}

	//All the required parameters are present
	return utils.Message(true, "success"), true
}

func (recipe *Recipe) Create() map[string]interface{} {

	if resp, ok := recipe.Validate(); !ok {
		return resp
	}

	db.Create(recipe)

	resp := utils.Message(true, "success")
	resp["recipe"] = recipe
	return resp
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
	if resp, ok := recipe.Validate(); !ok {
		return resp
	}

	db.Save(&recipe)
	resp := utils.Message(true, "success")
	resp["recipe"] = recipe

	return resp
}

func DeleteRecipe(recipe *Recipe) map[string]interface{} {
	if resp, ok := recipe.Validate(); !ok {
		return resp
	}

	db.Delete(&recipe)
	resp := utils.Message(true, "success")

	return resp
}
