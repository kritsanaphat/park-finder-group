package message

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageServices struct {
	MessageStorage *storage.LogStorage
}

type IMessageServices interface {
	RetriveChatList(ctx context.Context, customer_id primitive.ObjectID) ([]models.MessageRoom, error)
	RetriveChatLog(ctx context.Context, reservation_id string, start, limit int) *models.MessageRoom
}

func NewMessageServices(
	db *mongo.Database,
) IMessageServices {
	message_storage := storage.NewMessageStorage(db)

	return MessageServices{
		MessageStorage: message_storage,
	}
}
