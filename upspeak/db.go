package upspeak

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoStore stores references to a mongodb client and mongodb collections that
// the rest of upspeak needs
type MongoStore struct {
	Client  *mongo.Client
	Nodes   *mongo.Collection
	Rooms   *mongo.Collection
	Threads *mongo.Collection
}

// Connect connects to a MongoDB instance and returns useful collections as
// `MongoStore`
func Connect(uri string) (MongoStore, error) {
	var db MongoStore

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return db, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return db, err
	}

	db.Client = client
	db.Nodes = client.Database("upspeak_dev").Collection("nodes")
	db.Rooms = client.Database("upspeak_dev").Collection("rooms")
	db.Threads = client.Database("upspeak_dev").Collection("threads")
	return db, nil
}
