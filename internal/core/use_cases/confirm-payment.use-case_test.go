package use_cases

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/value_objects"
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConfirmPaymentUseCase(t *testing.T) {
	t.Run("should create use case with repository and message publisher", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)

		assert.NotNil(t, useCase)
		assert.Equal(t, mockRepo, useCase.repository)
		assert.Equal(t, mockPublisher, useCase.messagePublisher)
	})
}

func TestConfirmPaymentUseCase_Execute(t *testing.T) {
	cleanup := env.SetupTestEnv()
	defer cleanup()

	t.Run("should confirm payment and publish message when event is payment.updated", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
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

		mockPublisher.On("PublishOnTopic", ctx, mock.Anything, mock.Anything).Return(nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("should return error when payment is not found", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
		ctx := context.Background()

		orderID := "order-not-found"
		eventDTO := dto.WebhookEventDTO{
			ID:      "event-123",
			Type:    "payment",
			Action:  "payment.updated",
			OrderID: orderID,
		}

		mockRepo.On("FindByOrderID", ctx, orderID).Return(nil, errors.New("not found"))

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.Error(t, err)
		assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
		assert.Contains(t, err.Error(), orderID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should not update payment when event type is not payment", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
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

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
		mockPublisher.AssertNotCalled(t, "Create")
	})

	t.Run("should not update payment when action is not payment.updated", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
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

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
		mockPublisher.AssertNotCalled(t, "Create")
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
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

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.Error(t, err)
		assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should execute successfully when payment found even with empty OrderID in event", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
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

		topic := env.GetConfig().AWS.SNS.Topics.PaymentProcessed
		message := dto.PaymentProcessedEventDTO{
			OrderID: "",
			Status:  constants.ORDER_STATUS_CONFIRMED,
		}

		messageJSON, _ := message.ToJSON()

		mockPublisher.On("PublishOnTopic", ctx, topic, messageJSON).Return(nil)

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockPublisher := new(mocks.MockQueuePublisher)
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

		useCase := NewConfirmPaymentUseCase(mockRepo, mockPublisher)
		err := useCase.Execute(ctx, eventDTO)

		assert.Error(t, err)
		assert.EqualError(t, err, "update failed")
		mockRepo.AssertExpectations(t)
		mockPublisher.AssertNotCalled(t, "Create")
	})
}
