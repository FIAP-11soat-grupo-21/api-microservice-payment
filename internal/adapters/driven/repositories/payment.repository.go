package repositories

import (
	"context"

	"gorm.io/gorm"

	"payment_microservice/internal/core/domain/entities"
)

type GormPaymentDataSource struct {
	db *gorm.DB
}

func NewGormPaymentDataSource(instance *gorm.DB) *GormPaymentDataSource {
	return &GormPaymentDataSource{
		db: instance,
	}
}

func (r *GormPaymentDataSource) Insert(ctx context.Context, payment entities.Payment) error {
	err := r.db.
		WithContext(ctx).
		Create(toPersistence(payment)).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r *GormPaymentDataSource) Update(ctx context.Context, payment entities.Payment) error {
	var model = toPersistence(payment)

	err := r.db.WithContext(ctx).Save(&model).Error

	if err != nil {
		return err
	}

	return nil
}

func (ds *GormPaymentDataSource) FindByOrderID(ctx context.Context, orderID string) (entities.Payment, error) {
	paymentModel := new(PaymentModel)

	err := ds.db.WithContext(ctx).First(&paymentModel, "order_id = ?", orderID).Error
	if err != nil {
		return entities.Payment{}, err
	}

	return toDomain(*paymentModel)
}

func (ds *GormPaymentDataSource) Delete(ctx context.Context, paymentID string) error {
	err := ds.db.
		WithContext(ctx).
		Delete(&PaymentModel{}, "id = ?", paymentID).
		Error

	if err != nil {
		return err
	}

	return nil
}
