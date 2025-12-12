package use_cases

import (
	"payment_microservice/internal/core/application/gateways"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
)

type FindPaymentByOrderIDUseCase struct {
	Gateway gateways.PaymentGateway
}

func NewFindPaymentByOrderIDUseCase(gateway gateways.PaymentGateway) *FindPaymentByOrderIDUseCase {
	return &FindPaymentByOrderIDUseCase{
		Gateway: gateway,
	}
}

func (uc *FindPaymentByOrderIDUseCase) Execute(orderID string) (entities.Payment, error) {
	payment, err := uc.Gateway.FindByOrderID(orderID)

	if err != nil {
		return entities.Payment{}, &exceptions.PaymentNotFoundException{}
	}

	return payment, nil
}
