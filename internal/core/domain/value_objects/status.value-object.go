package value_objects

import (
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/exceptions"
	"slices"
)

var AllowedPaymentStatus = []string{
	constants.PAYMENT_STATUS_FAILED,
	constants.PAYMENT_STATUS_PENDING,
	constants.PAYMENT_STATUS_PAID,
	constants.PAYMENT_STATUS_REFUNDED,
}

type Status struct {
	value string
}

func NewStatus(value string) (Status, error) {
	if !isValidStatus(value) {
		return Status{}, &exceptions.InvalidPaymentDataException{
			Message: "invalid payment status",
		}
	}

	return Status{value: value}, nil
}

func NewStatusDefault() Status {
	return Status{value: constants.PAYMENT_STATUS_PENDING}
}

func (s Status) Value() string {
	return s.value
}

func isValidStatus(status string) bool {
	return slices.Contains(AllowedPaymentStatus, status)
}

func (s Status) IsPending() bool {
	return s.value == constants.PAYMENT_STATUS_PENDING
}

func (s *Status) SetPending() {
	s.value = constants.PAYMENT_STATUS_PENDING
}

func (s Status) IsPaid() bool {
	return s.value == constants.PAYMENT_STATUS_PAID
}

func (s *Status) SetPaid() {
	s.value = constants.PAYMENT_STATUS_PAID
}

func (s Status) IsFailed() bool {
	return s.value == constants.PAYMENT_STATUS_FAILED
}

func (s *Status) SetFailed() {
	s.value = constants.PAYMENT_STATUS_FAILED
}

func (s Status) IsRefunded() bool {
	return s.value == constants.PAYMENT_STATUS_REFUNDED
}

func (s *Status) SetRefunded() {
	s.value = constants.PAYMENT_STATUS_REFUNDED
}
