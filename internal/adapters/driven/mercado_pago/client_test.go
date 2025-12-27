package mercado_pago

import (
	"context"
	"errors"
	"io"
	"net/http"
	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/common/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	cleanup := env.SetupTestEnv(t)
	defer cleanup()

	t.Run("should create client successfully with valid config", func(t *testing.T) {
		client := NewClient()

		assert.NotNil(t, client)
		assert.NotNil(t, client.cfg)
		assert.Equal(t, "TEST_ACCESS_TOKEN", client.cfg.AccessToken)
		assert.Equal(t, int64(1234567890), client.cfg.CollectorID)
		assert.Equal(t, "test_pos_id", client.cfg.ExternalPosID)
		assert.Equal(t, "https://api.test.com", client.cfg.ApiBaseURL)
	})
}

func TestNewClientWithHTTPClient(t *testing.T) {
	t.Run("should create client with custom http client", func(t *testing.T) {
		cfg := &MercadoPagoConfig{
			AccessToken:   "test_token",
			CollectorID:   123456,
			ExternalPosID: "POS001",
			ApiBaseURL:    "https://api.test.com",
		}

		mockHTTP := mocks.NewMockHTTPClient(nil)
		client := NewClientWithHTTPClient(cfg, mockHTTP)

		assert.NotNil(t, client)
		assert.Equal(t, cfg, client.cfg)
		assert.Equal(t, mockHTTP, client.httpClient)
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("should return client config", func(t *testing.T) {
		cfg := &MercadoPagoConfig{
			AccessToken:   "test_token",
			CollectorID:   123456,
			ExternalPosID: "POS001",
			ApiBaseURL:    "https://api.test.com",
		}

		client := &MercadoPagoClient{cfg: cfg}
		returnedCfg := client.GetConfig()

		assert.Equal(t, cfg, returnedCfg)
	})
}

func TestDoRequest(t *testing.T) {
	ctx := context.Background()
	cfg := &MercadoPagoConfig{
		AccessToken:   "test_token",
		CollectorID:   123456,
		ExternalPosID: "POS001",
		ApiBaseURL:    "https://api.test.com",
	}

	t.Run("should execute request successfully", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://api.test.com/test/path", req.URL.String())
			assert.Equal(t, "POST", req.Method)
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer test_token", req.Header.Get("Authorization"))
			assert.Equal(t, "idempotency-123", req.Header.Get("X-Idempotency-Key"))

			return mocks.CreateSuccessResponse(`{"success": true}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		body := []byte(`{"test": "data"}`)

		result, err := client.DoRequest(ctx, "POST", "/test/path", body, "idempotency-123")

		assert.NoError(t, err)
		assert.Equal(t, []byte(`{"success": true}`), result)
	})

	t.Run("should return error when request creation fails", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(nil)
		client := NewClientWithHTTPClient(cfg, mockHTTP)

		ctxCanceled, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := client.DoRequest(ctxCanceled, "INVALID\nMETHOD", "/test/path", nil, "key")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create request")
	})

	t.Run("should return error when http client fails", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("connection error")
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		body := []byte(`{"test": "data"}`)

		_, err := client.DoRequest(ctx, "POST", "/test/path", body, "key")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute request")
	})

	t.Run("should return error for non-2xx status code", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateErrorResponse(400, `{"error": "bad request"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		body := []byte(`{"test": "data"}`)

		_, err := client.DoRequest(ctx, "POST", "/test/path", body, "key")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request failed with status code 400")
	})

	t.Run("should return error for 500 status code", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateErrorResponse(500, `{"error": "internal server error"}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		body := []byte(`{"test": "data"}`)

		_, err := client.DoRequest(ctx, "POST", "/test/path", body, "key")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request failed with status code 500")
	})

	t.Run("should return error when reading response body fails", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&errorReader{}),
			}, nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		body := []byte(`{"test": "data"}`)

		_, err := client.DoRequest(ctx, "POST", "/test/path", body, "key")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read response body")
	})

	t.Run("should handle empty request body", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateSuccessResponse(`{"success": true}`), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)

		result, err := client.DoRequest(ctx, "GET", "/test/path", nil, "key")

		assert.NoError(t, err)
		assert.Equal(t, []byte(`{"success": true}`), result)
	})

	t.Run("should handle empty response body", func(t *testing.T) {
		mockHTTP := mocks.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
			return mocks.CreateSuccessResponse(""), nil
		})

		client := NewClientWithHTTPClient(cfg, mockHTTP)
		body := []byte(`{"test": "data"}`)

		result, err := client.DoRequest(ctx, "POST", "/test/path", body, "key")

		assert.NoError(t, err)
		assert.Equal(t, []byte(""), result)
	})
}

type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read error")
}
