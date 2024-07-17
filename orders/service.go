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

func (s *orderService) CreateOrder(context.Context) error {
	return nil
}

func (s *orderService) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) error {

	if len(p.Items) == 0 {
		return common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Print(mergedItems)

	return nil
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
