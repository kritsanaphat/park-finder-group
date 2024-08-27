package connector

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client  *mongo.Client
	Context context.Context
}

func NewMongoDBClient(mongodbURI string, maxPoolSize uint64) *MongoDB {
	option := options.Client().ApplyURI(mongodbURI).SetMaxPoolSize(maxPoolSize)
	clinet, err := mongo.Connect(context.Background(), option)
	if err != nil {
		log.Fatal(err)
	}
	return &MongoDB{
		Client:  clinet,
		Context: context.Background(),
	}
}

func (conn *MongoDB) Disconnect() error {
	return conn.Client.Disconnect(conn.Context)
}

func (conn *MongoDB) Ping() error {
	return conn.Client.Ping(conn.Context, nil)
}

func (conn *MongoDB) SelectDB(db_name string) *mongo.Database {
	return conn.Client.Database(db_name)
}
