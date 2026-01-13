package factory

import (
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaymentRepository(t *testing.T) {
	cleanup := env.SetupTestEnv()
	defer cleanup()

	t.Run("should create payment repository successfully", func(t *testing.T) {
		repository := NewPaymentRepository()

		assert.NotNil(t, repository)
	})

	t.Run("should return IPaymentRepository interface", func(t *testing.T) {
		repository := NewPaymentRepository()

		assert.NotNil(t, repository)
		assert.Implements(t, (*ports.IPaymentRepository)(nil), repository)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		repository1 := NewPaymentRepository()
		repository2 := NewPaymentRepository()

		assert.NotNil(t, repository1)
		assert.NotNil(t, repository2)

		assert.NotSame(t, repository1, repository2)
	})
}
