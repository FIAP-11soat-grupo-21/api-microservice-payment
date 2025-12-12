package value_objects

import (
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/exceptions"
	"slices"
)

var AllowedPaymentMethods = []string{
	constants.PIX_PAYMENT_METHOD,
}

type Method struct {
	value string
}

func NewMethod(value string) (Method, error) {
	if !isValidMethod(value) {
		return Method{}, &exceptions.InvalidPaymentDataException{
			Message: "invalid payment method",
		}
	}

	return Method{value: value}, nil
}

func (m Method) Value() string {
	return m.value
}

func isValidMethod(method string) bool {
	return slices.Contains(AllowedPaymentMethods, method)
}

func (m Method) IsPix() bool {
	return m.value == constants.PIX_PAYMENT_METHOD
}
