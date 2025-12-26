package use_cases

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/value_objects"
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConfirmPaymentUseCase(t *testing.T) {
	t.Run("should create use case with repository and kitchen order service", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)

		assert.NotNil(t, useCase)
		assert.Equal(t, mockRepo, useCase.repository)
		assert.Equal(t, mockKitchenService, useCase.kitchenOrderService)
	})
}

func TestConfirmPaymentUseCase_Execute(t *testing.T) {
	t.Run("should confirm payment and create kitchen order when event is payment.updated", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx := context.Background()

		orderID := "order-123"
		eventDTO := dto.WebhookEventDTO{
			ID:          "event-123",
			Type:        "payment",
			Action:      "payment.updated",
			OrderID:     orderID,
			LiveMode:    true,
			DateCreated: "2025-12-21T10:00:00Z",
			APIVersion:  "v1",
		}

		amount, _ := value_objects.NewAmount(100.50)
		status, _ := value_objects.NewStatus("pending")
		method, _ := value_objects.NewMethod("pix")

		existingPayment := entities.Payment{
			ID:      "payment-123",
			OrderID: orderID,
			Amount:  amount,
			Status:  status,
			Method:  method,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(existingPayment, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(payment entities.Payment) bool {
			return payment.ID == existingPayment.ID && payment.Status.IsPaid()
		})).Return(nil)

		mockKitchenService.On("Create", ctx, dto.CreateKitchenOrderDTO{
			OrderID: orderID,
		}).Return(nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockKitchenService.AssertExpectations(t)
	})

	t.Run("should return error when payment is not found", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx := context.Background()

		orderID := "order-not-found"
		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "payment",
			Action:  "payment.updated",
			OrderID: orderID,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, errors.New("not found"))

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.Error(t, err)
		assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
		assert.Contains(t, err.Error(), orderID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should not update payment when event type is not payment", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx := context.Background()

		orderID := "order-123"
		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "order",
			Action:  "order.created",
			OrderID: orderID,
		}

		amount, _ := value_objects.NewAmount(100.50)
		status, _ := value_objects.NewStatus("pending")
		method, _ := value_objects.NewMethod("pix")

		existingPayment := entities.Payment{
			ID:      "payment-123",
			OrderID: orderID,
			Amount:  amount,
			Status:  status,
			Method:  method,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(existingPayment, nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
		mockKitchenService.AssertNotCalled(t, "Create")
	})

	t.Run("should not update payment when action is not payment.updated", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx := context.Background()

		orderID := "order-123"
		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "payment",
			Action:  "payment.created",
			OrderID: orderID,
		}

		amount, _ := value_objects.NewAmount(100.50)
		status, _ := value_objects.NewStatus("pending")
		method, _ := value_objects.NewMethod("pix")

		existingPayment := entities.Payment{
			ID:      "payment-123",
			OrderID: orderID,
			Amount:  amount,
			Status:  status,
			Method:  method,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(existingPayment, nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
		mockKitchenService.AssertNotCalled(t, "Create")
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx, cancel := context.WithCancel(context.Background())

		orderID := "order-123"
		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "payment",
			Action:  "payment.updated",
			OrderID: orderID,
		}

		cancel()

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, context.Canceled)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.Error(t, err)
		assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should execute successfully when payment found even with empty OrderID in event", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx := context.Background()

		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "payment",
			Action:  "payment.updated",
			OrderID: "",
		}

		amount, _ := value_objects.NewAmount(100.50)
		status, _ := value_objects.NewStatus("pending")
		method, _ := value_objects.NewMethod("pix")

		existingPayment := entities.Payment{
			ID:      "payment-123",
			OrderID: "",
			Amount:  amount,
			Status:  status,
			Method:  method,
		}

		mockRepo.On("FindByOrderID", ctx, "").Return(existingPayment, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)
		mockKitchenService.On("Create", ctx, dto.CreateKitchenOrderDTO{
			OrderID: "",
		}).Return(nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockKitchenService.AssertExpectations(t)
	})

	t.Run("should still execute kitchen service even if update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockKitchenService := new(mocks.MockKitchenOrderService)
		ctx := context.Background()

		orderID := "order-123"
		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "payment",
			Action:  "payment.updated",
			OrderID: orderID,
		}

		amount, _ := value_objects.NewAmount(100.50)
		status, _ := value_objects.NewStatus("pending")
		method, _ := value_objects.NewMethod("pix")

		existingPayment := entities.Payment{
			ID:      "payment-123",
			OrderID: orderID,
			Amount:  amount,
			Status:  status,
			Method:  method,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(existingPayment, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("update failed"))
		mockKitchenService.On("Create", ctx, dto.CreateKitchenOrderDTO{
			OrderID: orderID,
		}).Return(nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockKitchenService)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockKitchenService.AssertExpectations(t)
	})
}
