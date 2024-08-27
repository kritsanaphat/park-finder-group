package storage

import (
	"context"
	"fmt"
	"os"

	"gitlab.comparking-finderpark-finder-process/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccoutStorage struct {
	Collection *mongo.Collection
}

func NewCustomerAccoutStorage(db *mongo.Database) *AccoutStorage {
	return &AccoutStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_CUSTOMER_ACCOUNT_NAME")),
	}
}
func NewProviderAccoutStorage(db *mongo.Database) *AccoutStorage {
	return &AccoutStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_PROVIDER_ACCOUNT_NAME")),
	}
}
func NewAdminAccoutStorage(db *mongo.Database) *AccoutStorage {
	return &AccoutStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_ADMIN_ACCOUNT_NAME")),
	}
}

func (cs AccoutStorage) InsertAccount(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := cs.Collection.InsertOne(ctx, data)
	return result, err
}

func (cs AccoutStorage) UpdateAccountInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return result, err
}

func (cs AccoutStorage) UpdateManyAccountInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return result, err
}

func (cs AccoutStorage) FindAccountInterface(ctx context.Context, filter interface{}, user interface{}) error {
	err := cs.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}
	return err
}

func (cs AccoutStorage) FindAccountByID(id string) *models.CustomerAccount {
	ctx := context.Background()

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"_id": _id,
	}
	user := new(models.CustomerAccount)
	err = cs.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return user
}
