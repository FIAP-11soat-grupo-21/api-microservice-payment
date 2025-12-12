package use_cases

import (
	"context"
	"fmt"

	"payment_microservice/internal/core/application/dtos"
	payment_gateway "payment_microservice/internal/core/application/gateways"
	"payment_microservice/internal/core/domain/exceptions"
)

type ConfirmPaymentUseCase struct {
	Gateway payment_gateway.PaymentGateway
}

func NewConfirmPaymentUseCase(gateway payment_gateway.PaymentGateway) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{
		Gateway: gateway,
	}
}

func (uc *ConfirmPaymentUseCase) Execute(ctx context.Context, dto dtos.WebhookEventDTO) error {
	payment, err := uc.Gateway.FindByOrderID(dto.OrderID)

	if err != nil {
		return &exceptions.PaymentNotFoundException{
			Message: fmt.Sprintf("Payment not found for OrderID: %s", dto.OrderID),
		}
	}

	if dto.Type == "payment" && dto.Action == "payment.updated" {
		payment.Status.SetPaid()

		uc.Gateway.Confirm(ctx, payment)
	}

	return nil
}
