package value_objects

import (
	"testing"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/exceptions"

	"github.com/stretchr/testify/assert"
)

func TestNewStatus(t *testing.T) {
	t.Run("should create status with pending value", func(t *testing.T) {
		status, err := NewStatus(constants.PAYMENT_STATUS_PENDING)

		assert.NoError(t, err)
		assert.Equal(t, constants.PAYMENT_STATUS_PENDING, status.Value())
	})

	t.Run("should create status with paid value", func(t *testing.T) {
		status, err := NewStatus(constants.PAYMENT_STATUS_PAID)

		assert.NoError(t, err)
		assert.Equal(t, constants.PAYMENT_STATUS_PAID, status.Value())
	})

	t.Run("should create status with failed value", func(t *testing.T) {
		status, err := NewStatus(constants.PAYMENT_STATUS_FAILED)

		assert.NoError(t, err)
		assert.Equal(t, constants.PAYMENT_STATUS_FAILED, status.Value())
	})

	t.Run("should return error for invalid status", func(t *testing.T) {
		invalidStatus := "cancelled"

		status, err := NewStatus(invalidStatus)

		assert.Error(t, err)
		assert.Equal(t, Status{}, status)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "invalid payment status")
	})

	t.Run("should return error for empty status", func(t *testing.T) {
		status, err := NewStatus("")

		assert.Error(t, err)
		assert.Equal(t, Status{}, status)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})
}

func TestNewStatusDefault(t *testing.T) {
	t.Run("should create default status as pending", func(t *testing.T) {
		status := NewStatusDefault()

		assert.Equal(t, constants.PAYMENT_STATUS_PENDING, status.Value())
		assert.True(t, status.IsPending())
	})
}

func TestStatus_Value(t *testing.T) {
	t.Run("should return the correct value", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PAID)

		result := status.Value()

		assert.Equal(t, constants.PAYMENT_STATUS_PAID, result)
	})
}

func TestIsValidStatus(t *testing.T) {
	t.Run("should return true for pending status", func(t *testing.T) {
		result := isValidStatus(constants.PAYMENT_STATUS_PENDING)

		assert.True(t, result)
	})

	t.Run("should return true for paid status", func(t *testing.T) {
		result := isValidStatus(constants.PAYMENT_STATUS_PAID)

		assert.True(t, result)
	})

	t.Run("should return true for failed status", func(t *testing.T) {
		result := isValidStatus(constants.PAYMENT_STATUS_FAILED)

		assert.True(t, result)
	})

	t.Run("should return false for invalid status", func(t *testing.T) {
		result := isValidStatus("invalid")

		assert.False(t, result)
	})

	t.Run("should return false for empty string", func(t *testing.T) {
		result := isValidStatus("")

		assert.False(t, result)
	})
}

func TestStatus_IsPending(t *testing.T) {
	t.Run("should return true when status is pending", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PENDING)

		result := status.IsPending()

		assert.True(t, result)
	})

	t.Run("should return false when status is not pending", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PAID)

		result := status.IsPending()

		assert.False(t, result)
	})
}

func TestStatus_SetPending(t *testing.T) {
	t.Run("should set status to pending", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PAID)

		status.SetPending()

		assert.Equal(t, constants.PAYMENT_STATUS_PENDING, status.Value())
		assert.True(t, status.IsPending())
	})
}

func TestStatus_IsPaid(t *testing.T) {
	t.Run("should return true when status is paid", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PAID)

		result := status.IsPaid()

		assert.True(t, result)
	})

	t.Run("should return false when status is not paid", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PENDING)

		result := status.IsPaid()

		assert.False(t, result)
	})
}

func TestStatus_SetPaid(t *testing.T) {
	t.Run("should set status to paid", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PENDING)

		status.SetPaid()

		assert.Equal(t, constants.PAYMENT_STATUS_PAID, status.Value())
		assert.True(t, status.IsPaid())
	})
}

func TestStatus_IsFailed(t *testing.T) {
	t.Run("should return true when status is failed", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_FAILED)

		result := status.IsFailed()

		assert.True(t, result)
	})

	t.Run("should return false when status is not failed", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PENDING)

		result := status.IsFailed()

		assert.False(t, result)
	})
}

func TestStatus_SetFailed(t *testing.T) {
	t.Run("should set status to failed", func(t *testing.T) {
		status, _ := NewStatus(constants.PAYMENT_STATUS_PENDING)

		status.SetFailed()

		assert.Equal(t, constants.PAYMENT_STATUS_FAILED, status.Value())
		assert.True(t, status.IsFailed())
	})
}
