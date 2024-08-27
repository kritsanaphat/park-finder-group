package models

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerAccount struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	SSN               string             `json:"ssn" bson:"ssn"`
	FirstName         string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Birthday          string             `json:"birth_day" bson:"birth_day"`
	Phone             string             `json:"phone" bson:"phone"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"password" bson:"password"`
	Point             int                `json:"point" bson:"point"`
	ProfilePictureURL string             `json:"profile_picture_url" bson:"profile_picture_url"`
	Address           []CustomerAddress  `json:"address" bson:"address"`
	Verify            bool               `json:"verify" bson:"verify"`
	Reward            []CustomerReward   `json:"reward" bson:"reward"`
	Cashback          int                `json:"cashback" bson:"cashback"`
	FavoriteArea      []string           `json:"favorite_area" bson:"favorite_area"`
	Fine              *CustomerFine      `json:"fine" bson:"fine"`
	TimeStamp         time.Time          `json:"time_stamp" bson:"time_stamp"`
}

type CustomerReward struct {
	ID          primitive.ObjectID `json:"reward_id" bson:"reward_id"`
	Name        string             `json:"name" bson:"name"`
	Point       int                `json:"point" bson:"point"`
	BarcodeURL  string             `json:"barcode_url" bson:"barcode_url"`
	TimeStamp   time.Time          `json:"time_stamp" bson:"time_stamp"`
	ExpiredDate time.Time          `json:"expired_date" bson:"expired_date"`
}

type CustomerFine struct {
	OrderID     string             `json:"order_id" bson:"order_id"`
	ParkingID   primitive.ObjectID `json:"parking_id" bson:"parking_id"`
	ProviderID  primitive.ObjectID `json:"provider_id" bson:"provider_id"`
	Quantity    float32            `json:"quantity" bson:"quantity"`
	Price       int                `json:"price" bson:"price"`
	ParkingName string             `json:"parking_name" bson:"parking_name"`
}

type CustomerAddress struct {
	AddressID    primitive.ObjectID `json:"address_id" bson:"address_id"`
	AddressText  string             `json:"address_text" bson:"address_text"`
	SubDistrict  string             `json:"sub_district" bson:"sub_district"`
	District     string             `json:"district" bson:"district"`
	Province     string             `json:"province" bson:"province"`
	Postal_code  string             `json:"postal_code" bson:"postal_code"`
	Latitude     float64            `json:"latitude" bson:"latitude"`
	Longitude    float64            `json:"longitude" bson:"longitude"`
	LocationName string             `json:"location_name" bson:"location_name"`
	Default      bool               `json:"default" bson:"default"`
	TimeStamp    time.Time          `json:"time_stamp" bson:"time_stamp"`
}

type UpdateCustomerProfile struct {
	SSN               string    `json:"ssn" bson:"ssn"`
	FirstName         string    `json:"first_name" bson:"first_name"`
	LastName          string    `json:"last_name" bson:"last_name"`
	Birthday          string    `json:"birth_day" bson:"birth_day"`
	Phone             string    `json:"phone" bson:"phone"`
	ProfilePictureURL string    `json:"profile_picture_url" bson:"profile_picture_url"`
	TimeStamp         time.Time `json:"time_stamp" bson:"time_stamp"`
}

type UpdateCustomerAddress struct {
	Address CustomerAddress `json:"address" bson:"address"`
}

type UpdateCustomerReward struct {
	Reward CustomerReward `json:"reward" bson:"reward"`
}
type UpdatePassword struct {
	Password  string    `json:"password" bson:"password"`
	TimeStamp time.Time `json:"time_stamp" bson:"time_stamp"`
}

type VerifyCustomer struct {
	Verify    bool      `json:"verify" bson:"verify"`
	TimeStamp time.Time `json:"time_stamp" bson:"time_stamp"`
}

type Car struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	CustomerEmail string             `json:"customer_email" bson:"customer_email"`
	Name          string             `json:"name" bson:"name"`
	LicensePlate  string             `json:"license_plate" bson:"license_plate"`
	Brand         string             `json:"brand" bson:"brand"`
	Model         string             `json:"model" bson:"model"`
	Color         string             `json:"color" bson:"color"`
	CarPictureURL string             `json:"car_picture_url" bson:"car_picture_url"`
	Default       bool               `json:"default" bson:"default"`
	TimeStamp     time.Time          `json:"time_stamp" bson:"time_stamp"`
}

type UpdateCustomerCar struct {
	Name          string    `json:"name" bson:"name"`
	LicensePlate  string    `json:"license_plate" bson:"license_plate"`
	Brand         string    `json:"brand" bson:"brand"`
	Model         string    `json:"model" bson:"model"`
	Color         string    `json:"color" bson:"color"`
	CarPictureURL string    `json:"car_picture_url" bson:"car_picture_url"`
	TimeStamp     time.Time `json:"time_stamp" bson:"time_stamp"`
}

type MyReserve struct {
	Status          string `json:"status"`
	Pirce           int    `json:"price"`
	ParkingImageUrl string `json:"parking_image_url"`
	ParkingAddress  string `json:"parking_address"`
	Date            string `json:"date"`
}

func (c CustomerAccount) IDToString() string {
	return c.ID.Hex()
}

func (c CustomerAccount) ToMap() *echo.Map {
	return &echo.Map{
		"id":         c.IDToString(),
		"first_name": c.FirstName,
		"last_name":  c.LastName,
		"email":      c.Email,
		"verify":     c.Verify,
		"time_stamp": c.TimeStamp,
	}
}

func (c CustomerAccount) ToMapProfile() *echo.Map {
	return &echo.Map{
		"_id":                 c.IDToString(),
		"ssn":                 c.SSN,
		"first_name":          c.FirstName,
		"last_name":           c.LastName,
		"email":               c.Email,
		"verify":              c.Verify,
		"birth_day":           c.Birthday,
		"phone":               c.Phone,
		"point":               c.Point,
		"cashback":            c.Cashback,
		"profile_picture_url": c.ProfilePictureURL,
	}
}
