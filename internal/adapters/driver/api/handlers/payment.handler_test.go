package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/value_objects"
	"payment_microservice/internal/core/dto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestNewPaymentHandler(t *testing.T) {
	cleanup := env.SetupTestEnv()
	defer cleanup()

	handler := NewPaymentHandler()

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.repository)
	assert.NotNil(t, handler.gateway)
	assert.NotNil(t, handler.queuePublisher)
}

func TestCreatePixBilling_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	router.POST("/payments/pix", handler.CreatePixBilling)

	qrCode := "00020101021243650016COM.MERCADOLIBRE"

	mockGateway.On("CreatePIXBilling", mock.MatchedBy(func(dto dto.CreatePIXBillingDTO) bool {
		return dto.Amount == 100.50
	})).Return(dto.PIXBillingResultDTO{
		QRData: qrCode,
	}, nil)

	mockRepo.On("Insert", mock.Anything, mock.MatchedBy(func(payment entities.Payment) bool {
		return payment.OrderID == "order-123"
	})).Return(nil)

	requestBody := map[string]interface{}{
		"order_id": "order-123",
		"amount":   100.50,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/payments/pix", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "order-123", response["order_id"])
	assert.Equal(t, 100.50, response["amount"])
	assert.Equal(t, constants.PAYMENT_STATUS_PENDING, response["status"])
	assert.Equal(t, constants.PIX_PAYMENT_METHOD, response["method"])
	assert.NotEmpty(t, response["qr_code"])

	mockGateway.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreatePixBilling_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	handler := &PaymentHandler{
		repository:     new(mocks.MockPaymentRepository),
		gateway:        new(mocks.MockPaymentGateway),
		queuePublisher: new(mocks.MockQueuePublisher),
	}

	router.POST("/payments/pix", handler.CreatePixBilling)

	req, _ := http.NewRequest("POST", "/payments/pix", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePixBilling_GatewayError(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	router.POST("/payments/pix", handler.CreatePixBilling)

	mockGateway.On("CreatePIXBilling", mock.Anything).
		Return(dto.PIXBillingResultDTO{}, errors.New("gateway error"))

	requestBody := map[string]interface{}{
		"order_id": "order-123",
		"amount":   100.50,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/payments/pix", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockGateway.AssertExpectations(t)
}

func TestConfirmPayment_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	router.POST("/payments/webhook", handler.ConfirmPayment)

	amount, _ := value_objects.NewAmount(100.50)
	status, _ := value_objects.NewStatus(constants.PAYMENT_STATUS_PENDING)
	method, _ := value_objects.NewMethod(constants.PIX_PAYMENT_METHOD)

	payment := entities.Payment{
		ID:        "payment-123",
		OrderID:   "12345678",
		Amount:    amount,
		Status:    status,
		Method:    method,
		CreatedAt: time.Now(),
	}

	mockRepo.On("FindByOrderID", mock.Anything, "12345678").Return(payment, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockQueuePublisher.On("PublishOnTopic", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	requestBody := map[string]interface{}{
		"id":           "webhook-id",
		"live_mode":    true,
		"type":         "payment",
		"date_created": "2024-01-01T00:00:00Z",
		"user_id":      123,
		"api_version":  "v1",
		"action":       "payment.updated",
		"data": map[string]interface{}{
			"id": "12345678",
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockRepo.AssertExpectations(t)
	mockQueuePublisher.AssertExpectations(t)
}

func TestConfirmPayment_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	handler := &PaymentHandler{
		repository:     new(mocks.MockPaymentRepository),
		gateway:        new(mocks.MockPaymentGateway),
		queuePublisher: new(mocks.MockQueuePublisher),
	}

	router.POST("/payments/webhook", handler.ConfirmPayment)

	req, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestConfirmPayment_PaymentNotFound(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	router.POST("/payments/webhook", handler.ConfirmPayment)

	mockRepo.On("FindByOrderID", mock.Anything, "12345678").
		Return(entities.Payment{}, &exceptions.PaymentNotFoundException{Message: "payment not found"})

	requestBody := map[string]interface{}{
		"id":           "webhook-id",
		"live_mode":    true,
		"type":         "payment",
		"date_created": "2024-01-01T00:00:00Z",
		"user_id":      123,
		"api_version":  "v1",
		"action":       "payment.updated",
		"data": map[string]interface{}{
			"id": "12345678",
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestFindByOrderID_Success(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	router.GET("/payments/order/:orderID", handler.FindByOrderID)

	amount, _ := value_objects.NewAmount(100.50)
	status, _ := value_objects.NewStatus(constants.PAYMENT_STATUS_PENDING)
	method, _ := value_objects.NewMethod(constants.PIX_PAYMENT_METHOD)
	qrCode, _ := value_objects.NewQRCode("00020101021243650016COM.MERCADOLIBRE")

	payment := entities.Payment{
		ID:        "payment-123",
		OrderID:   "order-123",
		Amount:    amount,
		Status:    status,
		Method:    method,
		QRCode:    &qrCode,
		CreatedAt: time.Now(),
	}

	mockRepo.On("FindByOrderID", mock.Anything, "order-123").Return(payment, nil)

	req, _ := http.NewRequest("GET", "/payments/order/order-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "payment-123", response["id"])
	assert.Equal(t, "order-123", response["order_id"])
	assert.Equal(t, 100.50, response["amount"])
	assert.Equal(t, constants.PAYMENT_STATUS_PENDING, response["status"])
	assert.Equal(t, constants.PIX_PAYMENT_METHOD, response["method"])

	mockRepo.AssertExpectations(t)
}

func TestFindByOrderID_NotFound(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	// Adicionar middleware de erro para testar ctx.Error
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": c.Errors[0].Error()})
		}
	})

	router.GET("/payments/order/:orderID", handler.FindByOrderID)

	mockRepo.On("FindByOrderID", mock.Anything, "order-999").
		Return(entities.Payment{}, &exceptions.PaymentNotFoundException{Message: "payment not found"})

	req, _ := http.NewRequest("GET", "/payments/order/order-999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestFindByOrderID_WithoutQRCode(t *testing.T) {
	router := setupTestRouter()
	mockRepo := new(mocks.MockPaymentRepository)
	mockGateway := new(mocks.MockPaymentGateway)
	mockQueuePublisher := new(mocks.MockQueuePublisher)

	handler := &PaymentHandler{
		repository:     mockRepo,
		gateway:        mockGateway,
		queuePublisher: mockQueuePublisher,
	}

	router.GET("/payments/order/:orderID", handler.FindByOrderID)

	amount, _ := value_objects.NewAmount(100.50)
	status, _ := value_objects.NewStatus(constants.PAYMENT_STATUS_PAID)
	method, _ := value_objects.NewMethod(constants.PIX_PAYMENT_METHOD)
	paidAt := time.Now()

	payment := entities.Payment{
		ID:        "payment-123",
		OrderID:   "order-123",
		Amount:    amount,
		Status:    status,
		Method:    method,
		QRCode:    nil,
		CreatedAt: time.Now(),
		PaidAt:    &paidAt,
	}

	mockRepo.On("FindByOrderID", mock.Anything, "order-123").Return(payment, nil)

	req, _ := http.NewRequest("GET", "/payments/order/order-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "payment-123", response["id"])
	// qr_code com omitempty não aparece quando nil
	_, hasQRCode := response["qr_code"]
	assert.False(t, hasQRCode)
	assert.NotNil(t, response["paid_at"])

	mockRepo.AssertExpectations(t)
}
