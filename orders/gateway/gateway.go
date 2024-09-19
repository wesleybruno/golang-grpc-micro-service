package gateway

import (
	"context"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
)

type StockGateway interface {
	CheckIfItemIsInStock(ctx context.Context, customerID string, items []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
}
