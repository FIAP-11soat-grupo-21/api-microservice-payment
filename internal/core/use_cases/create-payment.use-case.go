package use_cases

import (
	"fmt"
	identity_manager "payment_microservice/internal/common/pkg/identity"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"
)

type CreatePaymentUseCase struct {
	repository ports.IPaymentRepository
	gateway    ports.IPaymentGateway
}

func NewCreatePaymentUseCase(repository ports.IPaymentRepository, gateway ports.IPaymentGateway) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		repository: repository,
		gateway:    gateway,
	}
}

func (uc *CreatePaymentUseCase) Execute(paymentDTO dto.CreatePaymentDTO) (entities.Payment, error) {
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

	if payment.Method.IsPix() {
		pixBillingResult, err := uc.gateway.CreatePIXBilling(dto.CreatePIXBillingDTO{
			Ctx:        paymentDTO.Ctx,
			ExternalID: payment.ID,
			Amount:     payment.Amount.Value(),
		})

		if err != nil {
			return entities.Payment{}, &exceptions.InvalidPaymentDataException{
				Message: fmt.Sprintf("Failed to create PIX billing: %v", err),
			}
		}

		err = payment.SetQrCode(pixBillingResult.QRData)

		if err != nil {
			return entities.Payment{}, &exceptions.InvalidPaymentDataException{}
		}

		err = uc.repository.Insert(paymentDTO.Ctx, *payment)

		if err != nil {
			return entities.Payment{}, &exceptions.InvalidPaymentDataException{
				Message: "Failed to insert payment on database",
			}
		}

		return *payment, nil
	}

	return entities.Payment{}, &exceptions.InvalidPaymentDataException{
		Message: "Invalid payment method - only PIX is supported",
	}
}
