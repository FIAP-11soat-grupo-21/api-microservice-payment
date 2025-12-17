package use_cases

import (
	"context"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
)

type FindPaymentByOrderIDUseCase struct {
	repository ports.IPaymentRepository
}

func NewFindPaymentByOrderIDUseCase(repository ports.IPaymentRepository) *FindPaymentByOrderIDUseCase {
	return &FindPaymentByOrderIDUseCase{
		repository: repository,
	}
}

func (uc *FindPaymentByOrderIDUseCase) Execute(ctx context.Context, orderID string) (entities.Payment, error) {
	payment, err := uc.repository.FindByOrderID(ctx, orderID)

	if err != nil {
		return entities.Payment{}, &exceptions.PaymentNotFoundException{}
	}

	return payment, nil
}
