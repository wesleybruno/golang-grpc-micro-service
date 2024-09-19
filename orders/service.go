package main

import (
	"context"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"github.com/wesleybruno/golang-grpc-micro-service/orders/gateway"
)

type orderService struct {
	store   OrdersStore
	gateway gateway.StockGateway
}

func NewOrderService(store OrdersStore, gateway gateway.StockGateway) *orderService {
	return &orderService{store, gateway}
}

func (s *orderService) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	o, err := s.store.GetOrder(ctx, p.OrderID, p.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *orderService) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, i []*pb.Item) (*pb.Order, error) {

	o, err := s.store.Create(ctx, p, i)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *orderService) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {

	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)

	inStock, items, err := s.gateway.CheckIfItemIsInStock(ctx, p.CustomerID, mergedItems)
	if err != nil {
		return nil, err
	}
	if !inStock {
		return items, common.ErrNoStock
	}

	return items, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {

	err := s.store.UpdateOrder(ctx, o.ID, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
