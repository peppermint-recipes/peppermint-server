package recipe

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getAllRecipes() ([]*Recipe, error) {
	var recipes []*Recipe

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("recipes")
	collection := db.Collection("recipes")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &recipes)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return recipes, nil
}

func getRecipeByID(id primitive.ObjectID) (*Recipe, error) {
	var recipe *Recipe

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("recipes")
	collection := db.Collection("recipes")
	result := collection.FindOne(ctx, bson.D{})
	if result == nil {
		return nil, errors.New("Could not find a Recipe")
	}
	err := result.Decode(&recipe)

	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	log.Printf("Recipis: %v", recipe)
	return recipe, nil
}

func createRecipe(recipe *Recipe) (primitive.ObjectID, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	recipe.ID = primitive.NewObjectID()

	result, err := client.Database("recipes").Collection("recipes").InsertOne(ctx, recipe)
	if err != nil {
		log.Printf("Could not create Recipe: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
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

	err := client.Database("recipes").Collection("recipes").FindOneAndUpdate(
		ctx, bson.M{"id": recipe.ID}, update, &opt).Decode(&updatedRecipe)
	if err != nil {
		log.Printf("Could not save Recipe: %v", err)
		return nil, err
	}
	return updatedRecipe, nil
}

func deleteRecipe(recipe *Recipe) error {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	_, err := client.Database("recipes").Collection("recipes").DeleteOne(ctx, recipe)
	if err != nil {
		log.Printf("Could not delete Recipe: %v", err)
		return err
	}
	return nil
}
