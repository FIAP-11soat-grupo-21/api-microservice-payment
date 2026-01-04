package repositories

import "time"

type PaymentModel struct {
	ID            string     `gorm:"column:id;primaryKey;size:36"`
	OrderID       string     `gorm:"column:order_id;uniqueIndex;size:36;not null"`
	Amount        float64    `gorm:"column:amount;not null"`
	Status        string     `gorm:"column:status;size:50;not null"`
	PaymentMethod string     `gorm:"column:payment_method;size:50;not null"`
	QRCodeURL     *string    `gorm:"column:qr_code_url"`
	PaidAt        *time.Time `gorm:"column:paid_at;type:timestamp(6)"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;type:timestamp(6)"`
}

func (PaymentModel) TableName() string { return "payments" }
