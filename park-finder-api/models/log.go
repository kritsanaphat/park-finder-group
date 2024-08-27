package models

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 1 Reservation
type Reservation struct {
	ID                primitive.ObjectID `bson:"_id"`
	ParkingID         primitive.ObjectID `bson:"parking_id"`
	ProviderID        primitive.ObjectID `bson:"provider_id"`
	CarID             primitive.ObjectID `bson:"car_id"`
	TransactionID     int                `bson:"transaction_id"`
	ParkingName       string             `bson:"parking_name"`
	CustomerEmail     string             `bson:"customer_email"`
	Price             int                `bson:"price"`
	Status            string             `bson:"status"`
	OrderID           string             `bson:"order_id"`
	DateStart         string             `bson:"date_start"`
	DateEnd           string             `bson:"date_end"`
	Type              string             `bson:"type"`
	HourStart         int                `bson:"hour_start"`
	MinStart          int                `bson:"min_start"`
	HourEnd           int                `bson:"hour_end"`
	MinEnd            int                `bson:"min_end"`
	PaymentChanel     string             `bson:"payment_chanel"`
	ParkingPictureUrl string             `bson:"parking_picture_url"`
	Address_Full      string             `bson:"address_full"`
	ModuleCode        string             `bson:"module_code"`
	IsExtend          bool               `bson:"is_extend"`
	TimeStamp         time.Time          `bson:"time_stamp"`
}

type UpdateStatusReservation struct {
	Status string `bson:"status"`
}

type InfoReserve struct {
	PaymentUrl         PaymentUrl `json:"paymentUrl" bson:"paymentUrl"`
	TransactionId      int64      `json:"transactionId" bson:"transactionId"`
	PaymentAccessToken string     `json:"paymentAccessToken" bson:"paymentAccessToken"`
}

type DetectStatus struct {
	RefID            string    `bson:"ref_id"`
	Entered          bool      `bson:"entered"`
	EnteredTimeStamp time.Time `bson:"entered_time_stamp"`
	Exited           bool      `bson:"exited"`
	ExitedTimeStamp  time.Time `bson:"exited_time_stamp"`
}

func (r Reservation) ToMapMyReservation() echo.Map {
	return echo.Map{
		"_id":                 r.ID,
		"address":             r.Address_Full,
		"parking_name":        r.ParkingName,
		"date_start":          r.DateStart,
		"date_end":            r.DateEnd,
		"hour_start":          r.HourStart,
		"min_start":           r.MinStart,
		"hour_end":            r.HourEnd,
		"min_end":             r.MinEnd,
		"price":               r.Price,
		"status":              r.Status,
		"parking_picture_url": r.ParkingPictureUrl,
		"date":                _converTime(r.TimeStamp),
	}
}

type Reservations []Reservation

func (r Reservations) Len() int {

	return len(r)
}
func (r Reservations) Less(i, j int) bool {
	dateStartI, _ := time.Parse("2006-01-02", r[i].DateStart)
	dateStartJ, _ := time.Parse("2006-01-02", r[j].DateStart)
	if dateStartI.Before(dateStartJ) {
		return true
	} else if dateStartI.After(dateStartJ) {
		return false
	}

	if r[i].HourStart < r[j].HourStart {
		return true
	} else if r[i].HourStart != r[j].HourStart {
		return false
	}

	if r[i].MinStart < r[j].MinStart {
		return true
	}
	return false
}

func (r Reservations) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// 2 Transaction
type TransactionLinePay struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	ParkingID     primitive.ObjectID `json:"parking_id" bson:"parking_id"`
	CustomerEmail string             `json:"customer_email" bson:"customer_email"`
	OrderId       string             `json:"order_Id" bson:"order_id"`
	Packages      []Package          `json:"packages" bson:"package"`
	Info          InfoReserve        `json:"info" bson:"info"`
	Status        string             `json:"status" bson:"status"`
	TimeStamp     time.Time          `json:"time_stamp" bson:"time_stamp"`
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

// 3 Message
type MessageRoom struct {
	ID            primitive.ObjectID `bson:"_id"`
	ReservationID primitive.ObjectID `bson:"reservation_id"`
	GroupList     []GroupList        `bson:"group_list"`
	MessageLog    []MessageLog       `bson:"message_log"`
}

func (r MessageRoom) ToMapMessageRoom() echo.Map {
	return echo.Map{
		"_id":            r.ID,
		"reservation_id": r.ReservationID,
		"group_list":     r.GroupList,
	}
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

// 4 Notification
type Notification struct {
	ID             primitive.ObjectID  `bson:"_id" json:"_id"`
	BroadcastType  string              `bson:"broadcast_type" json:"broadcast_type"`
	ReceiverID     *primitive.ObjectID `bson:"receiver_id"  json:"receiver_id"`
	Title          string              `bson:"title" json:"title"`
	Description    string              `bson:"description" json:"description"`
	CallbackMethod []CallbackMethod    `bson:"callback_method" json:"callback_method"`
	TimeStamp      time.Time           `json:"time_stamp" bson:"time_stamp"`
}

type CallbackMethod struct {
	Action      string `bson:"action" json:"action"`
	CallBackURL string `bson:"call_back_url" json:"call_back_url"`
}

func _converTime(ts time.Time) string {
	newTime := ts.Add(7 * time.Hour)

	formattedTime := newTime.Format("02 ก.ค 06, 15:04")
	return formattedTime
}

// 5 Report
type Report struct {
	ID         primitive.ObjectID `bson:"_id"`
	CustomerID primitive.ObjectID `bson:"customer_id"`
	ProviderID primitive.ObjectID `bson:"provider_id"`
	OrderID    string             `bson:"order_id"`
	Content    string             `bson:"content"`
	TimeStamp  time.Time          `bson:"time_stamp"`
}

// 6 Receipt
type Receipt struct {
	ID              primitive.ObjectID `bson:"_id"`
	ProviderID      primitive.ObjectID `bson:"provider_id"`
	ReceiptImageUrl string             `json:"receipt_image_url"`
	Price           int                `json:"price"`
	Month           string             `json:"month"`
	Year            string             `json:"year"`
	TimeStamp       time.Time          `bson:"time_stamp"`
}
