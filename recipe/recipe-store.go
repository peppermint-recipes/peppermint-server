package recipe

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	recipeCollectionName = "recipes"
)

var (
	errCouldNotMarhsallJSON   = errors.New("could not marshall recipe to json")
	errCouldNotFindRecipe     = errors.New("could not find recipe")
	errCouldNotCreateRecipe   = errors.New("could not create recipe")
	errCouldNotSaveRecipe     = errors.New("could not save recipe")
	errCouldNotDeleteRecipe   = errors.New("could not delete recipe")
	errCouldNotCreateObjectID = errors.New("could not create object id")
)

func getAllRecipes() ([]*Recipe, error) {
	var recipes []*Recipe

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database(database.DatabaseName)
	collection := db.Collection(recipeCollectionName)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &recipes)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}
	return recipes, nil
}

func getRecipeByID(id string) (*Recipe, error) {
	var recipe *Recipe

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	mongoObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Could not create object id from string. %v", err)

		return nil, errCouldNotCreateObjectID
	}

	db := client.Database(database.DatabaseName)
	collection := db.Collection(recipeCollectionName)
	result := collection.FindOne(ctx, bson.D{{"id", mongoObjectID}})
	if result == nil {
		return nil, errCouldNotFindRecipe
	}

	err = result.Decode(&recipe)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}

	return recipe, nil
}

func createRecipe(recipe *Recipe) (*Recipe, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	recipe.ID = primitive.NewObjectID()

	_, err := client.Database(database.DatabaseName).Collection(recipeCollectionName).InsertOne(ctx, recipe)
	if err != nil {
		log.Printf("Could not create Recipe: %v", err)
		return recipe, errCouldNotCreateRecipe
	}

	return recipe, nil
}

func updateRecipe(recipe *Recipe) (*Recipe, error) {
	var updatedRecipe *Recipe

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": recipe,
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database(database.DatabaseName).Collection(recipeCollectionName).FindOneAndUpdate(
		ctx, bson.M{"id": recipe.ID}, update, &opt).Decode(&updatedRecipe)
	if err != nil {
		log.Printf("Could not save Recipe: %v", err)

		return nil, errCouldNotSaveRecipe
	}

	return updatedRecipe, nil
}

func deleteRecipe(id string) (*Recipe, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	foundRecipe, err := getRecipeByID(id)
	if err != nil {
		return nil, err
	}
	foundRecipe.Deleted = true

	_, err = updateRecipe(foundRecipe)
	if err != nil {
		return nil, err
	}

	return foundRecipe, nil
}
