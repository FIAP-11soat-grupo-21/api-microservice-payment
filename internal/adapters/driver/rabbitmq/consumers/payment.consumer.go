package consumers

import (
	"payment_microservice/internal/adapters/driver/rabbitmq/handlers"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/infra/queue"
)

func RegisterConsumers() {
	cfg := env.GetConfig()
	go queue.RegisterConsumer("create-payment", cfg.RabbitMQ.Topics.CreatePayment, handlers.CreatePayment)
	go queue.RegisterConsumer("payment-order-error", cfg.RabbitMQ.Topics.OrderError, handlers.RollbackPayment)
}
