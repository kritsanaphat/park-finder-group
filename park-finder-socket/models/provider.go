package models

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProviderAccount struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	SSN               string             `json:"ssn" bson:"ssn"`
	FirstName         string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Birthday          string             `json:"birth_day" bson:"birth_day"`
	Phone             string             `json:"phone" bson:"phone"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"password" bson:"password"`
	ProfilePictureURL string             `json:"profile_picture_url" bson:"profile_picture_url"`
	Verify            bool               `json:"verify" bson:"verify"`
	Fine              int                `json:"fine" bson:"fine"`
	TimeStamp         time.Time          `json:"time_stamp" bson:"time_stamp"`
	BankAccount       BankAccount        `json:"bank_account" bson:"bank_account"`
}
type BankAccount struct {
	AccountBookImageUrl string `json:"account_book_image_url" bson:"account_book_image_url"`
	BankName            string `json:"bank_name" bson:"bank_name"`
	AccountName         string `json:"account_name" bson:"account_name"`
	AccountNumber       string `json:"account_number" bson:"account_number"`
}

type VerifyProvider struct {
	Verify    bool      `json:"verify" bson:"verify"`
	TimeStamp time.Time `json:"time_stamp" bson:"time_stamp"`
}

type ParkingArea struct {
	ID                     primitive.ObjectID  `json:"_id" bson:"_id"`
	ProviderID             primitive.ObjectID  `json:"provider_id" bson:"provider_id"`
	ParkingName            string              `json:"parking_name" bson:"parking_name"`
	Address                ParkingAddress      `json:"address" bson:"address"`
	Tag                    []string            `json:"tag" bson:"tag"`
	Price                  int16               `json:"price" bson:"price"`
	ParkingPictureUrl      string              `json:"parking_picture_url" bson:"parking_picture_url"`
	TitleDeedUrl           string              `json:"title_deed_url" bson:"title_deed_url"`
	LandCertificateUrl     string              `json:"land_certificate_url" bson:"land_certificate_url"`
	OverviewPictureUrl     []string            `json:"over_view_picture_url" bson:"over_view_picture_url"`
	MeasurementPictureUrl  []string            `json:"measurement_picture_url" bson:"measurement_picture_url"`
	ToatalParkingCount     int                 `json:"total_parking_count" bson:"total_parking_count"`
	IDCardUrl              string              `json:"id_card_url" bson:"id_card_url"`
	OpenStatus             bool                `json:"open_status" bson:"open_status"`
	OpenDetail             Daily               `json:"open_detail" bson:"open_detail"`
	TimeStampClose         *time.Time          `json:"time_stamp_close" bson:"time_stamp_close"`
	StatusApply            string              `json:"status_apply" bson:"status_apply"`
	StatusApplyDescription string              `json:"status_apply_description" bson:"status_apply_description"`
	Review                 []ReviewParkingArea `json:"review" bson:"review"`
	Distance               float32             `json:"distance" bson:"distance"`
	ReserveLog             []ReserveLog        `json:"reserve_log" bson:"reserve_log"`
	TimeStamp              time.Time           `json:"time_stamp" bson:"time_stamp"`
}

type ParkingAreaOpenDetailUpdate struct {
	OpenDetail Daily `json:"open_detail" bson:"open_detail"`
}

type ReserveLog struct {
	CustomerEmail string `bson:"customer_email"`
	DateStart     string `bson:"date_start"`
	DateEnd       string `bson:"date_end"`
	HourStart     int    `bson:"hour_start"`
	MinStart      int    `bson:"min_start"`
	HourEnd       int    `bson:"hour_end"`
	MinEnd        int    `bson:"min_end"`
}

type ReviewParkingArea struct {
	ReviewID    primitive.ObjectID `json:"review_id" bson:"review_id"`
	CustomerID  primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	ReviewScore int                `json:"review_score" bson:"review_score"`
	TimeStamp   time.Time          `json:"time_stamp" bson:"time_stamp"`
	Comment     string             `json:"comment" bson:"comment"`
}

type ParkingAddress struct {
	AddressText  string  `json:"address_text" bson:"address_text"`
	Sub_district string  `json:"sub_district" bson:"sub_district"`
	District     string  `json:"district" bson:"district"`
	Province     string  `json:"province" bson:"province"`
	Postal_code  string  `json:"postal_code" bson:"postal_code"`
	Latitude     float64 `json:"latitude" bson:"latitude"`
	Longitude    float64 `json:"longitude" bson:"longitude"`
}

type Daily struct {
	Monday    OpenTimeDetail `json:"monday" bson:"monday"`
	Tuesday   OpenTimeDetail `json:"tuesday" bson:"tuesday"`
	Wednesday OpenTimeDetail `json:"wednesday" bson:"wednesday"`
	Thursday  OpenTimeDetail `json:"thursday" bson:"thursday"`
	Friday    OpenTimeDetail `json:"friday" bson:"friday"`
	Saturday  OpenTimeDetail `json:"saturday" bson:"saturday"`
	Sunday    OpenTimeDetail `json:"sunday" bson:"sunday"`
}

type OpenTimeDetail struct {
	OpenTime  int `json:"open_time" bson:"open_time"`
	CloseTime int `json:"close_time" bson:"close_time"`
}

type Profit struct {
	ParkingName string `json:"parking_name"`
	Address     string `json:"address"`
	Count       int    `json:"count"`
	Profit      int    `json:"profit"`
}

func (p ProviderAccount) IDToString() string {
	return p.ID.Hex()
}

func (p ProviderAccount) ToMap() *echo.Map {
	return &echo.Map{
		"id":         p.IDToString(),
		"first_name": p.FirstName,
		"last_name":  p.LastName,
		"email":      p.Email,
		"verify":     p.Verify,
		"time_stamp": p.TimeStamp,
	}
}

func (p ProviderAccount) ToMapProfile() *echo.Map {
	return &echo.Map{
		"_id":                 p.ID,
		"ssn":                 p.SSN,
		"first_name":          p.FirstName,
		"last_name":           p.LastName,
		"email":               p.Email,
		"verify":              p.Verify,
		"birth_day":           p.Birthday,
		"phone":               p.Phone,
		"profile_picture_url": p.ProfilePictureURL,
	}
}
