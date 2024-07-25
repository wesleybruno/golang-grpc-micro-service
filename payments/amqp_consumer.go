package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"github.com/wesleybruno/golang-grpc-micro-service/common/broker"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func (c consumer) Listen(ch *amqp.Channel) {

	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				log.Fatal(err)
				continue
			}

			log.Printf("Message received: %v", o)

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("Failed to create paymanet: %v", err)
				continue
			}

			log.Printf("PaymentLink Created: %s", paymentLink)
		}
	}()

	<-forever

}
