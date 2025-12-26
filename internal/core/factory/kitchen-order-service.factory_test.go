package factory

import (
	"payment_microservice/internal/common/config/env"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKitchenOrderService(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	t.Run("should create kitchen order service successfully", func(t *testing.T) {
		service := NewKitchenOrderService()

		assert.NotNil(t, service)
	})

	t.Run("should return IKitchenOrderService interface", func(t *testing.T) {
		service := NewKitchenOrderService()

		assert.NotNil(t, service)
		assert.Implements(t, (*interface{})(nil), service)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		service1 := NewKitchenOrderService()
		service2 := NewKitchenOrderService()

		assert.NotNil(t, service1)
		assert.NotNil(t, service2)
		assert.NotSame(t, service1, service2)
	})
}
