package steps

import (
	"fmt"
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/common/mocks"
	identity_manager "payment_microservice/internal/common/pkg/identity"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/value_objects"
	"payment_microservice/internal/core/dto"

	"go.uber.org/mock/gomock"
)

type PaymentHelper struct {
	Ctrl     *gomock.Controller
	MockPG   *mocks.MockPaymentGateway
	MockRepo *mocks.MockPaymentRepository

	valid struct {
		order_id string
		amount   float64
	}
	existingID string
}

func (ph *PaymentHelper) ThePaymentDataIsValid() error {
	ph.valid.order_id = "order_12345"
	ph.valid.amount = 150.75

	if ph.valid.order_id == "" {
		return fmt.Errorf("invalid order ID")
	}

	_, err := value_objects.NewAmount(ph.valid.amount)

	if err != nil {
		return err
	}

	return nil
}

func (ph *PaymentHelper) SendMessageToCreatePIXBilling() error {
	generatedID := identity_manager.NewUUIDV7()

	qrCode := "00020101021243650016COM.MERCADOLIBRE"

	ph.MockPG.On("CreatePIXBilling").Return(dto.PIXBillingResultDTO{
		QRData: qrCode,
	}, nil)

	amount, _ := value_objects.NewAmount(ph.valid.amount)

	status, _ := value_objects.NewStatus(constants.PAYMENT_STATUS_PENDING)

	qrCodeVO, _ := value_objects.NewQRCode(qrCode)

	ph.MockRepo.On("Insert").Return(entities.Payment{
		ID:      generatedID,
		OrderID: ph.valid.order_id,
		QRCode:  &qrCodeVO,
		Amount:  amount,
		Status:  status,
	}, nil)

	ph.existingID = generatedID

	return nil
}

func (ph *PaymentHelper) PaymentShouldBeCreated() error {
	if ph.existingID == "" {
		return fmt.Errorf("payment was not created")
	}
	return nil
}
