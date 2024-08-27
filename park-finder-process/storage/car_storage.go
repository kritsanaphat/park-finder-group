package storage

import (
	"context"
	"fmt"
	"os"

	"gitlab.comparking-finderpark-finder-process/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarStorage struct {
	Collection *mongo.Collection
}

func NewCarStorage(db *mongo.Database) *CarStorage {
	return &CarStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_CAR_NAME")),
	}
}

func (cs CarStorage) InsertCar(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := cs.Collection.InsertOne(ctx, data)
	fmt.Println(result)
	return result, err
}

func (cs CarStorage) FindCarByEmail(ctx context.Context, email string) ([]models.Car, error) {
	filter := bson.M{"customer_email": email}

	cursor, err := cs.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []models.Car
	for cursor.Next(ctx) {
		var car models.Car
		if err := cursor.Decode(&car); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}

func (cs CarStorage) FindCarByInterface(ctx context.Context, filter interface{}, car interface{}) error {

	err := cs.Collection.FindOne(ctx, filter).Decode(car)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}
	return err
}

func (cs CarStorage) UpdateCar(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	fmt.Println("Modified count:", result.ModifiedCount)
	return result, err
}

func (cs CarStorage) DeleteCar(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := cs.Collection.DeleteOne(ctx, filter)
	fmt.Println("Delete count:", result.DeletedCount)
	return result, err
}

func (cs CarStorage) UpdateManyAccountInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return result, err
}
