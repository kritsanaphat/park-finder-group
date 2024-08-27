package models

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminAccount struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName         string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Phone             string             `json:"phone" bson:"phone"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"password" bson:"password"`
	ProfilePictureURL string             `json:"profile_picture_url" bson:"profile_picture_url"`
	TimeStamp         time.Time          `json:"time_stamp" bson:"time_stamp"`
}

type Reward struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string             `bson:"name"`
	Point           int                `bson:"point"`
	Title           string             `bson:"title"`
	Description     string             `bson:"description"`
	ExpiredDate     time.Time          `bson:"expired_date"`
	TimeStamp       time.Time          `bson:"time_stamp"`
	PreviewImageURL string             `bson:"preview_url"`
	Webhook         string             `bson:"webhook"`
	Condition       []string           `bson:"condition"`
	QuotaCount      int                `bson:"quota_count"`
	CreateBy        string             `bson:"create_by"`
}

func (c AdminAccount) IDToString() string {
	return c.ID.Hex()
}

func (c AdminAccount) ToMap() *echo.Map {
	return &echo.Map{
		"id":         c.IDToString(),
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"email":      c.Email,
		"time_stamp": c.TimeStamp,
	}
}
