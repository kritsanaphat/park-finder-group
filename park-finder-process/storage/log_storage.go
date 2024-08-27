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

type LogStorage struct {
	Collection *mongo.Collection
}

func NewTransactionStorage(db *mongo.Database) *LogStorage {
	return &LogStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_TRANSACTION_NAME")),
	}
}

func NewReserveStorage(db *mongo.Database) *LogStorage {
	return &LogStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_RESERVE_NAME")),
	}
}

func NewMessageStorage(db *mongo.Database) *LogStorage {
	return &LogStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_MESSAGE_NAME")),
	}
}

func NewNotificationStorage(db *mongo.Database) *LogStorage {
	return &LogStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_NOTIFICATION_NAME")),
	}
}

func (ls LogStorage) InsertLog(ctx context.Context, data interface{}) error {
	result, err := ls.Collection.InsertOne(ctx, data)
	fmt.Println("Insert count:", result.InsertedID)
	return err
}
func (ls LogStorage) InsertLogNotification(ctx context.Context, notification models.Notification) error {
	result, err := ls.Collection.InsertOne(ctx, notification)
	fmt.Println("Insert count:", result.InsertedID)
	return err
}
func (ls LogStorage) InsertLogReservation(ctx context.Context, data models.Reservation) (*mongo.InsertOneResult, error) {
	result, err := ls.Collection.InsertOne(ctx, data)
	fmt.Println("Insert count:", result.InsertedID)
	return result, err
}

func (ls LogStorage) UpdateLogByInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := ls.Collection.UpdateOne(ctx, filter, update)
	fmt.Println("Modified count:", result.ModifiedCount)
	return result, err
}

func (ls LogStorage) FindLogInterface(ctx context.Context, filter interface{}, user interface{}) error {
	err := ls.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}

	return err
}

func (ls LogStorage) FindMessgeExist(ctx context.Context, reservation_id primitive.ObjectID) *models.MessageRoom {

	filter := bson.M{
		"reservation_id": reservation_id,
	}
	msg := new(models.MessageRoom)
	err := ls.Collection.FindOne(ctx, filter).Decode(msg)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return msg

}

func (ls LogStorage) InsertMessageRoom(ctx context.Context, msg models.MessageRoom) error {
	_, err := ls.Collection.InsertOne(ctx, msg)
	return err

}

func (ls LogStorage) PushMessageLog(ctx context.Context, reservation_id primitive.ObjectID, msg models.MessageLog) error {
	filter := bson.M{"reservation_id": reservation_id}
	update := bson.M{
		"$push": bson.M{
			"message_log": msg,
		},
	}
	result, err := ls.Collection.UpdateOne(ctx, filter, update)
	fmt.Println("Modified count:", result.ModifiedCount)
	return err
}
