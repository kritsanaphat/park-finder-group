package services

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ats AutomateTeserServices) ResetData() error {
	customer := ats.CheckExistEmailCustomer()
	provider := ats.CheckExistEmailProvider()
	ctx := context.Background()
	_, err := ats.AdminAccoutStorage.DeleteManyAccountInterface(ctx, bson.M{"email": "automate_test@gmail.com"})
	if err != nil {
		return err
	}
	ctx = context.Background()
	_, err = ats.CustomerAccoutStorage.DeleteManyAccountInterface(ctx, bson.M{"email": "automate_test@gmail.com"})
	if err != nil {
		return err
	}
	ctx = context.Background()
	_, err = ats.ProviderAccountStorage.DeleteManyAccountInterface(ctx, bson.M{"email": "automate_test@gmail.com"})
	if err != nil {
		return err
	}

	ctx = context.Background()
	_, err = ats.CarStorage.DeleteMany(ctx, bson.M{"name": "automate_test"})
	if err != nil {
		return err
	}

	ctx = context.Background()
	_, err = ats.ParkingAreaStorage.DeleteManyAreaInterface(ctx, bson.M{"parking_name": "automate_test"})
	if err != nil {
		return err
	}

	ctx = context.Background()
	_, err = ats.RewardStorage.DeleteManyRewardInterface(ctx, bson.M{"name": "automate_test"})
	if err != nil {
		return err
	}

	ctx = context.Background()
	if provider != nil {
		_, err = ats.ReserveStorage.DeleteMany(ctx, bson.M{"parking_name": "automate_test"})
		if err != nil {
			return err
		}
	}

	ctx = context.Background()
	if customer != nil {
		_, err = ats.TransactionStorage.DeleteMany(ctx, bson.M{"customer_email": customer.Email})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ats AutomateTeserServices) AddCashback() error {
	ctx := context.Background()

	_, err := ats.CustomerAccoutStorage.UpdateAccountInterface(ctx, bson.M{"email": "automate_test@gmail.com"}, bson.M{"$set": bson.M{"cashback": 100}})
	if err != nil {
		return err
	}
	return nil
}

func (ats AutomateTeserServices) CheckExistEmailCustomer() *models.CustomerAccount {
	ctx := context.Background()

	filter := bson.M{"email": "automate_test@gmail.com"}

	user := new(models.CustomerAccount)
	err := ats.CustomerAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return user
}

func (ats AutomateTeserServices) CheckExistEmailProvider() *models.ProviderAccount {
	ctx := context.Background()

	filter := bson.M{"email": "automate_test@gmail.com"}

	user := new(models.ProviderAccount)
	err := ats.CustomerAccoutStorage.FindAccountInterface(ctx, filter, user)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return user
}
