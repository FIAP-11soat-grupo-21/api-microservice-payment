package use_cases

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/exceptions"
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCreatePaymentUseCase(t *testing.T) {
	t.Run("should create use case with repository and gateway", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)

		assert.NotNil(t, useCase)
		assert.Equal(t, mockRepo, useCase.repository)
		assert.Equal(t, mockGateway, useCase.gateway)
	})
}

func TestCreatePaymentUseCase_Execute(t *testing.T) {
	t.Run("should create payment with PIX successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  100.50,
			Method:  "pix",
		}

		pixBillingResult := dto.PIXBillingResultDTO{
			QRData: "00020126580014br.gov.bcb.pix",
		}

		mockGateway.On("CreatePIXBilling", mock.MatchedBy(func(billing dto.CreatePIXBillingDTO) bool {
			return billing.Amount == paymentDTO.Amount && billing.Ctx == ctx
		})).Return(pixBillingResult, nil)

		mockRepo.On("Insert", ctx, mock.MatchedBy(func(payment entities.Payment) bool {
			return payment.OrderID == paymentDTO.OrderID &&
				payment.Amount.Value() == paymentDTO.Amount &&
				payment.QRCode != nil
		})).Return(nil)

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.NoError(t, err)
		assert.Equal(t, paymentDTO.OrderID, payment.OrderID)
		assert.Equal(t, paymentDTO.Amount, payment.Amount.Value())
		assert.NotNil(t, payment.QRCode)
		mockRepo.AssertExpectations(t)
		mockGateway.AssertExpectations(t)
	})

	t.Run("should return error when amount is invalid", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  -10.00,
			Method:  "pix",
		}

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
	})

	t.Run("should return error when method is not PIX", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  100.50,
			Method:  "credit_card",
		}

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
	})

	t.Run("should return error when gateway fails to create PIX billing", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  100.50,
			Method:  "pix",
		}

		gatewayError := errors.New("gateway connection failed")
		mockGateway.On("CreatePIXBilling", mock.Anything).Return(dto.PIXBillingResultDTO{}, gatewayError)

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "Failed to create PIX billing")
		mockGateway.AssertExpectations(t)
	})

	t.Run("should return error when repository fails to insert payment", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  100.50,
			Method:  "pix",
		}

		pixBillingResult := dto.PIXBillingResultDTO{
			QRData: "00020126580014br.gov.bcb.pix",
		}

		mockGateway.On("CreatePIXBilling", mock.Anything).Return(pixBillingResult, nil)
		mockRepo.On("Insert", ctx, mock.Anything).Return(errors.New("database error"))

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		assert.Contains(t, err.Error(), "Failed to insert payment on database")
		mockRepo.AssertExpectations(t)
		mockGateway.AssertExpectations(t)
	})

	t.Run("should return error for invalid payment method", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  100.50,
			Method:  "invalid_method",
		}

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
	})

	t.Run("should return error when OrderID is empty", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "",
			Amount:  100.50,
			Method:  "pix",
		}

		pixBillingResult := dto.PIXBillingResultDTO{
			QRData: "00020126580014br.gov.bcb.pix",
		}

		mockGateway.On("CreatePIXBilling", mock.Anything).Return(pixBillingResult, nil)
		mockRepo.On("Insert", ctx, mock.Anything).Return(nil)

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.NoError(t, err)
		assert.Equal(t, "", payment.OrderID)
		mockRepo.AssertExpectations(t)
		mockGateway.AssertExpectations(t)
	})

	t.Run("should return error when QR code data is empty", func(t *testing.T) {
		mockRepo := new(mocks.MockPaymentRepository)
		mockGateway := new(mocks.MockPaymentGateway)
		ctx := context.Background()

		paymentDTO := dto.CreatePaymentDTO{
			Ctx:     ctx,
			OrderID: "order-123",
			Amount:  100.50,
			Method:  "pix",
		}

		pixBillingResult := dto.PIXBillingResultDTO{
			QRData: "",
		}

		mockGateway.On("CreatePIXBilling", mock.Anything).Return(pixBillingResult, nil)

		useCase := NewCreatePaymentUseCase(mockRepo, mockGateway)
		payment, err := useCase.Execute(paymentDTO)

		assert.Error(t, err)
		assert.Equal(t, entities.Payment{}, payment)
		assert.IsType(t, &exceptions.InvalidPaymentDataException{}, err)
		mockGateway.AssertExpectations(t)
	})
}
