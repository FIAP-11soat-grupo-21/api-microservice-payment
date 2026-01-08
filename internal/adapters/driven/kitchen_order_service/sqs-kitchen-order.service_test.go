package kitchen_order_service

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/assert"
)

type mockQueuePublisher struct {
	queueName string
	message   []byte
	called    bool
	err       error
}

func (m *mockQueuePublisher) PublishOnQueue(ctx context.Context, queueName string, message []byte) error {
	m.called = true
	m.queueName = queueName
	m.message = message
	return m.err
}

func (m *mockQueuePublisher) PublishOnTopic(ctx context.Context, topic string, message []byte) error {
	return nil
}

func TestSQSKitchenOrderServiceCreateSuccess(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	publisher := &mockQueuePublisher{}
	service := NewSQSKitchenOrderService(publisher)

	dtoInput := dto.CreateKitchenOrderDTO{OrderID: "order-123"}

	err := service.Create(context.Background(), dtoInput)

	assert.NoError(t, err)
	assert.True(t, publisher.called)
	assert.Equal(t, env.GetConfig().AWS.SQS.Queues.CreateKitchenOrder, publisher.queueName)
	assert.JSONEq(t, `{"OrderID":"order-123"}`, string(publisher.message))
}

func TestSQSKitchenOrderServiceCreatePublisherError(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	expectedErr := errors.New("publish error")
	publisher := &mockQueuePublisher{err: expectedErr}
	service := NewSQSKitchenOrderService(publisher)

	err := service.Create(context.Background(), dto.CreateKitchenOrderDTO{OrderID: "order-456"})

	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, publisher.called)
}
