package dtos

import (
	"context"
	"time"
)

type CreatePaymentDTO struct {
	Ctx     context.Context
	OrderID string
	Amount  float64
	Method  string
}

type PaymentResultDTO struct {
	ID              string
	OrderID         string
	Amount          float64
	Status          string
	Method          string
	TransactionCode *string
	QRCode          *string
	CreatedAt       time.Time
	PaidAt          *time.Time
}
