package main

import (
	"log"
	"os"

	broker "github.com/MauricioMilano/worker_stock_app/client"
	"github.com/MauricioMilano/worker_stock_app/utils"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

var br broker.Broker

func main() {
	log.Println("Stock bot service starting ...")

	err := godotenv.Load()
	utils.ErrorPanic(err, "Error loading .env file")

	rmqHost := os.Getenv("RMQ_HOST")
	rmqUserName := os.Getenv("RMQ_USERNAME")
	rmqPassword := os.Getenv("RMQ_PASSWORD")
	rmqPort := os.Getenv("RMQ_PORT")
	dsn := "amqp://" + rmqUserName + ":" + rmqPassword + "@" + rmqHost + ":" + rmqPort + "/"

	conn, err := amqp.Dial(dsn)
	utils.ErrorPanic(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.ErrorPanic(err, "Failed to open a channel")
	defer ch.Close()

	br.SetQueue(ch)
	go br.ReadMessages()
	select {}
}
