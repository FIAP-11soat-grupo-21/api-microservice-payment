package kitchen_order_service

import (
	"context"
	"encoding/json"
	"errors"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQKitchenOrderService(t *testing.T) {
	t.Run("should create service with default publisher", func(t *testing.T) {
		service := NewRabbitMQKitchenOrderService()

		assert.NotNil(t, service)
		assert.NotNil(t, service.publisher)
		assert.IsType(t, &DefaultQueuePublisher{}, service.publisher)
	})
}

func TestNewRabbitMQKitchenOrderServiceWithPublisher(t *testing.T) {
	t.Run("should create service with custom publisher", func(t *testing.T) {
		mockPublisher := mocks.NewMockQueuePublisher(nil)
		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		assert.NotNil(t, service)
		assert.Equal(t, mockPublisher, service.publisher)
	})
}

func TestCreate(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("should create kitchen order successfully", func(t *testing.T) {
		var capturedRoutingKey string
		var capturedBody []byte

		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			capturedRoutingKey = routingKey
			capturedBody = body
			return nil
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "order-123",
		}

		err := service.Create(ctx, orderDTO)

		assert.NoError(t, err)
		assert.Equal(t, "create.kitchen-order", capturedRoutingKey)

		var unmarshaledDTO dto.CreateKitchenOrderDTO
		err = json.Unmarshal(capturedBody, &unmarshaledDTO)
		assert.NoError(t, err)
		assert.Equal(t, "order-123", unmarshaledDTO.OrderID)
	})

	t.Run("should return error when publisher fails", func(t *testing.T) {
		expectedError := errors.New("publish error")

		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			return expectedError
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "order-456",
		}

		err := service.Create(ctx, orderDTO)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should handle empty order id", func(t *testing.T) {
		var capturedBody []byte

		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			capturedBody = body
			return nil
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "",
		}

		err := service.Create(ctx, orderDTO)

		assert.NoError(t, err)

		var unmarshaledDTO dto.CreateKitchenOrderDTO
		err = json.Unmarshal(capturedBody, &unmarshaledDTO)
		assert.NoError(t, err)
		assert.Equal(t, "", unmarshaledDTO.OrderID)
	})

	t.Run("should use correct routing key from config", func(t *testing.T) {
		var capturedRoutingKey string

		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			capturedRoutingKey = routingKey
			return nil
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "order-789",
		}

		err := service.Create(ctx, orderDTO)

		assert.NoError(t, err)
		cfg := env.GetConfig()
		assert.Equal(t, cfg.RabbitMQ.Topics.CreateKitchenOrder, capturedRoutingKey)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			return ctx.Err()
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "order-999",
		}

		err := service.Create(cancelledCtx, orderDTO)

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("should marshal dto correctly", func(t *testing.T) {
		var capturedBody []byte

		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			capturedBody = body
			return nil
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "order-special-123",
		}

		err := service.Create(ctx, orderDTO)

		assert.NoError(t, err)

		expectedJSON, _ := json.Marshal(orderDTO)
		assert.JSONEq(t, string(expectedJSON), string(capturedBody))
	})

	t.Run("should verify json marshalling is called", func(t *testing.T) {
		mockPublisher := mocks.NewMockQueuePublisher(func(ctx context.Context, routingKey string, body []byte) error {
			return nil
		})

		service := NewRabbitMQKitchenOrderServiceWithPublisher(mockPublisher)

		orderDTO := dto.CreateKitchenOrderDTO{
			OrderID: "test-order-json",
		}

		err := service.Create(ctx, orderDTO)

		assert.NoError(t, err)
	})
}

func TestDefaultQueuePublisher(t *testing.T) {
	t.Run("should implement QueuePublisher interface", func(t *testing.T) {
		publisher := &DefaultQueuePublisher{}

		var _ QueuePublisher = publisher
	})
}
