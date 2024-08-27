package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reward struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string             `bson:"name"`
	Point           float32            `bson:"point"`
	Description     string             `bson:"description"`
	ExpiredDate     time.Time          `bson:"expired_date"`
	TimeStamp       time.Time          `bson:"time_stamp"`
	PreviewImageURL string             `bson:"preview_url"`
	Webhook         string             `bson:"webhook"`
}
