package handlers

import (
	"net/http"
	"payment_microservice/internal/adapters/driver/api/schemas"
	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"
	"payment_microservice/internal/core/factory"
	"payment_microservice/internal/core/use_cases"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	repository ports.IPaymentRepository
	gateway    ports.IPaymentGateway
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		repository: factory.NewPaymentRepository(),
		gateway:    factory.NewPaymentGateway(),
	}
}

// @Summary Create Pix Billing
// @Tags Payments
// @Accept json
// @Produce json
// @Param order body schemas.CreatePixBillingSchema true "Pix billing to create"
// @Success 201 {object} schemas.PaymentResponseSchema
// @Failure 500
// @Router /payments/pix [post]
func (ph *PaymentHandler) CreatePixBilling(ctx *gin.Context) {
	var createPixBillingBody schemas.CreatePixBillingSchema

	if err := ctx.ShouldBindJSON(&createPixBillingBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentDTO := dto.CreatePaymentDTO{
		Ctx:     ctx,
		OrderID: createPixBillingBody.OrderID,
		Amount:  createPixBillingBody.Amount,
		Method:  constants.PIX_PAYMENT_METHOD,
	}

	createPaymentUseCase := use_cases.NewCreatePaymentUseCase(ph.repository, ph.gateway)

	newPayment, err := createPaymentUseCase.Execute(paymentDTO)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var qrCode *string
	if newPayment.QRCode != nil {
		qrCodeValue := newPayment.QRCode.Value()
		qrCode = &qrCodeValue
	}

	responseBody := schemas.PaymentResponseSchema{
		ID:              newPayment.ID,
		OrderID:         newPayment.OrderID,
		Amount:          newPayment.Amount.Value(),
		Status:          newPayment.Status.Value(),
		Method:          newPayment.Method.Value(),
		TransactionCode: newPayment.TransactionCode,
		QRCode:          qrCode,
		CreatedAt:       newPayment.CreatedAt,
		PaidAt:          newPayment.PaidAt,
	}

	ctx.JSON(http.StatusCreated, responseBody)
}

// @Summary Confirm Payment
// @Tags Payments
// @Accept json
// @Produce json
// @Param webhook body schemas.ConfirmPaymentSchema true "Payment confirmation webhook"
// @Success 204
// @Failure 500
// @Router /payments/webhook [post]
func (ph *PaymentHandler) ConfirmPayment(ctx *gin.Context) {
	var confirmPaymentBody schemas.ConfirmPaymentSchema

	if err := ctx.ShouldBindJSON(&confirmPaymentBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventDTO := dto.WebhookEventDTO{
		ID:          confirmPaymentBody.ID,
		LiveMode:    confirmPaymentBody.LiveMode,
		Type:        confirmPaymentBody.Type,
		DateCreated: confirmPaymentBody.DateCreated,
		UserID:      confirmPaymentBody.UserID,
		APIVersion:  confirmPaymentBody.APIVersion,
		Action:      confirmPaymentBody.Action,
		OrderID:     confirmPaymentBody.Data.ID,
	}

	confirmPaymentUseCase := use_cases.NewConfirmPaymentUseCase(ph.repository)

	err := confirmPaymentUseCase.Execute(ctx, eventDTO)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// @Summary Get a payment by OrderID
// @Tags Payments
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} schemas.PaymentResponseSchema
// @Failure 500
// @Router /payments/order/{orderID} [get]
func (h *PaymentHandler) FindByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	findPaymentUseCase := use_cases.NewFindPaymentByOrderIDUseCase(h.repository)

	payment, err := findPaymentUseCase.Execute(ctx, orderID)

	if err != nil {
		ctx.Error(err)
		return
	}

	var qrCode *string
	if payment.QRCode != nil {
		qrCodeValue := payment.QRCode.Value()
		qrCode = &qrCodeValue
	}

	responseBody := schemas.PaymentResponseSchema{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount.Value(),
		Status:          payment.Status.Value(),
		Method:          payment.Method.Value(),
		TransactionCode: payment.TransactionCode,
		QRCode:          qrCode,
		CreatedAt:       payment.CreatedAt,
		PaidAt:          payment.PaidAt,
	}

	ctx.JSON(http.StatusOK, responseBody)
}
