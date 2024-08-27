package services

import (
	"context"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminServices struct {
	AdminAccoutStorage *storage.AccoutStorage
	RewardStorage      *storage.RewardStorage
	ProviderStorage    *storage.AccoutStorage
	TransactionStorage *storage.LogStorage
	ReceiptStorage     *storage.LogStorage
	ParkingAreaStorage *storage.ParkingAreaStorage
	TokenStorage       *storage.TokenStorage
}

type IAdminServices interface {
	AddReward(ctx context.Context, user *models.AddRewardRequest, email string) (string, error)
	AdminRegister(ctx context.Context, user *models.RegisterAccountRequest) error

	CheckExistEmail(ctx context.Context, email string) *models.AdminAccount
	CheckVerifyEmail(ctx context.Context, email string) *models.AdminAccount
	CheckPassword(user models.AdminAccount, password string) bool

	AdminGetParkingArea(ctx context.Context, status string) *[]models.ParkingArea
	AdminUpdateParkingArea(ctx context.Context, parking_id, status, description string) error
	AdminCheckParkingAreaDetail(ctx context.Context, parking_id string) *models.ParkingArea

	AdminGetTransaction(ctx context.Context, year, month string) (*[]models.AdminTransactionResponse, error)
	AdminSubmitReceipt(ctx context.Context, image_url, provider_id, month, year string, price int) error

	AddToken(ctx context.Context, user *models.AdminAccount, token string) error
	CheckExistToken(ctx context.Context, tk string) *models.Token
	RevokeToken(ctx context.Context, user *models.AdminAccount, token string) error
	RevokeExpireToken(ctx context.Context, token string, expireDate time.Time) error
}

func NewAdminServices(
	db *mongo.Database,
) IAdminServices {
	admin_account_storage := storage.NewAdminAccoutStorage(db)
	reward_storage := storage.NewRewardStorage(db)
	parking_area_storage := storage.NewParkingAreaStorage(db)
	transaction_storage := storage.NewTransactionStorage(db)
	token_storage := storage.NewTokenStorage(db)
	provider_storage := storage.NewProviderAccoutStorage(db)
	receipt_storage := storage.NewReceiptStorage(db)
	return AdminServices{
		AdminAccoutStorage: admin_account_storage,
		RewardStorage:      reward_storage,
		ParkingAreaStorage: parking_area_storage,
		TokenStorage:       token_storage,
		TransactionStorage: transaction_storage,
		ProviderStorage:    provider_storage,
		ReceiptStorage:     receipt_storage,
	}
}
