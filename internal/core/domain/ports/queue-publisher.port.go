package ports

import "context"

type IQueuePublisher interface {
	PublishOnQueue(ctx context.Context, queueName string, message []byte) error
	PublishOnTopic(ctx context.Context, topic string, message []byte) error
}
