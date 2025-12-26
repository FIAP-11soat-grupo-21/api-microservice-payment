package kitchen_order_service

import (
	"context"
	"encoding/json"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/infra/queue"
	"payment_microservice/internal/core/dto"
)

// Refact this to another package
type QueuePublisher interface {
	PublishMessageWithContext(ctx context.Context, routingKey string, body []byte) error
}

type DefaultQueuePublisher struct{}

func (d *DefaultQueuePublisher) PublishMessageWithContext(ctx context.Context, routingKey string, body []byte) error {
	return queue.PublishMessageWithContext(ctx, routingKey, body)
}

type RabbitMQKitchenOrderService struct {
	publisher QueuePublisher
}

func NewRabbitMQKitchenOrderService() *RabbitMQKitchenOrderService {
	return &RabbitMQKitchenOrderService{
		publisher: &DefaultQueuePublisher{},
	}
}

func NewRabbitMQKitchenOrderServiceWithPublisher(publisher QueuePublisher) *RabbitMQKitchenOrderService {
	return &RabbitMQKitchenOrderService{
		publisher: publisher,
	}
}

func (s *RabbitMQKitchenOrderService) Create(ctx context.Context, dto dto.CreateKitchenOrderDTO) error {
	dtoJSON, err := json.Marshal(dto)

	if err != nil {
		return err
	}

	cfg := env.GetConfig()

	routingKey := cfg.RabbitMQ.Topics.CreateKitchenOrder

	return s.publisher.PublishMessageWithContext(
		ctx,
		routingKey,
		dtoJSON,
	)
}
