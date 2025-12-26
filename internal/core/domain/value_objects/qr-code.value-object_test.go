package value_objects

import (
	"testing"

	"payment_microservice/internal/core/domain/exceptions"

	"github.com/stretchr/testify/assert"
)

func TestNewQRCode(t *testing.T) {
	t.Run("should create QR code with valid URL", func(t *testing.T) {
		url := "00020126580014br.gov.bcb.pix"

		qrCode, err := NewQRCode(url)

		assert.NoError(t, err)
		assert.Equal(t, url, qrCode.Value())
	})

	t.Run("should create QR code with complex PIX code", func(t *testing.T) {
		url := "00020126580014br.gov.bcb.pix013603206e1e-cd05-4622-ad93-21406e17e0475204000053039865802BR5913FULANO DE TAL6008BRASILIA62070503***63041D3D"

		qrCode, err := NewQRCode(url)

		assert.NoError(t, err)
		assert.Equal(t, url, qrCode.Value())
	})

	t.Run("should return error when URL is empty", func(t *testing.T) {
		qrCode, err := NewQRCode("")

		assert.Error(t, err)
		assert.Equal(t, QRCode{}, qrCode)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "QR code URL and copy-paste code cannot be empty")
	})

	t.Run("should create QR code with minimum valid string", func(t *testing.T) {
		url := "a"

		qrCode, err := NewQRCode(url)

		assert.NoError(t, err)
		assert.Equal(t, url, qrCode.Value())
	})

	t.Run("should create QR code with special characters", func(t *testing.T) {
		url := "https://example.com/qr?code=123&ref=abc"

		qrCode, err := NewQRCode(url)

		assert.NoError(t, err)
		assert.Equal(t, url, qrCode.Value())
	})
}

func TestQRCode_Value(t *testing.T) {
	t.Run("should return the correct value", func(t *testing.T) {
		expectedURL := "00020126580014br.gov.bcb.pix"
		qrCode, _ := NewQRCode(expectedURL)

		result := qrCode.Value()

		assert.Equal(t, expectedURL, result)
	})
}
