package consumers

import "payment_microservice/internal/core/factory"

func RegisterConsumers() {
	queueConsumer := factory.NewQueueConsumer()

	paymentConsumer := NewPaymentConsumer(queueConsumer)
	paymentConsumer.RegisterConsumers()
}
