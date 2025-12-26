package mocks

import (
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/mock"
)

// MockPaymentGateway é um mock da interface IPaymentGateway
type MockPaymentGateway struct {
	mock.Mock
}

func (m *MockPaymentGateway) CreatePIXBilling(pixBilling dto.CreatePIXBillingDTO) (dto.PIXBillingResultDTO, error) {
	args := m.Called(pixBilling)
	return args.Get(0).(dto.PIXBillingResultDTO), args.Error(1)
}
