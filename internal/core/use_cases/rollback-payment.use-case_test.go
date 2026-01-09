package use_cases

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestNewRollbackPaymentUseCase(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)

	uc := NewRollbackPaymentUseCase(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.repository)
}

func TestRollbackPaymentUseCase_Execute_FindByOrderIDError(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRollbackPaymentUseCase(repo)

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

func TestRollbackPaymentUseCase_Execute_PaymentNotFound(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRollbackPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	emptyPayment := entities.Payment{}

	repo.On("FindByOrderID", ctx, orderID).
		Return(emptyPayment, nil).
		Once()

	err := uc.Execute(ctx, orderID)

	// Deve retornar sem erro, pois não há pagamento para reverter
	assert.NoError(t, err)
}

func TestRollbackPaymentUseCase_Execute_Success(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRollbackPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	// cria um payment não vazio para não cair em PaymentNotFound
	payment := entities.Payment{ID: "payment-1", OrderID: orderID}

	repo.On("FindByOrderID", ctx, orderID).
		Return(payment, nil).
		Once()

	repo.On("Delete", ctx, payment.ID).
		Return(nil).
		Once()

	err := uc.Execute(ctx, orderID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRollbackPaymentUseCase_Execute_UpdateError(t *testing.T) {
	repo := new(mocks.MockPaymentRepository)
	uc := NewRollbackPaymentUseCase(repo)

	ctx := context.Background()
	orderID := "order-123"

	payment := entities.Payment{ID: "payment-1", OrderID: orderID}

	expectedErr := errors.New("delete failed")

	repo.On("FindByOrderID", ctx, orderID).
		Return(payment, nil).
		Once()

	repo.On("Delete", ctx, payment.ID).
		Return(expectedErr).
		Once()

	err := uc.Execute(ctx, orderID)

	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}
