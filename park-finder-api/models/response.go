package models

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ReserveDetailResponse struct {
	Data   ReserveDetail `json:"data"`
	Status int           `json:"status"`
}

type ReserveDetail struct {
	ID           primitive.ObjectID `json:"reserve_id"`
	ProviderID   primitive.ObjectID `json:"provider_id"`
	ProviderName string             `json:"provider_name'`
	OrderId      string             `json:"order_id"`
	ParkingName  string             `json:"parking_name"`
	DateStart    string             `json:"date_start"`
	DateEnd      string             `json:"date_end"`
	HourStart    int                `json:"hour_start"`
	MinStart     int                `json:"min_start"`
	HourEnd      int                `json:"hour_end"`
	MinEnd       int                `json:"min_end"`
	Latitude     float64            `json:"latitude" bson:"latitude"`
	Longitude    float64            `json:"longitude" bson:"longitude"`
}

type CutomserRewardList struct {
	Data   []CustomerRedeemReward `json:"data"`
	Status int                    `json:"status"`
}

type CutomserHistoryPointResponse struct {
	Data   []CutomserHistoryPoint `json:"data"`
	Status int                    `json:"status"`
}

type CutomserHistoryPoint struct {
	Content         string    `json:"content"`
	Type            string    `json:"type"`
	Point           int       `json:"point"`
	TimeStampString string    `json:"time_stamp_string"`
	TimeStamp       time.Time `json:"time_stamp"`
}

type CustomerRedeemReward struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Point           int                `json:"point" bson:"point"`
	Title           string             `json:"title" bson:"title"`
	Description     string             `json:"description" bson:"description"`
	PreviewImageURL string             `json:"preview_url" bson:"preview_url"`
	Webhook         string             `json:"webhook" bson:"webhook"`
	Condition       []string           `json:"condition" bson:"condition"`
	QuotaCount      int                `json:"quota_count" bson:"quota_count"`
	CreateBy        string             `json:"create_by" bson:"create_by"`
	BarcodeURL      string             `json:"barcode_url" bson:"barcode_url"`
	ExpiredHour     int                `json:"customer_expired_date" bson:"customer_expired_date"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type StatusResponseFine struct {
	Data   *CustomerFine `json:"data"`
	Status int           `json:"status"`
}
type StatusResponseBool struct {
	Message bool `json:"message"`
	Status  int  `json:"status"`
}

type StatusReviewResponse struct {
	Message ReviewParkingArea `json:"message"`
	Status  int               `json:"status"`
}

type RegisterAreaLocationResponse struct {
	ParkingID interface{} `json:"parking_id"`
	Message   string      `json:"message"`
	Status    int         `json:"status"`
}

type StatusResponsePayment struct {
	Message       string `json:"message"`
	ReservationID string `json:"reservation_id"`
	Status        int    `json:"status"`
}
type LoginResponse struct {
	AccessToken string    `json:"access_token" validate:"required"`
	Data        *echo.Map `json:"data,omitempty"`
}

type CustomerProfileResponse struct {
	Profile *echo.Map `json:"profile"`
	Status  int       `json:"status"`
}

type CustomerAddressResponse struct {
	Data   []CustomerAddress `json:"address"`
	Status int               `json:"status"`
}

type CustomerCarResponse struct {
	Data   []Car `json:"data"`
	Status int   `json:"status"`
}

type SearchAreaResponse struct {
	Data   *[]ParkingArea `json:"data"`
	Status int            `json:"status"`
}

type RewardResponse struct {
	Data   *[]Reward `json:"data"`
	Status int       `json:"status"`
}
type RewardDetailResponse struct {
	Data   *Reward `json:"data"`
	Status int     `json:"status"`
}

type ParkingDetailResponse struct {
	Data   *ParkingArea `json:"data"`
	Status int          `json:"status"`
}

type LineReserveAPIResponse struct {
	ReturnCode    string      `json:"returnCode"`
	ReturnMessage string      `json:"returnMessage"`
	Info          InfoReserve `json:"info"`
}

type LineConfirmAPIResponse struct {
	ReturnCode    string      `json:"returnCode"`
	ReturnMessage string      `json:"returnMessage"`
	Info          InfoConfirm `json:"info"`
}

type ReserveResponse struct {
	ID      primitive.ObjectID `json:"_id"`
	OrderID string             `json:"order_id"`
	Status  int                `json:"status"`
}

type ExtendReserveResponse struct {
	Data   LineReserveRequest `json:"data"`
	Status int                `json:"status"`
}

type MyReserveResponse struct {
	Reservation []echo.Map `json:"data"`
	Status      int        `json:"status"`
}

type ListMessageRoomResponse struct {
	Data   *[]map[string]interface{} `json:"data"`
	Status int                       `json:"status"`
}

type ListNotificationResponse struct {
	Data   *[]Notification `json:"data"`
	Status int             `json:"status"`
}

type ListMessageLogResponse struct {
	Data   *MessageRoom `json:"data"`
	Status int          `json:"status"`
}

type DailyProfitResponse struct {
	Date string `json:"date"`
	Sum  int    `json:"sum"`
}

type WeeklyProfitResponse struct {
	Data  []DailyProfitResponse `json:"date"`
	Count int                   `json:"count"`
}

type AdminTransactionResponse struct {
	ProviderID  string      `json:"provider_id"`
	BankAccount BankAccount `json:"bank_information"`
	Count       int         `json:"count"`
	Sum         int         `json:"sum"`
	IsPay       bool        `json:"is_pay"`
}

type AccountAndBank struct {
	ID          string      `json:"id"`
	BankAccount BankAccount `json:"bank_information"`
}

type AdminTransactionResponseList struct {
	Data   *[]AdminTransactionResponse `json:"data"`
	Status int                         `json:"status"`
}
