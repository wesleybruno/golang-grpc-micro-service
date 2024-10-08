package gateway

import (
	"context"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrderById(ctx context.Context, customerId, orderId string) (*pb.Order, error)
}
