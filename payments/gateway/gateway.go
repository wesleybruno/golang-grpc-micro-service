package gateway

import (
	"context"
)

type OrdersGateway interface {
	UpdateOrderAfterPaymentLink(ctx context.Context, orderId, paymentLink string) error
}
