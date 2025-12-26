package handlers

import (
	"context"
	"log"
	"payment_microservice/internal/adapters/driver/rabbitmq/message"
	"payment_microservice/internal/core/factory"
	"payment_microservice/internal/core/use_cases"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RefoundPayment(d amqp.Delivery) {
	paymentRepository := factory.NewPaymentRepository()

	refoundPaymentUseCase := use_cases.NewRefoundPaymentUseCase(paymentRepository)

	messageParsed, err := message.NewPaymentMessageFromJSON(d.Body)

	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	ctx := context.Background()

	err = refoundPaymentUseCase.Execute(ctx, messageParsed.OrderID)

	if err != nil {
		log.Printf("Error executing refound payment use case: %v", err)
		return
	}

	log.Printf("Payment refunded successfully for order ID: %s", messageParsed.OrderID)

	err = d.Ack(false)

	if err != nil {
		log.Printf("Error acknowledging message: %v", err)
		return
	}
}
