package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	s "gitlab.comparking-finderpark-finder-socket/namespace"
	sv "gitlab.comparking-finderpark-finder-socket/src/services"

	"gitlab.comparking-finderpark-finder-socket/pkg/connector"
)

func main() {
	wd, _ := os.Getwd()
	if err := godotenv.Load(wd + "/.env"); err != nil {
		fmt.Println("Error loading .env file.")
	}

	mongodb := connector.NewMongoDBClient(os.Getenv("MONGODB_URI"), 100)
	if err := mongodb.Ping(); err != nil {
		log.Fatalln("Error connecting to MongoDB cluster")
	}

	// Select DB
	db := mongodb.SelectDB(os.Getenv("DATABASE_NAME"))
	// Create a WaitGroup to wait for both goroutines to finish
	var wg sync.WaitGroup
	socketService := s.NewSocketService(db)
	server := socketService.StartSocket()

	wg.Add(1)
	go func() {
		defer wg.Done()

		fs := http.FileServer(http.Dir("static"))
		http.Handle("/", fs)
		http.Handle("/socket.io/", server)

		log.Fatal(http.ListenAndServe(":4700", nil))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		webhook := sv.NewWebhookService(server)
		webhook.StartKafkaConsumer(200)
	}()

	wg.Wait()
}
