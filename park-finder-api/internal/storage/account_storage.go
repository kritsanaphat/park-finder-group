package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccoutStorage struct {
	Collection *mongo.Collection
}

func NewCustomerAccoutStorage(db *mongo.Database) *AccoutStorage {
	return &AccoutStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_CUSTOMER_ACCOUNT_NAME")),
	}
}
func NewProviderAccoutStorage(db *mongo.Database) *AccoutStorage {
	return &AccoutStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_PROVIDER_ACCOUNT_NAME")),
	}
}
func NewAdminAccoutStorage(db *mongo.Database) *AccoutStorage {
	return &AccoutStorage{
		Collection: db.Collection(os.Getenv("COLLECTION_ADMIN_ACCOUNT_NAME")),
	}
}

func (cs AccoutStorage) InsertAccount(ctx context.Context, data interface{}) (*mongo.InsertOneResult, error) {
	result, err := cs.Collection.InsertOne(ctx, data)
	return result, err
}

func (cs AccoutStorage) UpdateAccountInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func (cs AccoutStorage) UpdateManyAccountInterface(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := cs.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func (cs AccoutStorage) DeleteManyAccountInterface(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := cs.Collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func (cs AccoutStorage) FindAccountInterface(ctx context.Context, filter interface{}, user interface{}) error {
	err := cs.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}
	return err
}
func (cs AccoutStorage) FindAllListProviderID(ctx context.Context) []models.AccountAndBank {
	cursor, err := cs.Collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	var account_and_banks []models.AccountAndBank

	for cursor.Next(ctx) {
		var account models.ProviderAccount
		err := cursor.Decode(&account)
		if err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		account_and_bank := models.AccountAndBank{
			ID:          account.ID.Hex(),
			BankAccount: account.BankAccount,
		}
		account_and_banks = append(account_and_banks, account_and_bank)
	}
	return account_and_banks
}

func (cs AccoutStorage) FindCustomerAccountIDByEmail(ctx context.Context, email string) string {
	user := new(models.CustomerAccount)
	filter := bson.M{
		"email": email,
	}
	err := cs.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return ""
	}
	return user.ID.Hex()

}

func (cs AccoutStorage) FindCustomerAccountByID(ctx context.Context, id primitive.ObjectID) *models.CustomerAccount {
	user := new(models.CustomerAccount)
	filter := bson.M{
		"_id": id,
	}
	err := cs.Collection.FindOne(ctx, filter).Decode(user)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return user
}

func (cs AccoutStorage) PushFavoriteArea(ctx context.Context, parking_area, email string) error {

	filter := bson.M{
		"email": email,
	}
	fmt.Println(email)
	update := bson.M{
		"$push": bson.M{
			"favorite_area": parking_area,
		},
	}
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) PullFavoriteArea(ctx context.Context, parking_area, email string) error {

	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$pull": bson.M{
			"favorite_area": parking_area,
		},
	}
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) PushRedeemReward(ctx context.Context, customer_id, reward_id primitive.ObjectID, barcode, name string, point int) error {
	currentTime := time.Now()
	expired_date := currentTime.Add(24 * time.Hour)
	reward := &models.CustomerReward{
		Name:        name,
		ID:          reward_id,
		BarcodeURL:  barcode,
		ExpiredDate: expired_date,
		TimeStamp:   currentTime,
		Point:       point,
	}
	filter := bson.M{
		"_id": customer_id,
	}
	update := bson.M{
		"$push": bson.M{
			"reward": *reward,
		},
		"$inc": bson.M{
			"point": -point,
		},
	}

	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) AddCashback(ctx context.Context, email string, cashback int) error {
	filter := bson.M{"email": email}
	update := bson.M{"$inc": bson.M{"cashback": cashback}}
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified cashback count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) ResetFineCustomer(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"fine": &models.CustomerFine{},
		},
	}

	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified struc count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) UpdateCashback(ctx context.Context, email string, cashback int) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"cashback": cashback}}
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified cashback count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) AddPoint(ctx context.Context, email string, point int) error {
	filter := bson.M{"email": email}
	update := bson.M{"$inc": bson.M{"point": point}}
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified cashback count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) UpdateFineCustomer(ctx context.Context, _id primitive.ObjectID, fine int, reserve models.Reservation) error {
	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"fine.order_id":     reserve.OrderID,
			"fine.parking_id":   reserve.ParkingID,
			"fine.provider_id":  reserve.ProviderID,
			"fine.quantity":     1,
			"fine.parking_name": reserve.ParkingName,
		},
	}
	incUpdate := bson.M{
		"$inc": bson.M{
			"fine.price": fine,
		},
	}

	result, err := cs.Collection.UpdateOne(ctx, filter, incUpdate)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified fine count:", result.ModifiedCount)
	result, err = cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified struc count:", result.ModifiedCount)
	return nil
}

func (cs AccoutStorage) UpdateFineProvider(ctx context.Context, _id primitive.ObjectID, fine int) error {
	filter := bson.M{"_id": _id}
	update := bson.M{"$inc": bson.M{"fine": fine}}
	result, err := cs.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Modified cashback count:", result.ModifiedCount)
	return nil
}
