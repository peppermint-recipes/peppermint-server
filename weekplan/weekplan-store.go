package weekplan

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	weekplanCollectionName = "weekplans"
)

var (
	errCouldNotMarhsallJSON   = errors.New("could not marshall weekplan to json")
	errCouldNotFindWeekplan   = errors.New("could not find weekplan")
	errCouldNotCreateWeekplan = errors.New("could not create weekplan")
	errCouldNotSaveWeekplan   = errors.New("could not save weekplan")
	errCouldNotDeleteWeekplan = errors.New("could not delete weekplan")
	errCouldNotCreateObjectID = errors.New("could not create object id")
)

func getAllWeekplans() ([]*weekPlan, error) {
	var weekplans []*weekPlan

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database(database.DatabaseName)
	collection := db.Collection(weekplanCollectionName)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &weekplans)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}
	return weekplans, nil
}

func getWeekplanByID(id string) (*weekPlan, error) {
	var weekplan *weekPlan

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	mongoObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Could not create object id from string. %v", err)

		return nil, errCouldNotCreateObjectID
	}

	db := client.Database(database.DatabaseName)
	collection := db.Collection(weekplanCollectionName)
	result := collection.FindOne(ctx, bson.D{{"id", mongoObjectID}})
	if result == nil {
		return nil, errCouldNotFindWeekplan
	}
	err = result.Decode(&weekplan)

	if err != nil {
		log.Printf("Failed marshalling. %v", err)

		return nil, errCouldNotMarhsallJSON
	}

	return weekplan, nil
}

func createWeekplan(weekplan *weekPlan) (primitive.ObjectID, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	weekplan.ID = primitive.NewObjectID()

	insertWeekplanResult, err := client.Database(database.DatabaseName).Collection(weekplanCollectionName).InsertOne(ctx, weekplan)
	if err != nil {
		log.Printf("Could not create Weekplan: %v", err)
		return primitive.NilObjectID, errCouldNotCreateWeekplan
	}

	databaseObjectID := insertWeekplanResult.InsertedID.(primitive.ObjectID)

	return databaseObjectID, nil
}

func updateWeekplan(weekplan *weekPlan) (*weekPlan, error) {
	var updatedWeekplan *weekPlan

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": weekplan,
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database(database.DatabaseName).Collection(weekplanCollectionName).FindOneAndUpdate(
		ctx, bson.M{"id": weekplan.ID}, update, &opt).Decode(&updatedWeekplan)
	if err != nil {
		log.Printf("Could not save Weekplan: %v", err)

		return nil, errCouldNotSaveWeekplan
	}

	return updatedWeekplan, nil
}

func deleteWeekplan(id string) error {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	mongoObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Could not create object id from string. %v", err)

		return errCouldNotCreateObjectID
	}

	_, err = client.Database(database.DatabaseName).Collection(weekplanCollectionName).DeleteOne(
		ctx, bson.D{{"id", mongoObjectID}},
	)
	if err != nil {
		log.Printf("Could not delete Weekplan: %v", err)

		return errCouldNotDeleteWeekplan
	}

	return nil
}
