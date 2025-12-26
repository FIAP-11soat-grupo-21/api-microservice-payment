package value_objects

import (
	"testing"

	"payment_microservice/internal/core/domain/exceptions"

	"github.com/stretchr/testify/assert"
)

func TestNewAmount(t *testing.T) {
	t.Run("should create amount with valid positive value", func(t *testing.T) {
		value := 100.50

		amount, err := NewAmount(value)

		assert.NoError(t, err)
		assert.Equal(t, value, amount.Value())
	})

	t.Run("should create amount with zero value", func(t *testing.T) {
		value := 0.0

		amount, err := NewAmount(value)

		assert.NoError(t, err)
		assert.Equal(t, value, amount.Value())
	})

	t.Run("should return error when amount is negative", func(t *testing.T) {
		value := -10.50

		amount, err := NewAmount(value)

		assert.Error(t, err)
		assert.Equal(t, Amount{}, amount)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "amount cannot be negative")
	})

	t.Run("should create amount with large value", func(t *testing.T) {
		value := 999999.99

		amount, err := NewAmount(value)

		assert.NoError(t, err)
		assert.Equal(t, value, amount.Value())
	})

	t.Run("should create amount with decimal precision", func(t *testing.T) {
		value := 0.01

		amount, err := NewAmount(value)

		assert.NoError(t, err)
		assert.Equal(t, value, amount.Value())
	})
}

func TestAmount_Value(t *testing.T) {
	t.Run("should return the correct value", func(t *testing.T) {
		expectedValue := 50.75
		amount, _ := NewAmount(expectedValue)

		result := amount.Value()

		assert.Equal(t, expectedValue, result)
	})
}
