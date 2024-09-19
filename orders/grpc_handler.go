package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"github.com/wesleybruno/golang-grpc-micro-service/common/broker"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type gRpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	ch      *amqp.Channel
}

func NewGrRpcHandler(grpcServer *grpc.Server, service OrdersService, ch *amqp.Channel) {

	handler := &gRpcHandler{
		service: service,
		ch:      ch,
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)

}

func (h gRpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received %v", p)

	q, err := h.ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", q.Name))
	defer messageSpan.End()

	items, err := h.service.ValidateOrder(amqpContext, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(amqpContext, p, items)
	if err != nil {
		return nil, err
	}

	marshelledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	headers := broker.InjectAMQPHeaders(amqpContext)

	h.ch.PublishWithContext(amqpContext, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshelledOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	return o, nil
}

func (h gRpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	log.Printf("New order requested %v", p)

	o, err := h.service.GetOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	return o, nil

}

func (h gRpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	log.Printf("Order requested update %v", p)

	_, err := h.service.UpdateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil

}
