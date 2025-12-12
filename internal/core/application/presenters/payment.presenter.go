package presenters

import (
	"payment_microservice/internal/core/application/dtos"
	"payment_microservice/internal/core/domain/entities"
)

func ToResponse(payment entities.Payment) dtos.PaymentResultDTO {
	qrCode := payment.QRCode.Value()

	return dtos.PaymentResultDTO{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount.Value(),
		Status:          payment.Status.Value(),
		Method:          payment.Method.Value(),
		TransactionCode: payment.TransactionCode,
		QRCode:          &qrCode,
		CreatedAt:       payment.CreatedAt,
		PaidAt:          payment.PaidAt,
	}
}
