package entities

import (
	"testing"
	"time"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/domain/value_objects"

	"github.com/stretchr/testify/assert"
)

// Helper para criar payment com método específico (apenas para testes de validação)
func createPaymentWithMethod(method value_objects.Method) *Payment {
	amount, _ := value_objects.NewAmount(100.0)
	status := value_objects.NewStatusDefault()

	return &Payment{
		ID:        "test-123",
		OrderID:   "order-123",
		Amount:    amount,
		Status:    status,
		Method:    method,
		CreatedAt: time.Now(),
	}
}

func TestNewPaymentDefault(t *testing.T) {
	t.Run("should create payment with PIX method successfully", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, err := NewPaymentDefault(id, orderID, amount, method, qrData)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, id, payment.ID)
		assert.Equal(t, orderID, payment.OrderID)
		assert.Equal(t, amount, payment.Amount.Value())
		assert.Equal(t, method, payment.Method.Value())
		assert.True(t, payment.Status.IsPending())
		assert.Nil(t, payment.QRCode)
		assert.Nil(t, payment.TransactionCode)
		assert.Nil(t, payment.PaidAt)
		assert.False(t, payment.CreatedAt.IsZero())
	})

	t.Run("should create payment with PIX and QR code", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		qrData := "00020126580014br.gov.bcb.pix"

		payment, err := NewPaymentDefault(id, orderID, amount, method, &qrData)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.NotNil(t, payment.QRCode)
		assert.Equal(t, qrData, payment.QRCode.Value())
	})

	t.Run("should return error for invalid method", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := "credit_card"
		var qrData *string

		payment, err := NewPaymentDefault(id, orderID, amount, method, qrData)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error for non-PIX method", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := "debit_card"
		var qrData *string

		payment, err := NewPaymentDefault(id, orderID, amount, method, qrData)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "invalid payment method")
	})

	t.Run("should return error for negative amount", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := -10.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, err := NewPaymentDefault(id, orderID, amount, method, qrData)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error for invalid QR code", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		qrData := ""

		payment, err := NewPaymentDefault(id, orderID, amount, method, &qrData)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should create payment with zero amount", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 0.0
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, err := NewPaymentDefault(id, orderID, amount, method, qrData)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, 0.0, payment.Amount.Value())
	})

	t.Run("should create payment with empty orderID", func(t *testing.T) {
		id := "payment-123"
		orderID := ""
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, err := NewPaymentDefault(id, orderID, amount, method, qrData)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, "", payment.OrderID)
	})
}

func TestNewPayment(t *testing.T) {
	t.Run("should create payment with all fields", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PAID
		method := constants.PIX_PAYMENT_METHOD
		transactionCode := "txn-123"
		qrCode := "00020126580014br.gov.bcb.pix"
		now := time.Now()
		paidAt := now

		payment, err := NewPayment(id, orderID, amount, status, method, &transactionCode, &qrCode, &paidAt, now)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, id, payment.ID)
		assert.Equal(t, orderID, payment.OrderID)
		assert.Equal(t, amount, payment.Amount.Value())
		assert.Equal(t, status, payment.Status.Value())
		assert.Equal(t, method, payment.Method.Value())
		assert.Equal(t, &transactionCode, payment.TransactionCode)
		assert.NotNil(t, payment.QRCode)
		assert.Equal(t, qrCode, payment.QRCode.Value())
		assert.Equal(t, &paidAt, payment.PaidAt)
		assert.Equal(t, now, payment.CreatedAt)
	})

	t.Run("should create payment without optional fields", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, err := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Nil(t, payment.TransactionCode)
		assert.Nil(t, payment.QRCode)
		assert.Nil(t, payment.PaidAt)
	})

	t.Run("should return error for invalid method", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := "invalid"
		now := time.Now()

		payment, err := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error for invalid status", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := "invalid"
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, err := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error for invalid amount", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := -100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, err := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error for invalid QR code", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		qrCode := ""
		now := time.Now()

		payment, err := NewPayment(id, orderID, amount, status, method, nil, &qrCode, nil, now)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should create payment with failed status", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_FAILED
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, err := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.True(t, payment.Status.IsFailed())
	})
}

func TestPayment_MarkAsPaid(t *testing.T) {
	t.Run("should mark pending payment as paid", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, _ := NewPaymentDefault(id, orderID, amount, method, qrData)
		transactionCode := "txn-123"

		payment.MarkAsPaid(transactionCode)

		assert.True(t, payment.Status.IsPaid())
		assert.Equal(t, transactionCode, *payment.TransactionCode)
		assert.NotNil(t, payment.PaidAt)
		assert.False(t, payment.PaidAt.IsZero())
	})

	t.Run("should not change already paid payment", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PAID
		method := constants.PIX_PAYMENT_METHOD
		originalTxn := "txn-original"
		now := time.Now()
		paidAt := now.Add(-1 * time.Hour)

		payment, _ := NewPayment(id, orderID, amount, status, method, &originalTxn, nil, &paidAt, now)
		newTransactionCode := "txn-new"

		payment.MarkAsPaid(newTransactionCode)

		assert.True(t, payment.Status.IsPaid())
		assert.Equal(t, originalTxn, *payment.TransactionCode)
		assert.Equal(t, paidAt, *payment.PaidAt)
	})

	t.Run("should not mark failed payment as paid", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_FAILED
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, _ := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)
		transactionCode := "txn-123"

		payment.MarkAsPaid(transactionCode)

		assert.True(t, payment.Status.IsFailed())
		assert.Nil(t, payment.TransactionCode)
		assert.Nil(t, payment.PaidAt)
	})
}

func TestPayment_MarkAsFailed(t *testing.T) {
	t.Run("should mark pending payment as failed", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, _ := NewPaymentDefault(id, orderID, amount, method, qrData)

		payment.MarkAsFailed()

		assert.True(t, payment.Status.IsFailed())
	})

	t.Run("should mark paid payment as failed", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PAID
		method := constants.PIX_PAYMENT_METHOD
		transactionCode := "txn-123"
		now := time.Now()
		paidAt := now

		payment, _ := NewPayment(id, orderID, amount, status, method, &transactionCode, nil, &paidAt, now)

		payment.MarkAsFailed()

		assert.True(t, payment.Status.IsFailed())
	})

	t.Run("should keep failed status when already failed", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_FAILED
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, _ := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		payment.MarkAsFailed()

		assert.True(t, payment.Status.IsFailed())
	})
}

func TestPayment_SetQrCode(t *testing.T) {
	t.Run("should set QR code for PIX payment", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, _ := NewPaymentDefault(id, orderID, amount, method, qrData)
		newQRCode := "00020126580014br.gov.bcb.pix"

		err := payment.SetQrCode(newQRCode)

		assert.NoError(t, err)
		assert.NotNil(t, payment.QRCode)
		assert.Equal(t, newQRCode, payment.QRCode.Value())
	})

	t.Run("should update existing QR code", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		oldQRCode := "00020126580014br.gov.bcb.pix.old"

		payment, _ := NewPaymentDefault(id, orderID, amount, method, &oldQRCode)
		newQRCode := "00020126580014br.gov.bcb.pix.new"

		err := payment.SetQrCode(newQRCode)

		assert.NoError(t, err)
		assert.NotNil(t, payment.QRCode)
		assert.Equal(t, newQRCode, payment.QRCode.Value())
	})

	t.Run("should return error for empty QR code", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		method := constants.PIX_PAYMENT_METHOD
		var qrData *string

		payment, _ := NewPaymentDefault(id, orderID, amount, method, qrData)

		err := payment.SetQrCode("")

		assert.Error(t, err)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error when setting QR code for non-PIX payment", func(t *testing.T) {
		id := "payment-123"
		orderID := "order-123"
		amount := 100.50
		status := constants.PAYMENT_STATUS_PENDING
		method := constants.PIX_PAYMENT_METHOD
		now := time.Now()

		payment, _ := NewPayment(id, orderID, amount, status, method, nil, nil, nil, now)

		qrCode := "00020126580014br.gov.bcb.pix"
		err := payment.SetQrCode(qrCode)

		assert.NoError(t, err)
		assert.NotNil(t, payment.QRCode)
		assert.Equal(t, qrCode, payment.QRCode.Value())
	})

	t.Run("should validate PIX requirement comprehensively", func(t *testing.T) {
		pixMethod, _ := value_objects.NewMethod(constants.PIX_PAYMENT_METHOD)
		payment := createPaymentWithMethod(pixMethod)

		qrCode := "00020126580014br.gov.bcb.pix"
		err := payment.SetQrCode(qrCode)

		assert.NoError(t, err)
		assert.NotNil(t, payment.QRCode)
	})
}
