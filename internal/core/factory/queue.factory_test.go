package factory

import (
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/ports"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueuePublisher(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	t.Run("should create queue publisher successfully", func(t *testing.T) {
		publisher := NewQueuePublisher()

		assert.NotNil(t, publisher)
	})

	t.Run("should return IQueuePublisher interface", func(t *testing.T) {
		publisher := NewQueuePublisher()

		assert.NotNil(t, publisher)
		assert.Implements(t, (*ports.IQueuePublisher)(nil), publisher)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		publisher1 := NewQueuePublisher()
		publisher2 := NewQueuePublisher()

		assert.NotNil(t, publisher1)
		assert.NotNil(t, publisher2)
		assert.NotSame(t, publisher1, publisher2)
	})
}

func TestNewQueueConsumer(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	t.Run("should create queue consumer successfully", func(t *testing.T) {
		consumer := NewQueueConsumer()

		assert.NotNil(t, consumer)
	})

	t.Run("should return IQueueConsumer interface", func(t *testing.T) {
		consumer := NewQueueConsumer()
		assert.NotNil(t, consumer)
		assert.Implements(t, (*ports.IQueueConsumer)(nil), consumer)
	})

	t.Run("should create new instance on each call", func(t *testing.T) {
		consumer1 := NewQueueConsumer()
		consumer2 := NewQueueConsumer()

		assert.NotNil(t, consumer1)
		assert.NotNil(t, consumer2)
		assert.NotSame(t, consumer1, consumer2)
	})
}
