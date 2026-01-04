package handlers

import (
	"context"
	"log"
	"payment_microservice/internal/adapters/driver/rabbitmq/message"
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/dto"
	"payment_microservice/internal/core/factory"
	"payment_microservice/internal/core/use_cases"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreatePayment(d amqp.Delivery) {
	paymentRepository := factory.NewPaymentRepository()
	paymentGateway := factory.NewPaymentGateway()

	createPaymentUseCase := use_cases.NewCreatePaymentUseCase(paymentRepository, paymentGateway)

	messageParsed, err := message.NewCreatePaymentMessageFromJSON(d.Body)

	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	ctx := context.Background()

	_, err = createPaymentUseCase.Execute(dto.CreatePaymentDTO{
		Ctx:     ctx,
		OrderID: messageParsed.OrderID,
		Amount:  messageParsed.Amount,
		Method:  constants.PIX_PAYMENT_METHOD,
	})

	if err != nil {
		log.Printf("Error executing create payment use case: %v", err)
		return
	}

	log.Printf("Payment created successfully for order ID: %s", messageParsed.OrderID)
	err = d.Ack(false)

	if err != nil {
		log.Printf("Error acknowledging message: %v", err)
		return
	}

}

func RollbackPayment(d amqp.Delivery) {
	paymentRepository := factory.NewPaymentRepository()

	rollbackPaymentUseCase := use_cases.NewRollbackPaymentUseCase(paymentRepository)

	messageParsed, err := message.NewRollbackPaymentMessageFromJSON(d.Body)

	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	ctx := context.Background()

	err = rollbackPaymentUseCase.Execute(ctx, messageParsed.OrderID)

	if err != nil {
		log.Printf("Error executing rollback payment use case: %v", err)
		return
	}

	log.Printf("Payment rolled back successfully for order ID: %s", messageParsed.OrderID)
	err = d.Ack(false)

	if err != nil {
		log.Printf("Error acknowledging message: %v", err)
		return
	}
}
