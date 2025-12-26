package use_cases

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRefoundPaymentUseCase(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)

	uc := NewRefoundPaymentUseCase(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.repository)
}

func TestRefoundPaymentUseCase_Execute_FindByOrderIDError(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRefoundPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	expectedErr := errors.New("db error")

	repo.On("FindByOrderID", ctx, orderID).
		Return(entities.Payment{}, expectedErr).
		Once()

	err := uc.Execute(ctx, orderID)

	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

func TestRefoundPaymentUseCase_Execute_PaymentNotFound(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRefoundPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	emptyPayment := entities.Payment{}

	repo.On("FindByOrderID", ctx, orderID).
		Return(emptyPayment, nil).
		Once()

	err := uc.Execute(ctx, orderID)

	assert.IsType(t, &exceptions.PaymentNotFoundException{}, err)
	repo.AssertExpectations(t)
}

func TestRefoundPaymentUseCase_Execute_Success(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRefoundPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	// cria um payment não vazio para não cair em PaymentNotFound
	payment := entities.Payment{ID: "payment-1", OrderID: orderID}

	repo.On("FindByOrderID", ctx, orderID).
		Return(payment, nil).
		Once()

	repo.On("Update", ctx, mock.MatchedBy(func(p entities.Payment) bool {
		// garante que o pagamento passado para Update é o mesmo pedido e foi marcado como refundado
		return p.OrderID == orderID && p.Status.IsRefunded()
	})).
		Return(nil).
		Once()

	err := uc.Execute(ctx, orderID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRefoundPaymentUseCase_Execute_UpdateError(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRefoundPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	payment := entities.Payment{ID: "payment-1", OrderID: orderID}

	expectedErr := errors.New("update failed")

	repo.On("FindByOrderID", ctx, orderID).
		Return(payment, nil).
		Once()

	repo.On("Update", ctx, mock.AnythingOfType("entities.Payment")).
		Return(expectedErr).
		Once()

	err := uc.Execute(ctx, orderID)

	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}
