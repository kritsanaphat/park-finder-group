package storage

import (
	"context"
	"os"

	"gitlab.comparking-finderpark-finder-process/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenStorage struct {
	Collection *mongo.Collection
}

func NewTokenStorage(db *mongo.Database) *TokenStorage {
	return &TokenStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_TOKEN_NAME")),
	}
}

func (cts TokenStorage) InsertToken(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := cts.Collection.InsertOne(ctx, data)
	return result, err
}

func (cts TokenStorage) FindToken(ctx context.Context, tk string) *models.Token {
	token := new(models.Token)

	filter := bson.M{"token": tk}
	err := cts.Collection.FindOne(ctx, filter).Decode(token)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return token
}

func (cts TokenStorage) UpdateToken(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cts.Collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	return result, err
}

func (cts TokenStorage) UpdateExpireToken(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cts.Collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	return result, err
}
