package processor

import (
	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
