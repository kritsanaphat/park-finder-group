package reserve

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/internal/httpclient"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReserveServices struct {
	ParkingAreaStorage  *storage.ParkingAreaStorage
	TransactionStorage  *storage.LogStorage
	ReserveStorage      *storage.LogStorage
	NotificationService ns.INotificationServices
	HttpClient          *httpclient.HTTPClient
	ReportStorage       *storage.LogStorage
}

type IReserveServices interface {
	ReserveParking(ctx context.Context, email, order_id string, req *models.ReserveRequest) (string, error, *primitive.ObjectID)
	ReserveParkingInAdvance(ctx context.Context, email, order_id string, req *models.ReserveRequest) (string, error, *primitive.ObjectID)
	MyReserve(ctx context.Context, email, parking_id, status string) ([]models.Reservation, error)
	MyReserveProvider(ctx context.Context, provider_id primitive.ObjectID, parking_id, status string) ([]models.Reservation, error)
	FindParkingDetail(ctx context.Context, parking_id string) *models.ParkingArea
	FindReserveDetail(ctx context.Context, reserve_id string) *models.Reservation
	FindReserveDetailByOrderID(ctx context.Context, order_id string) *models.Reservation
	CheckExistOrderID(ctx context.Context, order_id string) *models.Reservation
	StartReserveParking(ctx context.Context, customer_id primitive.ObjectID, order_id, module_code, parking_id, license_plate, action, parking_name string) (error, bool)
	CreateReview(ctx context.Context, customer_id, fn, ln, comment, parking_id, order_id string, review_score int) error
	CheckReserveInNextHour(ctx context.Context, parking_id primitive.ObjectID, hour_start, hour_end int, date_end string) (*[]models.Reservation, error)
	ConfirmReservationInAdvance(ctx context.Context, order_id string) error
	CheckOrderIDByParkingIDAndCustomerEmail(ctx context.Context, provider_id, customer_email string) *models.Reservation
	ExtendReserve(ctx context.Context, order_id, action, date_end string, hour_end int, min_end int) error
	UpdateReserveStatus(ctx context.Context, status, order_id string) error
	MyReservePaymentComplete(ctx context.Context, email string) ([]models.CutomserHistoryPoint, error)
	CaptureCarReserve(ctx context.Context, receiver_id primitive.ObjectID, module_code, parking_name string) (string, error)
	ReportReservation(ctx context.Context, customer_id, provider_id primitive.ObjectID, content, order_id string) error
	UpdateReserveStatusAndRemoveJob(ctx context.Context, order_id string) error
}

func NewReserveServices(
	db *mongo.Database,
	ht *httpclient.HTTPClient,
	ns ns.INotificationServices,
) IReserveServices {
	parking_area_storage := storage.NewParkingAreaStorage(db)
	transaction_storage := storage.NewTransactionStorage(db)
	reserve_storage := storage.NewReserveStorage(db)
	report_storage := storage.NewReportStorage(db)

	return ReserveServices{
		ParkingAreaStorage:  parking_area_storage,
		TransactionStorage:  transaction_storage,
		ReserveStorage:      reserve_storage,
		NotificationService: ns,
		HttpClient:          ht,
		ReportStorage:       report_storage,
	}
}
