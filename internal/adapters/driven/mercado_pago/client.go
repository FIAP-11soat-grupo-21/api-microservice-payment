package mercado_pago

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"payment_microservice/internal/common/config/env"
	"strconv"
)

type MercadoPagoConfig struct {
	AccessToken   string
	CollectorID   int64
	ExternalPosID string
	ApiBaseURL    string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MercadoPagoClient struct {
	cfg        *MercadoPagoConfig
	httpClient HTTPClient
}

func (c *MercadoPagoClient) GetConfig() *MercadoPagoConfig {
	return c.cfg
}

func NewClient() *MercadoPagoClient {
	applicationConfig := env.GetConfig()

	collectorID, err := strconv.ParseInt(applicationConfig.MercadoPago.CollectorID, 10, 64)

	if err != nil {
		panic(fmt.Sprintf("Invalid CollectorID: %s", applicationConfig.MercadoPago.CollectorID))
	}

	cfg := &MercadoPagoConfig{
		AccessToken:   applicationConfig.MercadoPago.AccessToken,
		CollectorID:   collectorID,
		ExternalPosID: applicationConfig.MercadoPago.ExternalPosID,
		ApiBaseURL:    applicationConfig.MercadoPago.ApiBaseURL,
	}

	return &MercadoPagoClient{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

func NewClientWithHTTPClient(cfg *MercadoPagoConfig, httpClient HTTPClient) *MercadoPagoClient {
	return &MercadoPagoClient{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (c *MercadoPagoClient) DoRequest(ctx context.Context, method string, path string, requestBody []byte, idempotencyKey string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.cfg.ApiBaseURL, path)

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewBuffer(requestBody),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.AccessToken))
	req.Header.Set("X-Idempotency-Key", idempotencyKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Println(string(responseBody))

		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBody, nil
}
