package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mineamihai2001/cc/tema_1/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	ctx      context.Context
	Client   *mongo.Client
	database string
}

func NewClient(database string) *Client {
	uri := os.Getenv("MONGO_URI")
	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return &Client{
		ctx:      ctx,
		Client:   client,
		database: database,
	}
}

func Find[T Schema](c *Client, collection string, filter bson.D) *[]T {
	coll := c.Client.Database(c.database).Collection(collection)

	cursor, err := coll.Find(c.ctx, filter)
	core.Check(err)

	var results []T
	if err = cursor.All(c.ctx, &results); err != nil {
		panic(err)
	}
	return &results
}

func Insert[T Schema](c *Client, collection string, doc T) string {
	coll := c.Client.Database(c.database).Collection(collection)

	result, err := coll.InsertOne(c.ctx, doc)
	core.Check(err)

	return result.InsertedID.(primitive.ObjectID).String()
}

func InsertMany[T []interface{}](c *Client, collection string, doc T) []string {
	coll := c.Client.Database(c.database).Collection(collection)

	result, err := coll.InsertMany(c.ctx, doc)
	core.Check(err)

	var ids []string
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID).String())
	}

	return ids
}

func Update(c *Client, collection string, id string, update interface{}) int64 {
	coll := c.Client.Database(c.database).Collection(collection)

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", _id}}
	fmt.Println(filter)

	result, err := coll.UpdateOne(c.ctx, filter, update)
	core.Check(err)

	return result.ModifiedCount
}

func DeleteOne(c *Client, collection string, id string) int64 {
	coll := c.Client.Database(c.database).Collection(collection)

	_id, _ := primitive.ObjectIDFromHex(id)
	result, err := coll.DeleteOne(c.ctx, bson.D{{"_id", _id}})
	core.Check(err)
	return result.DeletedCount
}

func DeleteMany(c *Client, collection string, filter interface{}) int64 {
	coll := c.Client.Database(c.database).Collection(collection)

	result, err := coll.DeleteMany(c.ctx, filter)
	core.Check(err)
	return result.DeletedCount
}
