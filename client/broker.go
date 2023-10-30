package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MauricioMilano/worker_stock_app/services"
	"github.com/MauricioMilano/worker_stock_app/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type StockRequest struct {
	ChatRoomName string `json:"chatRoomName"`
	ChatRoomId   uint   `json:"chatRoomId"`
	ChatMessage  string `json:"chatMessage"`
}

type StockReponse struct {
	RoomId  uint   `json:"RoomId"`
	Message string `json:"Message"`
}

type Broker struct {
	ReceiverQueue  amqp.Queue
	PublisherQueue amqp.Queue
	Channel        *amqp.Channel
}

func (b *Broker) SetQueue(ch *amqp.Channel) {
	receiverQueue := os.Getenv("RECEIVER_QUEUE")
	publisherQueue := os.Getenv("PUBLISHER_QUEUE")

	q1, err := ch.QueueDeclare(
		receiverQueue, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	utils.ErrorPanic(err, "Failed to declare"+receiverQueue+" queue")

	q2, err := ch.QueueDeclare(
		publisherQueue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	utils.ErrorPanic(err, "Failed to declare "+publisherQueue+" queue")

	b.ReceiverQueue = q1
	b.PublisherQueue = q2
	b.Channel = ch
}

func (b *Broker) PublishMessage(sr StockReponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(sr)
	if err != nil {
		log.Printf("Response structure error %s ", err)
	}

	err = b.Channel.PublishWithContext(ctx,
		"",
		b.ReceiverQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	utils.ErrorPanic(err, "Failed to publish a message")
	log.Printf("Sent %s\n", body)
}

func (b *Broker) ReadMessages() {
	msgs, err := b.Channel.Consume(
		b.PublisherQueue.Name, // queue
		"",                    // consumer
		true,                  // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	utils.ErrorPanic(err, "Failed to register a consumer")

	rsvdMsgs := make(chan StockRequest)
	// forever := make(chan bool)

	// go func() {
	// 	for d := range msgs {
	// 		fmt.Printf("Mensagem recebida: %s\n", d.Body)

	// 	}
	// }()

	// fmt.Println("Pressione Ctrl+C para sair")
	// <-forever
	go messageTransformer(msgs, rsvdMsgs)
	go processRequest(rsvdMsgs, b)
	log.Printf(" Waiting for messages. To exit press CTRL+C")
}

func messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan StockRequest) {
	var sr StockRequest
	for d := range entries {
		log.Println("d.Body", string(d.Body))
		err := json.Unmarshal([]byte(d.Body), &sr)
		if err != nil {
			log.Printf("Error on received request : %s ", err)
			continue
		}
		log.Println("Received a request")
		receivedMessages <- sr
	}
}

func processRequest(stocksRequests <-chan StockRequest, b *Broker) {

	for request := range stocksRequests {
		log.Println("processing stock request for ", request.ChatRoomId)
		stock_name := strings.Replace(request.ChatMessage, "/stock=", "", 1)
		sr := StockReponse{
			RoomId:  request.ChatRoomId,
			Message: fmt.Sprintf("Processing: %s", stock_name),
		}
		go b.PublishMessage(sr)
		msg := services.EvalStock(stock_name)
		sr2 := StockReponse{
			RoomId:  request.ChatRoomId,
			Message: msg,
		}
		go b.PublishMessage(sr2)
		log.Println("processed", request.ChatMessage)
	}
}

func consume(receivedMessages <-chan amqp.Delivery, b *Broker) {

	forever := make(chan bool)

	go func() {
		for d := range receivedMessages {
			fmt.Printf("Mensagem recebida: %s\n", d.Body)

		}
	}()

	fmt.Println("Pressione Ctrl+C para sair")
	<-forever
}
