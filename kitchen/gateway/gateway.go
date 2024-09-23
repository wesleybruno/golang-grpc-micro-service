package gateway

import (
	"context"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
