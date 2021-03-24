package user

import (
	"errors"
	"log"

	"github.com/peppermint-recipes/peppermint-server/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	errCouldNotDeleteUser                    = errors.New("could not delete user")
	errCouldNotCreateObjectID                = errors.New("could not create object id")
	errCouldNotCreateHashedAndSaltedPassword = errors.New("could not create hashed and salted password")
	errUserNotAuthorized                     = errors.New("user not authorized")
)

func hashAndSalt(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
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

// func updateRecipe(recipe *Recipe) (*Recipe, error) {
// 	var updatedRecipe *Recipe

// 	client, ctx, cancel := database.GetConnection()
// 	defer cancel()
// 	defer client.Disconnect(ctx)

// 	update := bson.M{
// 		"$set": recipe,
// 	}

// 	upsert := false
// 	after := options.After
// 	opt := options.FindOneAndUpdateOptions{
// 		Upsert:         &upsert,
// 		ReturnDocument: &after,
// 	}

// 	err := client.Database(database.DatabaseName).
// 		Collection(userCollectionName).
// 		FindOneAndUpdate(
// 			ctx, bson.M{"id": recipe.ID}, update, &opt).
// 		Decode(&updatedRecipe)
// 	if err != nil {
// 		log.Printf("Could not save Recipe: %v", err)

// 		return nil, errCouldNotSaveRecipe
// 	}

// 	return updatedRecipe, nil
// }

// func deleteRecipe(id string) (*Recipe, error) {
// 	client, ctx, cancel := database.GetConnection()
// 	defer cancel()
// 	defer client.Disconnect(ctx)

// 	foundRecipe, err := getRecipeByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	foundRecipe.Deleted = true

// 	_, err = updateRecipe(foundRecipe)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return foundRecipe, nil
// }
