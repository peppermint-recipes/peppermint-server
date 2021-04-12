package user

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	userCollectionName = "users"
)

var (
	errCouldNotMarhsallJSON                  = errors.New("could not marshall user to json")
	errCouldNotFindUser                      = errors.New("could not find user")
	errCouldNotCreateUser                    = errors.New("could not create user")
	errCouldNotSaveUser                      = errors.New("could not save user")
	errCouldNotCreateObjectID                = errors.New("could not create object id")
	errCouldNotCreateHashedAndSaltedPassword = errors.New("could not create hashed and salted password")
	errUserNotAuthorized                     = errors.New("user not authorized")
)

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", errCouldNotCreateHashedAndSaltedPassword
	}

	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(errUserNotAuthorized)
		return false
	}

	return true
}

func IsUserAuthorized(userName string, userPassword string) (*User, error) {
	var user *User

	user, err := getUserByName(userName)
	if err != nil {
		return nil, err
	}

	if comparePasswords(user.Password, []byte(userPassword)) {
		return user, nil
	}

	return nil, errUserNotAuthorized
}

func getUserByID(id string) (*User, error) {
	var user *User

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	mongoObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Could not create object id from string. %v", err)

		return nil, errCouldNotCreateObjectID
	}

	db := client.Database(database.DatabaseName)
	collection := db.Collection(userCollectionName)
	result := collection.FindOne(ctx, bson.D{{"id", mongoObjectID}})
	if result == nil {
		return nil, errCouldNotFindUser
	}

	err = result.Decode(&user)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}

	return user, nil
}

func getUserByName(name string) (*User, error) {
	var user *User

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(database.DatabaseName)
	collection := db.Collection(userCollectionName)
	result := collection.FindOne(ctx, bson.D{{"name", name}})
	if result == nil {
		return nil, errCouldNotFindUser
	}

	err := result.Decode(&user)
	if err != nil {
		log.Printf("Failed marshalling %v", err)

		return nil, errCouldNotMarhsallJSON
	}

	return user, nil
}

func createUser(user *User) (*User, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	hashedAndSaltedPassword, err := hashAndSalt([]byte(user.Password))
	if err != nil {
		return nil, errCouldNotCreateHashedAndSaltedPassword
	}
	user.Password = hashedAndSaltedPassword

	user.ID = primitive.NewObjectID()

	_, err = client.Database(database.DatabaseName).
		Collection(userCollectionName).
		InsertOne(ctx, user)
	if err != nil {
		log.Printf("Could not create User: %v", err)
		return user, errCouldNotCreateUser
	}

	return user, nil
}

func updateUser(user *User) (*User, error) {
	var updatedUser *User

	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": user,
	}

	// TODO: handle password change

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database(database.DatabaseName).
		Collection(userCollectionName).
		FindOneAndUpdate(
			ctx, bson.M{"id": user.ID}, update, &opt).
		Decode(&updatedUser)
	if err != nil {
		log.Printf("Could not save User: %v", err)

		return nil, errCouldNotSaveUser
	}

	return updatedUser, nil
}

func deleteUser(id string) (*User, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	foundUser, err := getUserByID(id)
	if err != nil {
		return nil, err
	}
	foundUser.Deleted = true

	_, err = updateUser(foundUser)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}
