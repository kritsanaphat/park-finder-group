package services

import (
	"context"
	"time"

	"gitlab.com/parking-finder/parking-finder-api/internal/httpclient"
	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"gitlab.com/parking-finder/parking-finder-api/pkg/connector"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProviderServices struct {
	ProviderAccoutStorage *storage.AccoutStorage
	TokenStorage          *storage.TokenStorage
	ParkingAreaStorage    *storage.ParkingAreaStorage
	ReserveStorage        *storage.LogStorage
	TransactionStorage    *storage.LogStorage
	Redis                 *connector.Redis
	HttpClient            *httpclient.HTTPClient
}

type IProviderServices interface {
	ProviderRegister(ctx context.Context, user *models.RegisterAccountRequest) error
	UpdateProviderVerify(ctx context.Context, email string) error
	UpdateProviderProfile(ctx context.Context, email string, user *models.UpdateProfileRequest) error
	UpdateProviderBankAccount(ctx context.Context, email string, bank *models.BankAccount) error
	UpdateProviderPassword(ctx context.Context, user *models.ResetPasswordRequest, email string) error
	CheckExistEmail(ctx context.Context, email string) *models.ProviderAccount
	CheckVerifyEmail(ctx context.Context, email string) *models.ProviderAccount
	CheckPassword(user models.ProviderAccount, password string) bool

	SaveOTP(email string, otp string) error
	CheckOTP(email, otp string) error
	RemoveOTP(email string, otp string) error
	SaveOTPForgot(email string, otp string) error
	CheckOTPForgot(email, otp string) error
	RemoveOTPForgot(email string, otp string) error
	UpdateCustomerPasswordByForgot(ctx context.Context, password string, email string) error

	AddToken(ctx context.Context, user *models.ProviderAccount, token string) error
	CheckExistToken(ctx context.Context, tk string) *models.Token
	RevokeToken(ctx context.Context, user *models.ProviderAccount, token string) error
	RevokeExpireToken(ctx context.Context, token string, expireDate time.Time) error

	ProviderRegisterAreaLocaion(ctx context.Context, area *models.RegisterParkingAreaFirstStepRequest, email string) (interface{}, error)
	ProviderRegisterAreaDocument(ctx context.Context, area *models.RegisterParkingAreaDocumentStepRequest, email, parking_id string) error

	GetProviderArea(ctx context.Context, id string) []models.ParkingArea
	UpdateOpenAreaDailyStatus(ctx context.Context, daily *models.UpdateOpenAreaDailyStatusRequest) error
	UpdateOpenAreaQuickStatus(ctx context.Context, id, status string, range_time int) error
	UpdateOpenAreaInAdvanceStatus(ctx context.Context, id, status string, date []string) error
	CheckQuickAvailabilityCloseProviderArea(ctx context.Context, parking_area_id string, range_time int) *[]models.Reservation
	CheckDailyAvailabilityCloseProviderArea(ctx context.Context, request models.UpdateOpenAreaDailyStatusRequest) *[]models.Reservation
	CheckUpdatePriceArea(ctx context.Context, parking_area_id string, price int16) error

	CheckValidProviderArea(ctx context.Context, provider_id, parking_area_id string) error
	FineProvider(ctx context.Context, _id primitive.ObjectID, fine int) error

	//profit
	GetProviderProfitDaily(ctx context.Context, list_parking_id []string) (*models.DailyProfitResponse, int, error)
	GetProviderProfitWeekly(ctx context.Context, list_parking_id []string) (*[]models.DailyProfitResponse, int, error)
	GetProviderProfitMontly(ctx context.Context, list_parking_id []string) (*[]models.DailyProfitResponse, int, error)
	GetProviderProfitYearly(ctx context.Context, list_parking_id []string) (*[]models.DailyProfitResponse, int, error)
}

func NewProviderServices(
	db *mongo.Database,
	rd *connector.Redis,
	ht *httpclient.HTTPClient,
) IProviderServices {
	provider_account_storage := storage.NewProviderAccoutStorage(db)
	token_storage := storage.NewTokenStorage(db)
	parking_area_storage := storage.NewParkingAreaStorage(db)
	reserve_storage := storage.NewReserveStorage(db)
	transaction_storage := storage.NewTransactionStorage(db)
	return ProviderServices{
		ProviderAccoutStorage: provider_account_storage,
		TokenStorage:          token_storage,
		ParkingAreaStorage:    parking_area_storage,
		ReserveStorage:        reserve_storage,
		TransactionStorage:    transaction_storage,
		Redis:                 rd,
		HttpClient:            ht,
	}
}
