package storage

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (cs CarStorage) FindCarById(ctx context.Context, id string) *models.Car {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	filter := bson.M{"_id": _id}
	car := new(models.Car)
	err = cs.Collection.FindOne(ctx, filter).Decode(car)
	if err != nil {
		return nil
	}

	return car
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
	return result, err
}

func (cs CarStorage) DeleteCar(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := cs.Collection.DeleteOne(ctx, filter)
	return result, err
}
func (cs CarStorage) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := cs.Collection.DeleteMany(ctx, filter)
	return result, err
}

func (cs CarStorage) UpdateManyAccountInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}
