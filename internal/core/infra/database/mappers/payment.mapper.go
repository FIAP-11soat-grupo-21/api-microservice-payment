package mappers

import (
	"payment_microservice/internal/core/daos"
	"payment_microservice/internal/core/infra/database/models"
)

func FromDAOToModel(paymentDAO daos.PaymentDAO) models.PaymentModel {
	return models.PaymentModel{
		ID:              paymentDAO.ID,
		OrderID:         paymentDAO.OrderID,
		Amount:          paymentDAO.Amount,
		Status:          paymentDAO.Status,
		PaymentMethod:   paymentDAO.Method,
		TransactionCode: paymentDAO.TransactionCode,
		QRCodeURL:       paymentDAO.QRCode,
		PaidAt:          paymentDAO.PaidAt,
		CreatedAt:       paymentDAO.CreatedAt,
	}
}

func FromModelToOrderDAO(p *models.PaymentModel) daos.PaymentDAO {

	return daos.PaymentDAO{
		ID:              p.ID,
		OrderID:         p.OrderID,
		Amount:          p.Amount,
		Status:          p.Status,
		Method:          p.PaymentMethod,
		TransactionCode: p.TransactionCode,
		QRCode:          p.QRCodeURL,
		PaidAt:          p.PaidAt,
		CreatedAt:       p.CreatedAt,
	}
}
