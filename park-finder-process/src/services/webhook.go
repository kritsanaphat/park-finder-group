package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	//"fmt"
	//"strconv"

	"os"

	"github.com/IBM/sarama"
	"gitlab.comparking-finderpark-finder-process/kafkago"
	"gitlab.comparking-finderpark-finder-process/models"
	"gitlab.comparking-finderpark-finder-process/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebhookService struct {
	KafkaProcess        kafkago.IProducer
	LogStorage          *storage.LogStorage
	NotificationStorage *storage.LogStorage

	CustomerAccoutStorage *storage.AccoutStorage
	ProviderAccoutStorage *storage.AccoutStorage
}

type IWebHookService interface {
	StartKafkaConsumer(poolSize int)
}

func NewWebhookService(
	db *mongo.Database,

) IWebHookService {
	message_log_storage := storage.NewMessageStorage(db)
	kafka_produce := kafkago.NewProducerProvider()
	notification_storage := storage.NewNotificationStorage(db)
	customer_account_storage := storage.NewCustomerAccoutStorage(db)
	provider_account_storage := storage.NewProviderAccoutStorage(db)

	return &WebhookService{
		KafkaProcess:          kafka_produce,
		LogStorage:            message_log_storage,
		NotificationStorage:   notification_storage,
		CustomerAccoutStorage: customer_account_storage,
		ProviderAccoutStorage: provider_account_storage,
	}
}

func (cusRepo WebhookService) StartKafkaConsumer(poolSize int) {
	kafkaBrokers := strings.Split(os.Getenv("CLOUDKARAFKA_BROKERS"), ",")
	fmt.Println("Kafka Brokers: ", kafkaBrokers)

	consumerSetting := kafkago.SettingKafkaConsumer()

	topic := consumerSetting.DefaultTopicEnum
	group := consumerSetting.GroupMainLogEnum

	fmt.Printf("Kafka Topic: %s | Kafka Group: %s \n", topic, group)

	kafkago.NewConsumer(
		kafkaBrokers,
		topic,
		kafkago.StrategyRoundRobinEnum,
		group,
		cusRepo.ConsumerEventHandler,
		poolSize,
	)
}

func (cusRepo WebhookService) ConsumerEventHandler(message *sarama.ConsumerMessage) {
	var keyMessage string
	if err := json.Unmarshal(message.Key, &keyMessage); err != nil {
		log.Printf("error marshalling message KEY: %v", err)
		return
	}

	var valueMessage map[string]interface{}
	if err := json.Unmarshal(message.Value, &valueMessage); err != nil {
		log.Printf("error marshalling message VALUE: %v", err)
		return
	}

	msg, status := cusRepo.ProcessingWebhook(keyMessage, valueMessage)
	if msg != "ok" || status != 200 {
		fmt.Println("Error :: Cannot process webhook message | ", valueMessage)
		return
	}

	return
}

func (cusRepo WebhookService) ProcessingWebhook(keyMessage string, jsonData map[string]interface{}) (string, int) {

	if keyMessage == "message" {
		fmt.Println("-------------- Received Message --------------")
		msg := cusRepo.mapToMessageRoom(jsonData)
		err := cusRepo.ProcessToSaveChatLog(*msg)

		if err == nil {
			byteValue, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}

			byteKey, err := json.Marshal("message")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}

			if err := cusRepo.KafkaProcess.ProduceMessage("socket", byteKey, byteValue); err != nil {
				fmt.Println("Error producing message:", err)
			}
		}

	}

	if keyMessage == "notification" {
		fmt.Println("-------------- Received Notification --------------")
		notification := cusRepo.mapToNotification(jsonData)
		err := cusRepo.ProcessToSaveNotification(*notification)

		if err == nil {
			byteValue, err := json.Marshal(jsonData)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}

			byteKey, err := json.Marshal("notification")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}
			if err := cusRepo.KafkaProcess.ProduceMessage("socket", byteKey, byteValue); err != nil {
				fmt.Println("Error producing message:", err)
			}
		}
	}

	return "ok", 200

}

func (cusRepo WebhookService) ProcessToSaveChatLog(msg models.MessageRoom) error {
	ctx := context.Background()
	msg_exist := cusRepo.LogStorage.FindMessgeExist(ctx, msg.ReservationID)
	if msg_exist != nil {
		fmt.Println("Push")
		err := cusRepo.LogStorage.PushMessageLog(ctx, msg_exist.ReservationID, msg.MessageLog[0])
		if err != nil {
			return err
		}

	} else {
		fmt.Println("Add")
		err := cusRepo.LogStorage.InsertMessageRoom(ctx, msg)
		if err != nil {
			return err
		}

	}
	return nil
}

func (cusRepo WebhookService) ProcessToSaveNotification(notification models.Notification) error {
	ctx := context.Background()

	err := cusRepo.NotificationStorage.InsertLog(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}

func (cusRepo WebhookService) mapToMessageRoom(data map[string]interface{}) *models.MessageRoom {
	var senderID primitive.ObjectID
	var receiverID primitive.ObjectID
	var reservationID primitive.ObjectID
	var message models.Message
	var err error

	if data["sender_id"] != nil {
		senderID, err = primitive.ObjectIDFromHex(data["sender_id"].(string))
		if err != nil {
			return nil
		}
	}

	if data["receiver_id"] != nil {
		receiverID, err = primitive.ObjectIDFromHex(data["receiver_id"].(string))
		if err != nil {
			return nil
		}
	}
	if data["reservation_id"] != nil {
		reservationID, err = primitive.ObjectIDFromHex(data["reservation_id"].(string))
		if err != nil {
			return nil
		}
	}

	bangkokLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil
	}
	if data["message"].(map[string]interface{})["Type"] != nil && data["message"].(map[string]interface{})["Text"] != nil {
		message = models.Message{
			Type:      data["message"].(map[string]interface{})["Type"].(string),
			Text:      data["message"].(map[string]interface{})["Text"].(string),
			TimeStamp: time.Now().In(bangkokLocation),
		}
	}
	messageLog := models.MessageLog{
		SenderID:  senderID,
		ReciverID: receiverID,
		Message:   message,
	}
	sender := cusRepo.CustomerAccoutStorage.FindAccountByID(senderID.Hex())
	if sender == nil {
		sender = cusRepo.ProviderAccoutStorage.FindAccountByID(senderID.Hex())
		if sender == nil {
			return nil

		}
	}
	var groupID []models.GroupList
	var messageRoom models.MessageRoom
	receiver := cusRepo.CustomerAccoutStorage.FindAccountByID(receiverID.Hex())
	if receiver == nil {
		receiver = cusRepo.ProviderAccoutStorage.FindAccountByID(receiverID.Hex())
		if receiver == nil {
			return nil

		}
	}

	groupID = []models.GroupList{
		{
			ID:       sender.ID,
			ImageURL: sender.ProfilePictureURL,
			FullName: sender.FirstName + " " + sender.LastName,
		},
		{
			ID:       receiver.ID,
			ImageURL: receiver.ProfilePictureURL,
			FullName: receiver.FirstName + " " + receiver.LastName,
		},
	}
	messageRoom = models.MessageRoom{
		ID:            primitive.NewObjectID(),
		ReservationID: reservationID,
		GroupList:     groupID,
		MessageLog:    []models.MessageLog{messageLog},
	}

	return &messageRoom
}

func (cusRepo WebhookService) mapToNotification(data map[string]interface{}) *models.Notification {
	var receiverID primitive.ObjectID
	var broadcast_type string
	var title string
	var description string
	var err error
	var ok bool

	if data["receiver_id"] != nil {
		receiverID, err = primitive.ObjectIDFromHex(data["receiver_id"].(string))
		if err != nil {
			fmt.Println("Error converting receiver_id:", err)
			return nil
		}
	}

	var callbacks []models.CallbackMethod
	var callBackURLs []interface{}
	if data["callback_method"] != nil {
		callBackURLs, ok = data["callback_method"].([]interface{})
		if !ok {
			fmt.Println("Invalid or missing callback_method")
			return nil
		}
	}

	for _, temp := range callBackURLs {
		tempData, ok := temp.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid callback_method element:", temp)
			return nil
		}

		callback := models.CallbackMethod{
			Action:      tempData["action"].(string),
			CallBackURL: tempData["call_back_url"].(string),
		}
		callbacks = append(callbacks, callback)
	}

	if data["broadcast_type"] != nil {
		broadcast_type = data["broadcast_type"].(string)
	}
	if data["title"] != nil {
		title = data["title"].(string)
	}
	if data["description"] != nil {
		description = data["description"].(string)
	}
	bangkokLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil
	}
	notification := models.Notification{
		ID:             primitive.NewObjectID(),
		BroadcastType:  broadcast_type,
		ReceiverID:     &receiverID,
		Title:          title,
		Description:    description,
		CallbackMethod: callbacks,
		TimeStamp:      time.Now().In(bangkokLocation),
	}

	return &notification
}
