package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout           = 5
	connectionStringTemplate = "mongodb://%s:%s@%s"
	DatabaseName             = "peppermint-recipes"
)

var (
	mongoUsername string
	mongoPassword string
	mongoEndpoint string
)

func RegisterConnection(username string, password string, endpoint string) {
	mongoUsername = username
	mongoPassword = password
	mongoEndpoint = endpoint
}

// GetConnection Retrieves a client to the MongoDB
func GetConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	username := mongoUsername
	password := mongoPassword
	clusterEndpoint := mongoEndpoint

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	return client, ctx, cancel
}
