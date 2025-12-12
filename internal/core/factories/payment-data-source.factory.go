package factories

import (
	"payment_microservice/internal/core/infra/database/data_sources"
	"payment_microservice/internal/core/interfaces"
)

func NewPaymentDataSource() interfaces.IPaymentDataSource {
	return data_sources.NewGormPaymentDataSource()
}
