package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (rs RewardStorage) DeleteManyRewardInterface(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := rs.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rs RewardStorage) FindRewardByExpiredDate(ctx context.Context) *mongo.Cursor {
	filter := bson.M{"expired_date": bson.M{"$gt": time.Now()}}

	cursor, err := rs.Collection.Find(ctx, filter)
	if err != nil {
		return nil
	}

	return cursor
}

func (rs RewardStorage) FindRewardByID(ctx context.Context, id string) *models.Reward {
	reward := new(models.Reward)
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	filter := bson.M{"_id": _id}

	err = rs.Collection.FindOne(ctx, filter).Decode(reward)
	if err != nil {
		return nil
	}

	return reward
}

func (cs RewardStorage) RemoveQuotaCount(ctx context.Context, reward_id primitive.ObjectID, quota_count int) error {

	filter := bson.M{
		"_id": reward_id,
	}
	update := bson.M{
		"$set": bson.M{
			"quota_count": quota_count - 1,
		},
	}

	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return nil
}
