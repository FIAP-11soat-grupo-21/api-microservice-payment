package consumers

import (
	"log"
	"payment_microservice/internal/adapters/driver/queue/handlers"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"
)

var paymentConsumerLogFatalf = func(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

type PaymentConsumer struct {
	consumer ports.IQueueConsumer
}

func NewPaymentConsumer(consumer ports.IQueueConsumer) *PaymentConsumer {
	return &PaymentConsumer{
		consumer: consumer,
	}
}

func (pc *PaymentConsumer) RegisterConsumers() {
	cfg := env.GetConfig()

	err := pc.consumer.ConsumeQueue(cfg.AWS.SQS.Queues.CreatePayment, handlers.CreatePayment)

	if err != nil {
		paymentConsumerLogFatalf("Failed to register create payment consumer: %v", err)
	}

	err = pc.consumer.ConsumeQueue(cfg.AWS.SQS.Queues.OrderError, handlers.RollbackPayment)

	if err != nil {
		paymentConsumerLogFatalf("Failed to register rollback payment consumer: %v", err)
	}

	log.Println("All queue consumers registered successfully")
}
