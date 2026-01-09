package handlers

import (
	"context"
	"encoding/json"
	"log"
	"payment_microservice/internal/adapters/driver/queue/message"
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/dto"
	"payment_microservice/internal/core/factory"
	"payment_microservice/internal/core/use_cases"
)

func sendRollbackEvent(orderID string) {
	queuePublisher := factory.NewQueuePublisher()

	ctx := context.Background()
	cfg := env.GetConfig()
	topic := cfg.AWS.SNS.Topics.OrderError

	var messagePayload = struct {
		OrderID string `json:"order_id"`
	}{
		OrderID: orderID,
	}

	messageJSON, err := json.Marshal(messagePayload)

	if err != nil {
		log.Printf("Error marshalling rollback event message: %v", err)
		return
	}

	log.Println("Sending rollback event to Create Order SAGA")
	err = queuePublisher.PublishOnTopic(ctx, topic, messageJSON)

	if err != nil {
		log.Printf("Error publishing rollback event to SNS topic %s: %v", topic, err)
	}
}

func CreatePayment(msgBody []byte) error {
	paymentRepository := factory.NewPaymentRepository()
	paymentGateway := factory.NewPaymentGateway()
	createPaymentUseCase := use_cases.NewCreatePaymentUseCase(paymentRepository, paymentGateway)

	messageParsed, err := message.NewCreatePaymentMessageFromJSON(msgBody)

	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return err
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
		sendRollbackEvent(messageParsed.OrderID)
		return err
	}

	log.Printf("Payment created successfully for order ID: %s", messageParsed.OrderID)
	return nil
}

func RollbackPayment(msgBody []byte) error {
	paymentRepository := factory.NewPaymentRepository()

	rollbackPaymentUseCase := use_cases.NewRollbackPaymentUseCase(paymentRepository)

	messageParsed, err := message.NewRollbackPaymentMessageFromJSON(msgBody)

	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return err
	}

	ctx := context.Background()

	err = rollbackPaymentUseCase.Execute(ctx, messageParsed.OrderID)

	if err != nil {
		log.Printf("Error executing rollback payment use case: %v", err)
		return err
	}

	log.Printf("Payment rolled back successfully for order ID: %s", messageParsed.OrderID)
	return nil
}
