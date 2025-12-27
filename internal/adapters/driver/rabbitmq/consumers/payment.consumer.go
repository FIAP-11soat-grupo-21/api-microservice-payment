package consumers

import (
	"payment_microservice/internal/adapters/driver/rabbitmq/handlers"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/infra/queue"
)

func RegisterConsumers() {
	cfg := env.GetConfig()
	go queue.RegisterConsumer("refound-payment", cfg.RabbitMQ.Topics.RefoundPayment, handlers.RefoundPayment)
}
