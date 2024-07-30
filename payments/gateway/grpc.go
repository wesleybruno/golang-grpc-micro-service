package gateway

import (
	"context"
	"log"

	pb "github.com/wesleybruno/golang-grpc-micro-service/common/api"
	discovery "github.com/wesleybruno/golang-grpc-micro-service/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) UpdateOrderAfterPaymentLink(ctx context.Context, orderId, paymentLink string) error {

	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		log.Fatalf("Error to start conn %s", err)
	}

	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)

	_, err = c.UpdateOrder(ctx, &pb.Order{
		ID:          orderId,
		Status:      "waiting payment",
		PaymentLink: paymentLink,
	})

	if err != nil {
		return err
	}

	return nil

}
