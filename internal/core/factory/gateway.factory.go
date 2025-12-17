package factory

import (
	"payment_microservice/internal/adapters/driven/mercado_pago"
	"payment_microservice/internal/core/domain/ports"
)

func NewPaymentGateway() ports.IPaymentGateway {
	return mercado_pago.NewMercadoPagoGateway()
}
