package factory

import (
	"payment_microservice/internal/adapters/driven/repositories"
	"payment_microservice/internal/common/infra/database"
	"payment_microservice/internal/core/domain/ports"
)

func NewPaymentRepository() ports.IPaymentRepository {
	return repositories.NewGormPaymentDataSource(database.GetDB())
}
