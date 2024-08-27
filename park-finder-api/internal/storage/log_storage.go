package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/utility"
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

func NewReceiptStorage(db *mongo.Database) *LogStorage {
	return &LogStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_RECEIPT_NAME")),
	}
}

func NewReportStorage(db *mongo.Database) *LogStorage {
	return &LogStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_REPORT_NAME")),
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

func (ls LogStorage) InsertLog(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := ls.Collection.InsertOne(ctx, data)
	return result, err
}

func (ls LogStorage) InsertLogReservation(ctx context.Context, data models.Reservation) (*mongo.InsertOneResult, error) {
	result, err := ls.Collection.InsertOne(ctx, data)

	return result, err
}

func (ls LogStorage) UpdateLogByInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := ls.Collection.UpdateOne(ctx, filter, update)
	return result, err
}

func (ls LogStorage) FindLogInterface(ctx context.Context, filter interface{}, user interface{}) error {
	err := ls.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}

	return err
}

func (ls LogStorage) DeleteMany(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := ls.Collection.DeleteMany(ctx, filter)
	return result, err
}

// Reserve
func (ls LogStorage) FindReservationByOrderID(ctx context.Context, order_id string) *models.Reservation {
	reserve := new(models.Reservation)
	filter := bson.M{"order_id": order_id}
	err := ls.Collection.FindOne(ctx, filter).Decode(reserve)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return reserve
}

func (ls LogStorage) FindReservationByReserveID(ctx context.Context, reserve_id string) *models.Reservation {
	_id, err := primitive.ObjectIDFromHex(reserve_id)
	if err != nil {
		return nil
	}
	reserve := new(models.Reservation)
	filter := bson.M{"_id": _id}
	err = ls.Collection.FindOne(ctx, filter).Decode(reserve)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return reserve
}

func (ls LogStorage) FindActiveReservationByParkingID(ctx context.Context, parking_id string) ([]models.Reservation, error) {

	_id, err := primitive.ObjectIDFromHex(parking_id)
	if err != nil {
		fmt.Println("Error:", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{
				"parking_id": _id},
			{
				"status": "Process",
			},
		},
	}
	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (ls LogStorage) FindIncommingReservationByParkingID(ctx context.Context, parking_id primitive.ObjectID) ([]models.Reservation, error) {

	filter := bson.M{
		"$and": []bson.M{
			{
				"parking_id": parking_id},
			{
				"$or": []bson.M{
					{"status": "Process"},
				}},
		},
	}
	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (ls LogStorage) FindStillParkingReservationByParkingID(ctx context.Context, parking_id primitive.ObjectID) ([]models.Reservation, error) {

	filter := bson.M{
		"$and": []bson.M{
			{
				"parking_id": parking_id},
			{
				"$or": []bson.M{
					{"status": "Parking"},
				}},
		},
	}
	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (ls LogStorage) FindProcessReservationByCustomerEmail(ctx context.Context, customer_email, parking_id string) ([]models.Reservation, error) {
	filter := bson.M{}
	if parking_id != "" {
		_id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error:", err)
		}
		filter = bson.M{
			"$and": []bson.M{
				{
					"customer_email": customer_email},
				{
					"parking_id": _id},
				{
					"$or": []bson.M{
						{"status": "Process"},
						{"status": "Pending"},
						{"status": "Pending Approval"},
						{"status": "Pending Approval Process"},
						{"status": "Parking"},
					}},
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{
					"customer_email": customer_email},
				{
					"$or": []bson.M{
						{"status": "Process"},
						{"status": "Pending"},
						{"status": "Pending Approval"},
						{"status": "Pending Approval Process"},
						{"status": "Parking"},
					}},
			},
		}
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindProcessReservationByProviderID(ctx context.Context, provider_id primitive.ObjectID, parking_id string) ([]models.Reservation, error) {
	filter := bson.M{}
	if parking_id != "" {
		_id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error:", err)
		}
		filter = bson.M{
			"$and": []bson.M{
				{
					"provider_id": provider_id},
				{
					"parking_id": _id},
				{
					"$or": []bson.M{
						{"status": "Process"},
						{"status": "Pending"},
						{"status": "Pending Approval"},
						{"status": "Pending Approval Process"},
						{"status": "Parking"},
					}},
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{
					"provider_id": provider_id},
				{
					"$or": []bson.M{
						{"status": "Process"},
						{"status": "Pending"},
						{"status": "Pending Approval"},
						{"status": "Pending Approval Process"},
						{"status": "Parking"},
					}},
			},
		}
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			fmt.Println("Error", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindComplePaymentReservationByCustomerEmail(ctx context.Context, customer_email string) ([]models.Reservation, error) {
	filter := bson.M{}

	filter = bson.M{
		"$and": []bson.M{
			{
				"customer_email": customer_email},
			{
				"$or": []bson.M{
					{"status": "Process"},
					{"status": "Pending Approval Process"},
					{"status": "Parking"},
				}},
		},
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindComplePaymentReservationByProviderID(ctx context.Context, provider_id primitive.ObjectID) ([]models.Reservation, error) {
	filter := bson.M{}

	filter = bson.M{
		"$and": []bson.M{
			{
				"provider_id": provider_id},
			{
				"$or": []bson.M{
					{"status": "Process"},
					{"status": "Pending Approval Process"},
					{"status": "Parking"},
				}},
		},
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindCancelReservationByCustomerEmail(ctx context.Context, customer_email, parking_id string) ([]models.Reservation, error) {
	filter := bson.M{}
	if parking_id != "" {
		_id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error:", err)
		}
		filter = bson.M{
			"$and": []bson.M{
				{
					"customer_email": customer_email},
				{
					"parking_id": _id},
				{
					"status": "Cancel"},
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{
					"customer_email": customer_email},
				{
					"status": "Cancel"},
			},
		}
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindCancelReservationByProviderID(ctx context.Context, provider_id primitive.ObjectID, parking_id string) ([]models.Reservation, error) {
	filter := bson.M{}
	if parking_id != "" {
		_id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error:", err)
		}
		filter = bson.M{
			"$and": []bson.M{
				{
					"provider_id": provider_id},
				{
					"parking_id": _id},
				{
					"status": "Cancel"},
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{
					"provider_id": provider_id},
				{
					"status": "Cancel"},
			},
		}
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindSuccessfulReservationByCustomerEmail(ctx context.Context, customer_email, parking_id string) ([]models.Reservation, error) {

	filter := bson.M{}
	if parking_id != "" {
		_id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error:", err)
		}
		filter = bson.M{
			"$and": []bson.M{
				{
					"customer_email": customer_email},
				{
					"parking_id": _id},
				{
					"status": "Successful"},
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{
					"customer_email": customer_email},
				{
					"status": "Successful"},
			},
		}
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) FindSuccessfulReservationByProviderID(ctx context.Context, provider_id primitive.ObjectID, parking_id string) ([]models.Reservation, error) {

	filter := bson.M{}
	if parking_id != "" {
		_id, err := primitive.ObjectIDFromHex(parking_id)
		if err != nil {
			fmt.Println("Error:", err)
		}
		filter = bson.M{
			"$and": []bson.M{
				{
					"provider_id": provider_id},
				{
					"parking_id": _id},
				{
					"status": "Successful"},
			},
		}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{
					"provider_id": provider_id},
				{
					"status": "Successful"},
			},
		}
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var reservations []models.Reservation
	for cursor.Next(ctx) {
		var reservation models.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func (ls LogStorage) UpdateConfirmStatusToProcessByOrderID(ctx context.Context, order_id string) error {
	filter := bson.M{
		"order_id": order_id,
	}
	update := bson.M{
		"$set": models.UpdateStatusReservation{
			Status: "process",
		}}
	_, err := ls.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (ls LogStorage) FindReservationIDByParkingIDAndCustomerEmail(ctx context.Context, provider_id, customer_email string) *models.Reservation {
	_id, err := primitive.ObjectIDFromHex(provider_id)
	if err != nil {
		return nil
	}
	reserve := new(models.Reservation)
	filter := bson.M{
		"parking_id":     _id,
		"customer_email": customer_email,
		"status":         "Process",
	}
	err = ls.Collection.FindOne(ctx, filter).Decode(reserve)
	if err != nil {
		return nil
	}

	return reserve
}

// Message
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

func (ls LogStorage) FindMessageLogWithLimit(ctx context.Context, reservationID primitive.ObjectID, start, limit int) *models.MessageRoom {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"reservation_id": reservationID,
			},
		},
		{
			"$unwind": "$message_log",
		},
		{
			"$sort": bson.M{
				"message_log.message.time_stamp": -1,
			},
		},
		{
			"$group": bson.M{
				"_id":            "$_id",
				"reservation_id": bson.M{"$first": "$reservation_id"},
				"group_list":     bson.M{"$first": "$group_list"},
				"message_log":    bson.M{"$push": "$message_log"},
			},
		},
		{
			"$project": bson.M{
				"_id":            1,
				"reservation_id": 1,
				"group_list":     1,
				"message_log": bson.M{
					"$slice": []interface{}{
						"$message_log",
						start,
						start + limit,
					},
				},
			},
		},
	}

	cursor, err := ls.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil
	}

	defer cursor.Close(ctx)

	var result models.MessageRoom
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil
		}
		return &result
	}

	return nil
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
	_, err := ls.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (ls LogStorage) FindMessageListByAccountID(ctx context.Context, provider_id primitive.ObjectID) ([]models.MessageRoom, error) {
	filter := bson.M{
		"group_list": bson.M{
			"$elemMatch": bson.M{
				"_id": provider_id,
			},
			"$exists": true,
		},
	}
	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var chat_lists []models.MessageRoom
	for cursor.Next(ctx) {
		var chat models.MessageRoom
		if err := cursor.Decode(&chat); err != nil {
			return nil, err
		}
		chat_lists = append(chat_lists, chat)
	}

	return chat_lists, nil
}

// Notification
func (ls LogStorage) FindNotification(ctx context.Context, receiver_id, broadcast_type string) ([]models.Notification, error) {
	_id, err := primitive.ObjectIDFromHex(receiver_id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"$or": []bson.M{
			{"receiver_id": _id},
			{"broadcast_type": broadcast_type},
		},
	}
	cursor, err := ls.Collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	var notifications []models.Notification
	for cursor.Next(ctx) {
		var notification models.Notification
		if err := cursor.Decode(&notification); err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// Transaction
func (ls LogStorage) FindProfitDaily(ctx context.Context, list_parking_id []string) (string, int, int, error) {

	var count = 0
	var primitiveIDs []primitive.ObjectID
	for _, id := range list_parking_id {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", 0, 0, err
		}
		primitiveIDs = append(primitiveIDs, objectID)
	}

	targetDate := time.Now().UTC()
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"parking_id": bson.M{"$in": primitiveIDs},
		"time_stamp": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
		"status": "Successful",
	}
	cursor, err := ls.Collection.Find(context.Background(), filter)
	if err != nil {
		return "", 0, 0, err
	}
	defer cursor.Close(context.Background())

	var results []models.TransactionLinePay
	for cursor.Next(context.Background()) {
		var result models.TransactionLinePay
		if err := cursor.Decode(&result); err != nil {
			return "", 0, 0, err
		}
		results = append(results, result)
	}
	sum := 0
	for _, transaction := range results {
		for _, amount := range transaction.Packages {
			sum += int(amount.Amount)
			count += 1
		}
	}

	return startOfDay.Format("2006-01-02"), sum, count, nil
}

func (ls LogStorage) FindProfitWeekly(ctx context.Context, list_parking_id []string) ([]models.DailyProfitResponse, int, error) {
	targetDate := time.Now().UTC()
	startOfWeek := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day()-6, 0, 0, 0, 0, time.UTC)
	endOfWeek := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 0, time.UTC)
	var count = 0

	var primitiveIDs []primitive.ObjectID
	for _, id := range list_parking_id {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, 0, err
		}
		primitiveIDs = append(primitiveIDs, objectID)
	}

	filter := bson.M{
		"parking_id": bson.M{"$in": primitiveIDs},
		"time_stamp": bson.M{
			"$gte": startOfWeek,
			"$lt":  endOfWeek,
		},
		"status": "Successful",
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	dailySums := make(map[string]int)

	for cursor.Next(ctx) {
		var transaction models.TransactionLinePay
		if err := cursor.Decode(&transaction); err != nil {
			return nil, 0, err
		}

		day := transaction.TimeStamp.UTC().Format("2006-01-02")

		for _, amount := range transaction.Packages {
			dailySums[day] += int(amount.Amount)
			count += 1

		}
	}

	var result []map[string]interface{}
	for day, sum := range dailySums {
		result = append(result, map[string]interface{}{
			"date": day,
			"sum":  sum,
		})
	}

	var response []models.DailyProfitResponse
	for d := startOfWeek; d.Before(endOfWeek); d = d.AddDate(0, 0, 1) {
		day := d.UTC().Format("2006-01-02")
		sum, ok := dailySums[day]
		if !ok {
			sum = 0
		}
		data := models.DailyProfitResponse{
			Date: day,
			Sum:  sum,
		}
		response = append(response, data)
	}
	return response, count, nil
}

func (ls LogStorage) FindProfitMontly(ctx context.Context, list_parking_id []string) ([]models.DailyProfitResponse, int, error) {
	targetDate := time.Now().UTC()
	var count = 0

	startOfMonth := time.Date(targetDate.Year(), targetDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(24 * time.Hour)

	var result []models.DailyProfitResponse
	var primitiveIDs []primitive.ObjectID
	for _, id := range list_parking_id {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, 0, err
		}
		primitiveIDs = append(primitiveIDs, objectID)
	}

	for startOfWeek := startOfMonth; startOfWeek.Before(endOfMonth); startOfWeek = startOfWeek.AddDate(0, 0, 7) {
		endOfWeek := startOfWeek.AddDate(0, 0, 5).Add(24 * time.Hour)

		if endOfWeek.After(endOfMonth) {
			endOfWeek = endOfMonth.Add(-24 * time.Hour)

		}

		filter := bson.M{
			"parking_id": bson.M{"$in": primitiveIDs},
			"time_stamp": bson.M{
				"$gte": startOfWeek,
				"$lt":  endOfWeek,
			},
			"status": "Successful",
		}

		cursor, err := ls.Collection.Find(ctx, filter)
		if err != nil {
			return nil, 0, err
		}
		defer cursor.Close(ctx)

		sum := 0
		for cursor.Next(ctx) {
			var transaction models.TransactionLinePay
			if err := cursor.Decode(&transaction); err != nil {
				return nil, 0, err
			}
			for _, amount := range transaction.Packages {
				sum += int(amount.Amount)
				count += 1
			}
		}

		result = append(result, models.DailyProfitResponse{
			Date: fmt.Sprintf("%s - %s", startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02")),
			Sum:  sum,
		})
	}

	return result, count, nil
}

func (ls LogStorage) FindProfitYearly(ctx context.Context, list_parking_id []string) ([]models.DailyProfitResponse, int, error) {
	currentTime := time.Now().UTC()
	startOfYear := time.Now().UTC().AddDate(0, -12, 0)
	startOfYear = time.Date(startOfYear.Year(), startOfYear.Month(), 1, 0, 0, 0, 0, time.UTC)
	var count = 0

	endOfYear := time.Date(startOfYear.Year()+1, currentTime.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-1 * time.Second)
	var result []models.DailyProfitResponse

	var primitiveIDs []primitive.ObjectID
	for _, id := range list_parking_id {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, 0, err
		}
		primitiveIDs = append(primitiveIDs, objectID)
	}

	for startOfMonth := startOfYear; startOfMonth.Before(endOfYear); startOfMonth = startOfMonth.AddDate(0, 1, 0) {
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

		filter := bson.M{
			"parking_id": bson.M{"$in": primitiveIDs},
			"time_stamp": bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			},
			"status": "Successful",
		}

		cursor, err := ls.Collection.Find(ctx, filter)
		if err != nil {
			return nil, 0, err
		}
		defer cursor.Close(ctx)

		sum := 0
		for cursor.Next(ctx) {
			var transaction models.TransactionLinePay
			if err := cursor.Decode(&transaction); err != nil {
				return nil, 0, err
			}
			for _, amount := range transaction.Packages {
				sum += int(amount.Amount)
				count += 1
			}
		}

		result = append(result, models.DailyProfitResponse{
			Date: fmt.Sprintf("%s", startOfMonth.Format("2006-Jan")),
			Sum:  sum,
		})
	}

	return result, count, nil
}

func (ls LogStorage) FindTransactionFromAdmin(ctx context.Context, list_parking_id []string, month, year string) (int, error) {
	targetDate := time.Now().UTC()
	month_time := utility.ParseMonth(month)
	startOfMonth := time.Date(targetDate.Year(), month_time, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(24 * time.Hour)

	var primitiveIDs []primitive.ObjectID
	for _, id := range list_parking_id {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		primitiveIDs = append(primitiveIDs, objectID)
	}

	filter := bson.M{
		"parking_id": bson.M{"$in": primitiveIDs},
		"time_stamp": bson.M{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
		"status": "Successful",
	}

	cursor, err := ls.Collection.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	sum := 0
	for cursor.Next(ctx) {
		var transaction models.TransactionLinePay
		if err := cursor.Decode(&transaction); err != nil {
			return 0, err
		}
		for _, amount := range transaction.Packages {
			sum += int(amount.Amount)
		}
	}

	return sum, nil

}

// Receipt
func (ls LogStorage) FindExistReceipt(ctx context.Context, month, year, provider_id string, price int) (bool, int) {
	_id, err := primitive.ObjectIDFromHex(provider_id)
	if err != nil {
		return false, 0
	}
	filter := bson.M{
		"provider_id": _id,
		"month":       month,
		"year":        year,
	}
	var sum int
	cursor, err := ls.Collection.Find(ctx, filter)
	if err != nil {
		return false, 0
	}
	for cursor.Next(ctx) {
		var transaction models.Receipt
		if err := cursor.Decode(&transaction); err != nil {
			return false, 0
		}
		sum += transaction.Price
	}
	if sum == price {
		return true, 0
	} else {
		return false, price - sum
	}

}
