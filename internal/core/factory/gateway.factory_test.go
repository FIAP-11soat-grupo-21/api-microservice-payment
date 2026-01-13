package factory

import (
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaymentGateway(t *testing.T) {
	cleanup := env.SetupTestEnv()
	defer cleanup()

	t.Run("should create payment gateway successfully", func(t *testing.T) {
		gateway := NewPaymentGateway()

		assert.NotNil(t, gateway)
	})

	t.Run("should return IPaymentGateway interface", func(t *testing.T) {
		gateway := NewPaymentGateway()

		assert.NotNil(t, gateway)
		assert.Implements(t, (*ports.IPaymentGateway)(nil), gateway)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		gateway1 := NewPaymentGateway()
		gateway2 := NewPaymentGateway()

		assert.NotNil(t, gateway1)
		assert.NotNil(t, gateway2)
		assert.NotSame(t, gateway1, gateway2)
	})
}
