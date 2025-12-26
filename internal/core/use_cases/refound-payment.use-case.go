package use_cases

import (
	"context"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
)

type RefoundPaymentUseCase struct {
	repository ports.IPaymentRepository
	gateway    ports.IPaymentGateway
}

func NewRefoundPaymentUseCase(repository ports.IPaymentRepository) *RefoundPaymentUseCase {
	return &RefoundPaymentUseCase{
		repository: repository,
	}
}

func (uc *RefoundPaymentUseCase) Execute(ctx context.Context, orderID string) error {
	payment, err := uc.repository.FindByOrderID(ctx, orderID)

	if err != nil {
		return err
	}

	if payment.IsEmpty() {
		return new(exceptions.PaymentNotFoundException)
	}

	payment.MarkAsRefunded()

	err = uc.repository.Update(ctx, payment)

	if err != nil {
		return err
	}

	return nil
}
