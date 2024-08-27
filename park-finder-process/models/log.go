package models

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID                primitive.ObjectID `bson:"_id"`
	ParkingID         primitive.ObjectID `bson:"parking_id"`
	ProviderID        primitive.ObjectID `bson:"provider_id"`
	CarID             primitive.ObjectID `bson:"car_id"`
	TransactionID     primitive.ObjectID `bson:"transaction_id"`
	ParkingName       string             `bson:"parking_name"`
	CustomerEmail     string             `bson:"customer_email"`
	Price             int                `bson:"price"`
	Status            string             `bson:"status"`
	OrderID           string             `bson:"order_id"`
	Date              string             `bson:"date"`
	HourStart         int                `bson:"hour_start"`
	MinStart          int                `bson:"min_start"`
	HourEnd           int                `bson:"hour_end"`
	MinEnd            int                `bson:"min_end"`
	PaymentChanel     string             `bson:"payment_chanel"`
	ParkingPictureUrl string             `bson:"parking_picture_url"`
	Address_Full      string             `bson:"address_full"`
	TimeStamp         time.Time          `bson:"time_stamp"`
}

type TransactionLinePay struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"` //ReserveAPI
	ParkingID     primitive.ObjectID `json:"parking_id" bson:"parking_id"`
	CustomerEmail string             `json:"customer_email" bson:"customer_email"` //ReserveAPI
	OrderId       string             `json:"order_Id" bson:"order_id"`             //ReserveAPI
	Packages      []Package          `json:"packages" bson:"package"`              //ReserveAPI
	Info          InfoReserve        `json:"info" bson:"info"`                     //res of ReserveAPI
	Status        string             `json:"status" bson:"status"`                 //con
	TimeStamp     time.Time          `json:"time_stamp" bson:"time_stamp"`         //ReserveAPI
}

type Package struct {
	PackageID string    `json:"id" bson:"package_id"`
	Amount    float32   `json:"amount" bson:"amount"`
	Name      string    `json:"name" bson:"name"`
	Products  []Product `json:"products" bson:"products"`
}

type Product struct {
	Name     string  `json:"name" bson:"name"`
	Quantity float32 `json:"quantity" bson:"quantity"`
	Price    int     `json:"price" bson:"price"`
	ImageURL string  `json:"imageUrl" bson:"imageUrl"`
}

type RedirectUrls struct {
	ConfirmUrl     string `json:"confirmUrl" bson:"confirmUrl"`
	CancelUrl      string `json:"cancelUrl" bson:"cancelUrl"`
	ConfirmUrlType string `json:"confirmUrlType" bson:"confirmUrlType"`
}

type PaymentUrl struct {
	Web string `json:"web" bson:"web"`
	App string `json:"app" bson:"app"`
}

type InfoReserve struct {
	PaymentUrl         PaymentUrl `json:"paymentUrl" bson:"paymentUrl"`
	TransactionId      int64      `json:"transactionId" bson:"transactionId"`
	PaymentAccessToken string     `json:"paymentAccessToken" bson:"paymentAccessToken"`
}

type PayInfo struct {
	Method string `json:"method"`
	Amount int    `json:"amount"`
}

type InfoConfirm struct {
	TransactionId int64     `json:"transactionId"`
	OrderId       string    `json:"order_id"`
	PayInfo       []PayInfo `json:"payInfo"`
	Packages      []Package `json:"packages"`
}

type MessageRoom struct {
	ID            primitive.ObjectID `bson:"_id"`
	ReservationID primitive.ObjectID `bson:"reservation_id"`
	GroupList     []GroupList        `bson:"group_list"`
	MessageLog    []MessageLog       `bson:"message_log"`
}

type GroupList struct {
	ID       primitive.ObjectID `bson:"_id"`
	ImageURL string             `bson:"image_url"`
	FullName string             `bson:"full_name"`
}
type MessageLog struct {
	SenderID  primitive.ObjectID `bson:"sender_id"`
	ReciverID primitive.ObjectID `bson:"reciver_id"`
	Message   Message            `bson:"message"`
}

type Message struct {
	Type      string    `json:"type" bson:"type"`
	Text      string    `json:"text" bson:"text"`
	ImageURL  string    `json:"img_url" bson:"img_url"`
	TimeStamp time.Time `json:"time_stamp" bson:"time_stamp"`
}

func (r Reservation) ToMapMyReservation() echo.Map {
	return echo.Map{
		"_id":             r.ID,
		"address":         r.Address_Full,
		"parking_name":    r.ParkingName,
		"price":           r.Price,
		"status":          r.Status,
		"profile_picture": r.ParkingPictureUrl,
		"date":            _converTime(r.TimeStamp),
	}
}

func _converTime(ts time.Time) string {
	newTime := ts.Add(7 * time.Hour)

	formattedTime := newTime.Format("02 ก.ค 06, 15:04")
	return formattedTime
}

// 4 Notification
type Notification struct {
	ID             primitive.ObjectID  `bson:"_id" json:"_id"`
	BroadcastType  string              `bson:"broadcast_type" json:"broadcast_type"`
	ReceiverID     *primitive.ObjectID `bson:"receiver_id"  json:"receiver_id"`
	Title          string              `bson:"title" json:"title"`
	Description    string              `bson:"description" json:"description"`
	CallbackMethod []CallbackMethod    `bson:"callback_method" json:"callback_method"`
	TimeStamp      time.Time           `bson:"time_stamp" json:"time_stamp"`
}

type CallbackMethod struct {
	Action      string `bson:"action" json:"action"`
	CallBackURL string `bson:"call_back_url" json:"call_back_url"`
}
