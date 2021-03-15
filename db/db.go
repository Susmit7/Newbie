package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDBCollection ...
func GetDBCollection(collectione string) (*mongo.Collection, *mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://ronnie:rSoK0JNHEkY3y2UR@cluster1-shard-00-00.but8y.mongodb.net:27017,cluster1-shard-00-01.but8y.mongodb.net:27017,cluster1-shard-00-02.but8y.mongodb.net:27017/RHT?ssl=true&replicaSet=atlas-130qu9-shard-0&authSource=admin&retryWrites=true&w=majority"))
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	collection := client.Database("RHT").Collection(collectione)
	return collection, client, nil
}
