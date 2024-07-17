package main

import (
	"context"
	"log"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	"google.golang.org/grpc"
)

type gRpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service orderService
}

func NewGrRpcHandler(grpcServer *grpc.Server, service orderService) {

	handler := &gRpcHandler{
		service: service,
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)

}

func (h gRpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received %v", p)

	if err := h.service.ValidateOrder(ctx, p); err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID: "10",
	}

	return o, nil
}
