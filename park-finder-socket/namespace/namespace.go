package namespace

import (
	"encoding/json"
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"gitlab.comparking-finderpark-finder-socket/kafkago"
	storage "gitlab.comparking-finderpark-finder-socket/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type SocketService struct {
	KafkaProcess          kafkago.IProducer
	CustomerAccoutStorage *storage.AccoutStorage
	ProviderAccoutStorage *storage.AccoutStorage
}

type ISocketService interface {
	StartSocket() *socketio.Server
}

func NewSocketService(
	db *mongo.Database,

) ISocketService {
	kafka_produce := kafkago.NewProducerProvider()
	customer_account_storage := storage.NewCustomerAccoutStorage(db)
	provider_account_storage := storage.NewProviderAccoutStorage(db)

	return &SocketService{
		KafkaProcess:          kafka_produce,
		CustomerAccoutStorage: customer_account_storage,
		ProviderAccoutStorage: provider_account_storage,
	}
}

func (s SocketService) StartSocket() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("error establishing new socketio server")
	}

	server.On("connection", func(so socketio.Socket) {
		fmt.Println("####################### User Connected #######################")
		userID := so.Request().URL.Query().Get("user_id")

		err := s.CustomerAccoutStorage.FindAccountByID(userID)
		if err != nil {
			fmt.Println(err)
			err := s.ProviderAccoutStorage.FindAccountByProviderID(userID)
			if err != nil {
				fmt.Println("User Id invalid", err)
				so.Disconnect()

			} else {
				fmt.Println("Connect Successfuly")
				fmt.Println("User ID :", userID)
				so.Join("provider")
				so.Join(userID)

			}
		} else {
			fmt.Println("Connect Successfuly")
			fmt.Println("User ID :", userID)
			so.Join("customer")
			so.Join(userID)

		}
	})

	server.On("message", func(so socketio.Socket, msg map[string]interface{}) {
		fmt.Println("####################### On Message #######################")

		byteValue, err := json.Marshal(msg)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			return
		}

		byteKey, err := json.Marshal("message")
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			return
		}

		if err := s.KafkaProcess.ProduceMessage("process", byteKey, byteValue); err != nil {
			log.Println("Error producing message:", err)
			return
		}
	})
	fmt.Println("Server running on localhost: 4700")
	return server
}
