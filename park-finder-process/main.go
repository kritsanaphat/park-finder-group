package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.comparking-finderpark-finder-process/pkg/connector"
	sv "gitlab.comparking-finderpark-finder-process/src/services"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webhook := sv.NewWebhookService(db)

	webhook.StartKafkaConsumer(200)
}
