package storage

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
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

func (pas ParkingAreaStorage) InsertParkingArea(ctx context.Context, data interface{}) (interface{}, error) {
	result, err := pas.Collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

func (pas ParkingAreaStorage) UpdateLogByInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) UpdateOpenAreaRegisterDocumentClose(ctx context.Context, provider_id primitive.ObjectID, area models.RegisterParkingAreaDocumentStepRequest, parking_id string) (*mongo.UpdateResult, error) {
	_id, err := primitive.ObjectIDFromHex(parking_id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id":          _id,
		"provider_id":  provider_id,
		"status_apply": "wait for document",
	}

	update := bson.M{
		"$set": models.RegisterParkingAreaDocumentStepRequest{
			ParkingPictureUrl:     area.ParkingPictureUrl,
			TitleDeedUrl:          area.TitleDeedUrl,
			LandCertificateUrl:    area.LandCertificateUrl,
			IDCardUrl:             area.IDCardUrl,
			ToatalParkingCount:    area.ToatalParkingCount,
			OverviewPictureUrl:    area.OverviewPictureUrl,
			MeasurementPictureUrl: area.MeasurementPictureUrl,
			Price:                 area.Price,
			StatusApply:           "apply completed",
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	if result.ModifiedCount == 0 {
		return result, nil
	}
	return result, err
}

func (pas ParkingAreaStorage) DeleteManyAreaInterface(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := pas.Collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func (pas ParkingAreaStorage) UpdateOpenAreaDailyStatus(ctx context.Context, daily models.UpdateOpenAreaDailyStatusRequest) (*mongo.UpdateResult, error) {

	_id, err := primitive.ObjectIDFromHex(daily.ParkingAreaID)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$set": models.ParkingAreaOpenDetailUpdate{
			OpenDetail: models.Daily{
				Monday: models.OpenTimeDetail{
					OpenTime:  daily.Monday.OpenTime,
					CloseTime: daily.Monday.CloseTime,
				},
				Tuesday: models.OpenTimeDetail{
					OpenTime:  daily.Tuesday.OpenTime,
					CloseTime: daily.Tuesday.CloseTime,
				},
				Wednesday: models.OpenTimeDetail{
					OpenTime:  daily.Wednesday.OpenTime,
					CloseTime: daily.Wednesday.CloseTime,
				},
				Thursday: models.OpenTimeDetail{
					OpenTime:  daily.Thursday.OpenTime,
					CloseTime: daily.Thursday.CloseTime,
				},
				Friday: models.OpenTimeDetail{
					OpenTime:  daily.Friday.OpenTime,
					CloseTime: daily.Friday.CloseTime,
				},
				Saturday: models.OpenTimeDetail{
					OpenTime:  daily.Saturday.OpenTime,
					CloseTime: daily.Saturday.CloseTime,
				},
				Sunday: models.OpenTimeDetail{
					OpenTime:  daily.Sunday.OpenTime,
					CloseTime: daily.Sunday.CloseTime,
				},
			},
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) UpdateOpenAreaQuickStatusClose(ctx context.Context, id string, diff int) (*mongo.UpdateResult, error) {
	time_current := time.Now()
	newTime := time_current.Add(time.Duration(int(math.Abs(float64(diff)))) * time.Minute)

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"_id": _id,
	}

	update := bson.M{
		"$set": bson.M{
			"open_status":      false,
			"time_stamp_close": newTime,
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) UpdateOpenAreaQuickStatusOpen(ctx context.Context, id string, range_time int) (*mongo.UpdateResult, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"_id": _id,
	}

	update := bson.M{
		"$set": bson.M{
			"open_status":      true,
			"time_stamp_close": nil,
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) PushOpenAreaDateClose(ctx context.Context, id string, date_close []string) (*mongo.UpdateResult, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$push": bson.M{
			"date_close": bson.M{"$each": date_close},
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) PushReview(ctx context.Context, customer_id, fn, ln, comment, parking_id, order_id string, review_score int) (*mongo.UpdateResult, error) {

	_id, err := primitive.ObjectIDFromHex(parking_id)
	if err != nil {
		fmt.Println("Error:", err)
	}
	c_id, err := primitive.ObjectIDFromHex(customer_id)
	if err != nil {
		fmt.Println("Error:", err)
	}
	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$push": bson.M{
			"review": bson.M{
				"review_id":    primitive.NewObjectID(),
				"customer_id":  c_id,
				"first_name":   fn,
				"last_name":    ln,
				"time_stamp":   time.Now(),
				"comment":      comment,
				"order_id":     order_id,
				"review_score": review_score,
			},
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) PullOpenAreaDateClose(ctx context.Context, id string, date_close []string) (*mongo.UpdateResult, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"_id": _id,
	}
	update := bson.M{
		"$pull": bson.M{
			"date_close": bson.M{"$in": date_close},
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
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
		return err
	}

	return nil
}

func (pas ParkingAreaStorage) FindParkingByParkingIDList(ctx context.Context, ids []string) *[]models.ParkingArea {
	var objectIDs []primitive.ObjectID

	for _, id := range ids {
		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		objectIDs = append(objectIDs, _id)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := pas.Collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer cursor.Close(ctx)
	var areas []models.ParkingArea
	for cursor.Next(ctx) {
		var area models.ParkingArea
		err := cursor.Decode(&area)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		areas = append(areas, area)
	}

	if err := cursor.Err(); err != nil {
		return nil
	}
	return &areas
}

func (pas ParkingAreaStorage) FindParkingByProviderID(ctx context.Context, id string) ([]models.ParkingArea, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{"provider_id": _id}
	cursor, err := pas.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var areas []models.ParkingArea
	for cursor.Next(ctx) {
		var area models.ParkingArea
		if err := cursor.Decode(&area); err != nil {
			return nil, err
		}
		areas = append(areas, area)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return areas, nil
}

func (pas ParkingAreaStorage) FindParkingAreaByProviderID(ctx context.Context, provider_id, parking_area_id string) error {
	var area models.ParkingArea
	pv_id, err := primitive.ObjectIDFromHex(provider_id)
	if err != nil {
		fmt.Println("Error:", err)
	}
	pa_id, err := primitive.ObjectIDFromHex(parking_area_id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{"_id": pa_id, "provider_id": pv_id}
	err = pas.Collection.FindOne(ctx, filter).Decode(area)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}
	return err
}

func (pas ParkingAreaStorage) FindListParkingIDByProviderID(ctx context.Context, provider_id string) []string {
	pv_id, err := primitive.ObjectIDFromHex(provider_id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{"provider_id": pv_id}
	cursor, err := pas.Collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer cursor.Close(ctx)
	var list_parking_id []string
	for cursor.Next(ctx) {
		var area models.ParkingArea
		err := cursor.Decode(&area)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		list_parking_id = append(list_parking_id, area.ID.Hex())
	}

	if err := cursor.Err(); err != nil {
		return nil
	}
	return list_parking_id
}

func (pas ParkingAreaStorage) FindParkingAreaByStatus(ctx context.Context, status string) *[]models.ParkingArea {

	filter := bson.M{"status_apply": status}
	cursor, err := pas.Collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer cursor.Close(ctx)
	var areas []models.ParkingArea
	for cursor.Next(ctx) {
		var area models.ParkingArea
		err := cursor.Decode(&area)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		areas = append(areas, area)
	}

	if err := cursor.Err(); err != nil {
		return nil
	}
	return &areas
}

func (pas ParkingAreaStorage) UpdateParkingAreaByStatus(ctx context.Context, status, parking_id, description string) (*mongo.UpdateResult, error) {
	_id, err := primitive.ObjectIDFromHex(parking_id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"status_apply":             status,
			"status_apply_description": description,
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (pas ParkingAreaStorage) UpdateParkingAreaByPrice(ctx context.Context, parking_id primitive.ObjectID, price int16) (*mongo.UpdateResult, error) {

	filter := bson.M{"_id": parking_id}
	update := bson.M{
		"$set": bson.M{
			"price": price,
		},
	}

	result, err := pas.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

// Search Engine
func (pas ParkingAreaStorage) FindParkingAvailableOneDay(ctx context.Context, province1, province2, Date, daliy_open, daliy_close string, MinPrice, MaxPrice, Review int16, HourStart, HourEnd, MinStart, difference int) (*mongo.Cursor, error) {

	HourStart_string := strconv.Itoa(HourStart - 7)
	if len(HourStart_string) == 1 {
		HourStart_string = "0" + HourStart_string
	}
	MinStart_string := strconv.Itoa(MinStart)
	if len(MinStart_string) == 1 {
		MinStart_string = "0" + MinStart_string
	}

	timeString := HourStart_string + ":" + MinStart_string

	date, _ := time.Parse("2006-01-02", Date)
	hour, _ := time.Parse("15:04", timeString)

	start_ts := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), hour.Minute(), 0, 0, time.UTC)
	isoTime := start_ts.Format(time.RFC3339)
	parsedTime, _ := time.Parse(time.RFC3339, isoTime)

	dateArray := []string{Date}

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
			{
				"status_apply": "accepted",
			},
			{
				"$or": []bson.M{
					{
						"$expr": bson.M{
							"$gte": []interface{}{
								bson.M{"$avg": "$review.review_score"},
								Review,
							},
						},
					},
					{"review": bson.M{"$size": 0}},
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

func (pas ParkingAreaStorage) FindParkingAvailableMoreOneDay(ctx context.Context, province1, province2 string, date_names, date_array []string, MinPrice, MaxPrice, Review int16, HourStart, HourEnd, MinStart, difference int) (*mongo.Cursor, error) {
	HourStart_string := strconv.Itoa(HourStart - 7)
	if len(HourStart_string) == 1 {
		HourStart_string = "0" + HourStart_string
	}
	MinStart_string := strconv.Itoa(MinStart)
	if len(MinStart_string) == 1 {
		MinStart_string = "0" + MinStart_string
	}

	timeString := HourStart_string + ":" + MinStart_string

	date, _ := time.Parse("2006-01-02", date_array[0])
	hour, _ := time.Parse("15:04", timeString)

	start_ts := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), hour.Minute(), 0, 0, time.UTC)
	isoTime := start_ts.Format(time.RFC3339)
	parsedTime, _ := time.Parse(time.RFC3339, isoTime)

	filter_datename := _generateMutiDateDailyCloseNameFilter(date_names, HourStart, HourEnd, difference)
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
				"$and": filter_datename,
			},
			{
				"date_close": bson.M{
					"$nin": date_array,
				},
			},
			{
				"$or": []bson.M{
					{"time_stamp_close": bson.M{"$lte": parsedTime}},
					{"time_stamp_close": nil},
				},
			},
			{
				"status_apply": "accepted",
			},
			{
				"$or": []bson.M{
					{
						"$expr": bson.M{
							"$gte": []interface{}{
								bson.M{"$avg": "$review.review_score"},
								Review,
							},
						},
					},
					{"review": bson.M{"$size": 0}},
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

func (pas ParkingAreaStorage) FindParkingAvailableWithTagOneDay(ctx context.Context, province1, province2, Date, daliy_open, daliy_close string, tags string, MinPrice, MaxPrice, Review int16, HourStart, HourEnd, MinStart, difference int) (*mongo.Cursor, error) {
	HourStart_string := strconv.Itoa(HourStart - 7)
	if len(HourStart_string) == 1 {
		HourStart_string = "0" + HourStart_string
	}
	MinStart_string := strconv.Itoa(MinStart)
	if len(MinStart_string) == 1 {
		MinStart_string = "0" + MinStart_string
	}

	timeString := HourStart_string + ":" + MinStart_string

	date, _ := time.Parse("2006-01-02", Date)
	hour, _ := time.Parse("15:04", timeString)

	start_ts := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), hour.Minute(), 0, 0, time.UTC)
	isoTime := start_ts.Format(time.RFC3339)
	parsedTime, _ := time.Parse(time.RFC3339, isoTime)

	dateArray := []string{Date}

	filter := bson.M{
		"$and": []bson.M{
			{
				"price": bson.M{
					"$gte": MinPrice,
					"$lte": MaxPrice,
				},
			},
			{
				"tag": bson.M{
					"$regex": tags,
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
			{
				"status_apply": "accepted",
			},
			{
				"$or": []bson.M{
					{
						"$expr": bson.M{
							"$gte": []interface{}{
								bson.M{"$avg": "$review.review_score"},
								Review,
							},
						},
					},
					{"review": bson.M{"$size": 0}},
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

func (pas ParkingAreaStorage) FindParkingAvailableWithTagMoreOneDay(ctx context.Context, tags string, date_names, date_array []string, MinPrice, MaxPrice, Review int16, HourStart, HourEnd, MinStart, difference int) (*mongo.Cursor, error) {
	HourStart_string := strconv.Itoa(HourStart - 7)
	if len(HourStart_string) == 1 {
		HourStart_string = "0" + HourStart_string
	}
	MinStart_string := strconv.Itoa(MinStart)
	if len(MinStart_string) == 1 {
		MinStart_string = "0" + MinStart_string
	}

	timeString := HourStart_string + ":" + MinStart_string

	date, _ := time.Parse("2006-01-02", date_array[0])
	hour, _ := time.Parse("15:04", timeString)

	start_ts := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), hour.Minute(), 0, 0, time.UTC)
	isoTime := start_ts.Format(time.RFC3339)
	parsedTime, _ := time.Parse(time.RFC3339, isoTime)

	filter_datename := _generateMutiDateDailyCloseNameFilter(date_names, HourStart, HourEnd, difference)
	filter := bson.M{
		"$and": []bson.M{
			{
				"tag": bson.M{
					"$regex": tags,
				},
			},
			{
				"price": bson.M{
					"$gte": MinPrice,
					"$lte": MaxPrice,
				},
			},
			{
				"$and": filter_datename,
			},
			{
				"date_close": bson.M{
					"$nin": date_array,
				},
			},
			{
				"$or": []bson.M{
					{"time_stamp_close": bson.M{"$lte": parsedTime}},
					{"time_stamp_close": nil},
				},
			},
			{
				"status_apply": "accepted",
			},
			{
				"$or": []bson.M{
					{
						"$expr": bson.M{
							"$gte": []interface{}{
								bson.M{"$avg": "$review.review_score"},
								Review,
							},
						},
					},
					{"review": bson.M{"$size": 0}},
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

func _generateMutiDateDailyCloseNameFilter(date_names []string, start, end, diff int) []bson.M {
	var conditions []bson.M
	for i := 0; i < len(date_names); i += 2 {
		if i == 0 {
			condition := bson.M{
				"$and": []bson.M{
					{
						date_names[i]: bson.M{
							"$lte": start,
						},
						date_names[i+1]: bson.M{
							"$gte": 23,
						},
					},
				},
			}
			conditions = append(conditions, condition)

		} else if i == len(date_names)-2 {
			condition := bson.M{
				"$and": []bson.M{
					{
						date_names[len(date_names)-2]: bson.M{
							"$lte": 0,
						},
						date_names[len(date_names)-1]: bson.M{
							"$gte": end + diff,
						},
					},
				},
			}
			conditions = append(conditions, condition)
		} else {
			condition := bson.M{
				"$and": []bson.M{
					{
						date_names[i]: bson.M{
							"$lte": 0,
						},
						date_names[i+1]: bson.M{
							"$gte": 23,
						},
					},
				},
			}
			conditions = append(conditions, condition)
		}
	}
	return conditions
}
