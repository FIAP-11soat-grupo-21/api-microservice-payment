package schemas

import "time"

type ConfirmPaymentSchema struct {
	ID          string `json:"id"`
	LiveMode    bool   `json:"live_mode"`
	Type        string `json:"type"`
	DateCreated string `json:"date_created"`
	UserID      int64  `json:"user_id"`
	APIVersion  string `json:"api_version"`
	Action      string `json:"action"`
	Data        struct {
		ID string `json:"id"`
	} `json:"data"`
}

type PaymentResponseSchema struct {
	ID              string     `json:"id"`
	OrderID         string     `json:"order_id"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"`
	Method          string     `json:"method"`
	TransactionCode *string    `json:"transaction_code,omitempty"`
	QRCode          *string    `json:"qr_code,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
}
