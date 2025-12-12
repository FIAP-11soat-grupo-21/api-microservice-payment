package value_objects

import "payment_microservice/internal/core/domain/exceptions"

type QRCode struct {
	url string
}

func NewQRCode(url string) (QRCode, error) {
	if url == "" {
		return QRCode{}, &exceptions.InvalidPaymentDataException{
			Message: "QR code URL and copy-paste code cannot be empty",
		}
	}

	return QRCode{
		url: url,
	}, nil

}

func (q QRCode) Value() string {
	return q.url
}
