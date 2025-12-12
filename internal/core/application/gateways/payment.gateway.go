package gateways

import (
	"context"
	"time"

	"payment_microservice/internal/core/daos"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/interfaces"
)

type PaymentGateway struct {
	dataSource interfaces.IPaymentDataSource
	provider   interfaces.IPaymentProvider
}

func NewPaymentGateway(dataSource interfaces.IPaymentDataSource, provider interfaces.IPaymentProvider) *PaymentGateway {
	return &PaymentGateway{
		dataSource: dataSource,
		provider:   provider,
	}
}

func (pg *PaymentGateway) Insert(ctx context.Context, payment entities.Payment) error {
	var qrCodeUrl *string

	if payment.QRCode != nil {
		qrCodeValue := payment.QRCode.Value()
		qrCodeUrl = &qrCodeValue
	}

	return pg.dataSource.Insert(ctx, daos.PaymentDAO{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount.Value(),
		Method:          payment.Method.Value(),
		QRCode:          qrCodeUrl,
		Status:          payment.Status.Value(),
		TransactionCode: payment.TransactionCode,
		CreatedAt:       payment.CreatedAt,
		PaidAt:          payment.PaidAt,
	})
}

func (pg *PaymentGateway) FindByOrderID(orderID string) (entities.Payment, error) {
	paymentDAO, err := pg.dataSource.FindByOrderID(orderID)

	if err != nil {
		return entities.Payment{}, err
	}

	payment, err := entities.NewPayment(
		paymentDAO.ID,
		paymentDAO.OrderID,
		paymentDAO.Amount,
		paymentDAO.Status,
		paymentDAO.Method,
		paymentDAO.TransactionCode,
		paymentDAO.QRCode,
		paymentDAO.PaidAt,
		paymentDAO.CreatedAt,
	)

	if err != nil {
		return entities.Payment{}, err
	}

	return *payment, nil
}

func (pg *PaymentGateway) CreatePIXBilling(ctx context.Context, payment entities.Payment) (daos.PIXBillingResultDAO, error) {
	return pg.provider.CreatePIXBilling(daos.CreatePIXBillingDAO{
		Ctx:        ctx,
		ExternalID: payment.OrderID,
		Amount:     payment.Amount.Value(),
	})
}

func (pg *PaymentGateway) Confirm(ctx context.Context, payment entities.Payment) {
	now := time.Now()

	var qrCodeUrl *string

	if payment.QRCode != nil {
		qrCodeValue := payment.QRCode.Value()
		qrCodeUrl = &qrCodeValue
	}

	pg.dataSource.Update(ctx, daos.PaymentDAO{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount.Value(),
		Method:          payment.Method.Value(),
		QRCode:          qrCodeUrl,
		Status:          payment.Status.Value(),
		TransactionCode: payment.TransactionCode,
		CreatedAt:       payment.CreatedAt,
		PaidAt:          &now,
	})
}

func (pg *PaymentGateway) Update(ctx context.Context, payment entities.Payment) error {
	var qrCodeUrl *string

	if payment.QRCode != nil {
		qrCodeValue := payment.QRCode.Value()
		qrCodeUrl = &qrCodeValue
	}

	return pg.dataSource.Update(ctx, daos.PaymentDAO{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount.Value(),
		Method:          payment.Method.Value(),
		QRCode:          qrCodeUrl,
		Status:          payment.Status.Value(),
		TransactionCode: payment.TransactionCode,
		CreatedAt:       payment.CreatedAt,
		PaidAt:          payment.PaidAt,
	})
}
