package http_errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"payment_microservice/internal/core/domain/exceptions"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	return c, w
}

func TestHandleDomainErrors(t *testing.T) {
	t.Run("should handle PaymentNotFoundException and return 404", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := &exceptions.PaymentNotFoundException{
			Message: "payment with order ID order-123 not found",
		}

		result := HandleDomainErrors(err, ctx)

		assert.True(t, result)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "payment with order ID order-123 not found")
	})

	t.Run("should handle PaymentNotFoundException with default message", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := &exceptions.PaymentNotFoundException{}

		result := HandleDomainErrors(err, ctx)

		assert.True(t, result)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "payment not found")
	})

	t.Run("should handle InvalidPaymentDataException and return 400", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := &exceptions.InvalidPaymentDataException{
			Message: "invalid amount value",
		}

		result := HandleDomainErrors(err, ctx)

		assert.True(t, result)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid amount value")
	})

	t.Run("should handle InvalidPaymentDataException with default message", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := &exceptions.InvalidPaymentDataException{}

		result := HandleDomainErrors(err, ctx)

		assert.True(t, result)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid payment data")
	})

	t.Run("should return false for non-domain errors", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := errors.New("generic error")

		result := HandleDomainErrors(err, ctx)

		assert.False(t, result)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return false for nil error", func(t *testing.T) {
		ctx, w := setupTestContext()

		result := HandleDomainErrors(nil, ctx)

		assert.False(t, result)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should verify json response format for PaymentNotFoundException", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := &exceptions.PaymentNotFoundException{
			Message: "payment not found by ID",
		}

		result := HandleDomainErrors(err, ctx)

		assert.True(t, result)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error": "payment not found by ID"}`, w.Body.String())
	})

	t.Run("should verify json response format for InvalidPaymentDataException", func(t *testing.T) {
		ctx, w := setupTestContext()

		err := &exceptions.InvalidPaymentDataException{
			Message: "amount must be positive",
		}

		result := HandleDomainErrors(err, ctx)

		assert.True(t, result)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "amount must be positive"}`, w.Body.String())
	})
}
