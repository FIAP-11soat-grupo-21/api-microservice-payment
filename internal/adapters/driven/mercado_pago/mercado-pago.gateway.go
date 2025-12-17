package mercado_pago

import (
	"encoding/json"
	"fmt"
	"net/http"
	dtos "payment_microservice/internal/core/dto"
)

type MercadoPagoGateway struct {
	client *MercadoPagoClient
	cfg    *MercadoPagoConfig
}

func NewMercadoPagoGateway() *MercadoPagoGateway {
	client := NewClient()
	cfg := client.GetConfig()

	return &MercadoPagoGateway{
		client: client,
		cfg:    cfg,
	}
}

type qrCodeRequestBody struct {
	ExternalReference string   `json:"external_reference"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	NotificationURL   string   `json:"notification_url"`
	TotalAmount       float64  `json:"total_amount"`
	Items             []qrItem `json:"items"`
	CashOut           cashOut  `json:"cash_out"`
}

type qrItem struct {
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	UnitMeasure string  `json:"unit_measure"`
	TotalAmount float64 `json:"total_amount"`
}

type cashOut struct {
	Amount float64 `json:"amount"`
}

type QRCodeAPIResponse struct {
	QRData string `json:"qr_data"`
}

func (c *MercadoPagoGateway) CreatePIXBilling(pixBilling dtos.CreatePIXBillingDTO) (dtos.PIXBillingResultDTO, error) {
	url := fmt.Sprintf(
		"/instore/orders/qr/seller/collectors/%d/pos/%s/qrs",
		c.cfg.CollectorID,
		c.cfg.ExternalPosID,
	)

	body := qrCodeRequestBody{
		ExternalReference: pixBilling.ExternalID,
		Title:             fmt.Sprintf("Pagamento Pedido %s", pixBilling.ExternalID[:8]),
		Description:       fmt.Sprintf("Pagamento de R$ %.2f para a ordem %s", pixBilling.Amount, pixBilling.ExternalID),
		NotificationURL:   "https://tech-challenge.com.br", // Sem definição temporariamente
		TotalAmount:       pixBilling.Amount,
		Items: []qrItem{
			{
				Title:       "Pagamento total",
				UnitPrice:   pixBilling.Amount,
				Quantity:    1,
				UnitMeasure: "unit",
				TotalAmount: pixBilling.Amount,
			},
		},
		CashOut: cashOut{Amount: 0},
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		return dtos.PIXBillingResultDTO{}, fmt.Errorf("erro ao serializar o corpo da requisição: %w", err)
	}

	resp, err := c.client.DoRequest(
		pixBilling.Ctx,
		http.MethodPost,
		url,
		jsonBody,
		pixBilling.ExternalID, // Usando ExternalID como chave de idempotência
	)

	if err != nil {
		return dtos.PIXBillingResultDTO{}, err
	}

	var apiResp QRCodeAPIResponse

	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return dtos.PIXBillingResultDTO{}, fmt.Errorf("erro ao realizar JSON parse da resposta do MP: %w", err)
	}

	return dtos.PIXBillingResultDTO{
		QRData: apiResp.QRData,
	}, nil
}
