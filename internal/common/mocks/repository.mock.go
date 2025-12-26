package mocks

import (
	"context"
	"payment_microservice/internal/core/domain/entities"

	"github.com/stretchr/testify/mock"
)

// MockPaymentRepository é um mock da interface IPaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Insert(ctx context.Context, payment entities.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByOrderID(ctx context.Context, orderID string) (entities.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return entities.Payment{}, args.Error(1)
	}
	return args.Get(0).(entities.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment entities.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, paymentID string) error {
	args := m.Called(ctx, paymentID)
	return args.Error(0)
}
