package exceptions

type InvalidPaymentDataException struct {
	Message string
}

func (e *InvalidPaymentDataException) Error() string {
	if e.Message == "" {
		return "invalid payment data"
	}
	return e.Message
}

type PaymentNotFoundException struct {
	Message string
}

func (e *PaymentNotFoundException) Error() string {
	if e.Message == "" {
		return "payment not found"
	}
	return e.Message
}
