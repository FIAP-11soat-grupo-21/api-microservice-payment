package mocks

import (
	"time"

	"payment_microservice/internal/core/domain/entities"
)

// GetValidPaymentEntity retorna uma entidade de pagamento válida para testes
func GetValidPaymentEntity() entities.Payment {
	qrCode := "00020101021243650016COM.MERCADOLIBRE"
	paidAt := time.Now()
	createdAt := time.Now().Add(-1 * time.Hour)

	payment, err := entities.NewPayment(
		"payment-uuid-123",
		"order-uuid-456",
		100.50,
		"pending",
		"pix",
		&qrCode,
		&paidAt,
		createdAt,
	)

	if err != nil {
		panic("Failed to create mock payment: " + err.Error())
	}

	return *payment
}

// GetPendingPaymentEntity retorna um pagamento pendente
func GetPendingPaymentEntity() entities.Payment {
	qrCode := "00020101021243650016COM.MERCADOLIBRE"
	createdAt := time.Now().Add(-30 * time.Minute)

	payment, err := entities.NewPayment(
		"payment-uuid-pending",
		"order-uuid-pending",
		250.00,
		"pending",
		"pix",
		&qrCode,
		nil,
		createdAt,
	)

	if err != nil {
		panic("Failed to create mock pending payment: " + err.Error())
	}

	return *payment
}

// PaymentModelData representa os dados para criar um PaymentModel
type PaymentModelData struct {
	ID            string
	OrderID       string
	Amount        float64
	Status        string
	PaymentMethod string
	QRCodeURL     *string
	PaidAt        *time.Time
	CreatedAt     time.Time
}

// GetPaymentModelData retorna dados de modelo de pagamento
func GetPaymentModelData() PaymentModelData {
	qrCode := "00020101021243650016COM.MERCADOLIBRE"
	paidAt := time.Now()
	createdAt := time.Now().Add(-1 * time.Hour)

	return PaymentModelData{
		ID:            "payment-uuid-123",
		OrderID:       "order-uuid-456",
		Amount:        100.50,
		Status:        "pending",
		PaymentMethod: "pix",
		QRCodeURL:     &qrCode,
		PaidAt:        &paidAt,
		CreatedAt:     createdAt,
	}
}
