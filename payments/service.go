package main

import (
	"context"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type service struct{}

func NewService() *service {
	return &service{}

}

func (s service) CreatePayment(context.Context, *pb.Order) (string, error) {
	return "", nil

}
