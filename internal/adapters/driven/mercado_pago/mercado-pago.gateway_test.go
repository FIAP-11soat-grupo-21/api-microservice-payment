package mercado_pago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/mocks"
	dtos "payment_microservice/internal/core/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMercadoPagoGateway(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	t.Run("should create gateway successfully", func(t *testing.T) {
		gateway := NewMercadoPagoGateway()

		assert.NotNil(t, gateway)
		assert.NotNil(t, gateway.client)
		assert.NotNil(t, gateway.cfg)
		assert.Equal(t, int64(1234567890), gateway.cfg.CollectorID)
	})
}

func TestNewMercadoPagoGatewayWithClient(t *testing.T) {
	t.Run("should create gateway with custom client", func(t *testing.T) {
		cfg := &MercadoPagoConfig{
			AccessToken:   "test_token",
			CollectorID:   123456,
			ExternalPosID: "POS001",
			ApiBaseURL:    "https://api.test.com",
		}

		mockHTTP := mocks.NewMockHTTPClient(nil)
		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		assert.NotNil(t, gateway)
		assert.Equal(t, client, gateway.client)
		assert.Equal(t, cfg, gateway.cfg)
	})
}

func TestCreatePIXBilling(t *testing.T) {
	cfg := &MercadoPagoConfig{
		AccessToken:   "test_token",
		CollectorID:   123456,
		ExternalPosID: "POS001",
		ApiBaseURL:    "https://api.test.com",
	}

	t.Run("should create pix billing successfully", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "POST", req.Method)
			assert.Contains(t, req.URL.Path, "/instore/orders/qr/seller/collectors/123456/pos/POS001/qrs")
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer test_token", req.Header.Get("Authorization"))
			assert.Equal(t, "test-external-id-123", req.Header.Get("X-Idempotency-Key"))

			body, _ := io.ReadAll(req.Body)
			assert.Contains(t, string(body), "test-external-id-123")
			assert.Contains(t, string(body), "Pagamento total")
			assert.Contains(t, string(body), "100")

			return mocks.CreateSuccessResponse(`{"qr_data": "00020126580014BR.GOV.BCB.PIX"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-external-id-123",
			Amount:     100.50,
		}

		result, err := gateway.CreatePIXBilling(dto)

		assert.NoError(t, err)
		assert.Equal(t, "00020126580014BR.GOV.BCB.PIX", result.QRData)
	})

	t.Run("should handle short external id", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.Contains(t, string(body), "Pagamento Pedido short")

			return mocks.CreateSuccessResponse(`{"qr_data": "test-qr-data"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "short",
			Amount:     50.00,
		}

		result, err := gateway.CreatePIXBilling(dto)

		assert.NoError(t, err)
		assert.Equal(t, "test-qr-data", result.QRData)
	})

	t.Run("should return error when client request fails", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("connection error")
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-id",
			Amount:     100.00,
		}

		_, err := gateway.CreatePIXBilling(dto)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute request")
	})

	t.Run("should return error when response parsing fails", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateSuccessResponse(`invalid json`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-id",
			Amount:     100.00,
		}

		_, err := gateway.CreatePIXBilling(dto)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao realizar JSON parse da resposta do MP")
	})

	t.Run("should return error for http error status", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateErrorResponse(400, `{"error": "bad request"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-id",
			Amount:     100.00,
		}

		_, err := gateway.CreatePIXBilling(dto)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request failed with status code 400")
	})

	t.Run("should handle zero amount", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.Contains(t, string(body), `"total_amount":0`)

			return mocks.CreateSuccessResponse(`{"qr_data": "test-qr-data"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-id-zero",
			Amount:     0.0,
		}

		result, err := gateway.CreatePIXBilling(dto)

		assert.NoError(t, err)
		assert.Equal(t, "test-qr-data", result.QRData)
	})

	t.Run("should handle large amount", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			assert.Contains(t, string(body), "9999.99")

			return mocks.CreateSuccessResponse(`{"qr_data": "test-qr-data"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-id-large",
			Amount:     9999.99,
		}

		result, err := gateway.CreatePIXBilling(dto)

		assert.NoError(t, err)
		assert.Equal(t, "test-qr-data", result.QRData)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("context canceled")
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        ctx,
			ExternalID: "test-id",
			Amount:     100.00,
		}

		_, err := gateway.CreatePIXBilling(dto)

		assert.Error(t, err)
	})

	t.Run("should verify request body structure", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var requestBody qrCodeRequestBody
			err := json.Unmarshal(bodyBytes, &requestBody)
			assert.NoError(t, err)

			assert.Equal(t, "external-123", requestBody.ExternalReference)
			assert.Equal(t, "Pagamento Pedido external", requestBody.Title)
			assert.Contains(t, requestBody.Description, "R$ 250.00")
			assert.Contains(t, requestBody.Description, "external-123")
			assert.Equal(t, "http://localhost:8080/v1/payments/webhook", requestBody.NotificationURL)
			assert.Equal(t, 250.00, requestBody.TotalAmount)
			assert.Len(t, requestBody.Items, 1)
			assert.Equal(t, "Pagamento total", requestBody.Items[0].Title)
			assert.Equal(t, 250.00, requestBody.Items[0].UnitPrice)
			assert.Equal(t, 1, requestBody.Items[0].Quantity)
			assert.Equal(t, "unit", requestBody.Items[0].UnitMeasure)
			assert.Equal(t, 250.00, requestBody.Items[0].TotalAmount)
			assert.Equal(t, 0.0, requestBody.CashOut.Amount)

			return mocks.CreateSuccessResponse(`{"qr_data": "test-qr-data"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "external-123",
			Amount:     250.00,
		}

		result, err := gateway.CreatePIXBilling(dto)

		assert.NoError(t, err)
		assert.Equal(t, "test-qr-data", result.QRData)
	})

	t.Run("should handle empty qr_data in response", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateSuccessResponse(`{"qr_data": ""}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		gateway := NewMercadoPagoGatewayWithClient(client)

		dto := dtos.CreatePIXBillingDTO{
			Ctx:        context.Background(),
			ExternalID: "test-id",
			Amount:     100.00,
		}

		result, err := gateway.CreatePIXBilling(dto)

		assert.NoError(t, err)
		assert.Equal(t, "", result.QRData)
	})
}
