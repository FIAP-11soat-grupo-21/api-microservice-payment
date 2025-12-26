package value_objects

import (
	"testing"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/exceptions"

	"github.com/stretchr/testify/assert"
)

func TestNewMethod(t *testing.T) {
	t.Run("should create method with PIX", func(t *testing.T) {
		method, err := NewMethod(constants.PIX_PAYMENT_METHOD)

		assert.NoError(t, err)
		assert.Equal(t, constants.PIX_PAYMENT_METHOD, method.Value())
	})

	t.Run("should return error for invalid payment method", func(t *testing.T) {
		invalidMethod := "credit_card"

		method, err := NewMethod(invalidMethod)

		assert.Error(t, err)
		assert.Equal(t, Method{}, method)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "invalid payment method")
	})

	t.Run("should return error for empty method", func(t *testing.T) {
		method, err := NewMethod("")

		assert.Error(t, err)
		assert.Equal(t, Method{}, method)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error for unknown method", func(t *testing.T) {
		unknownMethod := "bitcoin"

		method, err := NewMethod(unknownMethod)

		assert.Error(t, err)
		assert.Equal(t, Method{}, method)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})
}

func TestMethod_Value(t *testing.T) {
	t.Run("should return the correct value", func(t *testing.T) {
		method, _ := NewMethod(constants.PIX_PAYMENT_METHOD)

		result := method.Value()

		assert.Equal(t, constants.PIX_PAYMENT_METHOD, result)
	})
}

func TestMethod_IsPix(t *testing.T) {
	t.Run("should return true for PIX method", func(t *testing.T) {
		method, _ := NewMethod(constants.PIX_PAYMENT_METHOD)

		result := method.IsPix()

		assert.True(t, result)
	})
}

func TestIsValidMethod(t *testing.T) {
	t.Run("should return true for valid method", func(t *testing.T) {
		result := isValidMethod(constants.PIX_PAYMENT_METHOD)

		assert.True(t, result)
	})

	t.Run("should return false for invalid method", func(t *testing.T) {
		result := isValidMethod("invalid")

		assert.False(t, result)
	})

	t.Run("should return false for empty string", func(t *testing.T) {
		result := isValidMethod("")

		assert.False(t, result)
	})
}
