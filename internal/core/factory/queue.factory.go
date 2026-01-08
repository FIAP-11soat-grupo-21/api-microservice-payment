package factory

import (
	"payment_microservice/internal/common/infra/queue"
	"payment_microservice/internal/core/domain/ports"
)

func NewQueuePublisher() ports.IQueuePublisher {
	return queue.NewSQSPublisher()
}

func NewQueueConsumer() ports.IQueueConsumer {
	return queue.NewSQSConsumer()
}
