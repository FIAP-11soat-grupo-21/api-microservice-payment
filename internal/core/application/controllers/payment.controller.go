package controllers

import (
	"context"

	"payment_microservice/internal/core/application/dtos"
	payment_gateways "payment_microservice/internal/core/application/gateways"
	"payment_microservice/internal/core/application/presenters"
	payment_interfaces "payment_microservice/internal/core/interfaces"
	"payment_microservice/internal/core/use_cases"
	payment_use_cases "payment_microservice/internal/core/use_cases"
	kitchen_order_gateways "tech_challenge/internal/kitchen-order/application/gateways"
	kitchen_order_interfaces "tech_challenge/internal/kitchen-order/interfaces"
	kitchen_order_use_cases "tech_challenge/internal/kitchen-order/use_cases"
)

type PaymentController struct {
	PaymentGateway      payment_gateways.PaymentGateway
	KitchenOrderGateway kitchen_order_gateways.KitchenOrderGateway
	OrderStatusGateway  kitchen_order_gateways.OrderStatusGateway
}

func NewPaymentController(
	paymentDataSource payment_interfaces.IPaymentDataSource,
	paymentProvider payment_interfaces.IPaymentProvider,
	kitchenOrderDataSource kitchen_order_interfaces.IKitchenOrderDataSource,
	orderStatusDataSource kitchen_order_interfaces.IOrderStatusDataSource,
) *PaymentController {
	return &PaymentController{
		PaymentGateway:      *payment_gateways.NewPaymentGateway(paymentDataSource, paymentProvider),
		KitchenOrderGateway: *kitchen_order_gateways.NewKitchenOrderGateway(kitchenOrderDataSource),
		OrderStatusGateway:  *kitchen_order_gateways.NewOrderStatusGateway(orderStatusDataSource),
	}
}

func (c *PaymentController) ConfirmPayment(ctx context.Context, webhookDTO dtos.WebhookEventDTO) error {
	confirmPaymentUseCase := payment_use_cases.NewConfirmPaymentUseCase(c.PaymentGateway)

	err := confirmPaymentUseCase.Execute(ctx, webhookDTO)

	if err != nil {
		return err
	}

	createKitchenOrderUseCase := kitchen_order_use_cases.NewCreateKitchenOrderUseCase(c.KitchenOrderGateway, c.OrderStatusGateway)

	_, err = createKitchenOrderUseCase.Execute(webhookDTO.OrderID)

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
