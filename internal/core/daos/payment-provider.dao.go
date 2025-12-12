package daos

import "context"

type CreatePIXBillingDAO struct {
	Ctx        context.Context
	ExternalID string
	Amount     float64
}

type PIXBillingResultDAO struct {
	QRData string
}
