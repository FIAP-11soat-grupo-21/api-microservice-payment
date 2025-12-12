package data_sources

import (
	"context"

	"gorm.io/gorm"

	"payment_microservice/internal/common/infra/database"
	"payment_microservice/internal/core/daos"
	"payment_microservice/internal/core/infra/database/mappers"
	"payment_microservice/internal/core/infra/database/models"
)

type GormPaymentDataSource struct {
	db *gorm.DB
}

func NewGormPaymentDataSource() *GormPaymentDataSource {
	return &GormPaymentDataSource{
		db: database.GetDB(),
	}
}

func (r *GormPaymentDataSource) Insert(ctx context.Context, paymentDAO daos.PaymentDAO) error {
	err := r.db.
		WithContext(ctx).
		Create(mappers.FromDAOToModel(paymentDAO)).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r *GormPaymentDataSource) Update(ctx context.Context, payment daos.PaymentDAO) error {
	var model = mappers.FromDAOToModel(payment)

	err := r.db.WithContext(ctx).Save(&model).Error

	if err != nil {
		return err
	}

	return nil
}

func (ds *GormPaymentDataSource) FindByOrderID(orderID string) (daos.PaymentDAO, error) {
	var paymentModel *models.PaymentModel

	err := ds.db.First(&paymentModel, "order_id = ?", orderID).Error
	if err != nil {
		return daos.PaymentDAO{}, err
	}

	return mappers.FromModelToOrderDAO(paymentModel), nil
}
