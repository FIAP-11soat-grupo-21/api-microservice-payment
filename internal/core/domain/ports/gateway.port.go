package ports

import (
	"payment_microservice/internal/core/dto"
)

type IPaymentGateway interface {
	CreatePIXBilling(pixBilling dto.CreatePIXBillingDTO) (dto.PIXBillingResultDTO, error)
}
