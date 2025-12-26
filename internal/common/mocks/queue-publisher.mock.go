package mocks

import (
	"context"
)

type MockQueuePublisher struct {
	PublishFunc func(ctx context.Context, routingKey string, body []byte) error
}

func (m *MockQueuePublisher) PublishMessageWithContext(ctx context.Context, routingKey string, body []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, routingKey, body)
	}
	return nil
}

func NewMockQueuePublisher(publishFunc func(ctx context.Context, routingKey string, body []byte) error) *MockQueuePublisher {
	return &MockQueuePublisher{
		PublishFunc: publishFunc,
	}
}
