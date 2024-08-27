package notification

import (
	"context"

	"gitlab.com/parking-finder/parking-finder-api/internal/storage"
	"gitlab.com/parking-finder/parking-finder-api/models"

	"gitlab.com/parking-finder/parking-finder-api/kafkago"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationServices struct {
	NotificationStorage *storage.LogStorage
	KafkaProcess        kafkago.IProducer
}

type INotificationServices interface {
	NotificationList(ctx context.Context, receiver_id, type_client string) ([]models.Notification, error)

	ConfirmReserveInAdvanceNotification(ctx context.Context, receiver_id primitive.ObjectID, customer_email string, reserve models.Reservation) error
	ProviderConfirmReserveInAdvanceNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation, full_address string) error
	ProviderCancelReserveInAdvanceNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation, full_address string) error
	ParkingAreaStatusUpdateNotification(ctx context.Context, receiver_id primitive.ObjectID, status, parking_name string) error
	BeforeTimeOutReserveNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error
	TimeOutReserveNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error
	AfterTimeOutReserveNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation, fine int) error
	LeaveTimeOutReserveNotification(ctx context.Context, receiver_id, parking_id primitive.ObjectID, order_id, parking_name string) error
	ReservationCancelNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error
	VertifyCustomerCarNotification(ctx context.Context, receiver_id primitive.ObjectID, license_plate, module_code, pic_url string, reserve *models.Reservation) error
	ReportParkingAreaNotification(ctx context.Context, receiver_id primitive.ObjectID, reserve *models.Reservation) error
	AddRewardNotification(ctx context.Context, tile, description, id, preview_url string) error
}

func NewNotificationServices(
	db *mongo.Database,
) INotificationServices {
	notification_storage := storage.NewNotificationStorage(db)
	kafka_produce := kafkago.NewProducerProvider()

	return NotificationServices{
		NotificationStorage: notification_storage,
		KafkaProcess:        kafka_produce,
	}
}
