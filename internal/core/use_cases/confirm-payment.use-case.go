package use_cases

import (
	"context"
	"fmt"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"
)

type ConfirmPaymentUseCase struct {
	repository       ports.IPaymentRepository
	messagePublisher ports.IQueuePublisher
}

func NewConfirmPaymentUseCase(repository ports.IPaymentRepository, messagePublisher ports.IQueuePublisher) *ConfirmPaymentUseCase {
	return &ConfirmPaymentUseCase{
		repository:       repository,
		messagePublisher: messagePublisher,
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
		payment.MarkAsPaid()

		if err := uc.repository.Update(ctx, payment); err != nil {
			return err
		}

		topic := env.GetConfig().AWS.SNS.Topics.PaymentProcessed

		message := dto.PaymentProcessedEventDTO{
			OrderID: payment.OrderID,
			Status:  constants.ORDER_STATUS_CONFIRMED,
		}

		messageJSON, err := message.ToJSON()

		if err != nil {
			return err
		}

		if err := uc.messagePublisher.PublishOnTopic(ctx, topic, messageJSON); err != nil {
			return err
		}
	}

	return nil
}
