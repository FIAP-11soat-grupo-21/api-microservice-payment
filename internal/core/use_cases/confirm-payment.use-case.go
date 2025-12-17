package use_cases

import (
	"context"
	"fmt"

	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"
)

type ConfirmPaymentUseCase struct {
	repository ports.IPaymentRepository
}

func NewConfirmPaymentUseCase(repository ports.IPaymentRepository) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{
		repository: repository,
	}
}

func (uc *ConfirmPaymentUseCase) Execute(ctx context.Context, dto dto.WebhookEventDTO) error {
	payment, err := uc.repository.FindByOrderID(ctx, dto.OrderID)

	if err != nil {
		return &exceptions.PaymentNotFoundException{
			Message: fmt.Sprintf("Payment not found for OrderID: %s", dto.OrderID),
		}
	}

	if dto.Type == "payment" && dto.Action == "payment.updated" {
		payment.Status.SetPaid()

		uc.repository.Update(ctx, payment)
	}

	return nil
}
