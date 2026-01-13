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

func (m *mockQueuePublisher) PublishOnQueue(_ context.Context, queueName string, message []byte) error {
	m.called = true
	m.queueName = queueName
	m.message = message
	return m.err
}

func (m *mockQueuePublisher) PublishOnTopic(_ context.Context, _ string, _ []byte) error {
	return nil
}

func TestSQSKitchenOrderServiceCreateSuccess(t *testing.T) {
	cleanup := env.SetupTestEnv()
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
	cleanup := env.SetupTestEnv()
	defer cleanup()

	expectedErr := errors.New("publish error")
	publisher := &mockQueuePublisher{err: expectedErr}
	service := NewSQSKitchenOrderService(publisher)

	err := service.Create(context.Background(), dto.CreateKitchenOrderDTO{OrderID: "order-456"})

	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, publisher.called)
}

func TestSQSKitchenOrderServiceCreateMarshalError(t *testing.T) {
	cleanup := env.SetupTestEnv()
	defer cleanup()

	originalMarshal := jsonMarshal
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("marshal error")
	}
	t.Cleanup(func() { jsonMarshal = originalMarshal })

	publisher := &mockQueuePublisher{}
	service := NewSQSKitchenOrderService(publisher)

	err := service.Create(context.Background(), dto.CreateKitchenOrderDTO{OrderID: "order-789"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal error")
	assert.False(t, publisher.called)
}
