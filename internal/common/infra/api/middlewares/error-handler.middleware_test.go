package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"payment_microservice/internal/core/domain/exceptions"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlerMiddleware())
	return router
}

func TestErrorHandlerMiddleware(t *testing.T) {
	t.Run("should handle domain error PaymentNotFoundException", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			err := &exceptions.PaymentNotFoundException{
				Message: "payment not found",
			}
			c.Error(err)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "payment not found")
	})

	t.Run("should handle domain error InvalidPaymentDataException", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			err := &exceptions.InvalidPaymentDataException{
				Message: "invalid data",
			}
			c.Error(err)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid data")
	})

	t.Run("should handle non-domain error with internal server error", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			err := errors.New("generic error")
			c.Error(err)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Internal server error"}`, w.Body.String())
	})

	t.Run("should not handle when no errors", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("should handle multiple errors and use last one", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			c.Error(errors.New("first error"))
			c.Error(errors.New("second error"))
			c.Error(&exceptions.PaymentNotFoundException{
				Message: "last error",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "last error")
	})

	t.Run("should call next and process handlers before checking errors", func(t *testing.T) {
		router := setupTestRouter()

		handlerCalled := false
		router.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.Error(errors.New("test error"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, handlerCalled)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should handle error and abort after middleware completes", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			c.Error(&exceptions.InvalidPaymentDataException{
				Message: "validation error",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation error")
	})

	t.Run("should call next when no errors", func(t *testing.T) {
		router := setupTestRouter()

		nextCalled := false

		router.GET("/test", func(c *gin.Context) {
			nextCalled = true
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, nextCalled)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should verify JSON format for internal server error", func(t *testing.T) {
		router := setupTestRouter()

		router.GET("/test", func(c *gin.Context) {
			c.Error(errors.New("database connection failed"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Internal server error"}`, w.Body.String())
	})

	t.Run("should process all handlers before middleware checks errors", func(t *testing.T) {
		router := setupTestRouter()

		var executionOrder []string

		router.GET("/test", func(c *gin.Context) {
			executionOrder = append(executionOrder, "handler")
			c.Error(&exceptions.PaymentNotFoundException{
				Message: "not found",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, []string{"handler"}, executionOrder)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
