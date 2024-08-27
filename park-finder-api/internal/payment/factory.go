package payment

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/internal/httpclient"
	ns "gitlab.com/parking-finder/parking-finder-api/internal/notification"
	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentServices struct {
	CustomerStorage     *storage.AccoutStorage
	TransactionStorage  *storage.LogStorage
	ReserveStorage      *storage.LogStorage
	NotificationService ns.INotificationServices
	ParkingAreaStorage  *storage.ParkingAreaStorage
	HttpClient          *httpclient.HTTPClient
}

type IPaymentServices interface {
	LinePayReserve(ctx context.Context, email, order_id, parking_id, action string, req *models.LineReserveRequest) (string, error)
	LinePayConfirm(ctx context.Context, transactionId, order_id, type_reserve string, is_extend bool) error
	LinePayCancel(ctx context.Context, transactionId, order_id string) error
	CashbackReserve(ctx context.Context, email, order_id, parking_id, type_reserve, action string, req *models.LineReserveRequest) error
	LinePayConfirmFine(ctx context.Context, transactionId, order_id string) error
}

func NewPaymentService(
	db *mongo.Database,
	ht *httpclient.HTTPClient,
	ns ns.INotificationServices,

) IPaymentServices {
	customer_storage := storage.NewCustomerAccoutStorage(db)
	transaction_storage := storage.NewTransactionStorage(db)
	parking_area_storage := storage.NewParkingAreaStorage(db)
	reserve_storage := storage.NewReserveStorage(db)

	return PaymentServices{
		CustomerStorage:     customer_storage,
		TransactionStorage:  transaction_storage,
		ReserveStorage:      reserve_storage,
		ParkingAreaStorage:  parking_area_storage,
		NotificationService: ns,
		HttpClient:          ht,
	}
}
