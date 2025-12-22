package factory

import (
	"payment_microservice/internal/adapters/driven/kitchen_order_service"
	"payment_microservice/internal/core/domain/ports"
)

func NewKitchenOrderService() ports.IKitchenOrderService {
	return kitchen_order_service.NewRabbitMQKitchenOrderService()
}
