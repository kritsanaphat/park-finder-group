package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (cs CustomerServices) UpdateCustomerProfile(ctx context.Context, email string, user *models.UpdateProfileRequest) error {
	location, _ := time.LoadLocation("Asia/Bangkok")
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": models.UpdateCustomerProfile{
			SSN:               user.SSN,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Birthday:          user.Birthday,
			Phone:             user.Phone,
			ProfilePictureURL: user.ProfilePictureURL,
			TimeStamp:         time.Now().In(location),
		}}

	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) CustomerRegisterCar(ctx context.Context, car *models.RegisterCarRequest, email string) error {
	filter := bson.M{"customer_email": email}
	c := new(models.Car)
	data := new(models.Car)

	err := cs.CarStorage.FindCarByInterface(ctx, filter, c)
	if err == mongo.ErrNoDocuments {
		data = &models.Car{
			ID:            primitive.NewObjectID(),
			CustomerEmail: email,
			Name:          car.Name,
			LicensePlate:  car.LicensePlate,
			Brand:         car.Brand,
			Model:         car.Model,
			Color:         car.Color,
			CarPictureURL: car.CarPictureURL,
			TimeStamp:     time.Now(),
			Default:       true,
		}
	} else {
		data = &models.Car{
			ID:            primitive.NewObjectID(),
			CustomerEmail: email,
			Name:          car.Name,
			LicensePlate:  car.LicensePlate,
			Brand:         car.Brand,
			Model:         car.Model,
			Color:         car.Color,
			CarPictureURL: car.CarPictureURL,
			TimeStamp:     time.Now(),
			Default:       false,
		}

	}

	cs.CarStorage.InsertCar(ctx, data)

	return nil
}

func (cs CustomerServices) CustomerCar(ctx context.Context, email string) []models.Car {
	var cars []models.Car
	cars, err := cs.CarStorage.FindCarByEmail(ctx, email)
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return cars
}

func (cs CustomerServices) CheckCustomerCarDetail(ctx context.Context, id string) *models.Car {
	car := cs.CarStorage.FindCarById(ctx, id)
	if car == nil {
		return nil
	}

	return car
}

func (cs CustomerServices) UpdateCustomerCar(ctx context.Context, car *models.UpdateCustomerCar, car_id string) error {
	location, _ := time.LoadLocation("Asia/Bangkok")

	id, err := primitive.ObjectIDFromHex(car_id)
	if err != nil {
		fmt.Println("Error:", err)

	}
	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$set": models.UpdateCustomerCar{
			Name:          car.Name,
			LicensePlate:  car.LicensePlate,
			Brand:         car.Brand,
			Model:         car.Model,
			Color:         car.Color,
			CarPictureURL: car.CarPictureURL,
			TimeStamp:     time.Now().In(location),
		}}

	_, err = cs.CarStorage.UpdateCar(ctx, filter, update)
	return err
}

func (cs CustomerServices) DeleteCustomerCar(ctx context.Context, car_id string) error {
	id, err := primitive.ObjectIDFromHex(car_id)
	if err != nil {
		fmt.Println("Error:", err)

	}
	filter := bson.M{
		"_id": id,
	}
	_, err = cs.CarStorage.DeleteCar(ctx, filter)
	return err
}

func (cs CustomerServices) CustomerRegisterAddress(ctx context.Context, request *models.RegisterAddressRequest, email string) error {

	address_text_split := strings.Split(request.Address, " ")

	filter := bson.M{
		"email": email,
	}

	var address_text string

	for i := 0; i < len(address_text_split); i++ {
		if i+6 <= len(address_text_split) {
			address_text += address_text_split[i] + " "
		}
	}

	address := models.CustomerAddress{
		AddressID:    primitive.NewObjectID(),
		AddressText:  address_text,
		SubDistrict:  address_text_split[len(address_text_split)-5],
		District:     address_text_split[len(address_text_split)-4],
		Province:     address_text_split[len(address_text_split)-3],
		Postal_code:  address_text_split[len(address_text_split)-2],
		Latitude:     request.Latitude,
		Longitude:    request.Longitude,
		LocationName: request.LocationName,
		TimeStamp:    time.Now(),
	}

	update := bson.M{
		"$push": models.UpdateCustomerAddress{
			Address: address,
		},
	}

	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) UpdateCustomerAddress(ctx context.Context, request *models.RegisterAddressRequest, email string, address_id string) error {

	id, err := primitive.ObjectIDFromHex(address_id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filter := bson.M{
		"email":              email,
		"address.address_id": id,
	}

	var address_text string
	address_text_split := strings.Split(request.Address, " ")

	for i := 0; i < len(address_text_split); i++ {
		if i+6 <= len(address_text_split) {
			address_text += address_text_split[i] + " "
		}
	}

	address := models.CustomerAddress{
		AddressID:    primitive.NewObjectID(),
		AddressText:  address_text,
		SubDistrict:  address_text_split[len(address_text_split)-5],
		District:     address_text_split[len(address_text_split)-4],
		Province:     address_text_split[len(address_text_split)-3],
		Postal_code:  address_text_split[len(address_text_split)-2],
		Latitude:     request.Latitude,
		Longitude:    request.Longitude,
		LocationName: request.LocationName,
		TimeStamp:    time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"address.$": address,
		},
	}
	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) DeleteCustomerAddress(ctx context.Context, address_id string, email string) error {
	id, err := primitive.ObjectIDFromHex(address_id)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	filter := bson.M{
		"email": email,
	}

	update := bson.M{
		"$pull": bson.M{
			"address": bson.M{"address_id": id},
		},
	}

	result, err := cs.CustomerAccoutStorage.UpdateAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) UpdateCustomerDefaultAddress(ctx context.Context, addressID string, email string) error {
	id, err := primitive.ObjectIDFromHex(addressID)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	filter := bson.M{
		"email": email,
	}

	update := bson.M{
		"$set": bson.M{
			"address.$[].default": false,
		},
	}
	_, err = cs.CustomerAccoutStorage.UpdateManyAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	filter = bson.M{
		"email":              email,
		"address.address_id": id,
	}

	update = bson.M{
		"$set": bson.M{
			"address.$.default": true,
		},
	}

	result, err := cs.CustomerAccoutStorage.UpdateManyAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) UpdateCustomerDefaultCar(ctx context.Context, CarID string, email string) error {
	id, err := primitive.ObjectIDFromHex(CarID)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	filter := bson.M{
		"customer_email": email,
	}

	update := bson.M{
		"$set": bson.M{
			"default": false,
		},
	}
	_, err = cs.CarStorage.UpdateManyAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	filter = bson.M{
		"_id": id,
	}

	update = bson.M{
		"$set": bson.M{
			"default": true,
		},
	}

	result, err := cs.CarStorage.UpdateManyAccountInterface(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no documents were modified")
	}

	return err
}

func (cs CustomerServices) CheckExistCarName(ctx context.Context, email, name string) *models.Car {
	filter := bson.M{
		"customer_email": email,
		"name":           name,
	}
	car := new(models.Car)
	err := cs.CarStorage.FindCarByInterface(ctx, filter, car)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return car
}

func (cs CustomerServices) CheckExistCarLicensePlate(ctx context.Context, email, license_plate string) *models.Car {
	filter := bson.M{
		"customer_email": email,
		"license_plate":  license_plate,
	}
	car := new(models.Car)
	err := cs.CarStorage.FindCarByInterface(ctx, filter, car)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return car
}

func (cs CustomerServices) CheckExistLocationName(ctx context.Context, email, location_name string) []models.CustomerAddress {
	filter := bson.M{
		"email":                 email,
		"address.location_name": location_name,
	}
	account := new(models.CustomerAccount)
	err := cs.CustomerAccoutStorage.FindAccountInterface(ctx, filter, account)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return account.Address
}

func (cs CustomerServices) UpdateCustomerFavoriteArea(ctx context.Context, email, parking_id, action string) error {
	if action == "push" {
		result := cs.CustomerAccoutStorage.PushFavoriteArea(ctx, parking_id, email)
		if result != nil {
			return result
		}

	} else if action == "pull" {
		result := cs.CustomerAccoutStorage.PullFavoriteArea(ctx, parking_id, email)
		if result != nil {
			return result
		}

	}
	return nil
}

func (cs CustomerServices) CustomerFavoriteArea(ctx context.Context, area_id []string) *[]models.ParkingArea {
	area := cs.ParkingAreaStorage.FindParkingByParkingIDList(ctx, area_id)
	if area == nil {
		return nil
	}
	return area
}

func (cs CustomerServices) CustomerReward(ctx context.Context) *[]models.Reward {
	var rewards []models.Reward
	cursor := cs.RewardStorage.FindRewardByExpiredDate(ctx)
	if cursor == nil {
		return nil
	}
	for cursor.Next(ctx) {
		var reward models.Reward
		if err := cursor.Decode(&reward); err == nil {
			rewards = append(rewards, reward)
		}
	}
	return &rewards
}

func (cs CustomerServices) CustomerRewardDetail(ctx context.Context, id string) *models.Reward {
	rewards := cs.RewardStorage.FindRewardByID(ctx, id)
	if rewards == nil {
		return nil
	}

	return rewards
}

func (cs CustomerServices) CustomerRedeemReward(ctx context.Context, id string, account models.CustomerAccount) (string, error) {
	rewards := cs.RewardStorage.FindRewardByID(ctx, id)
	if rewards == nil {
		return "", errors.New("Reward doesn't exist")
	}
	if account.Point < int(rewards.Point) {
		return "", errors.New("Your Point doesn't enough")
	}
	if rewards.QuotaCount < 1 {
		return "", errors.New("The quota is full")
	}
	barcode := cs.HttpClient.SendToPTTReward(rewards.Webhook)
	if barcode == "" {
		return "", errors.New("Error form third party")
	}
	err := cs.CustomerAccoutStorage.PushRedeemReward(ctx, account.ID, rewards.ID, barcode, rewards.Name, rewards.Point)
	if err != nil {
		return "", err
	}
	err = cs.RewardStorage.RemoveQuotaCount(ctx, rewards.ID, rewards.QuotaCount)
	if err != nil {
		return "", err
	}

	return barcode, nil
}

func (cs CustomerServices) CustomerRefund(ctx context.Context, email string, cashback int) error {
	err := cs.CustomerAccoutStorage.AddCashback(ctx, email, cashback)
	if err != nil {
		return err
	}
	return nil
}

func (cs CustomerServices) CustomerFine(ctx context.Context, _id primitive.ObjectID, fine int, reserve *models.Reservation) error {
	err := cs.CustomerAccoutStorage.UpdateFineCustomer(ctx, _id, fine, *reserve)
	if err != nil {
		return err
	}
	return nil
}

func (cs CustomerServices) CustomerUpdateCashback(ctx context.Context, email string, cashback int) error {
	err := cs.CustomerAccoutStorage.UpdateCashback(ctx, email, cashback)
	if err != nil {
		return err
	}
	return nil
}

func (cs CustomerServices) CheckExistReview(order_id string) bool {
	result := cs.Redis.Get(order_id)
	_, err := result.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false
		}
		return false
	}
	return true
}

func (cs CustomerServices) RemoveReviewCache(order_id string) error {
	ctx := context.Background()

	result := cs.Redis.Client.Del(ctx, order_id)
	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}

func (cs CustomerServices) RemoveCacheReview(order_id string) error {
	ctx := context.Background()

	result := cs.Redis.Client.Del(ctx, order_id)
	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}
