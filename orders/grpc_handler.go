package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"github.com/wesleybruno/golang-grpc-micro-service/common/broker"
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

	items, err := h.service.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(ctx, p, items)
	if err != nil {
		return nil, err
	}

	marshelledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	q, err := h.ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	h.ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshelledOrder,
		DeliveryMode: amqp.Persistent,
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
