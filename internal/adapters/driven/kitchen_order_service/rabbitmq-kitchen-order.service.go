package kitchen_order_service

import (
	"context"
	"encoding/json"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/infra/queue"
	"payment_microservice/internal/core/dto"
)

type RabbitMQKitchenOrderService struct{}

func NewRabbitMQKitchenOrderService() *RabbitMQKitchenOrderService {
	return &RabbitMQKitchenOrderService{}
}
func (s *RabbitMQKitchenOrderService) Create(ctx context.Context, dto dto.CreateKitchenOrderDTO) error {
	dtoJSON, err := json.Marshal(dto)

	if err != nil {
		return err
	}

	cfg := env.GetConfig()

	kitchenOrderQueueName := cfg.RabbitMQ.CreateKitchenOrderTopic

	return queue.PublishMessageWithContext(
		ctx,
		kitchenOrderQueueName,
		dtoJSON,
	)
}
