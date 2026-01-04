package use_cases

import (
	"context"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
)

type RollbackPaymentUseCase struct {
	repository ports.IPaymentRepository
}

func NewRollbackPaymentUseCase(repository ports.IPaymentRepository) *RollbackPaymentUseCase {
	return &RollbackPaymentUseCase{
		repository: repository,
	}
}

func (uc *RollbackPaymentUseCase) Execute(ctx context.Context, orderID string) error {
	payment, err := uc.repository.FindByOrderID(ctx, orderID)

	if err != nil {
		return err
	}

	if payment.IsEmpty() {
		return new(exceptions.PaymentNotFoundException)
	}

	err = uc.repository.Delete(ctx, payment.ID)

	if err != nil {
		return err
	}

	return nil
}
