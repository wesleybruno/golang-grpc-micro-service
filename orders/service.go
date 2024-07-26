package main

import (
	"context"
	"log"

	"github.com/wesleybruno/golang-grpc-micro-service/common"
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type orderService struct {
	store OrdersStore
}

func NewOrderService(store OrdersStore) *orderService {
	return &orderService{store}
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
	log.Print(mergedItems)

	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {

		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PrinceID: "",
			ID:       i.ID,
			Quantity: i.Quantity,
		})

	}

	return itemsWithPrice, nil
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
