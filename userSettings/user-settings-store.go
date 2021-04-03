package userSettings

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userSettingsCollectionName = "user-settings"
)

var (
	errCouldNotMarhsallJSON       = errors.New("could not marshall user to json")
	errCouldNotFindUserSettings   = errors.New("could not find user settings")
	errCouldNotCreateUserSettings = errors.New("could not create user settings")
	errCouldNotSaveUserSettings   = errors.New("could not save user settings")
	errCouldNotDeleteUserSettings = errors.New("could not delete user settings")
)

func getUserSettingsOfUser(userID string) (*UserSettings, error) {
	var userSettings *UserSettings

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(database.DatabaseName)
	collection := db.Collection(userSettingsCollectionName)
	result := collection.FindOne(ctx, bson.M{"userid": userID})
	if result == nil {
		return nil, errCouldNotCreateUserSettings
	}

	err := result.Decode(&userSettings)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}

	return userSettings, nil
}

func createUserSettings(userSettings *UserSettings) (*UserSettings, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	userSettings.ID = primitive.NewObjectID()

	_, err := client.Database(database.DatabaseName).
		Collection(userSettingsCollectionName).
		InsertOne(ctx, userSettings)
	if err != nil {
		log.Printf("Could not create User: %v", err)
		return userSettings, errCouldNotCreateUserSettings
	}

	return userSettings, nil
}

func updateUserSettings(userSettings *UserSettings) (*UserSettings, error) {
	var updatedUserSettings *UserSettings

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": userSettings,
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database(database.DatabaseName).
		Collection(userSettingsCollectionName).
		FindOneAndUpdate(
			ctx, bson.M{"id": userSettings.ID}, update, &opt).
		Decode(&updatedUserSettings)
	if err != nil {
		log.Printf("Could not save User: %v", err)

		return nil, errCouldNotSaveUserSettings
	}

	return updatedUserSettings, nil
}

func deleteUserSettingsForUser(userID string) (*UserSettings, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	foundUserSettings, err := getUserSettingsOfUser(userID)
	if err != nil {
		return nil, err
	}
	foundUserSettings.Deleted = true

	_, err = updateUserSettings(foundUserSettings)
	if err != nil {
		return nil, err
	}

	return foundUserSettings, nil
}
