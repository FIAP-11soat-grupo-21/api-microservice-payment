package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockQueuePublisher struct {
	mock.Mock
}

func (m *MockQueuePublisher) PublishOnQueue(ctx context.Context, queueName string, body []byte) error {
	args := m.Called(ctx, queueName, body)
	return args.Error(0)
}

func (m *MockQueuePublisher) PublishOnTopic(ctx context.Context, topic string, body []byte) error {
	args := m.Called(ctx, topic, body)
	return args.Error(0)
}
