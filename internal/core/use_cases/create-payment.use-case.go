package use_cases

import (
	"fmt"
	identity_manager "payment_microservice/internal/common/pkg/identity"
	"payment_microservice/internal/core/application/dtos"
	"payment_microservice/internal/core/application/gateways"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
)

type CreatePaymentUseCase struct {
	gateway gateways.PaymentGateway
}

func NewCreatePaymentUseCase(gateway gateways.PaymentGateway) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		gateway: gateway,
	}
}

func (uc *CreatePaymentUseCase) Execute(paymentDTO dtos.CreatePaymentDTO) (entities.Payment, error) {
	var qrData *string

	payment, err := entities.NewPaymentDefault(
		identity_manager.NewUUIDV4(),
		paymentDTO.OrderID,
		paymentDTO.Amount,
		paymentDTO.Method,
		qrData,
	)

	if err != nil {
		return entities.Payment{}, err
	}

	err = uc.gateway.Insert(paymentDTO.Ctx, *payment)

	if err != nil {
		return entities.Payment{}, &exceptions.InvalidPaymentDataException{
			Message: "Failed to insert payment",
		}
	}

	if payment.Method.IsPix() {
		pixBillingResult, err := uc.gateway.CreatePIXBilling(paymentDTO.Ctx, *payment)

		if err != nil {
			return entities.Payment{}, &exceptions.InvalidPaymentDataException{
				Message: fmt.Sprintf("Failed to create PIX billing: %v", err),
			}
		}

		payment.SetQrCode(pixBillingResult.QRData)

		err = uc.gateway.Update(paymentDTO.Ctx, *payment)

		if err != nil {
			return entities.Payment{}, &exceptions.InvalidPaymentDataException{
				Message: "Failed to update payment with QR data",
			}
		}

		return *payment, nil
	}

	return entities.Payment{}, &exceptions.InvalidPaymentDataException{
		Message: "Invalid payment method - only PIX is supported",
	}
}
