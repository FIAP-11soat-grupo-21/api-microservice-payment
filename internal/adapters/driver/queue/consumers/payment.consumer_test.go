package consumers

import (
	"errors"
	"testing"

	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"

	"github.com/stretchr/testify/assert"
)

type mockQueueConsumer struct {
	queues []mockQueueExpectation
}

type mockQueueExpectation struct {
	queue string
	err   error
}

func (m *mockQueueConsumer) ConsumeQueue(queueName string, handler ports.MessageHandler) error {
	if len(m.queues) == 0 {
		return nil
	}
	expectation := m.queues[0]
	m.queues = m.queues[1:]
	if expectation.queue != queueName {
		return errors.New("unexpected queue")
	}
	return expectation.err
}

func TestPaymentConsumerRegisterSuccess(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	mockConsumer := &mockQueueConsumer{
		queues: []mockQueueExpectation{
			{queue: env.GetConfig().AWS.SQS.Queues.CreatePayment},
			{queue: env.GetConfig().AWS.SQS.Queues.OrderError},
		},
	}

	paymentConsumer := NewPaymentConsumer(mockConsumer)

	paymentConsumer.RegisterConsumers()

	assert.Empty(t, mockConsumer.queues)
}

func TestPaymentConsumerRegisterCreatePaymentError(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	mockConsumer := &mockQueueConsumer{
		queues: []mockQueueExpectation{
			{queue: env.GetConfig().AWS.SQS.Queues.CreatePayment, err: errors.New("consume error")},
		},
	}

	paymentConsumer := NewPaymentConsumer(mockConsumer)

	called := stubPaymentConsumerFatal(t)
	assert.Panics(t, func() {
		paymentConsumer.RegisterConsumers()
	})
	assert.True(t, *called)
}

func TestPaymentConsumerRegisterRollbackError(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	mockConsumer := &mockQueueConsumer{
		queues: []mockQueueExpectation{
			{queue: env.GetConfig().AWS.SQS.Queues.CreatePayment},
			{queue: env.GetConfig().AWS.SQS.Queues.OrderError, err: errors.New("rollback error")},
		},
	}

	paymentConsumer := NewPaymentConsumer(mockConsumer)

	called := stubPaymentConsumerFatal(t)
	assert.Panics(t, func() {
		paymentConsumer.RegisterConsumers()
	})
	assert.True(t, *called)
}

func stubPaymentConsumerFatal(t *testing.T) *bool {
	original := paymentConsumerLogFatalf
	called := new(bool)
	paymentConsumerLogFatalf = func(format string, v ...interface{}) {
		*called = true
		panic("fatal")
	}
	t.Cleanup(func() {
		paymentConsumerLogFatalf = original
	})
	return called
}
