package services

import (
	"gitlab.com/parking-finder/parking-finder-api/internal/storage"

	"go.mongodb.org/mongo-driver/mongo"
)

type AutomateTeserServices struct {
	CustomerAccoutStorage  *storage.AccoutStorage
	ParkingAreaStorage     *storage.ParkingAreaStorage
	TokenStorage           *storage.TokenStorage
	CarStorage             *storage.CarStorage
	RewardStorage          *storage.RewardStorage
	ReserveStorage         *storage.LogStorage
	TransactionStorage     *storage.LogStorage
	NotificationStorage    *storage.LogStorage
	AdminAccoutStorage     *storage.AccoutStorage
	ProviderAccountStorage *storage.AccoutStorage
}

type IAutomateTeserServices interface {
	ResetData() error
	AddCashback() error
}

func NewAutomateTeserServices(
	db *mongo.Database,
) IAutomateTeserServices {
	customer_account_storage := storage.NewCustomerAccoutStorage(db)
	parking_area_storage := storage.NewParkingAreaStorage(db)
	token_storage := storage.NewTokenStorage(db)
	car_storage := storage.NewCarStorage(db)
	reward_storage := storage.NewRewardStorage(db)
	reserve_storage := storage.NewReserveStorage(db)
	transaction_storage := storage.NewTransactionStorage(db)
	notification_storage := storage.NewNotificationStorage(db)
	admin_account_storage := storage.NewAdminAccoutStorage(db)
	provider_account_storage := storage.NewProviderAccoutStorage(db)

	return AutomateTeserServices{
		CustomerAccoutStorage:  customer_account_storage,
		TokenStorage:           token_storage,
		CarStorage:             car_storage,
		ParkingAreaStorage:     parking_area_storage,
		RewardStorage:          reward_storage,
		ReserveStorage:         reserve_storage,
		TransactionStorage:     transaction_storage,
		NotificationStorage:    notification_storage,
		AdminAccoutStorage:     admin_account_storage,
		ProviderAccountStorage: provider_account_storage,
	}
}
