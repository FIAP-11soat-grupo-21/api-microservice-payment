package dto

import (
	"context"
	"encoding/json"
	"time"
)

type CreatePaymentDTO struct {
	Ctx     context.Context
	OrderID string
	Amount  float64
	Method  string
}

type PaymentResultDTO struct {
	ID        string
	OrderID   string
	Amount    float64
	Status    string
	Method    string
	QRCode    *string
	CreatedAt time.Time
	PaidAt    *time.Time
}

type CreatePIXBillingDTO struct {
	Ctx        context.Context
	ExternalID string
	Amount     float64
}

type PIXBillingResultDTO struct {
	QRData string
}

type PaymentProcessedEventDTO struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func (dto *PaymentProcessedEventDTO) ToJSON() ([]byte, error) {
	return json.Marshal(dto)
}
