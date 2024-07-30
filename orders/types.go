package main

import (
	"context"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	GetOrder(ctx context.Context, orderId, customerId string) (*pb.Order, error)
	UpdateOrder(ctx context.Context, orderId string, order *pb.Order) error
}
