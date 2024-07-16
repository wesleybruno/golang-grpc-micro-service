package main

import "context"

type orderService struct {
	store OrdersStore
}

func NewOrderService(store OrdersStore) *orderService {
	return &orderService{store}
}

func (s *orderService) CreateOrder(context.Context) error {
	return nil
}
