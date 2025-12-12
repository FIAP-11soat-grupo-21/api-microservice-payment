package factories

import (
	"payment_microservice/internal/core/infra/provider/mercado_pago"
	"payment_microservice/internal/core/interfaces"
)

func NewPaymentProvider() interfaces.IPaymentProvider {
	return mercado_pago.NewMercadoPagoProvider()
}
