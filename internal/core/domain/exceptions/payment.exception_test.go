package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidPaymentDataException_Error(t *testing.T) {
	t.Run("should return custom message when provided", func(t *testing.T) {
		message := "custom error message"
		exception := &InvalidPaymentDataException{
			Message: message,
		}

		result := exception.Error()

		assert.Equal(t, message, result)
	})

	t.Run("should return default message when no message provided", func(t *testing.T) {
		exception := &InvalidPaymentDataException{}

		result := exception.Error()

		assert.Equal(t, "invalid payment data", result)
	})

	t.Run("should return default message when empty message provided", func(t *testing.T) {
		exception := &InvalidPaymentDataException{
			Message: "",
		}

		result := exception.Error()

		assert.Equal(t, "invalid payment data", result)
	})

	t.Run("should return message with special characters", func(t *testing.T) {
		message := "error: amount cannot be negative (value: -10.5)"
		exception := &InvalidPaymentDataException{
			Message: message,
		}

		result := exception.Error()

		assert.Equal(t, message, result)
	})
}

func TestPaymentNotFoundException_Error(t *testing.T) {
	t.Run("should return custom message when provided", func(t *testing.T) {
		message := "payment with ID 123 not found"
		exception := &PaymentNotFoundException{
			Message: message,
		}

		result := exception.Error()

		assert.Equal(t, message, result)
	})

	t.Run("should return default message when no message provided", func(t *testing.T) {
		exception := &PaymentNotFoundException{}

		result := exception.Error()

		assert.Equal(t, "payment not found", result)
	})

	t.Run("should return default message when empty message provided", func(t *testing.T) {
		exception := &PaymentNotFoundException{
			Message: "",
		}

		result := exception.Error()

		assert.Equal(t, "payment not found", result)
	})

	t.Run("should return message with order reference", func(t *testing.T) {
		message := "Payment not found for OrderID: order-123"
		exception := &PaymentNotFoundException{
			Message: message,
		}

		result := exception.Error()

		assert.Equal(t, message, result)
	})
}
