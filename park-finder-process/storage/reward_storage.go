package storage

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type RewardStorage struct {
	Collection *mongo.Collection
}

func NewRewardStorage(db *mongo.Database) *RewardStorage {
	return &RewardStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_REWARD_NAME")),
	}
}

func (rs RewardStorage) InsertReward(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := rs.Collection.InsertOne(ctx, data)
	return result, err
}
