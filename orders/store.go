package main

import (
	"context"
	"errors"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

var orders = make([]*pb.Order, 0)

type orderStore struct {
}

func NewOrderStore() *orderStore {
	return &orderStore{}
}

func (s *orderStore) Create(ctx context.Context, o *pb.CreateOrderRequest, i []*pb.Item) (*pb.Order, error) {

	order := &pb.Order{
		ID:          "123",
		CustomerID:  o.CustomerID,
		Status:      "pending",
		PaymentLink: "",
		Items:       i,
	}

	orders = append(orders, order)

	return order, nil
}

func (s *orderStore) GetOrder(ctx context.Context, orderId, customerId string) (*pb.Order, error) {

	for _, o := range orders {

		if o.ID == orderId && o.CustomerID == customerId {
			return o, nil
		}

	}

	return nil, errors.New("order not found")
}

func (s *orderStore) UpdateOrder(ctx context.Context, orderId string, newOrder *pb.Order) error {

	for i, order := range orders {

		if order.ID == orderId {
			orders[i].Status = order.Status
			orders[i].PaymentLink = order.PaymentLink

			return nil

		}

	}

	return nil

}
