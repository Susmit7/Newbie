package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDBCollection ...
func GetDBCollection(collectione string) (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://susmit:6rusThrCaEhRyEx9@cluster0.whsab.mongodb.net/RHT?retryWrites=true&w=majority"))
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
