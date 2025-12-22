package ports

import (
	"context"
	"payment_microservice/internal/core/dto"
)

type IKitchenOrderService interface {
	Create(ctx context.Context, dto dto.CreateKitchenOrderDTO) error
}
