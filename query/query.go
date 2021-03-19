package query

import (
	"Newbie/db"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connection(val string) (*mongo.Collection, *mongo.Client) {

	collection, client, err := db.GetDBCollection(val)
	if err != nil {
		log.Fatal(err)
	}
	return collection, client

}

func Endconn(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func FindoneID(val string, id primitive.ObjectID, key string) *mongo.SingleResult {

	collection, client := Connection(val)
	result := collection.FindOne(context.TODO(), bson.M{key: id})
	defer Endconn(client)

	return result
}

func UpdateOne(val string, filter primitive.M, update primitive.M) {
	collection, client := Connection(val)
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	defer Endconn(client)
}

func InsertOne(val string, doc interface{}) *mongo.InsertOneResult {
	collection, client := Connection(val)
	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}
	defer Endconn(client)
	return result
}

func FindAll(val string, filter primitive.M) *mongo.Cursor {
	collection, client := Connection(val)
	skip := int64(0)
	limit := int64(10)
	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetLimit(limit), options.Find().SetSkip(skip))
	if err != nil {
		log.Fatal(err)
	}
	defer Endconn(client)
	return cursor
}
