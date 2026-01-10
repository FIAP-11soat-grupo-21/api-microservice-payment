package entities

import (
	"time"

	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/value_objects"
)

type Payment struct {
	ID        string
	OrderID   string
	Amount    value_objects.Amount
	Status    value_objects.Status
	Method    value_objects.Method
	QRCode    *value_objects.QRCode
	CreatedAt time.Time
	PaidAt    *time.Time
}

func NewPaymentDefault(id, orderID string, amount float64, method string, qrData *string) (*Payment, error) {
	paymentMethod, err := value_objects.NewMethod(method)

	if err != nil {
		return nil, err
	}

	paymentStatus := value_objects.NewStatusDefault()

	paymentAmount, err := value_objects.NewAmount(amount)

	if err != nil {
		return nil, err
	}

	var paymentQrCode *value_objects.QRCode

	if qrData != nil {
		qrCode, err := value_objects.NewQRCode(*qrData)

		if err != nil {
			return nil, err
		}

		paymentQrCode = &qrCode
	}

	return &Payment{
		ID:        id,
		OrderID:   orderID,
		Amount:    paymentAmount,
		Status:    paymentStatus,
		Method:    paymentMethod,
		QRCode:    paymentQrCode,
		CreatedAt: time.Now(),
	}, nil
}

func NewPayment(
	id,
	orderID string,
	amount float64,
	status,
	method string,
	qrCode *string,
	paidAt *time.Time,
	createdAt time.Time,
) (*Payment, error) {
	paymentMethod, err := value_objects.NewMethod(method)
	if err != nil {
		return nil, err
	}

	paymentStatus, err := value_objects.NewStatus(status)
	if err != nil {
		return nil, err
	}

	paymentAmount, err := value_objects.NewAmount(amount)
	if err != nil {
		return nil, err
	}

	var paymentQrCode *value_objects.QRCode
	if qrCode != nil {
		qrCodeVO, err := value_objects.NewQRCode(*qrCode)
		if err != nil {
			return nil, err
		}
		paymentQrCode = &qrCodeVO
	}

	return &Payment{
		ID:        id,
		OrderID:   orderID,
		Amount:    paymentAmount,
		Status:    paymentStatus,
		Method:    paymentMethod,
		QRCode:    paymentQrCode,
		PaidAt:    paidAt,
		CreatedAt: createdAt,
	}, nil
}

func (p *Payment) IsEmpty() bool {
	return p.ID == ""
}

func (p *Payment) MarkAsPaid() {
	if p.Status.IsPending() {
		p.Status.SetPaid()

		now := time.Now()
		p.PaidAt = &now
	}
}

func (p *Payment) MarkAsFailed() {
	p.Status.SetFailed()
}

func (p *Payment) MarkAsRefunded() {
	p.Status.SetRefunded()
}

func (p *Payment) SetQrCode(qrCode string) error {
	if p.Method.IsPix() {
		qrCodeVO, err := value_objects.NewQRCode(qrCode)

		if err != nil {
			return err
		}

		p.QRCode = &qrCodeVO
		return nil
	}

	return &exceptions.InvalidPaymentDataException{
		Message: "QR Code can only be set for Pix payments",
	}
}
