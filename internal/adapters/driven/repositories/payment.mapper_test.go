package repositories

import (
	"testing"
	"time"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/value_objects"

	"github.com/stretchr/testify/assert"
)

func TestToDomain(t *testing.T) {
	t.Run("should convert PaymentModel to Payment entity successfully", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		qrCodeURL := "00020101021243650016COM.MERCADOLIBRE"
		now := time.Now()
		paidAt := time.Now().Add(5 * time.Minute)

		paymentModel := PaymentModel{
			ID:            id,
			OrderID:       orderID,
			Amount:        amount,
			Status:        status,
			PaymentMethod: method,
			QRCodeURL:     &qrCodeURL,
			PaidAt:        &paidAt,
			CreatedAt:     now,
		}

		payment, err := toDomain(paymentModel)

		assert.NoError(t, err)
		assert.Equal(t, id, payment.ID)
		assert.Equal(t, orderID, payment.OrderID)
		assert.Equal(t, amount, payment.Amount.Value())
		assert.Equal(t, status, payment.Status.Value())
		assert.Equal(t, method, payment.Method.Value())
		assert.NotNil(t, payment.QRCode)
		assert.Equal(t, qrCodeURL, payment.QRCode.Value())
		assert.NotNil(t, payment.PaidAt)
		assert.Equal(t, paidAt, *payment.PaidAt)
		assert.Equal(t, now, payment.CreatedAt)
	})

	t.Run("should convert PaymentModel with nil optional fields", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		paymentModel := PaymentModel{
			ID:            id,
			OrderID:       orderID,
			Amount:        amount,
			Status:        status,
			PaymentMethod: method,
			QRCodeURL:     nil,
			PaidAt:        nil,
			CreatedAt:     now,
		}

		payment, err := toDomain(paymentModel)

		assert.NoError(t, err)
		assert.Equal(t, id, payment.ID)
		assert.Equal(t, orderID, payment.OrderID)
		assert.Nil(t, payment.QRCode)
		assert.Nil(t, payment.PaidAt)
	})

	t.Run("should convert PaymentModel with PAID status", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PAID
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		paymentModel := PaymentModel{
			ID:            id,
			OrderID:       orderID,
			Amount:        amount,
			Status:        status,
			PaymentMethod: method,
			CreatedAt:     now,
		}

		payment, err := toDomain(paymentModel)

		assert.NoError(t, err)
		assert.Equal(t, status, payment.Status.Value())
	})

	t.Run("should convert PaymentModel with FAILED status", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_FAILED
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		paymentModel := PaymentModel{
			ID:            id,
			OrderID:       orderID,
			Amount:        amount,
			Status:        status,
			PaymentMethod: method,
			CreatedAt:     now,
		}

		payment, err := toDomain(paymentModel)

		assert.NoError(t, err)
		assert.Equal(t, status, payment.Status.Value())
	})

	t.Run("should return error when PaymentModel has invalid amount", func(t *testing.T) {
		paymentModel := PaymentModel{
			ID:            "payment-123",
			OrderID:       "order-456",
			Amount:        -10.0,
			Status:        constants.PAYMENT_STATUS_PENDING,
			PaymentMethod: constants.PIX_PAYMENT_METHOD,
			CreatedAt:     time.Now(),
		}

		_, err := toDomain(paymentModel)

		assert.Error(t, err)
	})

	t.Run("should return error when PaymentModel has invalid status", func(t *testing.T) {
		paymentModel := PaymentModel{
			ID:            "payment-123",
			OrderID:       "order-456",
			Amount:        100.0,
			Status:        "INVALID_STATUS",
			PaymentMethod: constants.PIX_PAYMENT_METHOD,
			CreatedAt:     time.Now(),
		}

		_, err := toDomain(paymentModel)

		assert.Error(t, err)
	})

	t.Run("should return error when PaymentModel has invalid method", func(t *testing.T) {
		paymentModel := PaymentModel{
			ID:            "payment-123",
			OrderID:       "order-456",
			Amount:        100.0,
			Status:        constants.PAYMENT_STATUS_PENDING,
			PaymentMethod: "INVALID_METHOD",
			CreatedAt:     time.Now(),
		}

		_, err := toDomain(paymentModel)

		assert.Error(t, err)
	})

	t.Run("should return error when PaymentModel has invalid QR code", func(t *testing.T) {
		invalidQRCode := ""
		paymentModel := PaymentModel{
			ID:            "payment-123",
			OrderID:       "order-456",
			Amount:        100.0,
			Status:        constants.PAYMENT_STATUS_PENDING,
			PaymentMethod: constants.PIX_PAYMENT_METHOD,
			QRCodeURL:     &invalidQRCode,
			CreatedAt:     time.Now(),
		}

		_, err := toDomain(paymentModel)

		assert.Error(t, err)
	})
}

func TestToPersistence(t *testing.T) {
	t.Run("should convert Payment entity to PaymentModel successfully", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PAID
		method := constants.PIX_PAYMENT_METHOD
		qrCodeData := "00020101021243650016COM.MERCADOLIBRE"
		now := time.Now()
		paidAt := time.Now().Add(5 * time.Minute)

		qrCode, _ := value_objects.NewQRCode(qrCodeData)
		paymentAmount, _ := value_objects.NewAmount(amount)
		paymentStatus, _ := value_objects.NewStatus(status)
		paymentMethod, _ := value_objects.NewMethod(method)

		payment := entities.Payment{
			ID:        id,
			OrderID:   orderID,
			Amount:    paymentAmount,
			Status:    paymentStatus,
			Method:    paymentMethod,
			QRCode:    &qrCode,
			CreatedAt: now,
			PaidAt:    &paidAt,
		}

		paymentModel := toPersistence(payment)

		assert.Equal(t, id, paymentModel.ID)
		assert.Equal(t, orderID, paymentModel.OrderID)
		assert.Equal(t, amount, paymentModel.Amount)
		assert.Equal(t, status, paymentModel.Status)
		assert.Equal(t, method, paymentModel.PaymentMethod)
		assert.NotNil(t, paymentModel.QRCodeURL)
		assert.Equal(t, qrCodeData, *paymentModel.QRCodeURL)
		assert.NotNil(t, paymentModel.PaidAt)
		assert.Equal(t, paidAt, *paymentModel.PaidAt)
		assert.Equal(t, now, paymentModel.CreatedAt)
	})

	t.Run("should convert Payment entity with nil QRCode", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		paymentAmount, _ := value_objects.NewAmount(amount)
		paymentStatus, _ := value_objects.NewStatus(status)
		paymentMethod, _ := value_objects.NewMethod(method)

		payment := entities.Payment{
			ID:        id,
			OrderID:   orderID,
			Amount:    paymentAmount,
			Status:    paymentStatus,
			Method:    paymentMethod,
			QRCode:    nil,
			CreatedAt: now,
			PaidAt:    nil,
		}

		paymentModel := toPersistence(payment)

		assert.Equal(t, id, paymentModel.ID)
		assert.Equal(t, orderID, paymentModel.OrderID)
		assert.Equal(t, amount, paymentModel.Amount)
		assert.Equal(t, status, paymentModel.Status)
		assert.Equal(t, method, paymentModel.PaymentMethod)
		assert.Nil(t, paymentModel.QRCodeURL)
		assert.Nil(t, paymentModel.PaidAt)
		assert.Equal(t, now, paymentModel.CreatedAt)
	})

	t.Run("should convert Payment entity with nil PaidAt", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-456"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		paymentAmount, _ := value_objects.NewAmount(amount)
		paymentStatus, _ := value_objects.NewStatus(status)
		paymentMethod, _ := value_objects.NewMethod(method)

		payment := entities.Payment{
			ID:        id,
			OrderID:   orderID,
			Amount:    paymentAmount,
			Status:    paymentStatus,
			Method:    paymentMethod,
			CreatedAt: now,
			PaidAt:    nil,
		}

		paymentModel := toPersistence(payment)

		assert.Nil(t, paymentModel.PaidAt)
	})

	t.Run("should extract QRCode value correctly", func(t *testing.T) {
		qrCodeData := "00020101021243650016COM.MERCADOLIBRE"
		qrCode, _ := value_objects.NewQRCode(qrCodeData)

		paymentAmount, _ := value_objects.NewAmount(100.0)
		paymentStatus, _ := value_objects.NewStatus(constants.PAYMENT_STATUS_PENDING)
		paymentMethod, _ := value_objects.NewMethod(constants.PIX_PAYMENT_METHOD)

		payment := entities.Payment{
			ID:        "payment-123",
			OrderID:   "order-456",
			Amount:    paymentAmount,
			Status:    paymentStatus,
			Method:    paymentMethod,
			QRCode:    &qrCode,
			CreatedAt: time.Now(),
		}

		paymentModel := toPersistence(payment)

		assert.NotNil(t, paymentModel.QRCodeURL)
		assert.Equal(t, qrCodeData, *paymentModel.QRCodeURL)
	})

	t.Run("should convert Payment with all value objects", func(t *testing.T) {
		paymentAmount, _ := value_objects.NewAmount(200.0)
		paymentStatus, _ := value_objects.NewStatus(constants.PAYMENT_STATUS_FAILED)
		paymentMethod, _ := value_objects.NewMethod(constants.PIX_PAYMENT_METHOD)
		qrCode, _ := value_objects.NewQRCode("test-qr-code")

		payment := entities.Payment{
			ID:        "payment-999",
			OrderID:   "order-999",
			Amount:    paymentAmount,
			Status:    paymentStatus,
			Method:    paymentMethod,
			QRCode:    &qrCode,
			CreatedAt: time.Now(),
		}

		paymentModel := toPersistence(payment)

		assert.Equal(t, 200.0, paymentModel.Amount)
		assert.Equal(t, constants.PAYMENT_STATUS_FAILED, paymentModel.Status)
		assert.Equal(t, constants.PIX_PAYMENT_METHOD, paymentModel.PaymentMethod)
		assert.NotNil(t, paymentModel.QRCodeURL)
	})
}
