package repositories

import (
	"payment_microservice/internal/core/domain/entities"
)

func toDomain(paymentModel PaymentModel) (entities.Payment, error) {
	payment, err := entities.NewPayment(
		paymentModel.ID,
		paymentModel.OrderID,
		paymentModel.Amount,
		paymentModel.Status,
		paymentModel.PaymentMethod,
		paymentModel.QRCodeURL,
		paymentModel.PaidAt,
		paymentModel.CreatedAt,
	)

	if err != nil {
		return entities.Payment{}, err
	}

	return *payment, nil
}

func toPersistence(payment entities.Payment) PaymentModel {
	var qrCodeURL *string

	if payment.QRCode != nil {
		url := payment.QRCode.Value()
		qrCodeURL = &url
	}

	return PaymentModel{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount.Value(),
		Status:        payment.Status.Value(),
		PaymentMethod: payment.Method.Value(),
		QRCodeURL:     qrCodeURL,
		PaidAt:        payment.PaidAt,
		CreatedAt:     payment.CreatedAt,
	}
}
