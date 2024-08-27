package message

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ms MessageServices) RetriveChatList(ctx context.Context, provider_id primitive.ObjectID) ([]models.MessageRoom, error) {
	message_list, err := ms.MessageStorage.FindMessageListByAccountID(ctx, provider_id)
	if err != nil {
		return nil, err
	}

	return message_list, nil
}

func (ms MessageServices) RetriveChatLog(ctx context.Context, reservation_id string, start, limit int) *models.MessageRoom {
	id, err := primitive.ObjectIDFromHex(reservation_id)
	if err != nil {
		return nil
	}
	message_list := ms.MessageStorage.FindMessageLogWithLimit(ctx, id, start, limit)
	if message_list == nil {
		return nil
	}

	return message_list
}
