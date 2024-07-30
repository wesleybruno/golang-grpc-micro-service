package main

import (
	"context"

	"testing"

	"github.com/wesleybruno/golang-grpc-micro-service/common/api"
	inmemRegistry "github.com/wesleybruno/golang-grpc-micro-service/common/discovery/inmem"
	"github.com/wesleybruno/golang-grpc-micro-service/payment/gateway"
	"github.com/wesleybruno/golang-grpc-micro-service/payment/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()
	gateway := gateway.NewGRPCGateway(registry)
	svc := NewService(processor, gateway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
