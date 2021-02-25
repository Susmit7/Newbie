package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDBCollection ...
func GetDBCollection(collectione string) (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://ronnie:rSoK0JNHEkY3y2UR@cluster1-shard-00-00.but8y.mongodb.net:27017,cluster1-shard-00-01.but8y.mongodb.net:27017,cluster1-shard-00-02.but8y.mongodb.net:27017/RHT?ssl=true&replicaSet=atlas-130qu9-shard-0&authSource=admin&retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	//defer client.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	collection := client.Database("RHT").Collection(collectione)
	return collection, nil
}
