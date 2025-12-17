package ports

import (
	"context"
	"payment_microservice/internal/core/domain/entities"
)

type IPaymentRepository interface {
	Insert(ctx context.Context, payment entities.Payment) error
	FindByOrderID(ctx context.Context, orderID string) (entities.Payment, error)
	Update(ctx context.Context, payment entities.Payment) error
	Delete(ctx context.Context, paymentID string) error
}
