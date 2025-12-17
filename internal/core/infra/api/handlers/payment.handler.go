package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"payment_microservice/internal/common/config/constants"
	"payment_microservice/internal/core/application/controllers"
	"payment_microservice/internal/core/application/dtos"
	payment_factories "payment_microservice/internal/core/factories"
	"payment_microservice/internal/core/infra/api/schemas"
)

type PaymentHandler struct {
	paymentController controllers.PaymentController
}

func NewPaymentHandler() *PaymentHandler {
	paymentDataSource := payment_factories.NewPaymentDataSource()
	paymentProvider := payment_factories.NewPaymentProvider()
	paymentController := controllers.NewPaymentController(paymentDataSource, paymentProvider)

	return &PaymentHandler{
		paymentController: *paymentController,
	}
}

// @Summary Create Pix Billing
// @Tags Payments
// @Accept json
// @Router /payments/pix [post]
func (ph *PaymentHandler) CreatePixBilling(ctx *gin.Context) {
	var createPixBillingBody schemas.CreatePixBillingSchema

	if err := ctx.ShouldBindJSON(&createPixBillingBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentDTO := dtos.CreatePaymentDTO{
		Ctx:     ctx,
		OrderID: createPixBillingBody.OrderID,
		Amount:  createPixBillingBody.Amount,
		Method:  constants.PIX_PAYMENT_METHOD,
	}

	res, err := ph.paymentController.CreatePayment(paymentDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

// @Summary Confirm Payment
// @Tags Payments
// @Accept json
// @Router /payments/webhook [post]
func (ph *PaymentHandler) ConfirmPayment(ctx *gin.Context) {
	var confirmPaymentBody schemas.ConfirmPaymentSchema

	if err := ctx.ShouldBindJSON(&confirmPaymentBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventDTO := dtos.WebhookEventDTO{
		ID:          confirmPaymentBody.ID,
		LiveMode:    confirmPaymentBody.LiveMode,
		Type:        confirmPaymentBody.Type,
		DateCreated: confirmPaymentBody.DateCreated,
		UserID:      confirmPaymentBody.UserID,
		APIVersion:  confirmPaymentBody.APIVersion,
		Action:      confirmPaymentBody.Action,
		OrderID:     confirmPaymentBody.Data.ID,
	}

	err := ph.paymentController.ConfirmPayment(ctx, eventDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// @Summary Get a payment by OrderID
// @Tags Payments
// @Produce json
// @Param orderID path string true "Order ID"
// @Router /payments/order/{orderID} [get]
func (h *PaymentHandler) FindByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	payment, err := h.paymentController.GetByOrderID(orderID)

	if err != nil {
		ctx.Error(err)
		return
	}

	response := schemas.PaymentResponseSchema{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount,
		Status:          payment.Status,
		Method:          payment.Method,
		TransactionCode: payment.TransactionCode,
		QRCode:          payment.QRCode,
		CreatedAt:       payment.CreatedAt,
		PaidAt:          payment.PaidAt,
	}

	ctx.JSON(http.StatusOK, response)
}
