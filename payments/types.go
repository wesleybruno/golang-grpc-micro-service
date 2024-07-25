package main

import (
	"context"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
