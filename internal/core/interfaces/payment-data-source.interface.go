package interfaces

import (
	"context"

	"payment_microservice/internal/core/daos"
)

type IPaymentDataSource interface {
	Insert(ctx context.Context, payment daos.PaymentDAO) error
	Update(ctx context.Context, payment daos.PaymentDAO) error
	FindByOrderID(orderId string) (daos.PaymentDAO, error)
}
