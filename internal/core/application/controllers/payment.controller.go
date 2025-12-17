package controllers

import (
	"context"

	"payment_microservice/internal/core/application/dtos"
	payment_gateways "payment_microservice/internal/core/application/gateways"
	"payment_microservice/internal/core/application/presenters"
	payment_interfaces "payment_microservice/internal/core/interfaces"
	"payment_microservice/internal/core/use_cases"
)

type PaymentController struct {
	PaymentGateway payment_gateways.PaymentGateway
}

func NewPaymentController(
	paymentDataSource payment_interfaces.IPaymentDataSource,
	paymentProvider payment_interfaces.IPaymentProvider,
) *PaymentController {
	return &PaymentController{
		PaymentGateway: *payment_gateways.NewPaymentGateway(paymentDataSource, paymentProvider),
	}
}

func (c *PaymentController) CreatePayment(paymentDTO dtos.CreatePaymentDTO) (dtos.PaymentResultDTO, error) {
	createPaymentUseCase := use_cases.NewCreatePaymentUseCase(c.PaymentGateway)
	payment, err := createPaymentUseCase.Execute(paymentDTO)

	if err != nil {
		return dtos.PaymentResultDTO{}, err
	}

	return presenters.ToResponse(payment), nil
}

func (c *PaymentController) ConfirmPayment(ctx context.Context, webhookDTO dtos.WebhookEventDTO) error {
	confirmPaymentUseCase := use_cases.NewConfirmPaymentUseCase(c.PaymentGateway)

	err := confirmPaymentUseCase.Execute(ctx, webhookDTO)

	if err != nil {
		return err
	}

	return nil
}

func (c *PaymentController) GetByOrderID(orderID string) (dtos.PaymentResultDTO, error) {
	GetPaymentByOrderIDUseCase := use_cases.NewFindPaymentByOrderIDUseCase(c.PaymentGateway)

	payment, err := GetPaymentByOrderIDUseCase.Execute(orderID)

	if err != nil {
		return dtos.PaymentResultDTO{}, err
	}

	return presenters.ToResponse(payment), nil

}
