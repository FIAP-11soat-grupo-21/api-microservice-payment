package interfaces

import (
	"payment_microservice/internal/core/daos"
)

type IPaymentProvider interface {
	CreatePIXBilling(pixBilling daos.CreatePIXBillingDAO) (daos.PIXBillingResultDAO, error)
}
