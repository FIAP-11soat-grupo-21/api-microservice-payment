package daos

import "time"

type PaymentDAO struct {
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
