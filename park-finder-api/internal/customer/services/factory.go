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

type CustomerServices struct {
	CustomerAccoutStorage *storage.AccoutStorage
	ParkingAreaStorage    *storage.ParkingAreaStorage
	TokenStorage          *storage.TokenStorage
	CarStorage            *storage.CarStorage
	RewardStorage         *storage.RewardStorage
	Redis                 *connector.Redis
	HttpClient            *httpclient.HTTPClient
}

type ICustomerServices interface {
	CustomerRegister(ctx context.Context, user *models.RegisterAccountRequest) error
	UpdateCustomerVerify(ctx context.Context, email string) error
	UpdateCustomerProfile(ctx context.Context, email string, user *models.UpdateProfileRequest) error
	UpdateCustomerPassword(ctx context.Context, user *models.ResetPasswordRequest, email string) error
	UpdateCustomerPasswordByForgot(ctx context.Context, password string, email string) error

	CheckExistEmail(ctx context.Context, email string) *models.CustomerAccount
	CheckVerifyEmail(ctx context.Context, email string) *models.CustomerAccount
	CheckPassword(user models.CustomerAccount, password string) bool

	AddToken(ctx context.Context, user *models.CustomerAccount, token string) error
	CheckExistToken(ctx context.Context, tk string) *models.Token
	RevokeToken(ctx context.Context, user *models.CustomerAccount, token string) error
	RevokeExpireToken(ctx context.Context, token string, expireDate time.Time) error

	SaveOTP(email string, otp string) error
	CheckOTP(email, otp string) error
	RemoveOTP(email string, otp string) error

	SaveOTPForgot(email string, otp string) error
	CheckOTPForgot(email, otp string) error
	RemoveOTPForgot(email string, otp string) error

	CustomerRegisterCar(ctx context.Context, car *models.RegisterCarRequest, email string) error
	CustomerCar(ctx context.Context, email string) []models.Car
	CheckCustomerCarDetail(ctx context.Context, id string) *models.Car
	UpdateCustomerCar(ctx context.Context, car *models.UpdateCustomerCar, car_id string) error
	DeleteCustomerCar(ctx context.Context, car_id string) error
	CheckExistCarName(ctx context.Context, email, name string) *models.Car
	CheckExistCarLicensePlate(ctx context.Context, email, license_plate string) *models.Car

	CustomerRegisterAddress(ctx context.Context, address *models.RegisterAddressRequest, email string) error
	CheckExistLocationName(ctx context.Context, email, location_name string) []models.CustomerAddress
	UpdateCustomerAddress(ctx context.Context, request *models.RegisterAddressRequest, email string, address_id string) error
	DeleteCustomerAddress(ctx context.Context, address_id string, email string) error
	UpdateCustomerDefaultAddress(ctx context.Context, addressID string, email string) error
	UpdateCustomerDefaultCar(ctx context.Context, carID string, email string) error

	UpdateCustomerFavoriteArea(ctx context.Context, email, parking_id, action string) error
	CustomerFavoriteArea(ctx context.Context, area_id []string) *[]models.ParkingArea

	CustomerReward(ctx context.Context) *[]models.Reward
	CustomerRewardDetail(ctx context.Context, id string) *models.Reward
	CustomerRedeemReward(ctx context.Context, id string, account models.CustomerAccount) (string, error)

	CustomerRefund(ctx context.Context, email string, cashback int) error
	CustomerFine(ctx context.Context, _id primitive.ObjectID, fine int, reserve *models.Reservation) error
	CustomerUpdateCashback(ctx context.Context, email string, cashback int) error

	CheckExistReview(order_id string) bool
	RemoveReviewCache(order_id string) error
}

func NewCustomerServices(
	db *mongo.Database,
	rd *connector.Redis,
	ht *httpclient.HTTPClient,

) ICustomerServices {
	customer_account_storage := storage.NewCustomerAccoutStorage(db)
	parking_area_storage := storage.NewParkingAreaStorage(db)
	token_storage := storage.NewTokenStorage(db)
	car_storage := storage.NewCarStorage(db)
	reward_storage := storage.NewRewardStorage(db)
	return CustomerServices{
		CustomerAccoutStorage: customer_account_storage,
		TokenStorage:          token_storage,
		CarStorage:            car_storage,
		ParkingAreaStorage:    parking_area_storage,
		RewardStorage:         reward_storage,
		Redis:                 rd,
		HttpClient:            ht,
	}
}
