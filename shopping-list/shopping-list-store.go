package shoppinglist

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	shoppingListsCollectionName = "shoppinglists"
)

var (
	errCouldNotMarhsallJSON       = errors.New("could not marshall shopping list to json")
	errCouldNotFindShoppingList   = errors.New("could not find shopping list")
	errCouldNotCreateShoppingList = errors.New("could not create shopping list")
	errCouldNotSaveShoppingList   = errors.New("could not save shopping list")
	errCouldNotDeleteShoppingList = errors.New("could not delete shopping list")
	errCouldNotCreateObjectID     = errors.New("could not create object id")
)

func getAllShoppingLists() ([]*shoppingList, error) {
	var shoppingLists []*shoppingList

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database(database.DatabaseName)
	collection := db.Collection(shoppingListsCollectionName)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &shoppingLists)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}
	return shoppingLists, nil
}

func getShoppingListByID(id string) (*shoppingList, error) {
	var sl *shoppingList

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	mongoObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Could not create object id from string. %v", err)

		return nil, errCouldNotCreateObjectID
	}

	db := client.Database(database.DatabaseName)
	collection := db.Collection(shoppingListsCollectionName)
	result := collection.FindOne(ctx, bson.D{{"id", mongoObjectID}})
	if result == nil {
		return nil, errCouldNotFindShoppingList
	}
	err = result.Decode(&sl)

	if err != nil {
		log.Printf("Failed marshalling. %v", err)

		return nil, errCouldNotMarhsallJSON
	}

	return sl, nil
}

func createShoppingList(sl *shoppingList) (primitive.ObjectID, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	sl.ID = primitive.NewObjectID()

	insertWeekplanResult, err := client.Database(database.DatabaseName).Collection(shoppingListsCollectionName).InsertOne(ctx, sl)
	if err != nil {
		log.Printf("Could not create shopping list: %v", err)
		return primitive.NilObjectID, errCouldNotCreateShoppingList
	}

	databaseObjectID := insertWeekplanResult.InsertedID.(primitive.ObjectID)

	return databaseObjectID, nil
}

func updateShoppingList(sl *shoppingList) (*shoppingList, error) {
	var updatedShoppingList *shoppingList

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": sl,
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database(database.DatabaseName).Collection(shoppingListsCollectionName).FindOneAndUpdate(
		ctx, bson.M{"id": sl.ID}, update, &opt).Decode(&updatedShoppingList)
	if err != nil {
		log.Printf("Could not save shopping list: %v", err)

		return nil, errCouldNotSaveShoppingList
	}

	return updatedShoppingList, nil
}

func deleteShoppingList(id string) error {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	mongoObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Could not create object id from string. %v", err)

		return errCouldNotCreateObjectID
	}

	_, err = client.Database(database.DatabaseName).Collection(shoppingListsCollectionName).DeleteOne(
		ctx, bson.D{{"id", mongoObjectID}},
	)
	if err != nil {
		log.Printf("Could not delete shopping list: %v", err)

		return errCouldNotDeleteShoppingList
	}

	return nil
}
