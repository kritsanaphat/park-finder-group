package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	//"fmt"
	//"strconv"

	"os"

	"github.com/IBM/sarama"
	socketio "github.com/googollee/go-socket.io"
	"gitlab.comparking-finderpark-finder-socket/kafkago"
)

type WebhookService struct {
	server *socketio.Server
}

type IWebHookService interface {
	StartKafkaConsumer(poolSize int)
}

func NewWebhookService(server_in *socketio.Server) IWebHookService {
	return &WebhookService{
		server: server_in,
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
		fmt.Println("Message")
		messageLogs, ok := jsonData["MessageLog"].([]interface{})
		if !ok {
			return "Invalid MessageLog type", 400
		}
		var senderID string
		var receiverID string
		// Access the first message log
		if len(messageLogs) > 0 {
			firstMessageLog, ok := messageLogs[0].(map[string]interface{})
			if !ok {
				return "Invalid message log type", 400
			}

			// Access the "SenderID" and "ReciverID" fields
			senderID, ok = firstMessageLog["SenderID"].(string)
			if !ok {
				return "Invalid SenderID type", 400
			}

			receiverID, ok = firstMessageLog["ReciverID"].(string)
			if !ok {
				return "Invalid ReciverID type", 400
			}
		}
		fmt.Println("-------------Message-------------")
		fmt.Println("Broadcast To", receiverID)
		cusRepo.server.BroadcastTo(receiverID, "message", jsonData)
		fmt.Println("-------------Message-------------")
		fmt.Println("Broadcast To", senderID)
		cusRepo.server.BroadcastTo(senderID, "message", jsonData)

	}

	if keyMessage == "notification" {
		if jsonData["broadcast_type"] == "Personal" {
			reciver_id := jsonData["receiver_id"]

			reciver_id_string := reciver_id.(string)
			fmt.Println("-------------Notification-------------")
			fmt.Println("Broadcast Specific To", reciver_id_string)
			cusRepo.server.BroadcastTo(reciver_id_string, "notification", jsonData)
		} else {
			fmt.Println("Broadcast Group To", jsonData["broadcast_type"].(string))
			cusRepo.server.BroadcastTo(jsonData["broadcast_type"].(string), "notification", jsonData)
		}
	}

	return "ok", 200

}
