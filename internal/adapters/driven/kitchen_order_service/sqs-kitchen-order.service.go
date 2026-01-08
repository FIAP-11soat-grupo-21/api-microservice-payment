package kitchen_order_service

import (
	"context"
	"encoding/json"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"
)

type SQSKitchenOrderService struct {
	publisher ports.IQueuePublisher
}

func NewSQSKitchenOrderService(publisher ports.IQueuePublisher) *SQSKitchenOrderService {
	return &SQSKitchenOrderService{publisher: publisher}
}

func (s *SQSKitchenOrderService) Create(ctx context.Context, dto dto.CreateKitchenOrderDTO) error {
	dtoJSON, err := json.Marshal(dto)

	if err != nil {
		return err
	}

	cfg := env.GetConfig()

	queueName := cfg.AWS.SQS.Queues.CreateKitchenOrder

	return s.publisher.PublishOnQueue(
		ctx,
		queueName,
		dtoJSON,
	)
}
