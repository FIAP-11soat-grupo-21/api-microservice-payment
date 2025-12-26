package mocks

import (
	"context"
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

// MockKitchenOrderService é um mock da interface IKitchenOrderService
type MockKitchenOrderService struct {
	mock.Mock
}

func (m *MockKitchenOrderService) Create(ctx context.Context, kitchenOrderDTO dto.CreateKitchenOrderDTO) error {
	args := m.Called(ctx, kitchenOrderDTO)
	return args.Error(0)
}
