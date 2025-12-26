package use_cases

import (
	"context"
	"testing"

	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"

	"github.com/stretchr/testify/assert"
)

func TestNewFindPaymentByOrderIDUseCase(t *testing.T) {
	t.Run("should create use case with repository", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)

		useCase := NewFindPaymentByOrderIDUseCase(mockRepo)

		assert.NotNil(t, useCase)
		assert.Equal(t, mockRepo, useCase.repository)
	})
}

func TestFindPaymentByOrderIDUseCase_Execute(t *testing.T) {
	t.Run("should return payment when found successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		ctx := context.Background()
		orderID := "order-123"

		expectedPayment := entities.Payment{
			ID:      "payment-123",
			OrderID: orderID,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(expectedPayment, nil)
		useCase := NewFindPaymentByOrderIDUseCase(mockRepo)

		payment, err := useCase.Execute(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPayment, payment)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return PaymentNotFoundException when payment not found", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		ctx := context.Background()
		orderID := "order-not-found"

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, &exceptions.PaymentNotFoundException{})
		useCase := NewFindPaymentByOrderIDUseCase(mockRepo)

		payment, err := useCase.Execute(ctx, orderID)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return PaymentNotFoundException when repository returns any error", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		ctx := context.Background()
		orderID := "order-error"

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, &exceptions.InvalidPaymentDataException{})
		useCase := NewFindPaymentByOrderIDUseCase(mockRepo)

		payment, err := useCase.Execute(ctx, orderID)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when context is cancelled", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		ctx, cancel := context.WithCancel(context.Background())
		orderID := "order-123"

		cancel()

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, context.Canceled)
		useCase := NewFindPaymentByOrderIDUseCase(mockRepo)

		payment, err := useCase.Execute(ctx, orderID)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when orderID is empty", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		ctx := context.Background()
		orderID := ""

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, &exceptions.PaymentNotFoundException{})
		useCase := NewFindPaymentByOrderIDUseCase(mockRepo)

		payment, err := useCase.Execute(ctx, orderID)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		mockRepo.AssertExpectations(t)
	})
}
