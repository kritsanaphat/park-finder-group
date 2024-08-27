package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.comparking-finderpark-finder-process/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParkingAreaStorage struct {
	Collection *mongo.Collection
}

func NewParkingAreaStorage(db *mongo.Database) *ParkingAreaStorage {
	return &ParkingAreaStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_PARKING_AREA_NAME")),
	}
}

func (pas ParkingAreaStorage) InsertParkingArea(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := pas.Collection.InsertOne(ctx, data)
	return result, err
}

func (pas ParkingAreaStorage) UpdateLogByInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	fmt.Println("Modified count:", result.ModifiedCount)
	return result, err
}

func (pas ParkingAreaStorage) FindParkingByParkingID(ctx context.Context, id string, area *models.ParkingArea) error {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{"_id": _id}
	err = pas.Collection.FindOne(ctx, filter).Decode(area)
	if err == mongo.ErrNoDocuments {
		fmt.Println(err)
		return err
	}

	return nil
}

func (pas ParkingAreaStorage) FindParkingAvailable(ctx context.Context, province1, province2, Date, daliy_open, daliy_close string, MinPrice, MaxPrice int16, HourStart, HourEnd, MinStart, difference int) (*mongo.Cursor, error) {

	HourStart_string := strconv.Itoa(HourStart - 7)
	if len(HourStart_string) == 1 {
		HourStart_string = "0" + HourStart_string
	}
	MinStart_string := strconv.Itoa(MinStart)
	if len(MinStart_string) == 1 {
		MinStart_string = "0" + MinStart_string
	}

	timeString := HourStart_string + ":" + MinStart_string

	// Parse date and time strings
	date, _ := time.Parse("2006-01-02", Date)
	hour, _ := time.Parse("15:04", timeString) // Use "15:04" for parsing the time component

	start_ts := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), hour.Minute(), 0, 0, time.UTC)
	isoTime := start_ts.Format(time.RFC3339)
	parsedTime, _ := time.Parse(time.RFC3339, isoTime)

	dateArray := []string{Date}
	fmt.Println(province1, province2)

	filter := bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"address.province": province1},
					{"address.province": province2},
				},
			},
			{
				"price": bson.M{
					"$gte": MinPrice,
					"$lte": MaxPrice,
				},
			},
			{
				"$and": []bson.M{
					{
						daliy_open: bson.M{
							"$lte": HourStart,
						},
						daliy_close: bson.M{
							"$gte": HourEnd + difference,
						},
					},
				},
			},
			{
				"date_close": bson.M{
					"$nin": dateArray,
				},
			},
			{
				"$or": []bson.M{
					{"time_stamp_close": bson.M{"$lte": parsedTime}},
					{"time_stamp_close": nil},
				},
			},
		},
	}

	cursor, err := pas.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}
