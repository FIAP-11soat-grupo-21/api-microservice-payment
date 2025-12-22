package use_cases

import (
	"context"
	"fmt"

	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"
)

type ConfirmPaymentUseCase struct {
	repository          ports.IPaymentRepository
	kitchenOrderService ports.IKitchenOrderService
}

func NewConfirmPaymentUseCase(repository ports.IPaymentRepository, kitchenOrderService ports.IKitchenOrderService) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{
		repository:          repository,
		kitchenOrderService: kitchenOrderService,
	}
}

func (uc *ConfirmPaymentUseCase) Execute(ctx context.Context, eventDTO dto.WebhookEventDTO) error {
	payment, err := uc.repository.FindByOrderID(ctx, eventDTO.OrderID)

	if err != nil {
		return &exceptions.PaymentNotFoundException{
			Message: fmt.Sprintf("Payment not found for OrderID: %s", eventDTO.OrderID),
		}
	}

	if eventDTO.Type == "payment" && eventDTO.Action == "payment.updated" {
		payment.Status.SetPaid()

		uc.repository.Update(ctx, payment)

		kitchenOrderDTO := dto.CreateKitchenOrderDTO{
			OrderID: eventDTO.OrderID,
		}
		uc.kitchenOrderService.Create(ctx, kitchenOrderDTO)
	}

	return nil
}
