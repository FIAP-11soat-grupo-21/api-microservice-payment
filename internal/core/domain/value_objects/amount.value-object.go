package value_objects

import "payment_microservice/internal/core/domain/exceptions"

type Amount struct {
	value float64
}

func NewAmount(value float64) (Amount, error) {
	if value < 0 {
		return Amount{}, &exceptions.InvalidPaymentDataException{
			Message: "amount cannot be negative",
		}
	}

	return Amount{value: value}, nil
}

func (a Amount) Value() float64 {
	return a.value
}
